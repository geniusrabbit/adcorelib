// Package batchclient provides high-performance batch execution capabilities for HTTP requests.
// It implements the httpclient.Driver interface and delegates requests to underlying driver implementations.
// This package is designed to work with any httpclient.Driver implementation and provides:
// - Concurrent batch execution of multiple requests
// - Rate limiting capabilities
// - Connection pooling management
// - Performance monitoring and statistics
package batchclient

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/geniusrabbit/adcorelib/net/httpclient"
)

// BatchRequest represents a single request in a batch operation.
type BatchRequest struct {
	Request  httpclient.Request
	Response httpclient.Response
	Error    error
	Index    int // Original index in the batch
}

// BatchResult contains the results of a batch operation.
type BatchResult struct {
	Requests  []BatchRequest
	Completed int
	Errors    int
	Duration  time.Duration
	StartTime time.Time
	EndTime   time.Time
}

// BatchExecutor provides high-performance batch execution of HTTP requests.
// It implements httpclient.Driver interface and can be used as a drop-in replacement
// for any driver while providing batch processing capabilities.
type BatchExecutor struct {
	driver     httpclient.Driver
	maxWorkers int

	// Statistics
	totalRequests  int64
	totalCompleted int64
	totalErrors    int64
	totalDuration  time.Duration
	statsMutex     sync.RWMutex
}

// NewBatchExecutor creates a new batch executor with the specified driver and worker count.
// The driver parameter can be any implementation of httpclient.Driver interface.
// If maxWorkers is <= 0, it defaults to 100 workers.
func NewBatchExecutor(driver httpclient.Driver, maxWorkers int) *BatchExecutor {
	if maxWorkers <= 0 {
		maxWorkers = 100 // Default to 100 workers
	}

	return &BatchExecutor{
		driver:     driver,
		maxWorkers: maxWorkers,
	}
}

// Request implements httpclient.Driver interface.
// Creates a new HTTP request using the underlying driver.
func (be *BatchExecutor) Request(method, url string, body io.Reader) (httpclient.Request, error) {
	return be.driver.Request(method, url, body)
}

// Do implements httpclient.Driver interface.
// Executes a single HTTP request using the underlying driver.
func (be *BatchExecutor) Do(req httpclient.Request) (httpclient.Response, error) {
	return be.driver.Do(req)
}

// ExecuteBatch executes multiple HTTP requests concurrently with optimal performance.
// This is the main batch processing method that provides high-throughput execution.
func (be *BatchExecutor) ExecuteBatch(ctx context.Context, requests []httpclient.Request) BatchResult {
	start := time.Now()

	// Create new channels for this batch to avoid reuse issues
	requestChan := make(chan BatchRequest, be.maxWorkers*2)
	resultChan := make(chan BatchRequest, be.maxWorkers*2)
	var wg sync.WaitGroup

	// Initialize result
	result := BatchResult{
		Requests:  make([]BatchRequest, len(requests)),
		StartTime: start,
	}

	// Update statistics
	be.statsMutex.Lock()
	be.totalRequests += int64(len(requests))
	be.statsMutex.Unlock()

	// Start workers
	for i := 0; i < be.maxWorkers; i++ {
		wg.Add(1)
		go be.batchWorker(ctx, requestChan, resultChan, &wg)
	}

	// Send requests
	go func() {
		for i, req := range requests {
			select {
			case requestChan <- BatchRequest{Request: req, Index: i}:
			case <-ctx.Done():
				return
			}
		}
		close(requestChan)
	}()

	// Collect results
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Process results
	for batchReq := range resultChan {
		result.Requests[batchReq.Index] = batchReq
		result.Completed++
		if batchReq.Error != nil {
			result.Errors++
		}
	}

	result.Duration = time.Since(start)
	result.EndTime = time.Now()

	// Update statistics
	be.statsMutex.Lock()
	be.totalCompleted += int64(result.Completed)
	be.totalErrors += int64(result.Errors)
	be.totalDuration += result.Duration
	be.statsMutex.Unlock()

	return result
}

// batchWorker processes requests from the request channel for a single batch.
func (be *BatchExecutor) batchWorker(ctx context.Context, requestChan <-chan BatchRequest, resultChan chan<- BatchRequest, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case batchReq, ok := <-requestChan:
			if !ok {
				return // Channel closed
			}

			// Execute the request using the underlying driver
			resp, err := be.driver.Do(batchReq.Request)
			batchReq.Response = resp
			batchReq.Error = err

			// Send result
			select {
			case resultChan <- batchReq:
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// Stats returns execution statistics.
func (be *BatchExecutor) Stats() BatchStats {
	be.statsMutex.RLock()
	defer be.statsMutex.RUnlock()

	var avgDuration time.Duration
	if be.totalCompleted > 0 {
		avgDuration = time.Duration(int64(be.totalDuration) / be.totalCompleted)
	}

	return BatchStats{
		TotalRequests:   be.totalRequests,
		TotalCompleted:  be.totalCompleted,
		TotalErrors:     be.totalErrors,
		TotalDuration:   be.totalDuration,
		AverageDuration: avgDuration,
		ErrorRate:       float64(be.totalErrors) / float64(be.totalRequests),
	}
}

// ResetStats resets all execution statistics.
func (be *BatchExecutor) ResetStats() {
	be.statsMutex.Lock()
	defer be.statsMutex.Unlock()

	be.totalRequests = 0
	be.totalCompleted = 0
	be.totalErrors = 0
	be.totalDuration = 0
}

// BatchStats contains execution statistics.
type BatchStats struct {
	TotalRequests   int64
	TotalCompleted  int64
	TotalErrors     int64
	TotalDuration   time.Duration
	AverageDuration time.Duration
	ErrorRate       float64
}
