# HTTP Client Package

The `httpclient` package provides an abstraction layer for HTTP client operations. It defines interfaces for making HTTP requests and handling responses, allowing different HTTP client implementations to be used interchangeably.

## Features

- **Two-phase request pattern**: Create a request, then execute it
- **Interface-based design**: Easy to mock and test
- **Flexible header management**: Set headers before execution
- **Resource management**: Proper cleanup with `Close()` method
- **Standard implementation**: Built-in implementation using Go's standard library

## Package Structure

```sh
httpclient/
├── driver.go          # Core interfaces
├── driver_test.go     # Interface tests
├── example_test.go    # Usage examples
├── batchclient/       # Batch processing extensions
│   ├── batch.go       # Batch executor
│   ├── ratelimiter.go # Rate limiting
│   ├── pool.go        # Connection pooling
│   ├── *_test.go      # Tests
│   └── README.md      # Batch client documentation
├── stdhttpclient/     # Standard implementation
│   ├── driver.go      # Driver implementation
│   ├── driver_test.go # Implementation tests
│   ├── request.go     # Request wrapper
│   ├── response.go    # Response wrapper
│   └── batch.go       # DEPRECATED: Use batchclient instead
└── README.md         # This file
```

## Interfaces

### Driver

The main interface for HTTP client operations:

```go
type Driver interface {
    Request(method, url string, body io.Reader) (Request, error)
    Do(req Request) (Response, error)
}
```

### Request

Represents a configurable HTTP request:

```go
type Request interface {
    SetHeader(key, value string)
}
```

### Response

Represents an HTTP response:

```go
type Response interface {
    io.Closer
    StatusCode() int
    Body() io.Reader
}
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "io"
    "log"
    
    "github.com/geniusrabbit/adcorelib/net/httpclient/stdhttpclient"
)

func main() {
    // Create a new HTTP client driver
    driver := stdhttpclient.NewDriver()
    
    // Phase 1: Create a request
    req, err := driver.Request("GET", "https://httpbin.org/get", nil)
    if err != nil {
        log.Fatal(err)
    }
    
    // Modify the request by setting headers
    req.SetHeader("User-Agent", "my-app/1.0")
    req.SetHeader("Accept", "application/json")
    
    // Phase 2: Execute the request
    resp, err := driver.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Close()
    
    // Process the response
    fmt.Printf("Status Code: %d\n", resp.StatusCode())
    
    // Read response body
    body, err := io.ReadAll(resp.Body())
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Response Body: %s\n", string(body))
}
```

### POST Request with JSON

```go
import (
    "strings"
    "github.com/geniusrabbit/adcorelib/net/httpclient/stdhttpclient"
)

func makePostRequest() {
    driver := stdhttpclient.NewDriver()
    
    // Create POST request with JSON body
    jsonData := `{"name":"John Doe","email":"john@example.com"}`
    req, err := driver.Request("POST", "https://httpbin.org/post", strings.NewReader(jsonData))
    if err != nil {
        log.Fatal(err)
    }
    
    // Set required headers
    req.SetHeader("Content-Type", "application/json")
    req.SetHeader("Authorization", "Bearer token123")
    
    // Execute the request
    resp, err := driver.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Close()
    
    fmt.Printf("Status Code: %d\n", resp.StatusCode())
}
```

### Custom HTTP Client

```go
import (
    "net/http"
    "time"
    "github.com/geniusrabbit/adcorelib/net/httpclient/stdhttpclient"
)

func useCustomClient() {
    // Create a custom HTTP client with timeout
    httpClient := &http.Client{
        Timeout: 30 * time.Second,
    }
    
    // Create driver with custom client
    driver := stdhttpclient.NewDriverWithHTTPClient(httpClient)
    
    // Use the driver as normal...
}
```

### Dependency Injection

```go
import "github.com/geniusrabbit/adcorelib/net/httpclient"

// This function accepts any implementation of httpclient.Driver
func makeRequest(client httpclient.Driver, url string) (int, error) {
    req, err := client.Request("GET", url, nil)
    if err != nil {
        return 0, err
    }
    
    req.SetHeader("User-Agent", "my-app/1.0")
    
    resp, err := client.Do(req)
    if err != nil {
        return 0, err
    }
    defer resp.Close()
    
    return resp.StatusCode(), nil
}
```

## Design Principles

### Two-Phase Request Pattern

The package uses a two-phase approach:

1. **Create**: Build a request object using `Driver.Request()`
2. **Execute**: Send the request using `Driver.Do()`

This pattern allows for:

- Request modification (headers, etc.) between creation and execution
- Middleware implementation
- Request reuse and modification
- Better testability

### Interface-Based Design

All functionality is exposed through interfaces, making it easy to:

- Mock for testing
- Implement custom HTTP clients
- Add middleware layers
- Switch between implementations

### Resource Management

The `Response` interface implements `io.Closer` to ensure proper cleanup:

- Always call `resp.Close()` when done
- Use `defer resp.Close()` for automatic cleanup
- Prevents resource leaks

## Testing

The package includes comprehensive tests:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run benchmarks
go test -bench=. ./...

# Run specific test
go test -run TestDriver_Request ./...
```

### Test Coverage

- Interface compliance tests
- Standard implementation tests
- Error handling tests
- Integration tests with real HTTP servers
- Benchmark tests
- Example tests

## Implementation Details

### Standard Implementation

The `stdhttpclient` package provides a standard implementation using Go's `net/http` package:

- `Driver`: Wraps `http.Client`
- `Request`: Wraps `http.Request`
- `Response`: Wraps `http.Response`

### Performance

The implementation is lightweight and efficient:

- Minimal overhead over standard `net/http`
- No unnecessary allocations
- Proper resource management
- Connection pooling (via `http.Client`)

## Extended Modules

### BatchClient

The `batchclient` package provides high-performance batch processing capabilities that work with any `httpclient.Driver` implementation:

- **BatchExecutor**: Concurrent batch execution with configurable worker pools
- **RateLimitedExecutor**: Rate-limited execution with token bucket algorithm
- **ConnectionPool**: Connection pooling with retry logic and health checks
- **Chainable**: All components implement the Driver interface and can be chained
- **Statistics**: Built-in monitoring and performance metrics

See [batchclient/README.md](batchclient/README.md) for detailed documentation and examples.

## Contributing

When contributing to this package:

1. Maintain interface compatibility
2. Add tests for new functionality
3. Follow Go conventions and best practices
4. Document public APIs with GoDoc comments
5. Ensure proper error handling
6. Add examples for new features

## License

This package is part of the adcorelib project and follows the same license terms.

## Batch Processing Examples

### Batch Processing

```go
import (
    "context"
    "time"
    "github.com/geniusrabbit/adcorelib/net/httpclient"
    "github.com/geniusrabbit/adcorelib/net/httpclient/batchclient"
    "github.com/geniusrabbit/adcorelib/net/httpclient/stdhttpclient"
)

func batchProcessing() {
    // Create base driver
    driver := stdhttpclient.NewDriver()

    // Wrap with batch execution (100 concurrent workers)
    batchExecutor := batchclient.NewBatchExecutor(driver, 100)

    // Create multiple requests
    var requests []httpclient.Request
    for i := 0; i < 50; i++ {
        req, err := batchExecutor.Request("GET", fmt.Sprintf("https://api.example.com/data/%d", i), nil)
        if err != nil {
            log.Printf("Failed to create request %d: %v", i, err)
            continue
        }
        requests = append(requests, req)
    }

    // Execute all requests concurrently
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    result := batchExecutor.ExecuteBatch(ctx, requests)

    fmt.Printf("Batch completed: %d/%d requests (errors: %d) in %v\n", 
        result.Completed, len(result.Requests), result.Errors, result.Duration)

    // Clean up
    for _, req := range result.Requests {
        if req.Response != nil {
            req.Response.Close()
        }
    }
}
```

### Chained Executors

```go
func chainedExecutors() {
    // Create base driver
    driver := stdhttpclient.NewDriver()

    // Chain executors: Pool -> Rate Limiter -> Batch Executor
    pool := batchclient.NewConnectionPool(driver)
    rateLimiter := batchclient.NewRateLimitedExecutor(pool, 10, 5) // 10 req/sec
    batchExecutor := batchclient.NewBatchExecutor(rateLimiter, 20)

    // Now you have connection pooling, rate limiting, and batch processing
    // Use batchExecutor for all operations
    
    // ... create and execute requests ...
    
    // Clean up
    rateLimiter.Stop()
    pool.Stop()
}
```
