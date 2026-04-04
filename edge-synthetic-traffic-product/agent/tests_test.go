package main

import (
	"net/http"
	"net/http/httptest"
	"net"
	"testing"

	"github.com/edgesynth/edgesynth/pkg/models"
)

func TestRunTask_HTTP_Success(t *testing.T) {
	// Create a dummy HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	task := models.Task{
		ID:      "t1",
		AgentID: "a1",
		Type:    models.TaskTypeHTTP,
		Target:  server.URL,
	}

	res := RunTask(task)
	if !res.Success {
		t.Errorf("Expected success, got error: %s", res.ErrorMsg)
	}
	if res.LatencyMs <= 0 {
		t.Errorf("Expected latency > 0, got %f", res.LatencyMs)
	}
}

func TestRunTask_HTTP_Fail(t *testing.T) {
	task := models.Task{
		ID:      "t2",
		AgentID: "a1",
		Type:    models.TaskTypeHTTP,
		Target:  "http://localhost:12345/nonexistent", // Assuming nothing is listening here
	}

	res := RunTask(task)
	if res.Success {
		t.Errorf("Expected failure for non-existent target")
	}
}

func TestRunTask_TCP_Success(t *testing.T) {
	// Start a dummy TCP server
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}
	defer l.Close()

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	task := models.Task{
		ID:      "t3",
		AgentID: "a1",
		Type:    models.TaskTypeTCP,
		Target:  l.Addr().String(),
	}

	res := RunTask(task)
	if !res.Success {
		t.Errorf("Expected success, got error: %s", res.ErrorMsg)
	}
}
