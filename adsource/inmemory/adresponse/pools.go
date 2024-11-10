package adresponse

import (
	"sync"

	"github.com/geniusrabbit/adcorelib/rand"
)

// Common pools
var (
	responseAdItemPool = sync.Pool{
		New: func() any { return new(ResponseAdItem) },
	}
)

///////////////////////////////////////////////////////////////////////////////
/// Response ad item sync pool
///////////////////////////////////////////////////////////////////////////////

// BorrowResponseAdItem object
func BorrowResponseAdItem() *ResponseAdItem {
	item := responseAdItemPool.Get().(*ResponseAdItem)
	item.ItemID = rand.UUID()
	return item
}

// ReturnResponseAdItem back to pool
func ReturnResponseAdItem(o *ResponseAdItem) {
	o.reset()
	responseAdItemPool.Put(o)
}
