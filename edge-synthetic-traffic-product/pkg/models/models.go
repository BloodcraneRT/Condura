package models

import "time"

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
	TaskTypePing       TaskType = "ping"
	TaskTypeHTTP       TaskType = "http"
	TaskTypeTCP        TaskType = "tcp"
	TaskTypeUDP        TaskType = "udp"
	TaskTypeDownload   TaskType = "download"
	TaskTypeUpload     TaskType = "upload"
	TaskTypeDNS        TaskType = "dns"
	TaskTypeTraceroute TaskType = "traceroute"
	TaskTypePCAPReplay TaskType = "pcap_replay"
)

// Task represents a synthetic traffic test.
type Task struct {
	ID       string   `json:"id"`
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
	Timestamp   time.Time `json:"timestamp"`
	Success     bool      `json:"success"`
	LatencyMs   float64   `json:"latencyMs"`
	ErrorMsg      string    `json:"errorMsg,omitempty"`
	BytesSent     int64     `json:"bytesSent"`
	BytesRecv     int64     `json:"bytesRecv"`
	PacketLoss    float64   `json:"packetLoss"` // For ping/udp tests
	TTFBMs        float64   `json:"ttfbMs,omitempty"`
	DNSTimeMs     float64   `json:"dnsTimeMs,omitempty"`
	ConnectTimeMs float64   `json:"connectTimeMs,omitempty"`
	Hops          int       `json:"hops,omitempty"`
}

// AggregatedMetrics represents summary statistics of tests over time.
type AggregatedMetrics struct {
	TotalTests      int64   `json:"totalTests"`
	TotalSuccesses  int64   `json:"totalSuccesses"`
	TotalFailures   int64   `json:"totalFailures"`
	AvgLatencyMs    float64 `json:"avgLatencyMs"`
	TotalBytesSent  int64   `json:"totalBytesSent"`
	TotalBytesRecv  int64   `json:"totalBytesRecv"`
}
