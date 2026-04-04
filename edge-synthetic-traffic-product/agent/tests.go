package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/edgesynth/edgesynth/pkg/models"
)

// RunTask executes the synthetic task based on its type.
func RunTask(t models.Task) models.TaskResult {
	res := models.TaskResult{
		TaskID:    t.ID,
		AgentID:   t.AgentID,
		Timestamp: time.Now(),
	}

	start := time.Now()
	var err error

	switch t.Type {
	case models.TaskTypeHTTP:
		err = runHTTPTest(t.Target)
	case models.TaskTypeTCP:
		err = runTCPTest(t.Target)
	case models.TaskTypeUDP:
		// Simplified UDP check, since UDP is connectionless
		err = runUDPTest(t.Target)
	case models.TaskTypePing:
		// Fallback to TCP ping if real ICMP needs root
		err = runTCPTest(t.Target + ":80") // Assuming target is an IP, fallback to 80
	case models.TaskTypeDownload:
		var bytesRecv int64
		bytesRecv, err = runDownloadTest(t.Target)
		res.BytesRecv = bytesRecv
	case models.TaskTypeUpload:
		var bytesSent int64
		bytesSent, err = runUploadTest(t.Target)
		res.BytesSent = bytesSent
	default:
		err = fmt.Errorf("unknown task type: %s", t.Type)
	}

	res.LatencyMs = float64(time.Since(start).Milliseconds())

	if err != nil {
		res.Success = false
		res.ErrorMsg = err.Error()
	} else {
		res.Success = true
	}

	return res
}

func runHTTPTest(target string) error {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(target)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}
	return nil
}

func runTCPTest(target string) error {
	conn, err := net.DialTimeout("tcp", target, 5*time.Second)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

func runUDPTest(target string) error {
	conn, err := net.DialTimeout("udp", target, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send a dummy payload
	_, err = conn.Write([]byte("ping"))
	return err
}

func runDownloadTest(target string) (int64, error) {
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(target)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}

	bytesRead, err := io.Copy(io.Discard, resp.Body)
	return bytesRead, err
}

func runUploadTest(target string) (int64, error) {
	client := http.Client{Timeout: 30 * time.Second}

	// Create a dummy 10MB payload
	size := 10 * 1024 * 1024
	payload := make([]byte, size)
	for i := range payload {
		payload[i] = 'B'
	}

	req, err := http.NewRequest("POST", target, bytes.NewReader(payload))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return 0, fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}

	return int64(size), nil
}
