package stdhttpclient

import "time"

// PerformanceConfig holds configuration for high-performance HTTP client settings.
type PerformanceConfig struct {
	// Connection pool settings
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	MaxConnsPerHost     int

	// Timeout settings
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ResponseHeaderTimeout time.Duration
	ExpectContinueTimeout time.Duration
	RequestTimeout        time.Duration

	// Performance settings
	DisableCompression bool
	ForceAttemptHTTP2  bool
	DisableKeepAlives  bool
}

// DefaultPerformanceConfig returns default high-performance configuration.
func DefaultPerformanceConfig() PerformanceConfig {
	return PerformanceConfig{
		MaxIdleConns:          1000,
		MaxIdleConnsPerHost:   100,
		MaxConnsPerHost:       200,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		RequestTimeout:        30 * time.Second,
		DisableCompression:    false,
		ForceAttemptHTTP2:     true,
		DisableKeepAlives:     false,
	}
}

// ExtremePerformanceConfig returns configuration optimized for maximum throughput.
// Use this configuration when you need to send thousands of requests per second.
func ExtremePerformanceConfig() PerformanceConfig {
	return PerformanceConfig{
		MaxIdleConns:          2000,
		MaxIdleConnsPerHost:   200,
		MaxConnsPerHost:       500,
		IdleConnTimeout:       120 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		ExpectContinueTimeout: 500 * time.Millisecond,
		RequestTimeout:        15 * time.Second,
		DisableCompression:    true, // Disable compression for better CPU performance
		ForceAttemptHTTP2:     true,
		DisableKeepAlives:     false,
	}
}
