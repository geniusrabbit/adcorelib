package batchclient

import (
	"context"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/geniusrabbit/adcorelib/net/httpclient"
)

// Mock implementations for testing
type mockRequest struct {
	method  string
	url     string
	body    io.Reader
	headers map[string]string
}

func (mr *mockRequest) SetHeader(key, value string) {
	if mr.headers == nil {
		mr.headers = make(map[string]string)
	}
	mr.headers[key] = value
}

type mockResponse struct {
	statusCode int
	body       io.Reader
	closed     bool
}

func (mr *mockResponse) Close() error {
	mr.closed = true
	return nil
}

func (mr *mockResponse) StatusCode() int {
	return mr.statusCode
}

func (mr *mockResponse) Body() io.Reader {
	return mr.body
}

type mockDriver struct {
	requestCount int
	requestDelay time.Duration
	shouldError  bool
	errorAfter   int
	mu           sync.Mutex
}

func (md *mockDriver) Request(method, url string, body io.Reader) (httpclient.Request, error) {
	if md.shouldError {
		return nil, errors.New("mock request creation error")
	}
	return &mockRequest{
		method: method,
		url:    url,
		body:   body,
	}, nil
}

func (md *mockDriver) Do(req httpclient.Request) (httpclient.Response, error) {
	md.mu.Lock()
	defer md.mu.Unlock()

	md.requestCount++

	if md.requestDelay > 0 {
		time.Sleep(md.requestDelay)
	}

	if md.shouldError || (md.errorAfter > 0 && md.requestCount > md.errorAfter) {
		return nil, errors.New("mock execution error")
	}

	return &mockResponse{
		statusCode: 200,
		body:       strings.NewReader("mock response"),
	}, nil
}

func (md *mockDriver) GetRequestCount() int {
	md.mu.Lock()
	defer md.mu.Unlock()
	return md.requestCount
}

func (md *mockDriver) ResetRequestCount() {
	md.mu.Lock()
	defer md.mu.Unlock()
	md.requestCount = 0
}

func TestBatchExecutor_Basic(t *testing.T) {
	driver := &mockDriver{}
	executor := NewBatchExecutor(driver, 10)

	// Test driver interface delegation
	req, err := executor.Request("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	resp, err := executor.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode() != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode())
	}

	resp.Close()
}

func TestBatchExecutor_ExecuteBatch(t *testing.T) {
	driver := &mockDriver{}
	executor := NewBatchExecutor(driver, 5)

	// Create multiple requests
	var requests []httpclient.Request
	for i := 0; i < 10; i++ {
		req, err := executor.Request("GET", "http://example.com", nil)
		if err != nil {
			t.Fatalf("Failed to create request %d: %v", i, err)
		}
		requests = append(requests, req)
	}

	// Execute batch
	ctx := context.Background()
	result := executor.ExecuteBatch(ctx, requests)

	// Verify results
	if result.Completed != 10 {
		t.Errorf("Expected 10 completed requests, got %d", result.Completed)
	}

	if result.Errors != 0 {
		t.Errorf("Expected 0 errors, got %d", result.Errors)
	}

	if len(result.Requests) != 10 {
		t.Errorf("Expected 10 results, got %d", len(result.Requests))
	}

	// Verify all requests were executed
	for i, req := range result.Requests {
		if req.Index != i {
			t.Errorf("Expected request index %d, got %d", i, req.Index)
		}
		if req.Response == nil {
			t.Errorf("Expected response for request %d", i)
		}
		if req.Error != nil {
			t.Errorf("Expected no error for request %d, got %v", i, req.Error)
		}
	}

	// Check statistics
	stats := executor.Stats()
	if stats.TotalRequests != 10 {
		t.Errorf("Expected 10 total requests in stats, got %d", stats.TotalRequests)
	}
	if stats.TotalCompleted != 10 {
		t.Errorf("Expected 10 completed requests in stats, got %d", stats.TotalCompleted)
	}
}

func TestBatchExecutor_WithErrors(t *testing.T) {
	driver := &mockDriver{errorAfter: 5}
	executor := NewBatchExecutor(driver, 3)

	// Create multiple requests
	var requests []httpclient.Request
	for i := 0; i < 10; i++ {
		req, err := executor.Request("GET", "http://example.com", nil)
		if err != nil {
			t.Fatalf("Failed to create request %d: %v", i, err)
		}
		requests = append(requests, req)
	}

	// Execute batch
	ctx := context.Background()
	result := executor.ExecuteBatch(ctx, requests)

	// Verify results
	if result.Completed != 10 {
		t.Errorf("Expected 10 completed requests, got %d", result.Completed)
	}

	if result.Errors != 5 {
		t.Errorf("Expected 5 errors, got %d", result.Errors)
	}

	// Check statistics
	stats := executor.Stats()
	if stats.TotalErrors != 5 {
		t.Errorf("Expected 5 total errors in stats, got %d", stats.TotalErrors)
	}
}

func TestBatchExecutor_ContextCancellation(t *testing.T) {
	driver := &mockDriver{requestDelay: 100 * time.Millisecond}
	executor := NewBatchExecutor(driver, 2)

	// Create multiple requests
	var requests []httpclient.Request
	for i := 0; i < 5; i++ {
		req, err := executor.Request("GET", "http://example.com", nil)
		if err != nil {
			t.Fatalf("Failed to create request %d: %v", i, err)
		}
		requests = append(requests, req)
	}

	// Execute batch with short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	result := executor.ExecuteBatch(ctx, requests)

	// Should complete fewer requests due to timeout
	if result.Completed >= 5 {
		t.Errorf("Expected fewer than 5 completed requests due to timeout, got %d", result.Completed)
	}
}

func TestBatchExecutor_StatsReset(t *testing.T) {
	driver := &mockDriver{}
	executor := NewBatchExecutor(driver, 10)

	// Create and execute some requests
	var requests []httpclient.Request
	for i := 0; i < 5; i++ {
		req, err := executor.Request("GET", "http://example.com", nil)
		if err != nil {
			t.Fatalf("Failed to create request %d: %v", i, err)
		}
		requests = append(requests, req)
	}

	ctx := context.Background()
	executor.ExecuteBatch(ctx, requests)

	// Check initial stats
	stats := executor.Stats()
	if stats.TotalRequests != 5 {
		t.Errorf("Expected 5 total requests, got %d", stats.TotalRequests)
	}

	// Reset stats
	executor.ResetStats()

	// Check reset stats
	stats = executor.Stats()
	if stats.TotalRequests != 0 {
		t.Errorf("Expected 0 total requests after reset, got %d", stats.TotalRequests)
	}
	if stats.TotalCompleted != 0 {
		t.Errorf("Expected 0 completed requests after reset, got %d", stats.TotalCompleted)
	}
}

func BenchmarkBatchExecutor_ExecuteBatch(b *testing.B) {
	driver := &mockDriver{}
	executor := NewBatchExecutor(driver, 100)

	// Create requests
	var requests []httpclient.Request
	for i := 0; i < 1000; i++ {
		req, err := executor.Request("GET", "http://example.com", nil)
		if err != nil {
			b.Fatalf("Failed to create request %d: %v", i, err)
		}
		requests = append(requests, req)
	}

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		driver.ResetRequestCount()
		executor.ExecuteBatch(ctx, requests)
	}
}

func BenchmarkBatchExecutor_SingleRequest(b *testing.B) {
	driver := &mockDriver{}
	executor := NewBatchExecutor(driver, 100)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req, err := executor.Request("GET", "http://example.com", nil)
		if err != nil {
			b.Fatalf("Failed to create request: %v", err)
		}

		resp, err := executor.Do(req)
		if err != nil {
			b.Fatalf("Failed to execute request: %v", err)
		}

		resp.Close()
	}
}
