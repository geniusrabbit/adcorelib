//
// @project GeniusRabbit 2016, 2021
// @author Dmitry Ponomarev <demdxx@gmail.com>
//

package utils

// CachePool object
type CachePool struct {
	pool chan *ResponseCache
}

// NewCachePool with max size
func NewCachePool(size int) *CachePool {
	return &CachePool{
		pool: make(chan *ResponseCache, size),
	}
}

// Borrow cache
func (p *CachePool) Borrow() (c *ResponseCache) {
	select {
	case c = <-p.pool:
		break
	default:
		c = &ResponseCache{pool: p}
	}
	return c
}

// Return returns a Cacher to the pool.
func (p *CachePool) Return(c *ResponseCache) {
	if c == nil {
		return
	}
	c.Reset()
	select {
	case p.pool <- c:
		break
	default:
		// let it go, let it go...
	}
}
