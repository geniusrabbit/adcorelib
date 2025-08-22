package stdhttpzeroclient

import "github.com/geniusrabbit/adcorelib/net/httpclient/stdhttpclient"

// PerformanceConfig is an alias for stdhttpclient.PerformanceConfig.
type PerformanceConfig = stdhttpclient.PerformanceConfig

// DefaultPerformanceConfig returns default high-performance configuration.
func DefaultPerformanceConfig() PerformanceConfig {
	return stdhttpclient.DefaultPerformanceConfig()
}

// ExtremePerformanceConfig returns configuration optimized for maximum throughput.
// Use this configuration when you need to send thousands of requests per second.
func ExtremePerformanceConfig() PerformanceConfig {
	return stdhttpclient.ExtremePerformanceConfig()
}
