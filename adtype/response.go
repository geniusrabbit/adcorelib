//
// @project GeniusRabbit corelib 2016 – 2019, 2025
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019, 2025
//

package adtype

import (
	"context"
	"iter"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// Response type definition
type Response interface {
	// AuctionID response
	AuctionID() string

	// AuctionType of request
	AuctionType() types.AuctionType

	// Source of response
	Source() Source

	// Request information
	Request() BidRequester

	// Ads list
	Ads() []ResponseItemCommon

	// IterAds returns an iterator over the response items.
	IterAds() iter.Seq[ResponseItem]

	// Item by impression code
	Item(impid string) ResponseItemCommon

	// Count of response items
	Count() int

	// Validate response
	Validate() error

	// Error of the response
	Error() error

	// Context value
	Context(ctx ...context.Context) context.Context

	// Get context item by key
	Get(key string) any

	// Release response and all linked objects
	Release()
}
