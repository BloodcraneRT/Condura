package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/edgesynth/edgesynth/pkg/api"
	"github.com/edgesynth/edgesynth/pkg/models"
)

var (
	dbPath   string
	port     int
	database *DB
)

func main() {
	flag.StringVar(&dbPath, "db", "edgesynth.db", "Path to SQLite database")
	flag.IntVar(&port, "port", 8080, "Port to run server on")
	flag.Parse()

	log.Printf("Starting EdgeSynth Control Plane Server on port %d...", port)

	var err error
	database, err = initDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// API Routes for Agents
	http.HandleFunc(api.RouteRegisterAgent, handleRegister)
	http.HandleFunc(api.RouteHeartbeat, handleHeartbeat)
	http.HandleFunc(api.RouteGetTasks, handleGetTasks)
	http.HandleFunc(api.RouteSubmitResults, handleSubmitResults)

	// API Routes for UI
	http.HandleFunc("/api/v1/ui/agents", handleUIAgents)
	http.HandleFunc("/api/v1/ui/results", handleUIResults)
	http.HandleFunc("/api/v1/ui/tasks", handleUITasks) // for creating tasks
	http.HandleFunc("/api/v1/ui/sources", handleUISources) // for getting/creating sources

	// Data endpoints for synthetic load
	http.HandleFunc("/api/v1/testdata/download", handleDownloadData)
	http.HandleFunc("/api/v1/testdata/upload", handleUploadData)

	// Serve Static UI
	// Fallback to searching relative path if running from root dir
	uiDir := "./server/ui-dist"
	if _, err := os.Stat(uiDir); os.IsNotExist(err) {
		uiDir = "./ui-dist"
	}
	fs := http.FileServer(http.Dir(uiDir))
	http.Handle("/", fs)

	// Start server
	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, enableCORS(http.DefaultServeMux)))
}

// enableCORS is a simple middleware to allow cross-origin requests for UI dev
func enableCORS(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req api.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ip := r.RemoteAddr
	id, err := database.RegisterAgent(req, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Registered new agent: %s (%s)", req.Hostname, id)
	json.NewEncoder(w).Encode(api.RegisterResponse{AgentID: id})
}

func handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req api.HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.UpdateHeartbeat(req.AgentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleGetTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	agentID := r.URL.Query().Get("agentId")
	if agentID == "" {
		http.Error(w, "Missing agentId", http.StatusBadRequest)
		return
	}

	tasks, err := database.GetTasksForAgent(agentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	json.NewEncoder(w).Encode(tasks)
}

func handleSubmitResults(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var res models.TaskResult
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := database.InsertResult(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UI Handlers

func handleUIAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := database.GetAllAgents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if agents == nil {
		agents = []models.Agent{}
	}
	json.NewEncoder(w).Encode(agents)
}

func handleUIResults(w http.ResponseWriter, r *http.Request) {
	results, err := database.GetRecentResults(100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if results == nil {
		results = []models.TaskResult{}
	}
	json.NewEncoder(w).Encode(results)
}

func handleUITasks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var t models.Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := database.CreateTask(t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// Data Endpoint Handlers

// handleDownloadData generates synthetic data of a requested size and format
func handleDownloadData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sizeStr := r.URL.Query().Get("sizeBytes")
	sizeBytes, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil || sizeBytes <= 0 {
		sizeBytes = 1024 * 1024 // Default to 1MB
	}

	// Limit to 5GB per request to avoid blowing up the server
	if sizeBytes > 5*1024*1024*1024 {
		sizeBytes = 5 * 1024 * 1024 * 1024
	}

	format := r.URL.Query().Get("format")

	if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":"`))

		// Write dummy A's to fill size
		remaining := sizeBytes - 12
		chunkSize := int64(1024 * 1024)
		chunk := make([]byte, chunkSize)
		for i := range chunk {
			chunk[i] = 'A'
		}
		for remaining > 0 {
			writeSize := remaining
			if writeSize > chunkSize {
				writeSize = chunkSize
			}
			w.Write(chunk[:writeSize])
			remaining -= writeSize
		}
		w.Write([]byte(`"}`))
	} else {
		// Default binary stream
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", strconv.FormatInt(sizeBytes, 10))

		chunkSize := int64(4 * 1024 * 1024) // 4MB chunks
		chunk := make([]byte, chunkSize)
		remaining := sizeBytes

		for remaining > 0 {
			writeSize := remaining
			if writeSize > chunkSize {
				writeSize = chunkSize
			}
			_, err := w.Write(chunk[:writeSize])
			if err != nil {
				// Client likely disconnected
				break
			}
			remaining -= writeSize
		}
	}
}

// handleUploadData accepts streaming uploads of arbitrary sizes
func handleUploadData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Discard the body to simulate upload processing
	bytesRead, err := io.Copy(io.Discard, r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Successfully received %d bytes", bytesRead)))
}

func handleUISources(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		sources, err := database.GetSources()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if sources == nil {
			sources = []models.Source{}
		}
		json.NewEncoder(w).Encode(sources)
		return
	}
	if r.Method == http.MethodPost {
		var s models.Source
		if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := database.CreateSource(s); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
