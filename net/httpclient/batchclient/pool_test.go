package batchclient

import (
	"context"
	"testing"
	"time"
)

func TestConnectionPool_Basic(t *testing.T) {
	driver := &mockDriver{}
	pool := NewConnectionPool(driver)

	// Test driver interface delegation
	req, err := pool.Request("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	resp, err := pool.Do(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.StatusCode() != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode())
	}

	resp.Close()
}

func TestConnectionPool_WarmupConnections(t *testing.T) {
	driver := &mockDriver{}
	pool := NewConnectionPool(driver)

	hosts := []string{
		"http://example.com",
		"http://test.com",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := pool.WarmupConnections(ctx, hosts, 2)
	if err != nil {
		t.Fatalf("Expected no error during warmup, got %v", err)
	}

	// Check that warmup requests were made
	if driver.GetRequestCount() != 4 { // 2 hosts * 2 connections each
		t.Errorf("Expected 4 warmup requests, got %d", driver.GetRequestCount())
	}

	// Check statistics
	stats := pool.Stats()
	if stats.TotalConnections != 4 {
		t.Errorf("Expected 4 total connections, got %d", stats.TotalConnections)
	}
}

func TestConnectionPool_WarmupWithErrors(t *testing.T) {
	driver := &mockDriver{shouldError: true}
	pool := NewConnectionPool(driver)

	hosts := []string{"http://example.com"}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := pool.WarmupConnections(ctx, hosts, 2)
	if err == nil {
		t.Error("Expected error during warmup with failing driver")
	}
}

func TestConnectionPool_WarmupWithTimeout(t *testing.T) {
	driver := &mockDriver{requestDelay: 100 * time.Millisecond}
	pool := NewConnectionPool(driver)

	hosts := []string{"http://example.com"}

	// Very short timeout to force cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := pool.WarmupConnections(ctx, hosts, 3)
	if err != nil && err != context.DeadlineExceeded {
		t.Errorf("Expected timeout error during warmup, got: %v", err)
	}
}

func TestConnectionPool_RetryLogic(t *testing.T) {
	driver := &mockDriver{errorAfter: 2} // Fail after 2 successful requests
	pool := NewConnectionPool(driver)

	// Configure retry policy
	pool.SetRetryPolicy(3, 10*time.Millisecond)

	// Execute requests - some should fail and retry
	req, err := pool.Request("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// First request should succeed
	resp, err := pool.Do(req)
	if err != nil {
		t.Fatalf("Expected first request to succeed, got %v", err)
	}
	resp.Close()

	// Second request should succeed
	resp, err = pool.Do(req)
	if err != nil {
		t.Fatalf("Expected second request to succeed, got %v", err)
	}
	resp.Close()

	// Third request should fail even with retries
	_, err = pool.Do(req)
	if err == nil {
		t.Error("Expected third request to fail after retries")
	}

	// Check statistics
	stats := pool.Stats()
	if stats.ConnectionFailures == 0 {
		t.Error("Expected some connection failures")
	}
}

func TestConnectionPool_HealthCheck(t *testing.T) {
	driver := &mockDriver{}
	pool := NewConnectionPool(driver)

	hosts := []string{"http://example.com"}

	// Enable health check with short interval
	pool.EnableHealthCheck(hosts, 50*time.Millisecond)

	// Wait for a couple of health check cycles
	time.Sleep(150 * time.Millisecond)

	// Disable health check
	pool.DisableHealthCheck()

	// Check that health check requests were made
	if driver.GetRequestCount() == 0 {
		t.Error("Expected health check requests to be made")
	}
}

func TestConnectionPool_StatsReset(t *testing.T) {
	driver := &mockDriver{}
	pool := NewConnectionPool(driver)

	// Execute some requests
	for i := 0; i < 3; i++ {
		req, err := pool.Request("GET", "http://example.com", nil)
		if err != nil {
			t.Fatalf("Failed to create request %d: %v", i, err)
		}
		resp, err := pool.Do(req)
		if err != nil {
			t.Fatalf("Failed to execute request %d: %v", i, err)
		}
		resp.Close()
	}

	// Check initial stats
	stats := pool.Stats()
	if stats.ConnectionAttempts != 3 {
		t.Errorf("Expected 3 connection attempts, got %d", stats.ConnectionAttempts)
	}

	// Reset stats
	pool.ResetStats()

	// Check reset stats
	stats = pool.Stats()
	if stats.ConnectionAttempts != 0 {
		t.Errorf("Expected 0 connection attempts after reset, got %d", stats.ConnectionAttempts)
	}
	if stats.ActiveConnections != 0 {
		t.Errorf("Expected 0 active connections after reset, got %d", stats.ActiveConnections)
	}
}

func TestConnectionPool_Stop(t *testing.T) {
	driver := &mockDriver{}
	pool := NewConnectionPool(driver)

	// Enable health check
	pool.EnableHealthCheck([]string{"http://example.com"}, 100*time.Millisecond)

	// Stop the pool
	pool.Stop()

	// Health check should be disabled
	stats := pool.Stats()
	if stats.HealthCheckEnabled {
		t.Error("Expected health check to be disabled after stop")
	}
}

func BenchmarkConnectionPool_Execute(b *testing.B) {
	driver := &mockDriver{}
	pool := NewConnectionPool(driver)

	req, err := pool.Request("GET", "http://example.com", nil)
	if err != nil {
		b.Fatalf("Failed to create request: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := pool.Do(req)
		if err != nil {
			b.Fatalf("Failed to execute request: %v", err)
		}
		resp.Close()
	}
}
