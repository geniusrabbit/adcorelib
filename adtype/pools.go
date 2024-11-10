package adtype

import (
	"sync"
)

// Common pools
var (
	responsePool = sync.Pool{
		New: func() any { return new(Response) },
	}
)

///////////////////////////////////////////////////////////////////////////////
/// Response sync pool
///////////////////////////////////////////////////////////////////////////////

// BorrowResponse object
func BorrowResponse(request *BidRequest, source Source, items []ResponserItemCommon, err error) *Response {
	resp := responsePool.Get().(*Response)
	resp.context = request.Ctx
	resp.request = request
	resp.source = source
	resp.items = items
	resp.err = err
	return resp
}

// ReturnResponse back to pool
func ReturnResponse(o *Response) {
	o.reset()
	responsePool.Put(o)
}
