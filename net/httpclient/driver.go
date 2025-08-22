// Package httpclient provides an abstraction layer for HTTP client operations.
// It defines interfaces for making HTTP requests and handling responses,
// allowing for different HTTP client implementations to be used interchangeably.
//
// The package follows a two-phase request pattern:
// 1. Create a request using Driver.Request()
// 2. Execute the request using Driver.Do()
//
// This design allows for request modification (headers, etc.) between
// creation and execution, providing flexibility for middleware and
// request customization.
package httpclient

import "io"

// Request represents an HTTP request that can be configured before execution.
// It provides methods to modify request properties such as headers.
type Request interface {
	// SetHeader sets a header key-value pair for the request.
	// If the header already exists, it will be replaced.
	// The key is case-insensitive as per HTTP specification.
	SetHeader(key, value string)
}

// Response represents an HTTP response received from a server.
// It provides access to the response status code and body,
// and implements io.Closer to ensure proper resource cleanup.
type Response interface {
	// Close closes the response body and releases associated resources.
	// It should be called when the response is no longer needed.
	io.Closer

	// StatusCode returns the HTTP status code of the response.
	// Standard HTTP status codes are defined in the net/http package.
	StatusCode() int

	// Body returns an io.Reader for the response body.
	// The caller is responsible for reading and closing the body.
	Body() io.Reader
}

// Driver is the main interface for HTTP client operations.
// It provides a two-phase approach to HTTP requests:
// first creating a configurable request, then executing it.
//
// Implementations of this interface should handle:
// - Connection management and pooling
// - Timeout handling
// - Error handling and retries
// - Security considerations (TLS, certificates, etc.)
type Driver interface {
	// Request creates a new HTTP request with the specified method, URL, and body.
	// The request is not executed immediately; it must be passed to Do() for execution.
	//
	// Parameters:
	//   method: HTTP method (GET, POST, PUT, DELETE, etc.)
	//   url: Target URL for the request
	//   body: Request body as an io.Reader, or nil for methods without body
	//
	// Returns:
	//   Request: A configurable request object that can be modified before execution
	//   error: An error if the request cannot be created
	Request(method, url string, body io.Reader) (Request, error)

	// Do executes the HTTP request and returns the response.
	// This method performs the actual network communication.
	//
	// Parameters:
	//   req: The request to execute, previously created by Request()
	//
	// Returns:
	//   Response: The HTTP response, which must be closed by the caller
	//   error: An error if the request fails or cannot be executed
	Do(req Request) (Response, error)
}
