package batchclient

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/geniusrabbit/adcorelib/net/httpclient"
)

// ConnectionPool provides connection pooling statistics and management.
// It wraps any httpclient.Driver implementation and provides connection management features.
type ConnectionPool struct {
	driver httpclient.Driver

	// Statistics
	totalConnections   int64
	activeConnections  int64
	connectionAttempts int64
	connectionFailures int64
	statsMutex         sync.RWMutex

	// Configuration
	maxRetries          int
	retryDelay          time.Duration
	healthCheckInterval time.Duration

	// Health checking
	healthCheckEnabled bool
	healthCheckHosts   []string
	stopChan           chan struct{}
	stopped            bool
	stopMutex          sync.RWMutex
}

// NewConnectionPool creates a new connection pool manager.
// The driver parameter can be any implementation of httpclient.Driver interface.
func NewConnectionPool(driver httpclient.Driver) *ConnectionPool {
	return &ConnectionPool{
		driver:              driver,
		maxRetries:          3,
		retryDelay:          100 * time.Millisecond,
		healthCheckInterval: 30 * time.Second,
		stopChan:            make(chan struct{}),
	}
}

// Request implements httpclient.Driver interface.
// Creates a new HTTP request using the underlying driver.
func (cp *ConnectionPool) Request(method, url string, body io.Reader) (httpclient.Request, error) {
	return cp.driver.Request(method, url, body)
}

// Do implements httpclient.Driver interface.
// Executes a HTTP request with connection pooling and retry logic.
func (cp *ConnectionPool) Do(req httpclient.Request) (httpclient.Response, error) {
	cp.statsMutex.Lock()
	cp.connectionAttempts++
	cp.statsMutex.Unlock()

	// Get retry configuration safely
	cp.stopMutex.RLock()
	maxRetries := cp.maxRetries
	retryDelay := cp.retryDelay
	cp.stopMutex.RUnlock()

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(retryDelay)
		}

		resp, err := cp.driver.Do(req)
		if err == nil {
			cp.statsMutex.Lock()
			cp.activeConnections++
			cp.statsMutex.Unlock()
			return resp, nil
		}

		lastErr = err
		cp.statsMutex.Lock()
		cp.connectionFailures++
		cp.statsMutex.Unlock()
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries+1, lastErr)
}

// WarmupConnections pre-establishes connections to the specified hosts.
// This can improve performance by avoiding connection overhead during actual requests.
func (cp *ConnectionPool) WarmupConnections(ctx context.Context, hosts []string, connectionsPerHost int) error {
	var wg sync.WaitGroup
	errorChan := make(chan error, len(hosts)*connectionsPerHost)

	for _, host := range hosts {
		for i := 0; i < connectionsPerHost; i++ {
			wg.Add(1)
			go func(host string) {
				defer wg.Done()

				select {
				case <-ctx.Done():
					errorChan <- ctx.Err()
					return
				default:
				}

				// Make a simple HEAD request to establish connection
				req, err := cp.driver.Request("HEAD", host, nil)
				if err != nil {
					errorChan <- fmt.Errorf("failed to create warmup request for %s: %w", host, err)
					return
				}

				resp, err := cp.driver.Do(req)
				if err != nil {
					errorChan <- fmt.Errorf("failed to warmup connection to %s: %w", host, err)
					return
				}

				if resp != nil {
					resp.Close()
				}

				cp.statsMutex.Lock()
				cp.totalConnections++
				cp.statsMutex.Unlock()
			}(host)
		}
	}

	// Wait for all warmup requests to complete
	wg.Wait()
	close(errorChan)

	// Collect any errors
	var errors []error
	for err := range errorChan {
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("warmup completed with %d errors: %v", len(errors), errors[0])
	}

	return nil
}

// EnableHealthCheck enables periodic health checking for the specified hosts.
func (cp *ConnectionPool) EnableHealthCheck(hosts []string, interval time.Duration) {
	cp.stopMutex.Lock()
	defer cp.stopMutex.Unlock()

	cp.healthCheckEnabled = true
	cp.healthCheckHosts = hosts
	cp.healthCheckInterval = interval

	go cp.healthCheckLoop()
}

// DisableHealthCheck disables periodic health checking.
func (cp *ConnectionPool) DisableHealthCheck() {
	cp.stopMutex.Lock()
	defer cp.stopMutex.Unlock()

	cp.healthCheckEnabled = false
}

// healthCheckLoop runs periodic health checks.
func (cp *ConnectionPool) healthCheckLoop() {
	ticker := time.NewTicker(cp.healthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cp.stopMutex.RLock()
			enabled := cp.healthCheckEnabled
			cp.stopMutex.RUnlock()

			if !enabled {
				return
			}
			cp.performHealthCheck()
		case <-cp.stopChan:
			return
		}
	}
}

// performHealthCheck performs health checks on configured hosts.
func (cp *ConnectionPool) performHealthCheck() {
	for _, host := range cp.healthCheckHosts {
		go func(host string) {
			req, err := cp.driver.Request("HEAD", host, nil)
			if err != nil {
				return
			}

			resp, err := cp.driver.Do(req)
			if err != nil {
				cp.statsMutex.Lock()
				cp.connectionFailures++
				cp.statsMutex.Unlock()
				return
			}

			if resp != nil {
				resp.Close()
			}
		}(host)
	}
}

// SetRetryPolicy configures the retry policy for failed connections.
func (cp *ConnectionPool) SetRetryPolicy(maxRetries int, retryDelay time.Duration) {
	cp.stopMutex.Lock()
	defer cp.stopMutex.Unlock()

	cp.maxRetries = maxRetries
	cp.retryDelay = retryDelay
}

// Stop stops all background operations.
func (cp *ConnectionPool) Stop() {
	cp.stopMutex.Lock()
	defer cp.stopMutex.Unlock()

	if !cp.stopped {
		close(cp.stopChan)
		cp.stopped = true
		cp.healthCheckEnabled = false
	}
}

// Stats returns connection pool statistics.
func (cp *ConnectionPool) Stats() ConnectionPoolStats {
	cp.statsMutex.RLock()
	totalConnections := cp.totalConnections
	activeConnections := cp.activeConnections
	connectionAttempts := cp.connectionAttempts
	connectionFailures := cp.connectionFailures
	cp.statsMutex.RUnlock()

	cp.stopMutex.RLock()
	healthCheckEnabled := cp.healthCheckEnabled
	healthCheckHosts := len(cp.healthCheckHosts)
	healthCheckInterval := cp.healthCheckInterval
	maxRetries := cp.maxRetries
	retryDelay := cp.retryDelay
	cp.stopMutex.RUnlock()

	var successRate float64
	if connectionAttempts > 0 {
		successRate = float64(connectionAttempts-connectionFailures) / float64(connectionAttempts)
	}

	return ConnectionPoolStats{
		TotalConnections:    totalConnections,
		ActiveConnections:   activeConnections,
		ConnectionAttempts:  connectionAttempts,
		ConnectionFailures:  connectionFailures,
		SuccessRate:         successRate,
		MaxRetries:          maxRetries,
		RetryDelay:          retryDelay,
		HealthCheckEnabled:  healthCheckEnabled,
		HealthCheckHosts:    healthCheckHosts,
		HealthCheckInterval: healthCheckInterval,
	}
}

// ResetStats resets all connection pool statistics.
func (cp *ConnectionPool) ResetStats() {
	cp.statsMutex.Lock()
	defer cp.statsMutex.Unlock()

	cp.totalConnections = 0
	cp.activeConnections = 0
	cp.connectionAttempts = 0
	cp.connectionFailures = 0
}

// ConnectionPoolStats contains connection pool statistics.
type ConnectionPoolStats struct {
	TotalConnections    int64
	ActiveConnections   int64
	ConnectionAttempts  int64
	ConnectionFailures  int64
	SuccessRate         float64
	MaxRetries          int
	RetryDelay          time.Duration
	HealthCheckEnabled  bool
	HealthCheckHosts    int
	HealthCheckInterval time.Duration
}
