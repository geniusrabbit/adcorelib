package batchclient_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/geniusrabbit/adcorelib/net/httpclient"
	"github.com/geniusrabbit/adcorelib/net/httpclient/batchclient"
	"github.com/geniusrabbit/adcorelib/net/httpclient/stdhttpclient"
)

// TestIntegrationWithStdHTTPClient demonstrates that batchclient works with stdhttpclient
func TestIntegrationWithStdHTTPClient(t *testing.T) {
	// Create a standard HTTP client driver
	driver := stdhttpclient.NewDriver()

	// Wrap it with batch execution
	batchExecutor := batchclient.NewBatchExecutor(driver, 10)

	// Test basic driver interface methods
	req, err := batchExecutor.Request("GET", "https://httpbin.org/json", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Execute single request
	resp, err := batchExecutor.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}

	if resp.StatusCode() != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode())
	}

	resp.Close()
}

// TestChainedExecutorsWithStdHTTPClient demonstrates chaining multiple executors
func TestChainedExecutorsWithStdHTTPClient(t *testing.T) {
	// Create base driver
	driver := stdhttpclient.NewDriver()

	// Chain executors: Pool -> Rate Limiter -> Batch Executor
	pool := batchclient.NewConnectionPool(driver)
	rateLimiter := batchclient.NewRateLimitedExecutor(pool, 5, 3)
	batchExecutor := batchclient.NewBatchExecutor(rateLimiter, 10)

	// Create multiple requests
	var requests []httpclient.Request
	for i := 0; i < 5; i++ {
		req, err := batchExecutor.Request("GET", "https://httpbin.org/json", nil)
		if err != nil {
			t.Fatalf("Failed to create request %d: %v", i, err)
		}
		requests = append(requests, req)
	}

	// Execute batch
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result := batchExecutor.ExecuteBatch(ctx, requests)

	// Check results
	if result.Completed != len(requests) {
		t.Errorf("Expected %d completed requests, got %d", len(requests), result.Completed)
	}

	if result.Errors > 0 {
		t.Errorf("Expected no errors, got %d", result.Errors)
	}

	// Clean up
	for _, req := range result.Requests {
		if req.Response != nil {
			req.Response.Close()
		}
	}

	// Print statistics
	fmt.Printf("Batch stats: %d requests, %d completed, %d errors\n",
		len(requests), result.Completed, result.Errors)

	batchStats := batchExecutor.Stats()
	fmt.Printf("Total batch stats: %d requests, %d completed, %d errors\n",
		batchStats.TotalRequests, batchStats.TotalCompleted, batchStats.TotalErrors)

	// Clean up
	rateLimiter.Stop()
	pool.Stop()
}

// BenchmarkIntegrationWithStdHTTPClient benchmarks the integration
func BenchmarkIntegrationWithStdHTTPClient(b *testing.B) {
	driver := stdhttpclient.NewDriver()
	batchExecutor := batchclient.NewBatchExecutor(driver, 50)

	// Create requests
	var requests []httpclient.Request
	for i := 0; i < 100; i++ {
		req, err := batchExecutor.Request("GET", "https://httpbin.org/json", nil)
		if err != nil {
			b.Fatalf("Failed to create request %d: %v", i, err)
		}
		requests = append(requests, req)
	}

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result := batchExecutor.ExecuteBatch(ctx, requests)

		// Clean up
		for _, req := range result.Requests {
			if req.Response != nil {
				req.Response.Close()
			}
		}
	}
}
