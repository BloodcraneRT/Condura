package main

import (
	"testing"
	"time"

	"github.com/edgesynth/edgesynth/pkg/models"
)

func TestDatabaseOperations(t *testing.T) {
	t.Skip("Skipping postgres DB test for now as we don't have a local postgres instance running in the test environment.")

	// Use an in-memory database for testing
	db, err := initDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to init db: %v", err)
	}
	defer db.Close()

	// Test Create Task
	task := models.Task{
		Type:    models.TaskTypeHTTP,
		Target:  "https://google.com",
		Enabled: true,
	}
	err = db.CreateTask(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Test Get Tasks
	tasks, err := db.GetEnabledTasks()
	if err != nil {
		t.Fatalf("Failed to get tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}

	// Test Insert Result
	res := models.TaskResult{
		TaskID:    tasks[0].ID,
		Timestamp: time.Now(),
		Success:   true,
		LatencyMs: 42.0,
	}
	err = db.InsertResult(res)
	if err != nil {
		t.Fatalf("Failed to insert result: %v", err)
	}

	// Test Get Results
	results, err := db.GetRecentResults(10)
	if err != nil {
		t.Fatalf("Failed to get results: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}
}
