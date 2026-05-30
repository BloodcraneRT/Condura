package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunDownloadTest(t *testing.T) {
	tests := []struct {
		name          string
		handler       http.HandlerFunc
		expectedBytes int64
		expectError   bool
	}{
		{
			name: "Success",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("hello world"))
			},
			expectedBytes: 11,
			expectError:   false,
		},
		{
			name: "HTTP Error 404",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			expectedBytes: 0,
			expectError:   true,
		},
		{
			name: "HTTP Error 500",
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedBytes: 0,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			bytes, err := runDownloadTest(server.URL)
			if (err != nil) != tt.expectError {
				t.Errorf("runDownloadTest() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if bytes != tt.expectedBytes {
				t.Errorf("runDownloadTest() bytes = %v, expected %v", bytes, tt.expectedBytes)
			}
		})
	}

	t.Run("Connection Error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		url := server.URL
		server.Close() // Close immediately to cause connection error

		_, err := runDownloadTest(url)
		if err == nil {
			t.Error("runDownloadTest() expected error on closed server, got nil")
		}
	})
}
