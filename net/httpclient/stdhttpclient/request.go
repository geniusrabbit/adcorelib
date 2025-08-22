package stdhttpclient

import "net/http"

// Request is a high-performance wrapper around http.Request that implements the httpclient.Request interface.
// It provides methods to set headers and access the underlying HTTP request.
// The request uses object pooling to minimize allocations.
type Request struct {
	HTTP   *http.Request
	driver *Driver // Reference to driver for object pooling
}

// SetHeader sets a header for the HTTP request.
func (r *Request) SetHeader(key, value string) {
	r.HTTP.Header.Set(key, value)
}

func (r *Request) clear() {
	r.HTTP = nil
	r.driver = nil
}
