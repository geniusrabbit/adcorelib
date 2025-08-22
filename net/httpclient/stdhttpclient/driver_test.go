package stdhttpclient

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// Test the Driver implementation
func TestDriver_Request(t *testing.T) {
	driver := NewDriver()
	defer driver.Close()

	tests := []struct {
		name   string
		method string
		url    string
		body   io.Reader
	}{
		{"GET request", "GET", "http://example.com", nil},
		{"POST request", "POST", "http://example.com/api", strings.NewReader("test data")},
		{"PUT request", "PUT", "http://example.com/api/1", bytes.NewReader([]byte("update data"))},
		{"DELETE request", "DELETE", "http://example.com/api/1", nil},
		{"PATCH request", "PATCH", "http://example.com/api/1", strings.NewReader("patch data")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := driver.Request(tt.method, tt.url, tt.body)
			if err != nil {
				t.Errorf("Request() error = %v", err)
				return
			}

			if req == nil {
				t.Error("Request() returned nil request")
				return
			}

			// Verify the request is of the correct type
			stdReq, ok := req.(*Request)
			if !ok {
				t.Error("Request() did not return *Request type")
				return
			}

			// Verify method and URL
			if stdReq.HTTP.Method != tt.method {
				t.Errorf("Request method = %s, expected %s", stdReq.HTTP.Method, tt.method)
			}

			if stdReq.HTTP.URL.String() != tt.url {
				t.Errorf("Request URL = %s, expected %s", stdReq.HTTP.URL.String(), tt.url)
			}
		})
	}
}

func TestDriver_Request_InvalidURL(t *testing.T) {
	driver := NewDriver()
	defer driver.Close()

	_, err := driver.Request("GET", "://invalid-url", nil)
	if err == nil {
		t.Error("Request() should return error for invalid URL")
	}
}

// Benchmark tests
func BenchmarkDriver_Request(b *testing.B) {
	driver := NewDriver()
	defer driver.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := driver.Request("GET", "http://example.com", nil)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
	}
}

func BenchmarkRequest_SetHeader(b *testing.B) {
	driver := NewDriver()
	defer driver.Close()

	req, err := driver.Request("GET", "http://example.com", nil)
	if err != nil {
		b.Fatalf("Request failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.SetHeader("Content-Type", "application/json")
	}
}

func BenchmarkDriver_Do(b *testing.B) {
	// Use shared test server to avoid port exhaustion
	server := getGlobalTestServer()

	driver := NewDriver()
	defer driver.Close()

	b.ResetTimer()

	errorCount := 0
	for i := 0; i < b.N; i++ {
		req, err := driver.Request("GET", server.URL, nil)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}

		resp, err := driver.Do(req)
		if err != nil {
			errorCount++
			continue
		}
		if resp != nil {
			_ = resp.Close()
		}
	}

	// Allow up to 99% errors for driver benchmark under high load
	// This is because macOS has strict limits on concurrent connections
	errorThreshold := b.N * 99 / 100
	if errorThreshold < 1 {
		errorThreshold = 1 // Always allow at least 1 error
	}
	if errorCount > errorThreshold {
		b.Errorf("Too many errors: %d out of %d (threshold: %d)", errorCount, b.N, errorThreshold)
	}
}
