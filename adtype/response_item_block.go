//
// @project GeniusRabbit rotator 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package adtype

import (
	"context"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

//
// When we need to combinate several AD items to one
//

// ResponseItemBlock group of simple items
type ResponseItemBlock struct {
	Items   []ResponserItem
	context context.Context
}

// ID of current response item (unique code of current response)
func (i *ResponseItemBlock) ID() string {
	return ""
}

// Impression place object
func (i *ResponseItemBlock) Impression() *Impression {
	return i.Items[0].Impression()
}

// ImpressionID code
func (i *ResponseItemBlock) ImpressionID() string {
	return i.Items[0].ImpressionID()
}

// ExtImpressionID it's unique code of the auction bid impression
func (i *ResponseItemBlock) ExtImpressionID() string {
	return ""
}

// ExtTargetID of the external network
func (i *ResponseItemBlock) ExtTargetID() string {
	return ""
}

// PriorityFormatType from current Ad
func (i *ResponseItemBlock) PriorityFormatType() types.FormatType {
	return types.FormatInvalidType
}

// Price of whole response
func (i *ResponseItemBlock) Price(action admodels.Action) (price billing.Money) {
	for _, it := range i.Items {
		price += it.Price(action)
	}
	return price
}

// AuctionCPMBid value price without any comission
func (i *ResponseItemBlock) AuctionCPMBid() (bid billing.Money) {
	for _, it := range i.Items {
		bid += it.AuctionCPMBid()
	}
	return bid
}

// Ads list
func (i *ResponseItemBlock) Ads() []ResponserItem {
	return i.Items
}

// Count of response items
func (i *ResponseItemBlock) Count() int {
	return len(i.Items)
}

// Validate response
func (i *ResponseItemBlock) Validate() (err error) {
	if len(i.Items) < 1 {
		return ErrResponseEmpty
	}
	for _, it := range i.Items {
		if err = it.Validate(); nil != err {
			return err
		}
	}
	return err
}

// Context value
func (i *ResponseItemBlock) Context(ctx ...context.Context) context.Context {
	if len(ctx) > 0 && ctx[0] != nil {
		i.context = ctx[0]
	}
	return i.context
}

// Get ext field
func (i *ResponseItemBlock) Get(key string) any {
	if i.context != nil {
		return i.context.Value(key)
	}
	return nil
}

var (
	_ ResponserMultipleItem = &ResponseItemBlock{}
)
