# BatchClient - Generic HTTP Client Batch Processing

The `batchclient` package provides high-performance batch processing capabilities for any `httpclient.Driver` implementation. It offers generic wrappers that can enhance any HTTP client with batch execution, rate limiting, and connection pooling features.

## Features

- **Interface-Based Design**: Works with any `httpclient.Driver` implementation
- **Batch Processing**: Concurrent execution of multiple HTTP requests
- **Rate Limiting**: Configurable request rate limits
- **Connection Pooling**: Connection management with retry logic and health checks
- **Chainable**: Components can be combined for complex request processing
- **Statistics**: Built-in monitoring and statistics collection
- **Context Support**: Full context cancellation support

## Components

### BatchExecutor

Provides concurrent batch execution of HTTP requests with configurable worker pools.

```go
// Create batch executor with 100 workers
executor := batchclient.NewBatchExecutor(driver, 100)

// Execute multiple requests concurrently
result := executor.ExecuteBatch(ctx, requests)
```

**Features:**

- Configurable worker pool size
- Concurrent request processing
- Request ordering preservation
- Execution statistics
- Context cancellation support

### RateLimitedExecutor

Provides rate-limited execution of HTTP requests with configurable request rates.

```go
// Create rate limiter (10 requests per second)
limiter := batchclient.NewRateLimitedExecutor(driver, 10, 5)

// Execute request with rate limiting
resp, err := limiter.Execute(ctx, req)
```

**Features:**

- Configurable requests per second limit
- Token bucket algorithm
- Dynamic rate adjustment
- Execution statistics
- Graceful shutdown

### ConnectionPool

Provides connection pooling with retry logic and health monitoring.

```go
// Create connection pool
pool := batchclient.NewConnectionPool(driver)

// Configure retry policy
pool.SetRetryPolicy(3, 100*time.Millisecond)

// Warmup connections
pool.WarmupConnections(ctx, hosts, 2)

// Enable health checking
pool.EnableHealthCheck(hosts, 30*time.Second)
```

**Features:**

- Connection warmup
- Configurable retry policies
- Health check monitoring
- Connection statistics
- Failure recovery

## Usage Examples

### Basic Batch Processing

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/geniusrabbit/adcorelib/net/httpclient"
    "github.com/geniusrabbit/adcorelib/net/httpclient/batchclient"
    "github.com/geniusrabbit/adcorelib/net/httpclient/stdhttpclient"
)

func main() {
    // Create base driver
    driver := stdhttpclient.NewDriver()

    // Wrap with batch executor
    executor := batchclient.NewBatchExecutor(driver, 50)

    // Create requests
    var requests []httpclient.Request
    for i := 0; i < 100; i++ {
        req, err := executor.Request("GET", fmt.Sprintf("https://api.example.com/data/%d", i), nil)
        if err != nil {
            log.Printf("Failed to create request %d: %v", i, err)
            continue
        }
        requests = append(requests, req)
    }

    // Execute batch
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    result := executor.ExecuteBatch(ctx, requests)
    
    fmt.Printf("Completed %d/%d requests in %v\n", 
        result.Completed, len(result.Requests), result.Duration)
    
    // Clean up
    for _, req := range result.Requests {
        if req.Response != nil {
            req.Response.Close()
        }
    }
}
```

### Rate Limited Execution

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/geniusrabbit/adcorelib/net/httpclient/batchclient"
    "github.com/geniusrabbit/adcorelib/net/httpclient/stdhttpclient"
)

func main() {
    // Create base driver
    driver := stdhttpclient.NewDriver()

    // Wrap with rate limiter (5 requests per second)
    limiter := batchclient.NewRateLimitedExecutor(driver, 5, 10)
    defer limiter.Stop()

    ctx := context.Background()
    
    // Execute requests with rate limiting
    for i := 0; i < 20; i++ {
        req, err := limiter.Request("GET", "https://api.example.com/limited", nil)
        if err != nil {
            log.Printf("Failed to create request %d: %v", i, err)
            continue
        }

        start := time.Now()
        resp, err := limiter.Execute(ctx, req)
        if err != nil {
            log.Printf("Failed to execute request %d: %v", i, err)
            continue
        }

        fmt.Printf("Request %d completed in %v (status: %d)\n", 
            i+1, time.Since(start), resp.StatusCode())
        resp.Close()
    }
}
```

### Chained Executors

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/geniusrabbit/adcorelib/net/httpclient"
    "github.com/geniusrabbit/adcorelib/net/httpclient/batchclient"
    "github.com/geniusrabbit/adcorelib/net/httpclient/stdhttpclient"
)

func main() {
    // Create base driver
    driver := stdhttpclient.NewDriver()

    // Chain executors: Pool -> Rate Limiter -> Batch Executor
    pool := batchclient.NewConnectionPool(driver)
    pool.SetRetryPolicy(3, 100*time.Millisecond)

    rateLimiter := batchclient.NewRateLimitedExecutor(pool, 20, 10)
    batchExecutor := batchclient.NewBatchExecutor(rateLimiter, 50)

    // Warmup connections
    hosts := []string{"https://api.example.com"}
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err := pool.WarmupConnections(ctx, hosts, 3)
    if err != nil {
        log.Printf("Warmup error: %v", err)
    }

    // Create requests
    var requests []httpclient.Request
    for i := 0; i < 100; i++ {
        req, err := batchExecutor.Request("GET", fmt.Sprintf("https://api.example.com/data/%d", i), nil)
        if err != nil {
            log.Printf("Failed to create request %d: %v", i, err)
            continue
        }
        requests = append(requests, req)
    }

    // Execute with all features
    start := time.Now()
    result := batchExecutor.ExecuteBatch(ctx, requests)
    elapsed := time.Since(start)

    fmt.Printf("Executed %d requests in %v\n", result.Completed, elapsed)
    fmt.Printf("Errors: %d\n", result.Errors)

    // Print statistics
    fmt.Println("Statistics:")
    batchStats := batchExecutor.Stats()
    fmt.Printf("  Batch: %d total, %d completed, %d errors\n", 
        batchStats.TotalRequests, batchStats.TotalCompleted, batchStats.TotalErrors)

    rateStats := rateLimiter.Stats()
    fmt.Printf("  Rate Limiter: %d total, %d executed, %d/s rate\n", 
        rateStats.TotalRequests, rateStats.TotalExecuted, rateStats.RequestsPerSecond)

    poolStats := pool.Stats()
    fmt.Printf("  Pool: %d connections, %.2f%% success rate\n", 
        poolStats.TotalConnections, poolStats.SuccessRate*100)

    // Clean up
    for _, req := range result.Requests {
        if req.Response != nil {
            req.Response.Close()
        }
    }
    
    rateLimiter.Stop()
    pool.Stop()
}
```

## Statistics and Monitoring

All components provide detailed statistics:

### BatchExecutor Stats

- Total requests processed
- Completed requests
- Error count
- Total execution time
- Average execution time
- Error rate

### RateLimitedExecutor Stats

- Current rate limit (requests/second)
- Total requests received
- Total requests executed
- Throttled requests
- Queue size and capacity

### ConnectionPool Stats

- Total connections created
- Active connections
- Connection attempts
- Connection failures
- Success rate
- Retry configuration
- Health check status

## Error Handling

The package provides robust error handling:

- **Context Cancellation**: All operations respect context cancellation
- **Timeout Handling**: Configurable timeouts for all operations
- **Retry Logic**: Configurable retry policies with exponential backoff
- **Error Aggregation**: Batch operations collect and report individual errors
- **Graceful Shutdown**: Clean shutdown with resource cleanup

## Performance Considerations

- **Worker Pool Sizing**: Configure worker pools based on your concurrency requirements
- **Rate Limiting**: Set appropriate rate limits to avoid overwhelming target servers
- **Connection Warmup**: Use connection warmup for better initial performance
- **Resource Cleanup**: Always close responses and stop executors when done
- **Memory Usage**: Monitor memory usage with large request batches

## Thread Safety

All components are thread-safe and can be used concurrently from multiple goroutines.

## Dependencies

- `github.com/geniusrabbit/adcorelib/net/httpclient`: Core HTTP client interfaces
- Standard library: `context`, `sync`, `time`

## License

This package is part of the adcorelib project and follows the same license terms.
