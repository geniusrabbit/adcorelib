package stdhttpclient

import (
	"crypto/tls"
	"io"
	"net/http"
	"sync"

	"github.com/geniusrabbit/adcorelib/net/httpclient"
)

// Driver is a high-performance HTTP client driver that implements the httpclient.Driver interface.
// It uses the standard library's http.Client with optimizations for high throughput scenarios.
// This driver is optimized for maximum performance with minimal allocations and persistent connections.
type Driver struct {
	HTTPClient *http.Client

	// Object pools for reducing allocations
	requestPool  sync.Pool
	responsePool sync.Pool
}

// NewDriver creates a new Driver with a high-performance HTTP client.
// The client is optimized for high throughput with:
// - Persistent connections and connection pooling
// - Aggressive keepalive settings
// - Optimized timeouts for low latency
// - Large connection pool limits
func NewDriver() *Driver {
	return NewDriverWithHTTPClient(newHighPerformanceClient())
}

// NewHighPerformanceDriver creates a new Driver with maximum performance settings.
// This function provides the most aggressive optimization for high-throughput scenarios.
func NewHighPerformanceDriver() *Driver {
	return NewDriverWithHTTPClient(newHighPerformanceClient())
}

// NewDriverWithConfig creates a new Driver with custom performance configuration.
func NewDriverWithConfig(config PerformanceConfig) *Driver {
	transport := &http.Transport{
		MaxIdleConns:          config.MaxIdleConns,
		MaxIdleConnsPerHost:   config.MaxIdleConnsPerHost,
		MaxConnsPerHost:       config.MaxConnsPerHost,
		IdleConnTimeout:       config.IdleConnTimeout,
		TLSHandshakeTimeout:   config.TLSHandshakeTimeout,
		ResponseHeaderTimeout: config.ResponseHeaderTimeout,
		ExpectContinueTimeout: config.ExpectContinueTimeout,
		DisableCompression:    config.DisableCompression,
		ForceAttemptHTTP2:     config.ForceAttemptHTTP2,
		DisableKeepAlives:     config.DisableKeepAlives,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // For testing purposes, disable TLS verification
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.RequestTimeout,
	}

	return NewDriverWithHTTPClient(client)
}

// NewExtremePerformanceDriver creates a driver optimized for extreme throughput.
func NewExtremePerformanceDriver() *Driver {
	return NewDriverWithConfig(ExtremePerformanceConfig())
}

// NewDriverWithHTTPClient creates a new Driver with the specified HTTP client.
// This allows for custom HTTP client configurations, such as timeouts or transport settings.
// The driver is initialized with object pools for optimal performance.
func NewDriverWithHTTPClient(client *http.Client) *Driver {
	d := &Driver{HTTPClient: client}

	// Initialize object pools to reduce allocations
	d.requestPool = sync.Pool{
		New: func() any { return &Request{driver: d} },
	}
	d.responsePool = sync.Pool{
		New: func() any { return &Response{driver: d} },
	}

	return d
}

// Request creates a new HTTP request with the specified method, URL, and body.
// It uses object pooling to minimize allocations and improve performance.
func (d *Driver) Request(method, url string, body io.Reader) (httpclient.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Get a pooled Request object
	pooledReq := d.requestPool.Get().(*Request)
	pooledReq.HTTP = req
	// driver already set in pool New function

	return pooledReq, nil
}

// Do sends the HTTP request and returns the response.
// It uses object pooling to minimize allocations and improve performance.
func (d *Driver) Do(req httpclient.Request) (httpclient.Response, error) {
	httpReq := req.(*Request)
	resp, err := d.HTTPClient.Do(httpReq.HTTP)
	if err != nil {
		// Return the request to the pool on error
		d.releaseRequest(httpReq)
		return nil, err
	}

	// Get a pooled Response object
	pooledResp := d.responsePool.Get().(*Response)
	pooledResp.HTTP = resp
	// driver already set in pool New function
	pooledResp.request = httpReq // Store reference to request for cleanup

	return pooledResp, nil
}

// Close releases resources held by the driver.
// It closes the HTTP client and clears the object pools.
// This should be called when the driver is no longer needed to prevent resource leaks.
// It is safe to call this multiple times.
func (d *Driver) Close() {
	// Close the HTTP client to release resources
	if d.HTTPClient != nil && d.HTTPClient.Transport != nil {
		if transport, ok := d.HTTPClient.Transport.(*http.Transport); ok {
			transport.CloseIdleConnections()
		}
	}

	// Clear object pools
	d.requestPool = sync.Pool{}
	d.responsePool = sync.Pool{}
}

// releaseRequest returns a Request object to the pool for reuse.
func (d *Driver) releaseRequest(req *Request) {
	if req != nil {
		req.clear()
		d.requestPool.Put(req)
	}
}

// releaseResponse returns a Response object to the pool for reuse.
func (d *Driver) releaseResponse(resp *Response) {
	if resp != nil {
		resp.clear()
		d.responsePool.Put(resp)
	}
}
