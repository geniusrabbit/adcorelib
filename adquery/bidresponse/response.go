//
// @project GeniusRabbit corelib 2017 - 2019, 2025
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019, 2025
//

package bidresponse

import (
	"context"
	"iter"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
)

// Response from different sources
type Response struct {
	request adtype.BidRequester
	source  adtype.Source
	items   []adtype.ResponseItemCommon
	err     error
	context context.Context
}

// NewResponse common object
func NewResponse(request adtype.BidRequester, source adtype.Source, items []adtype.ResponseItemCommon, err error) *Response {
	return &Response{
		request: request,
		source:  source,
		items:   items,
		err:     err,
		context: request.Context(),
	}
}

// AuctionID response
func (r *Response) AuctionID() string {
	if r == nil || r.request == nil {
		return ""
	}
	return r.request.AuctionID()
}

// AuctionType of request
func (r *Response) AuctionType() types.AuctionType {
	if r == nil || r.request == nil {
		return types.UndefinedAuctionType
	}
	return r.request.AuctionType()
}

// Source of response
func (r *Response) Source() adtype.Source {
	return r.source
}

// Request information
func (r *Response) Request() adtype.BidRequester {
	return r.request
}

// AddItem to response
func (r *Response) AddItem(it adtype.ResponseItemCommon) {
	r.items = append(r.items, it)
}

// Item by impression code
func (r *Response) Item(impid string) adtype.ResponseItemCommon {
	for _, it := range r.items {
		if it.ImpressionID() == impid {
			return it
		}
	}
	return nil
}

// Ads list
func (r *Response) Ads() []adtype.ResponseItemCommon {
	return r.items
}

// IterAds returns an iterator over the response items.
func (r *Response) IterAds() iter.Seq[adtype.ResponseItem] {
	return func(yield func(adtype.ResponseItem) bool) {
		for _, it := range r.items {
			switch itV := it.(type) {
			case nil:
			case adtype.ResponseItem:
				if !yield(itV) {
					return
				}
			case adtype.ResponseMultipleItem:
				for _, mit := range itV.Ads() {
					if !yield(mit) {
						return
					}
				}
			default:
				// do nothing
			}
		}
	}
}

// Count of response items
func (r *Response) Count() int {
	return len(r.items)
}

// Validate response
func (r *Response) Validate() (err error) {
	if r.err != nil {
		return r.err
	}
	if r.Count() < 1 {
		return adtype.ErrResponseEmpty
	}
	for _, it := range r.items {
		if err = it.Validate(); err != nil {
			break
		}
	}
	return
}

// Error of the response
func (r *Response) Error() error {
	return r.err
}

// Context value
func (r *Response) Context(ctx ...context.Context) context.Context {
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
func (r *Response) Get(key string) any {
	if r.context == nil {
		return nil
	}
	return r.context.Value(key)
}

// Release response and all linked objects
func (r *Response) Release() {
	if r != nil {
		r.reset()
	}
}

func (r *Response) reset() {
	r.items = r.items[:0]
	r.err = nil
	r.context = nil
	r.source = nil
	r.request = nil
}

var (
	_ adtype.Response = &Response{}
)
