package stdhttpzeroclient

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestZeroAllocDriver tests the zero-allocation driver
func TestZeroAllocDriver(t *testing.T) {
	driver := NewDriver()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	// Test basic request
	req, err := driver.Request("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Request creation failed: %v", err)
	}

	resp, err := driver.Do(req)
	if err != nil {
		t.Fatalf("Request execution failed: %v", err)
	}
	defer resp.Close()

	if resp.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode())
	}
}

// TestZeroAllocDriverWithHeaders tests header handling
func TestZeroAllocDriverWithHeaders(t *testing.T) {
	driver := NewDriver()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer token123" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	req, err := driver.Request("GET", server.URL, nil)
	if err != nil {
		t.Fatalf("Request creation failed: %v", err)
	}

	req.SetHeader("Authorization", "Bearer token123")

	resp, err := driver.Do(req)
	if err != nil {
		t.Fatalf("Request execution failed: %v", err)
	}
	defer resp.Close()

	if resp.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode())
	}
}

// TestZeroAllocDriverStringRequest tests string request helper
func TestZeroAllocDriverStringRequest(t *testing.T) {
	driver := NewDriver()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	req, err := driver.StringRequest("POST", server.URL, "test body")
	if err != nil {
		t.Fatalf("String request creation failed: %v", err)
	}

	resp, err := driver.Do(req)
	if err != nil {
		t.Fatalf("Request execution failed: %v", err)
	}
	defer resp.Close()

	if resp.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode())
	}
}

// TestZeroAllocDriverJSONRequest tests JSON request helper
func TestZeroAllocDriverJSONRequest(t *testing.T) {
	driver := NewDriver()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	jsonData := []byte(`{"key":"value"}`)
	req, err := driver.JSONRequest("POST", server.URL, jsonData)
	if err != nil {
		t.Fatalf("JSON request creation failed: %v", err)
	}

	resp, err := driver.Do(req)
	if err != nil {
		t.Fatalf("Request execution failed: %v", err)
	}
	defer resp.Close()

	if resp.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode())
	}
}

// TestZeroAllocDriverFormRequest tests form request helper
func TestZeroAllocDriverFormRequest(t *testing.T) {
	driver := NewDriver()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	formData := url.Values{
		"key1": []string{"value1"},
		"key2": []string{"value2"},
	}

	req, err := driver.FormRequest("POST", server.URL, formData)
	if err != nil {
		t.Fatalf("Form request creation failed: %v", err)
	}

	resp, err := driver.Do(req)
	if err != nil {
		t.Fatalf("Request execution failed: %v", err)
	}
	defer resp.Close()

	if resp.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode())
	}
}

// TestFastRequest tests FastRequest functionality
func TestFastRequest(t *testing.T) {
	req := NewFastRequest("GET", "http://example.com", nil)

	// Test header setting
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Authorization", "Bearer token")

	// Test conversion to HTTP request
	httpReq, err := req.ToHTTPRequest()
	if err != nil {
		t.Fatalf("ToHTTPRequest failed: %v", err)
	}

	if httpReq.Method != "GET" {
		t.Errorf("Expected method GET, got %s", httpReq.Method)
	}

	if httpReq.URL.String() != "http://example.com" {
		t.Errorf("Expected URL http://example.com, got %s", httpReq.URL.String())
	}

	if httpReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", httpReq.Header.Get("Content-Type"))
	}

	if httpReq.Header.Get("Authorization") != "Bearer token" {
		t.Errorf("Expected Authorization Bearer token, got %s", httpReq.Header.Get("Authorization"))
	}

	// Test reset
	req.Reset()
	if req.method != "" || req.url != "" || req.body != nil {
		t.Error("Reset did not clear request properly")
	}
}

// BenchmarkZeroAllocDriver benchmarks the zero-allocation driver
func BenchmarkZeroAllocDriver(b *testing.B) {
	driver := NewDriver()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, err := driver.Request("GET", server.URL, nil)
			if err != nil {
				b.Fatalf("Request creation failed: %v", err)
			}

			resp, err := driver.Do(req)
			if err != nil {
				b.Fatalf("Request execution failed: %v", err)
			}

			resp.Close()
		}
	})
}

// BenchmarkZeroAllocDriverWithHeaders benchmarks with headers
func BenchmarkZeroAllocDriverWithHeaders(b *testing.B) {
	driver := NewDriver()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, err := driver.Request("POST", server.URL, strings.NewReader("test"))
		if err != nil {
			b.Fatalf("Request creation failed: %v", err)
		}

		req.SetHeader("Content-Type", "text/plain")
		req.SetHeader("Authorization", "Bearer token")

		resp, err := driver.Do(req)
		if err != nil {
			b.Fatalf("Request execution failed: %v", err)
		}

		resp.Close()
	}
}

// BenchmarkZeroAllocDriverStringRequest benchmarks string request helper
func BenchmarkZeroAllocDriverStringRequest(b *testing.B) {
	driver := NewDriver()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, err := driver.StringRequest("POST", server.URL, "test body")
		if err != nil {
			b.Fatalf("String request creation failed: %v", err)
		}

		resp, err := driver.Do(req)
		if err != nil {
			b.Fatalf("Request execution failed: %v", err)
		}

		resp.Close()
	}
}

// BenchmarkZeroAllocDriverJSONRequest benchmarks JSON request helper
func BenchmarkZeroAllocDriverJSONRequest(b *testing.B) {
	driver := NewDriver()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	jsonData := []byte(`{"key":"value"}`)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, err := driver.JSONRequest("POST", server.URL, jsonData)
		if err != nil {
			b.Fatalf("JSON request creation failed: %v", err)
		}

		resp, err := driver.Do(req)
		if err != nil {
			b.Fatalf("Request execution failed: %v", err)
		}

		resp.Close()
	}
}
