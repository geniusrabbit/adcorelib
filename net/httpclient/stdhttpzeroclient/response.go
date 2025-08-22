package stdhttpzeroclient

import (
	"io"
	"net/http"
)

// Response is a high-performance wrapper around http.Response that implements the httpclient.Response interface.
// It provides methods to access the response status code, body, and to close the response body.
// The response uses object pooling to minimize allocations.
type Response struct {
	HTTP   *http.Response
	driver *ZeroAllocDriver // Reference to driver for object pooling
}

// StatusCode returns the HTTP status code of the response.
// If the response is nil, it returns 0.
func (r *Response) StatusCode() int {
	if r == nil {
		return 0
	}
	return r.HTTP.StatusCode
}

// Body returns the body of the HTTP response as an io.Reader.
// If the response is nil, it returns nil.
func (r *Response) Body() io.Reader {
	return r.HTTP.Body
}

// Close closes the response body and returns objects to the pool for reuse.
// This method should always be called when the response is no longer needed
// to prevent memory leaks and to return pooled objects for reuse.
func (r *Response) Close() error {
	var err error
	if r.HTTP != nil && r.HTTP.Body != nil {
		err = r.HTTP.Body.Close()
	}

	// Return this response to the pool
	if r.driver != nil {
		r.driver.releaseResponse(r)
	}

	return err
}

func (r *Response) clear() {
	r.HTTP = nil
	r.driver = nil
}
