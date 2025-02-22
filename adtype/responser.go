//
// @project GeniusRabbit corelib 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package adtype

import (
	"context"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// Responser type
type Responser interface {
	// AuctionID response
	AuctionID() string

	// AuctionType of request
	AuctionType() types.AuctionType

	// Source of response
	Source() Source

	// Request information
	Request() *BidRequest

	// Ads list
	Ads() []ResponserItemCommon

	// Item by impression code
	Item(impid string) ResponserItemCommon

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
