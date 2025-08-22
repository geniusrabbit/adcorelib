package stdhttpzeroclient

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/geniusrabbit/adcorelib/net/httpclient"
)

// ZeroAllocDriver is a driver optimized for zero allocations in hot paths.
type ZeroAllocDriver struct {
	HTTPClient   *http.Client
	requestPool  sync.Pool
	responsePool sync.Pool
	bufferPool   sync.Pool
}

// NewDriver creates a new zero-allocation driver.
func NewDriver() *ZeroAllocDriver {
	return NewDriverWithHTTPClient(newHighPerformanceClient())
}

// NewDriverWithConfig creates a new Driver with custom performance configuration.
func NewDriverWithConfig(config PerformanceConfig) *ZeroAllocDriver {
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
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   config.RequestTimeout,
	}

	return NewDriverWithHTTPClient(client)
}

// NewExtremePerformanceDriver creates a driver optimized for extreme throughput.
func NewExtremePerformanceDriver() *ZeroAllocDriver {
	return NewDriverWithConfig(ExtremePerformanceConfig())
}

// NewDriverWithHTTPClient creates a new Driver with the specified HTTP client.
// This allows for custom HTTP client configurations, such as timeouts or transport settings.
// The driver is initialized with object pools for optimal performance.
func NewDriverWithHTTPClient(client *http.Client) *ZeroAllocDriver {
	d := &ZeroAllocDriver{HTTPClient: client}

	d.requestPool = sync.Pool{
		New: func() any {
			return &FastRequest{
				headers: make(map[string]string),
			}
		},
	}

	d.responsePool = sync.Pool{
		New: func() any { return &Response{} },
	}

	d.bufferPool = sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, 1024))
		},
	}

	return d
}

// Request creates a new request with minimal allocations.
func (d *ZeroAllocDriver) Request(method, url string, body io.Reader) (httpclient.Request, error) {
	fastReq := d.requestPool.Get().(*FastRequest)
	fastReq.method = method
	fastReq.url = url
	fastReq.body = body

	return fastReq, nil
}

// Do executes the request.
func (d *ZeroAllocDriver) Do(req httpclient.Request) (httpclient.Response, error) {
	fastReq := req.(*FastRequest)

	// Convert to http.Request
	httpReq, err := fastReq.ToHTTPRequest()
	if err != nil {
		d.releaseRequest(fastReq)
		return nil, err
	}

	// Execute request
	resp, err := d.HTTPClient.Do(httpReq)
	if err != nil {
		d.releaseRequest(fastReq)
		return nil, err
	}

	// Get pooled response
	pooledResp := d.responsePool.Get().(*Response)
	pooledResp.HTTP = resp
	pooledResp.driver = d // Set driver reference for pooling

	// Return request to pool
	d.releaseRequest(fastReq)

	return pooledResp, nil
}

// GetBuffer gets a buffer from the pool.
func (d *ZeroAllocDriver) GetBuffer() *bytes.Buffer {
	return d.bufferPool.Get().(*bytes.Buffer)
}

// PutBuffer returns a buffer to the pool.
func (d *ZeroAllocDriver) PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	d.bufferPool.Put(buf)
}

// StringRequest creates a request with string body using pooled buffer.
func (d *ZeroAllocDriver) StringRequest(method, url, body string) (httpclient.Request, error) {
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	return d.Request(method, url, bodyReader)
}

// JSONRequest creates a request with JSON body using pooled buffer.
func (d *ZeroAllocDriver) JSONRequest(method, url string, jsonBytes []byte) (httpclient.Request, error) {
	var bodyReader io.Reader
	if len(jsonBytes) > 0 {
		bodyReader = bytes.NewReader(jsonBytes)
	}

	req, err := d.Request(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.SetHeader("Content-Type", "application/json")
	return req, nil
}

// FormRequest creates a request with form data using pooled buffer.
func (d *ZeroAllocDriver) FormRequest(method, url string, formData url.Values) (httpclient.Request, error) {
	var bodyReader io.Reader
	if len(formData) > 0 {
		bodyReader = strings.NewReader(formData.Encode())
	}

	req, err := d.Request(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.SetHeader("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

// releaseRequest returns a Request object to the pool for reuse.
func (d *ZeroAllocDriver) releaseRequest(req *FastRequest) {
	if req != nil {
		req.clear()
		d.requestPool.Put(req)
	}
}

// releaseResponse returns a Response object to the pool for reuse.
func (d *ZeroAllocDriver) releaseResponse(resp *Response) {
	if resp != nil {
		resp.clear()
		d.responsePool.Put(resp)
	}
}
