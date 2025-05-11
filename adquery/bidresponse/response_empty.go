package bidresponse

import (
	"context"

	"github.com/bsm/openrtb"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// ResponseEmpty object represents empty response and response item
type ResponseEmpty struct {
	ItemID  string
	Req     *adtype.BidRequest
	Src     adtype.Source
	Imp     *adtype.Impression
	Bid     *openrtb.Bid
	Err     error
	context context.Context
}

// NewEmptyResponse by bid request
func NewEmptyResponse(request *adtype.BidRequest, src adtype.Source, err error) *ResponseEmpty {
	return &ResponseEmpty{Req: request, Src: src, Err: err, context: request.Ctx}
}

// ID of current response item (unique code of current response)
func (r ResponseEmpty) ID() string {
	return r.ItemID
}

// AuctionID returns ID of the current auction
func (r ResponseEmpty) AuctionID() string {
	if r.Req == nil {
		return ""
	}
	return r.Req.ID
}

// AuctionType returns type of the auction (first price, second price, etc)
func (r ResponseEmpty) AuctionType() types.AuctionType {
	if r.Req == nil {
		return types.UndefinedAuctionType
	}
	return r.Req.AuctionType
}

// Source returns source of the response
func (r ResponseEmpty) Source() adtype.Source {
	return r.Src
}

// PriorityFormatType from current Ad
func (r ResponseEmpty) PriorityFormatType() types.FormatType {
	return types.FormatUndefinedType
}

// Request returns the original request object
func (r ResponseEmpty) Request() *adtype.BidRequest {
	return r.Req
}

// Impression place object
func (r ResponseEmpty) Impression() *adtype.Impression {
	return r.Imp
}

// ImpressionID unique code string
func (r ResponseEmpty) ImpressionID() string {
	if r.Imp == nil {
		return ""
	}
	return r.Imp.ID
}

// ExtImpressionID it's unique code of the auction bid impression
func (r ResponseEmpty) ExtImpressionID() string {
	if r.Bid == nil {
		return ""
	}
	return r.Bid.ImpID
}

// ExtTargetID of the external network
func (r ResponseEmpty) ExtTargetID() string {
	if r.Imp == nil {
		return ""
	}
	return r.Imp.ExternalTargetID
}

// Ads list
func (r ResponseEmpty) Ads() []adtype.ResponserItemCommon {
	return nil
}

// Item by impression code
func (r ResponseEmpty) Item(impid string) adtype.ResponserItemCommon {
	return nil
}

// InternalAuctionCPMBid value provides maximal possible price without any comission
// According to this value the system can choice the best item for the auction
func (r ResponseEmpty) InternalAuctionCPMBid() billing.Money {
	return 0
}

// Count of response items
func (r ResponseEmpty) Count() int {
	return 0
}

// Validate response
func (r ResponseEmpty) Validate() error {
	if r.Err != nil {
		return r.Err
	}
	return adtype.ErrResponseEmpty
}

// Error of the response
func (r ResponseEmpty) Error() error {
	return r.Err
}

// Context value
func (r *ResponseEmpty) Context(ctx ...context.Context) context.Context {
	if r != nil && len(ctx) > 0 && ctx[0] != nil {
		r.context = ctx[0]
	}
	if r.context == nil && r.Req != nil {
		return r.Req.Ctx
	}
	return r.context
}

// Get ext field
func (r *ResponseEmpty) Get(key string) (res any) {
	if r != nil && r.context != nil {
		res = r.context.Value(key)
	}
	return res
}

// Release response and all linked objects
func (r *ResponseEmpty) Release() {}

var (
	_ adtype.ResponserItemCommon = (*ResponseEmpty)(nil)
	_ adtype.Responser           = (*ResponseEmpty)(nil)
)
