package batchclient_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/geniusrabbit/adcorelib/net/httpclient"
	"github.com/geniusrabbit/adcorelib/net/httpclient/batchclient"
	"github.com/geniusrabbit/adcorelib/net/httpclient/stdhttpclient"
)

func ExampleBatchExecutor() {
	// Create a base HTTP client driver
	driver := stdhttpclient.NewDriver()

	// Wrap it with batch execution capabilities
	batchExecutor := batchclient.NewBatchExecutor(driver, 100) // 100 concurrent workers

	// Create multiple requests
	var requests []httpclient.Request
	for i := 0; i < 10; i++ {
		req, err := batchExecutor.Request("GET", "https://httpbin.org/delay/1", nil)
		if err != nil {
			log.Printf("Failed to create request %d: %v", i, err)
			continue
		}
		requests = append(requests, req)
	}

	// Execute all requests concurrently
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	start := time.Now()
	result := batchExecutor.ExecuteBatch(ctx, requests)
	elapsed := time.Since(start)

	// Print results
	fmt.Printf("Batch execution completed in %v\n", elapsed)
	fmt.Printf("Total requests: %d\n", len(result.Requests))
	fmt.Printf("Completed: %d\n", result.Completed)
	fmt.Printf("Errors: %d\n", result.Errors)

	// Clean up responses
	for _, req := range result.Requests {
		if req.Response != nil {
			req.Response.Close()
		}
	}

	// Print statistics
	stats := batchExecutor.Stats()
	fmt.Printf("Statistics - Total: %d, Completed: %d, Errors: %d\n",
		stats.TotalRequests, stats.TotalCompleted, stats.TotalErrors)
}

func ExampleRateLimitedExecutor() {
	// Create a base HTTP client driver
	driver := stdhttpclient.NewDriver()

	// Wrap it with rate limiting (10 requests per second)
	rateLimiter := batchclient.NewRateLimitedExecutor(driver, 10, 5)

	// Execute requests - will be rate limited
	ctx := context.Background()

	start := time.Now()
	for i := 0; i < 5; i++ {
		req, err := rateLimiter.Request("GET", "https://httpbin.org/get", nil)
		if err != nil {
			log.Printf("Failed to create request %d: %v", i, err)
			continue
		}

		resp, err := rateLimiter.Execute(ctx, req)
		if err != nil {
			log.Printf("Failed to execute request %d: %v", i, err)
			continue
		}

		fmt.Printf("Request %d completed with status %d\n", i+1, resp.StatusCode())
		resp.Close()
	}
	elapsed := time.Since(start)

	fmt.Printf("5 requests completed in %v (rate limited)\n", elapsed)

	// Print statistics
	stats := rateLimiter.Stats()
	fmt.Printf("Rate limiter stats - Requests: %d, Executed: %d, Rate: %d/s\n",
		stats.TotalRequests, stats.TotalExecuted, stats.RequestsPerSecond)

	// Clean up
	rateLimiter.Stop()
}

func ExampleConnectionPool() {
	// Create a base HTTP client driver
	driver := stdhttpclient.NewDriver()

	// Wrap it with connection pool management
	pool := batchclient.NewConnectionPool(driver)

	// Configure retry policy
	pool.SetRetryPolicy(3, 100*time.Millisecond)

	// Warmup connections to target hosts
	hosts := []string{
		"https://httpbin.org",
		"https://jsonplaceholder.typicode.com",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("Warming up connections...")
	err := pool.WarmupConnections(ctx, hosts, 2)
	if err != nil {
		log.Printf("Warmup error: %v", err)
	}

	// Enable health checking
	pool.EnableHealthCheck(hosts, 30*time.Second)

	// Execute requests with connection pooling
	for i := 0; i < 3; i++ {
		req, err := pool.Request("GET", "https://httpbin.org/get", nil)
		if err != nil {
			log.Printf("Failed to create request %d: %v", i, err)
			continue
		}

		resp, err := pool.Do(req)
		if err != nil {
			log.Printf("Failed to execute request %d: %v", i, err)
			continue
		}

		fmt.Printf("Request %d completed with status %d\n", i+1, resp.StatusCode())
		resp.Close()
	}

	// Print statistics
	stats := pool.Stats()
	fmt.Printf("Pool stats - Connections: %d, Attempts: %d, Failures: %d, Success Rate: %.2f%%\n",
		stats.TotalConnections, stats.ConnectionAttempts, stats.ConnectionFailures, stats.SuccessRate*100)

	// Clean up
	pool.Stop()
}

func Example_chainedExecutors() {
	// Create a base HTTP client driver
	driver := stdhttpclient.NewDriver()

	// Chain multiple executors together
	// First: Connection pool management
	pool := batchclient.NewConnectionPool(driver)
	pool.SetRetryPolicy(2, 50*time.Millisecond)

	// Second: Rate limiting
	rateLimiter := batchclient.NewRateLimitedExecutor(pool, 20, 10)

	// Third: Batch execution
	batchExecutor := batchclient.NewBatchExecutor(rateLimiter, 50)

	// Create requests
	var requests []httpclient.Request
	for i := 0; i < 10; i++ {
		req, err := batchExecutor.Request("GET", "https://httpbin.org/get", nil)
		if err != nil {
			log.Printf("Failed to create request %d: %v", i, err)
			continue
		}
		requests = append(requests, req)
	}

	// Execute with all features: connection pooling, rate limiting, and batch processing
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("Executing with chained executors (pool + rate limiter + batch)...")
	start := time.Now()
	result := batchExecutor.ExecuteBatch(ctx, requests)
	elapsed := time.Since(start)

	fmt.Printf("Chained execution completed in %v\n", elapsed)
	fmt.Printf("Completed: %d/%d requests\n", result.Completed, len(result.Requests))

	// Clean up responses
	for _, req := range result.Requests {
		if req.Response != nil {
			req.Response.Close()
		}
	}

	// Print all statistics
	fmt.Println("\nStatistics:")

	batchStats := batchExecutor.Stats()
	fmt.Printf("  Batch - Total: %d, Completed: %d, Errors: %d\n",
		batchStats.TotalRequests, batchStats.TotalCompleted, batchStats.TotalErrors)

	rateStats := rateLimiter.Stats()
	fmt.Printf("  Rate Limiter - Total: %d, Executed: %d, Rate: %d/s\n",
		rateStats.TotalRequests, rateStats.TotalExecuted, rateStats.RequestsPerSecond)

	poolStats := pool.Stats()
	fmt.Printf("  Pool - Connections: %d, Attempts: %d, Success Rate: %.2f%%\n",
		poolStats.TotalConnections, poolStats.ConnectionAttempts, poolStats.SuccessRate*100)

	// Clean up
	rateLimiter.Stop()
	pool.Stop()
}

// Example of using batchclient with a custom driver
func Example_withCustomDriver() {
	// Create a custom driver that logs requests
	customDriver := &loggingDriver{
		underlying: stdhttpclient.NewDriver(),
	}

	// Use batch execution with the custom driver
	batchExecutor := batchclient.NewBatchExecutor(customDriver, 10)

	// Create and execute requests
	var requests []httpclient.Request
	for i := 0; i < 3; i++ {
		req, err := batchExecutor.Request("GET", "https://httpbin.org/get", nil)
		if err != nil {
			log.Printf("Failed to create request %d: %v", i, err)
			continue
		}
		requests = append(requests, req)
	}

	ctx := context.Background()
	result := batchExecutor.ExecuteBatch(ctx, requests)

	fmt.Printf("Custom driver executed %d requests\n", result.Completed)

	// Clean up
	for _, req := range result.Requests {
		if req.Response != nil {
			req.Response.Close()
		}
	}
}

// loggingDriver is a custom driver that logs all requests
type loggingDriver struct {
	underlying httpclient.Driver
}

func (ld *loggingDriver) Request(method, url string, body io.Reader) (httpclient.Request, error) {
	log.Printf("Creating request: %s %s", method, url)
	return ld.underlying.Request(method, url, body)
}

func (ld *loggingDriver) Do(req httpclient.Request) (httpclient.Response, error) {
	log.Printf("Executing request")
	return ld.underlying.Do(req)
}
