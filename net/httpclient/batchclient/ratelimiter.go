package batchclient

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/geniusrabbit/adcorelib/net/httpclient"
)

// RateLimitedExecutor provides rate-limited execution of HTTP requests.
// It implements httpclient.Driver interface and can wrap any underlying driver
// to provide rate limiting capabilities.
type RateLimitedExecutor struct {
	driver            httpclient.Driver
	rateLimiter       chan struct{}
	maxWorkers        int
	requestsPerSecond int

	// Statistics
	totalRequests  int64
	totalExecuted  int64
	totalThrottled int64
	statsMutex     sync.RWMutex

	// Control
	stopChan  chan struct{}
	stopped   bool
	stopMutex sync.RWMutex
}

// NewRateLimitedExecutor creates a new rate-limited executor.
// The driver parameter can be any implementation of httpclient.Driver interface.
// requestsPerSecond controls the maximum rate of request execution.
func NewRateLimitedExecutor(driver httpclient.Driver, requestsPerSecond int, maxWorkers int) *RateLimitedExecutor {
	if requestsPerSecond <= 0 {
		requestsPerSecond = 100 // Default rate limit
	}

	if maxWorkers <= 0 {
		maxWorkers = 10 // Default workers for rate limiter
	}

	rle := &RateLimitedExecutor{
		driver:            driver,
		rateLimiter:       make(chan struct{}, requestsPerSecond),
		maxWorkers:        maxWorkers,
		requestsPerSecond: requestsPerSecond,
		stopChan:          make(chan struct{}),
	}

	// Fill the rate limiter initially
	for i := 0; i < requestsPerSecond; i++ {
		rle.rateLimiter <- struct{}{}
	}

	// Start the refill goroutine
	go rle.refillRateLimiter()

	return rle
}

// Request implements httpclient.Driver interface.
// Creates a new HTTP request using the underlying driver.
func (rle *RateLimitedExecutor) Request(method, url string, body io.Reader) (httpclient.Request, error) {
	return rle.driver.Request(method, url, body)
}

// Do implements httpclient.Driver interface.
// Executes a HTTP request with rate limiting applied.
func (rle *RateLimitedExecutor) Do(req httpclient.Request) (httpclient.Response, error) {
	return rle.Execute(context.Background(), req)
}

// Execute executes a request with rate limiting.
// This method respects the configured rate limit and will block if necessary.
func (rle *RateLimitedExecutor) Execute(ctx context.Context, req httpclient.Request) (httpclient.Response, error) {
	// Update request statistics
	rle.statsMutex.Lock()
	rle.totalRequests++
	rle.statsMutex.Unlock()

	// Wait for rate limiter
	select {
	case <-rle.rateLimiter:
		// Got permission to proceed
		rle.statsMutex.Lock()
		rle.totalExecuted++
		rle.statsMutex.Unlock()
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-rle.stopChan:
		return nil, context.Canceled
	}

	// Execute the request using the underlying driver
	return rle.driver.Do(req)
}

// ExecuteWithTimeout executes a request with both rate limiting and timeout.
func (rle *RateLimitedExecutor) ExecuteWithTimeout(req httpclient.Request, timeout time.Duration) (httpclient.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return rle.Execute(ctx, req)
}

// refillRateLimiter periodically refills the rate limiter tokens.
func (rle *RateLimitedExecutor) refillRateLimiter() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rle.refillTokens()
		case <-rle.stopChan:
			return
		}
	}
}

// refillTokens adds tokens to the rate limiter up to the configured limit.
func (rle *RateLimitedExecutor) refillTokens() {
	for i := 0; i < rle.requestsPerSecond; i++ {
		select {
		case rle.rateLimiter <- struct{}{}:
			// Token added successfully
		default:
			// Rate limiter is full, stop adding tokens
			return
		}
	}
}

// UpdateRate changes the rate limit at runtime.
func (rle *RateLimitedExecutor) UpdateRate(requestsPerSecond int) {
	if requestsPerSecond <= 0 {
		return
	}

	rle.stopMutex.Lock()
	defer rle.stopMutex.Unlock()

	// Create new rate limiter with updated capacity
	newRateLimiter := make(chan struct{}, requestsPerSecond)

	// Fill the new rate limiter
	for i := 0; i < requestsPerSecond; i++ {
		newRateLimiter <- struct{}{}
	}

	// Replace the old rate limiter
	rle.rateLimiter = newRateLimiter
	rle.requestsPerSecond = requestsPerSecond
}

// Stop stops the rate limiter and all background goroutines.
func (rle *RateLimitedExecutor) Stop() {
	rle.stopMutex.Lock()
	defer rle.stopMutex.Unlock()

	if !rle.stopped {
		close(rle.stopChan)
		rle.stopped = true
	}
}

// Stats returns rate limiter statistics.
func (rle *RateLimitedExecutor) Stats() RateLimiterStats {
	rle.statsMutex.RLock()
	defer rle.statsMutex.RUnlock()

	var throttleRate float64
	if rle.totalRequests > 0 {
		throttleRate = float64(rle.totalThrottled) / float64(rle.totalRequests)
	}

	return RateLimiterStats{
		RequestsPerSecond: rle.requestsPerSecond,
		TotalRequests:     rle.totalRequests,
		TotalExecuted:     rle.totalExecuted,
		TotalThrottled:    rle.totalThrottled,
		ThrottleRate:      throttleRate,
		QueueSize:         len(rle.rateLimiter),
		QueueCapacity:     cap(rle.rateLimiter),
	}
}

// ResetStats resets all rate limiter statistics.
func (rle *RateLimitedExecutor) ResetStats() {
	rle.statsMutex.Lock()
	defer rle.statsMutex.Unlock()

	rle.totalRequests = 0
	rle.totalExecuted = 0
	rle.totalThrottled = 0
}

// RateLimiterStats contains rate limiter statistics.
type RateLimiterStats struct {
	RequestsPerSecond int
	TotalRequests     int64
	TotalExecuted     int64
	TotalThrottled    int64
	ThrottleRate      float64
	QueueSize         int
	QueueCapacity     int
}
