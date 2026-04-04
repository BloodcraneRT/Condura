package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/edgesynth/edgesynth/pkg/api"
	"github.com/edgesynth/edgesynth/pkg/models"
)

var (
	serverURL string
	agentID   string
)

func main() {
	flag.StringVar(&serverURL, "server", "http://localhost:8080", "Control plane server URL")
	flag.Parse()

	log.Printf("Starting Edge Agent, connecting to %s\n", serverURL)

	err := register()
	if err != nil {
		log.Fatalf("Failed to register agent: %v", err)
	}

	log.Printf("Successfully registered with Agent ID: %s\n", agentID)

	ctx, cancel := contextWithSigterm()
	defer cancel()

	var wg sync.WaitGroup

	// Start Heartbeat
	wg.Add(1)
	go func() {
		defer wg.Done()
		heartbeatLoop(ctx.Done())
	}()

	// Start Task Poller
	wg.Add(1)
	go func() {
		defer wg.Done()
		taskPollerLoop(ctx.Done())
	}()

	<-ctx.Done()
	log.Println("Shutting down agent...")
	wg.Wait()
	log.Println("Agent stopped.")
}

func contextWithSigterm() (contextType, func()) {
	// Simple channel-based context for termination
	done := make(chan struct{})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		close(done)
	}()

	return doneCtx{done}, func() {
		// allow manual cancel if needed, though typically signal driven here
		select {
		case <-done:
		default:
			close(done)
		}
	}
}

type contextType interface {
	Done() <-chan struct{}
}
type doneCtx struct {
	done chan struct{}
}
func (d doneCtx) Done() <-chan struct{} { return d.done }

func register() error {
	hostname, _ := os.Hostname()
	reqData := api.RegisterRequest{
		Hostname: hostname,
		OS:       runtime.GOOS,
	}
	body, _ := json.Marshal(reqData)

	resp, err := http.Post(serverURL+api.RouteRegisterAgent, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status: %s", resp.Status)
	}

	var resData api.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&resData); err != nil {
		return err
	}

	agentID = resData.AgentID
	return nil
}

func heartbeatLoop(stop <-chan struct{}) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			reqData := api.HeartbeatRequest{AgentID: agentID}
			body, _ := json.Marshal(reqData)
			resp, err := http.Post(serverURL+api.RouteHeartbeat, "application/json", bytes.NewBuffer(body))
			if err != nil {
				log.Printf("Heartbeat failed: %v", err)
				continue
			}
			resp.Body.Close()
		case <-stop:
			return
		}
	}
}

func taskPollerLoop(stop <-chan struct{}) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fetchAndRunTasks()
		case <-stop:
			return
		}
	}
}

func fetchAndRunTasks() {
	resp, err := http.Get(fmt.Sprintf("%s%s?agentId=%s", serverURL, api.RouteGetTasks, agentID))
	if err != nil {
		log.Printf("Failed to fetch tasks: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read body for error
		b, _ := io.ReadAll(resp.Body)
		log.Printf("Error fetching tasks, status: %d, body: %s", resp.StatusCode, string(b))
		return
	}

	var tasks []models.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		log.Printf("Failed to decode tasks: %v", err)
		return
	}

	for _, task := range tasks {
		if !task.Enabled {
			continue
		}
		// In a real app, we would schedule these based on Interval.
		// For simplicity, we just run them sequentially on each poll right now.
		go func(t models.Task) {
			res := RunTask(t)
			submitResult(res)
		}(task)
	}
}

func submitResult(res models.TaskResult) {
	body, _ := json.Marshal(res)
	resp, err := http.Post(serverURL+api.RouteSubmitResults, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to submit result for task %s: %v", res.TaskID, err)
		return
	}
	defer resp.Body.Close()
}
