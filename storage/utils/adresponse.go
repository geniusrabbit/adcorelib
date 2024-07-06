package utils

import "github.com/geniusrabbit/adcorelib/admodels"

// ResponseCache object
type ResponseCache struct {
	pool *CachePool
	data []*admodels.VirtualAd
}

// Data of the response
func (rs *ResponseCache) Data() []*admodels.VirtualAd {
	return rs.data
}

// Reset response cache
func (rs *ResponseCache) Reset() {
	if len(rs.data) > 0 {
		rs.data = rs.data[:0]
	}
}

// Release the cache response
func (rs *ResponseCache) Release() {
	rs.pool.Return(rs)
}
