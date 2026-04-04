package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/edgesynth/edgesynth/pkg/models"
)

//go:embed ui-dist/*
var staticFS embed.FS

var (
	connStr  string
	port     int
	database *DB
)

func main() {
	defaultConnStr := os.Getenv("DATABASE_URL")
	if defaultConnStr == "" {
		defaultConnStr = "clickhouse://default:@localhost:9000/default?dial_timeout=10s&read_timeout=20s"
	}
	flag.StringVar(&connStr, "db", defaultConnStr, "ClickHouse connection string")
	flag.IntVar(&port, "port", 8080, "Port to run server on")
	flag.Parse()

	log.Printf("Starting EdgeSynth Control Plane Server on port %d...", port)
	log.Printf("Connecting to Database...")

	var err error
	// Retry connection for 30s to allow docker container startup
	for i := 0; i < 30; i++ {
		database, err = initDB(connStr)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Start internal task runner
	go startTaskRunner(database)

	// API Routes for UI
	http.HandleFunc("/api/v1/ui/metrics", handleUIMetrics)
	http.HandleFunc("/api/v1/ui/results", handleUIResults)
	http.HandleFunc("/api/v1/ui/tasks", handleUITasks) // for creating tasks
	http.HandleFunc("/api/v1/ui/sources", handleUISources) // for getting/creating sources

	// Data endpoints for synthetic load
	http.HandleFunc("/api/v1/testdata/download", handleDownloadData)
	http.HandleFunc("/api/v1/testdata/upload", handleUploadData)

	// Serve Static UI
	// Serve directly from embedded FS
	subFS, err := fs.Sub(staticFS, "ui-dist")
	if err != nil {
		log.Fatalf("Failed to initialize embedded UI: %v", err)
	}
	fileServer := http.FileServer(http.FS(subFS))
	http.Handle("/", fileServer)

	// Start server
	addr := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(addr, enableCORS(http.DefaultServeMux)))
}

var (
	activeTasks     = make(map[string]context.CancelFunc)
	activeTasksLock sync.Mutex
)

func startTaskRunner(db *DB) {
	// Sync loop checks every 5 seconds for new or removed tasks
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		tasks, err := db.GetEnabledTasks()
		if err != nil {
			log.Printf("Failed to get tasks: %v", err)
			continue
		}

		currentTaskIDs := make(map[string]bool)

		activeTasksLock.Lock()
		for _, task := range tasks {
			currentTaskIDs[task.ID] = true
			if _, exists := activeTasks[task.ID]; !exists {
				// Start a new goroutine for this task
				ctx, cancel := context.WithCancel(context.Background())
				activeTasks[task.ID] = cancel
				go runTaskLoop(ctx, db, task)
			}
		}

		// Stop and clean up any tasks that are no longer enabled
		for id, cancel := range activeTasks {
			if !currentTaskIDs[id] {
				cancel()
				delete(activeTasks, id)
			}
		}
		activeTasksLock.Unlock()
	}
}

func runTaskLoop(ctx context.Context, db *DB, t models.Task) {
	// Execute immediately once
	executeTask(db, t)

	interval := t.Interval
	if interval < 1 {
		interval = 10 // Minimum 10 seconds
	}

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			executeTask(db, t)
		case <-ctx.Done():
			return
		}
	}
}

func executeTask(db *DB, t models.Task) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic executing task %s: %v", t.ID, r)
		}
	}()

	res := RunTask(t)
	err := db.InsertResult(res)
	if err != nil {
		log.Printf("Failed to insert result for task %s: %v", t.ID, err)
	}
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

// UI Handlers

func handleUIMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	metrics, err := database.GetAggregatedMetrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(metrics)
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
	if r.Method == http.MethodGet {
		tasks, err := database.GetEnabledTasks()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if tasks == nil {
			tasks = []models.Task{}
		}
		json.NewEncoder(w).Encode(tasks)
		return
	}
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
	if r.Method == http.MethodDelete {
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing task id", http.StatusBadRequest)
			return
		}
		if err := database.DeleteTask(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
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
