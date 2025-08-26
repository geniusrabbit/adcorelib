//
// @project GeniusRabbit corelib 2017 - 2019, 2025
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019, 2025
//

package adtype

import (
	"context"
	"iter"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// ResponseError from different sources
type ResponseError struct {
	context context.Context
	request *BidRequest
	err     error
}

// NewErrorResponse object
func NewErrorResponse(request *BidRequest, err error) *ResponseError {
	return &ResponseError{
		request: request,
		err:     err,
		context: request.Ctx,
	}
}

// AuctionID response
func (r *ResponseError) AuctionID() string {
	if r == nil || r.request == nil {
		return ""
	}
	return r.request.ID
}

// AuctionType of request
func (r *ResponseError) AuctionType() types.AuctionType {
	if r == nil || r.request == nil {
		return types.UndefinedAuctionType
	}
	return r.request.AuctionType
}

// Source of response
func (r *ResponseError) Source() Source { return nil }

// Request information
func (r *ResponseError) Request() *BidRequest { return r.request }

// AddItem to response
func (r *ResponseError) AddItem(it ResponserItemCommon) {
	panic("error response can't add item")
}

// Item by impression code
func (r *ResponseError) Item(impid string) ResponserItemCommon { return nil }

// Ads list
func (r *ResponseError) Ads() []ResponserItemCommon { return nil }

// IterAds returns an iterator over the response items.
func (r *ResponseError) IterAds() iter.Seq[ResponserItem] {
	return func(yield func(ResponserItem) bool) {}
}

// Count of response items
func (r *ResponseError) Count() int { return 0 }

// Validate response
func (r *ResponseError) Validate() (err error) {
	if r.err != nil {
		return r.err
	}
	return ErrResponseEmpty
}

// Error of the response
func (r *ResponseError) Error() error {
	return r.err
}

// Context value
func (r *ResponseError) Context(ctx ...context.Context) context.Context {
	if r != nil && len(ctx) > 0 {
		oldContext := r.context
		r.context = ctx[0]
		return oldContext
	}
	if r == nil || r.context == nil {
		return context.Background()
	}
	return r.context
}

// Get context item by key
func (r *ResponseError) Get(key string) any {
	if r.context == nil {
		return nil
	}
	return r.context.Value(key)
}

// Release response and all linked objects
func (r *ResponseError) Release() {
	if r != nil {
		r.reset()
	}
}

func (r *ResponseError) reset() {
	r.err = nil
	r.context = nil
	r.request = nil
}

var (
	_ Responser = &ResponseError{}
)
