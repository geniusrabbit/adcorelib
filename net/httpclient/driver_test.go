package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockRequest is a mock implementation of the Request interface for testing
type mockRequest struct {
	headers map[string]string
	method  string
	url     string
	body    io.Reader
}

func newMockRequest(method, url string, body io.Reader) *mockRequest {
	return &mockRequest{
		headers: make(map[string]string),
		method:  method,
		url:     url,
		body:    body,
	}
}

func (m *mockRequest) SetHeader(key, value string) {
	m.headers[key] = value
}

func (m *mockRequest) SetContentType(contentType string) {
	m.headers["Content-Type"] = contentType
}

// mockResponse is a mock implementation of the Response interface for testing
type mockResponse struct {
	statusCode int
	body       io.Reader
	closed     bool
}

func newMockResponse(statusCode int, body string) *mockResponse {
	return &mockResponse{
		statusCode: statusCode,
		body:       strings.NewReader(body),
		closed:     false,
	}
}

func (m *mockResponse) Close() error {
	m.closed = true
	return nil
}

func (m *mockResponse) StatusCode() int {
	return m.statusCode
}

func (m *mockResponse) Body() io.Reader {
	return m.body
}

// mockDriver is a mock implementation of the Driver interface for testing
type mockDriver struct {
	requests  []mockRequest
	responses []mockResponse
	errors    []error
	callCount int
}

func newMockDriver() *mockDriver {
	return &mockDriver{
		requests:  make([]mockRequest, 0),
		responses: make([]mockResponse, 0),
		errors:    make([]error, 0),
		callCount: 0,
	}
}

func (m *mockDriver) Request(method, url string, body io.Reader) (Request, error) {
	defer func() { m.callCount++ }()

	req := newMockRequest(method, url, body)
	m.requests = append(m.requests, *req)

	if m.callCount < len(m.errors) && m.errors[m.callCount] != nil {
		return nil, m.errors[m.callCount]
	}

	return req, nil
}

func (m *mockDriver) Do(req Request) (Response, error) {
	defer func() { m.callCount++ }()

	if m.callCount < len(m.errors) && m.errors[m.callCount] != nil {
		return nil, m.errors[m.callCount]
	}

	if m.callCount < len(m.responses) {
		return &m.responses[m.callCount], nil
	}

	// Default response
	return newMockResponse(200, "OK"), nil
}

// Test the Request interface
func TestRequest_SetHeader(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{"Standard header", "Content-Type", "application/json"},
		{"Authorization header", "Authorization", "Bearer token123"},
		{"Custom header", "X-Custom-Header", "custom-value"},
		{"Case insensitive key", "content-type", "text/plain"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newMockRequest("GET", "http://example.com", nil)
			req.SetHeader(tt.key, tt.value)

			if req.headers[tt.key] != tt.value {
				t.Errorf("SetHeader() failed: expected %s=%s, got %s", tt.key, tt.value, req.headers[tt.key])
			}
		})
	}
}

func TestRequest_SetHeader_Overwrite(t *testing.T) {
	req := newMockRequest("GET", "http://example.com", nil)

	// Set initial value
	req.SetHeader("Content-Type", "application/json")
	if req.headers["Content-Type"] != "application/json" {
		t.Errorf("Initial SetHeader() failed")
	}

	// Overwrite with new value
	req.SetHeader("Content-Type", "text/plain")
	if req.headers["Content-Type"] != "text/plain" {
		t.Errorf("SetHeader() overwrite failed: expected text/plain, got %s", req.headers["Content-Type"])
	}
}

// Test the Response interface
func TestResponse_StatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"Success", 200},
		{"Created", 201},
		{"Bad Request", 400},
		{"Unauthorized", 401},
		{"Not Found", 404},
		{"Internal Server Error", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := newMockResponse(tt.statusCode, "test body")

			if resp.StatusCode() != tt.statusCode {
				t.Errorf("StatusCode() = %d, expected %d", resp.StatusCode(), tt.statusCode)
			}
		})
	}
}

func TestResponse_Body(t *testing.T) {
	testBody := "This is a test response body"
	resp := newMockResponse(200, testBody)

	body := resp.Body()
	if body == nil {
		t.Fatal("Body() returned nil")
	}

	// Read the body
	data, err := io.ReadAll(body)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	if string(data) != testBody {
		t.Errorf("Body() content = %q, expected %q", string(data), testBody)
	}
}

func TestResponse_Close(t *testing.T) {
	resp := newMockResponse(200, "test body")

	if resp.closed {
		t.Error("Response should not be closed initially")
	}

	err := resp.Close()
	if err != nil {
		t.Errorf("Close() returned error: %v", err)
	}

	if !resp.closed {
		t.Error("Response should be closed after calling Close()")
	}
}

// Test the Driver interface
func TestDriver_Request(t *testing.T) {
	driver := newMockDriver()

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
			}

			// Verify the request was recorded
			if len(driver.requests) == 0 {
				t.Error("Driver did not record the request")
			}
		})
	}
}

func TestDriver_Request_Error(t *testing.T) {
	driver := newMockDriver()
	expectedErr := fmt.Errorf("request creation failed")
	driver.errors = append(driver.errors, expectedErr)

	req, err := driver.Request("GET", "http://example.com", nil)
	if err != expectedErr {
		t.Errorf("Request() error = %v, expected %v", err, expectedErr)
	}

	if req != nil {
		t.Error("Request() should return nil request on error")
	}
}

func TestDriver_Do(t *testing.T) {
	driver := newMockDriver()

	// Create a request
	req := newMockRequest("GET", "http://example.com", nil)

	// Add a response
	expectedResp := newMockResponse(200, "success")
	driver.responses = append(driver.responses, *expectedResp)

	resp, err := driver.Do(req)
	if err != nil {
		t.Errorf("Do() error = %v", err)
		return
	}

	if resp == nil {
		t.Error("Do() returned nil response")
		return
	}

	if resp.StatusCode() != 200 {
		t.Errorf("Do() response status = %d, expected 200", resp.StatusCode())
	}
}

func TestDriver_Do_Error(t *testing.T) {
	driver := newMockDriver()
	expectedErr := fmt.Errorf("request execution failed")

	// Set up errors: no error for Request(), then error for Do()
	driver.errors = []error{nil, expectedErr}

	// Create request successfully (uses first error slot)
	req, err := driver.Request("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Request() should succeed, got error: %v", err)
	}

	// Execute request (uses second error slot)
	resp, err := driver.Do(req)
	if err != expectedErr {
		t.Errorf("Do() error = %v, expected %v", err, expectedErr)
	}

	if resp != nil {
		t.Error("Do() should return nil response on error")
	}
}

// Integration test demonstrating the two-phase request pattern
func TestDriver_IntegrationFlow(t *testing.T) {
	driver := newMockDriver()

	// Phase 1: Create request
	req, err := driver.Request("POST", "http://api.example.com/users", strings.NewReader(`{"name":"John"}`))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Modify request (add headers)
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Authorization", "Bearer token123")

	// Verify headers were set
	mockReq := req.(*mockRequest)
	if mockReq.headers["Content-Type"] != "application/json" {
		t.Error("Content-Type header not set correctly")
	}
	if mockReq.headers["Authorization"] != "Bearer token123" {
		t.Error("Authorization header not set correctly")
	}

	// Phase 2: Execute request
	expectedResp := newMockResponse(201, `{"id":1,"name":"John"}`)
	driver.responses = []mockResponse{*expectedResp}
	driver.callCount = 0 // Reset counter for Do() call

	resp, err := driver.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	// Verify response
	if resp.StatusCode() != 201 {
		t.Errorf("Expected status 201, got %d", resp.StatusCode())
	}

	// Read response body
	body, err := io.ReadAll(resp.Body())
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	expectedBody := `{"id":1,"name":"John"}`
	if string(body) != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, string(body))
	}

	// Close response
	err = resp.Close()
	if err != nil {
		t.Errorf("Failed to close response: %v", err)
	}
}

// Benchmark tests
func BenchmarkRequest_SetHeader(b *testing.B) {
	req := newMockRequest("GET", "http://example.com", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.SetHeader("Content-Type", "application/json")
	}
}

func BenchmarkDriver_Request(b *testing.B) {
	driver := newMockDriver()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := driver.Request("GET", "http://example.com", nil)
		if err != nil {
			b.Fatalf("Request failed: %v", err)
		}
	}
}

// Test with real HTTP server (integration test)
func TestDriver_WithRealServer(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Echo back the request method and headers
		w.Header().Set("Content-Type", "application/json")

		if r.Header.Get("Authorization") != "" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"method":"%s","auth":"present"}`, r.Method)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"error":"unauthorized"}`)
		}
	}))
	defer server.Close()

	// This test would need the actual stdhttpclient implementation
	// For now, we'll just demonstrate the test structure
	t.Skip("Skipping real server test - requires actual Driver implementation")
}
