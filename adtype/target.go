package adtype

import (
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

// Target type object represents the placement of the ad
// in the ad system. It can be a zone, smartlink, or any other
type Target interface {
	// ID of object (Zone OR SmartLink only)
	ID() uint64

	// Codename of the target (equal to tagid)
	Codename() string

	// ObjectKey of the target
	ObjectKey() string

	// PricingModel of the target
	// Undefined as any priceing model type
	PricingModel() types.PricingModel

	// AlternativeAdCode returns URL or any code (HTML, XML, etc)
	// for alternative ad code for the target.
	// Key represents adformat code (banner, video, direct, etc)
	AlternativeAdCode(key string) string

	// PurchasePrice gives the price of view from external resource
	PurchasePrice(action Action) billing.Money

	// CommissionShareFactor of current target from 0 to 1
	CommissionShareFactor() float64

	// RevenueShareReduceFactor correction factor to reduce target proce of the access point to avoid descrepancy
	// Returns percent from 0 to 1 for reducing of the value
	// If there is 10% of price correction, it means that 10% of the final price must be ignored
	RevenueShareReduceFactor() float64

	// Account object
	Account() Account
}
