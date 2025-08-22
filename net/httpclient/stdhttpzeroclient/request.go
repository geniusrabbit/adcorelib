package stdhttpzeroclient

import (
	"io"
	"net/http"
	"sync"

	"golang.org/x/exp/maps"
)

// FastRequest is a zero-allocation request implementation for high-performance scenarios.
type FastRequest struct {
	method    string
	url       string
	body      io.Reader
	headers   map[string]string
	headersMu sync.RWMutex
}

// NewFastRequest creates a new FastRequest with minimal allocations.
func NewFastRequest(method, urlStr string, body io.Reader) *FastRequest {
	return &FastRequest{
		method:  method,
		url:     urlStr,
		body:    body,
		headers: make(map[string]string),
	}
}

// SetHeader sets a header for the request.
func (r *FastRequest) SetHeader(key, value string) {
	r.headersMu.Lock()
	r.headers[key] = value
	r.headersMu.Unlock()
}

// SetContentType sets the Content-Type header for the request.
func (r *FastRequest) SetContentType(contentType string) {
	r.SetHeader("Content-Type", contentType)
}

// ToHTTPRequest converts FastRequest to http.Request.
func (r *FastRequest) ToHTTPRequest() (*http.Request, error) {
	req, err := http.NewRequest(r.method, r.url, r.body)
	if err != nil {
		return nil, err
	}

	r.headersMu.RLock()
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}
	r.headersMu.RUnlock()

	return req, nil
}

// Reset resets the request for reuse.
func (r *FastRequest) Reset() {
	r.method = ""
	r.url = ""
	r.body = nil
	r.headersMu.Lock()
	maps.Clear(r.headers)
	r.headersMu.Unlock()
}

func (r *FastRequest) clear() {
	r.Reset()
}
