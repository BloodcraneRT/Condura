package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"os/exec"
	"runtime"
	"time"

	"github.com/edgesynth/edgesynth/pkg/models"
)

// RunTask executes the synthetic task based on its type.
func RunTask(t models.Task) models.TaskResult {
	res := models.TaskResult{
		TaskID:    t.ID,
		Timestamp: time.Now(),
	}

	start := time.Now()
	var err error

	switch t.Type {
	case models.TaskTypeHTTP:
		var ttfb, dnsTime, connectTime float64
		err = runHTTPTest(t.Target, &ttfb, &dnsTime, &connectTime)
		// Store detailed metrics inside ErrorMsg temporarily or extended fields later.
		// For now, we'll append to ErrorMsg if successful just to show it, or we could add columns.
		// We'll add columns in a moment. Let's add them to the result struct.
		res.TTFBMs = ttfb
		res.DNSTimeMs = dnsTime
		res.ConnectTimeMs = connectTime
	case models.TaskTypeTCP:
		err = runTCPTest(t.Target)
	case models.TaskTypeUDP:
		err = runUDPTest(t.Target)
	case models.TaskTypePing:
		err = runTCPTest(t.Target + ":80")
	case models.TaskTypeDownload:
		var bytesRecv int64
		bytesRecv, err = runDownloadTest(t.Target)
		res.BytesRecv = bytesRecv
	case models.TaskTypeUpload:
		var bytesSent int64
		bytesSent, err = runUploadTest(t.Target)
		res.BytesSent = bytesSent
	case models.TaskTypeDNS:
		var dnsTime float64
		err = runDNSTest(t.Target, &dnsTime)
		res.DNSTimeMs = dnsTime
	case models.TaskTypeTraceroute:
		var hops int
		var output string
		hops, output, err = runTracerouteTest(t.Target)
		res.Hops = hops
		if err == nil {
			// Overload ErrorMsg to store traceroute output if successful for UI display
			res.ErrorMsg = output
		}
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

func runHTTPTest(target string, ttfb, dnsTime, connectTime *float64) error {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		return err
	}

	var dnsStart, dnsDone, connStart, connDone, gotFirstByte time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { dnsDone = time.Now() },
		ConnectStart: func(_, _ string) { connStart = time.Now() },
		ConnectDone: func(net, addr string, err error) { connDone = time.Now() },
		GotFirstResponseByte: func() { gotFirstByte = time.Now() },
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !dnsStart.IsZero() && !dnsDone.IsZero() {
		*dnsTime = float64(dnsDone.Sub(dnsStart).Milliseconds())
	}
	if !connStart.IsZero() && !connDone.IsZero() {
		*connectTime = float64(connDone.Sub(connStart).Milliseconds())
	}
	if !gotFirstByte.IsZero() {
		*ttfb = float64(gotFirstByte.Sub(connDone).Milliseconds())
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}
	return nil
}

func runDNSTest(target string, dnsTime *float64) error {
	start := time.Now()
	// Using the default resolver
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := net.DefaultResolver.LookupHost(ctx, target)
	*dnsTime = float64(time.Since(start).Milliseconds())
	return err
}

func runTracerouteTest(target string) (int, string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("tracert", "-d", "-h", "15", "-w", "500", target)
	} else {
		cmd = exec.Command("traceroute", "-n", "-m", "15", "-w", "1", target)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	output := out.String()

	// Hop count parsing is complex and OS-dependent.
	// For this edge tool, we'll store the raw output for the user to view.
	// Assume successful execution means path traced.
	return 0, output, err
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
	// A 60-second timeout allows for downloads up to ~2GB to complete depending on network speed
	client := http.Client{Timeout: 60 * time.Second}
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
	// A 60-second timeout allows for slow uploads
	client := http.Client{Timeout: 60 * time.Second}

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
