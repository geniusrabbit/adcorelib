package batchclient

import (
	"context"
	"testing"
	"time"

	"github.com/geniusrabbit/adcorelib/net/httpclient"
)

func TestRateLimitedExecutor_Basic(t *testing.T) {
	driver := &mockDriver{}
	executor := NewRateLimitedExecutor(driver, 10, 5) // 10 requests/second, 5 workers

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

func TestRateLimitedExecutor_RateLimiting(t *testing.T) {
	driver := &mockDriver{}
	executor := NewRateLimitedExecutor(driver, 2, 5) // 2 requests/second

	// Create multiple requests
	var requests []httpclient.Request
	for i := 0; i < 5; i++ {
		req, err := executor.Request("GET", "http://example.com", nil)
		if err != nil {
			t.Fatalf("Failed to create request %d: %v", i, err)
		}
		requests = append(requests, req)
	}

	// Execute requests and measure time
	start := time.Now()
	ctx := context.Background()

	for _, req := range requests {
		resp, err := executor.Execute(ctx, req)
		if err != nil {
			t.Fatalf("Failed to execute request: %v", err)
		}
		resp.Close()
	}

	elapsed := time.Since(start)

	// Should take at least 2 seconds for 5 requests at 2 requests/second
	// (allowing some tolerance for timing)
	if elapsed < 1*time.Second {
		t.Errorf("Expected at least 1 second for rate limiting, got %v", elapsed)
	}

	// Check statistics
	stats := executor.Stats()
	if stats.TotalRequests != 5 {
		t.Errorf("Expected 5 total requests, got %d", stats.TotalRequests)
	}
	if stats.TotalExecuted != 5 {
		t.Errorf("Expected 5 executed requests, got %d", stats.TotalExecuted)
	}
}

func TestRateLimitedExecutor_WithTimeout(t *testing.T) {
	driver := &mockDriver{}
	executor := NewRateLimitedExecutor(driver, 1, 5) // 1 request/second

	req, err := executor.Request("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Execute with timeout
	resp, err := executor.ExecuteWithTimeout(req, 500*time.Millisecond)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	resp.Close()
}

func TestRateLimitedExecutor_ContextCancellation(t *testing.T) {
	driver := &mockDriver{}
	executor := NewRateLimitedExecutor(driver, 1, 5) // 1 request/second

	req, err := executor.Request("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Execute first request to consume the token
	ctx := context.Background()
	resp, err := executor.Execute(ctx, req)
	if err != nil {
		t.Fatalf("Failed to execute first request: %v", err)
	}
	resp.Close()

	// Create second request
	req2, err := executor.Request("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create second request: %v", err)
	}

	// Execute second request with cancelled context
	ctx2, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err = executor.Execute(ctx2, req2)
	if err == nil {
		t.Error("Expected error due to cancelled context")
	}
}

func TestRateLimitedExecutor_UpdateRate(t *testing.T) {
	driver := &mockDriver{}
	executor := NewRateLimitedExecutor(driver, 2, 5) // Start with 2 requests/second

	// Check initial rate
	stats := executor.Stats()
	if stats.RequestsPerSecond != 2 {
		t.Errorf("Expected 2 requests/second, got %d", stats.RequestsPerSecond)
	}

	// Update rate
	executor.UpdateRate(5)

	// Check updated rate
	stats = executor.Stats()
	if stats.RequestsPerSecond != 5 {
		t.Errorf("Expected 5 requests/second after update, got %d", stats.RequestsPerSecond)
	}
}

func TestRateLimitedExecutor_StatsReset(t *testing.T) {
	driver := &mockDriver{}
	executor := NewRateLimitedExecutor(driver, 10, 5)

	// Execute some requests
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		req, err := executor.Request("GET", "http://example.com", nil)
		if err != nil {
			t.Fatalf("Failed to create request %d: %v", i, err)
		}
		resp, err := executor.Execute(ctx, req)
		if err != nil {
			t.Fatalf("Failed to execute request %d: %v", i, err)
		}
		resp.Close()
	}

	// Check initial stats
	stats := executor.Stats()
	if stats.TotalRequests != 3 {
		t.Errorf("Expected 3 total requests, got %d", stats.TotalRequests)
	}

	// Reset stats
	executor.ResetStats()

	// Check reset stats
	stats = executor.Stats()
	if stats.TotalRequests != 0 {
		t.Errorf("Expected 0 total requests after reset, got %d", stats.TotalRequests)
	}
	if stats.TotalExecuted != 0 {
		t.Errorf("Expected 0 executed requests after reset, got %d", stats.TotalExecuted)
	}
}

func TestRateLimitedExecutor_Stop(t *testing.T) {
	driver := &mockDriver{}
	executor := NewRateLimitedExecutor(driver, 10, 5)

	// Stop the executor
	executor.Stop()

	// Give it a moment to actually stop
	time.Sleep(10 * time.Millisecond)

	// Try to execute a request after stopping
	req, err := executor.Request("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err = executor.Execute(ctx, req)
	if err == nil {
		t.Error("Expected error when executing after stop")
	}
}

func BenchmarkRateLimitedExecutor_Execute(b *testing.B) {
	driver := &mockDriver{}
	executor := NewRateLimitedExecutor(driver, 1000, 10) // High rate limit for benchmarking

	req, err := executor.Request("GET", "http://example.com", nil)
	if err != nil {
		b.Fatalf("Failed to create request: %v", err)
	}

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := executor.Execute(ctx, req)
		if err != nil {
			b.Fatalf("Failed to execute request: %v", err)
		}
		resp.Close()
	}
}
