package stdhttpzeroclient

import (
	"net/http"
	"time"
)

// newHighPerformanceClient creates an optimized HTTP client for high throughput.
func newHighPerformanceClient() *http.Client {
	transport := &http.Transport{
		// Connection pooling settings
		MaxIdleConns:        1000, // Maximum idle connections across all hosts
		MaxIdleConnsPerHost: 100,  // Maximum idle connections per host
		MaxConnsPerHost:     200,  // Maximum connections per host

		// Keepalive settings for persistent connections
		IdleConnTimeout:     90 * time.Second, // How long idle connections stay open
		TLSHandshakeTimeout: 10 * time.Second, // TLS handshake timeout

		// Response settings
		ResponseHeaderTimeout: 10 * time.Second, // Time to wait for response headers
		ExpectContinueTimeout: 1 * time.Second,  // Time to wait for 100-continue

		// Disable compression for better performance if not needed
		DisableCompression: false, // Keep compression enabled by default

		// Force HTTP/2 for better multiplexing
		ForceAttemptHTTP2: true,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Overall request timeout
	}
}
