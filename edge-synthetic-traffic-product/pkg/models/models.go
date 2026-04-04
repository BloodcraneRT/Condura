package models

import "time"

// Agent represents a synthetic traffic agent running at the edge.
type Agent struct {
	ID        string    `json:"id"`
	Hostname  string    `json:"hostname"`
	OS        string    `json:"os"`
	IP        string    `json:"ip"`
	Status    string    `json:"status"` // online, offline
	LastSeen  time.Time `json:"lastSeen"`
}

// Source represents a target for synthetic traffic (e.g. an IP or URL).
type Source struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Target    string `json:"target"` // The IP or URL
	IsDefault bool   `json:"isDefault"`
}

// TaskType represents the kind of synthetic test to run.
type TaskType string

const (
	TaskTypePing TaskType = "ping"
	TaskTypeHTTP TaskType = "http"
	TaskTypeTCP  TaskType = "tcp"
	TaskTypeUDP  TaskType = "udp"
)

// Task represents a synthetic traffic test assigned to an agent.
type Task struct {
	ID       string   `json:"id"`
	AgentID  string   `json:"agentId"`
	Type     TaskType `json:"type"`
	Target   string   `json:"target"`             // e.g., "8.8.8.8", "https://example.com"
	Interval int      `json:"interval"`           // in seconds
	Config   string   `json:"config,omitempty"`   // JSON string of specific test config
	Enabled  bool     `json:"enabled"`
}

// TaskResult contains the metrics gathered from a single execution of a Task.
type TaskResult struct {
	ID          string    `json:"id"`
	TaskID      string    `json:"taskId"`
	AgentID     string    `json:"agentId"`
	Timestamp   time.Time `json:"timestamp"`
	Success     bool      `json:"success"`
	LatencyMs   float64   `json:"latencyMs"`
	ErrorMsg    string    `json:"errorMsg,omitempty"`
	BytesSent   int64     `json:"bytesSent,omitempty"`
	BytesRecv   int64     `json:"bytesRecv,omitempty"`
	PacketLoss  float64   `json:"packetLoss,omitempty"` // For ping/udp tests
}
