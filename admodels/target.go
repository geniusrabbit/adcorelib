//
// @project GeniusRabbit corelib 2017 - 2019, 2021
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019, 2021
//

package admodels

import (
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/models"
)

// Target type object
type Target interface {
	// ID of object (Zone OR SmartLink only)
	ID() uint64

	// Codename of the target (equal to tagid)
	Codename() string

	// PricingModel of the target
	// Undefined as any priceing model type
	PricingModel() types.PricingModel

	// AlternativeAdCode returns URL or any code (HTML, XML, etc)
	// for alternative ad code for the target.
	// Key represents adformat code (banner, video, direct, etc)
	AlternativeAdCode(key string) string

	// PurchasePrice gives the price of view from external resource
	PurchasePrice(action Action) billing.Money

	// RevenueShareFactor of current target from 0 to 1
	RevenueShareFactor() float64

	// ComissionShareFactor of current target from 0 to 1
	ComissionShareFactor() float64

	// RevenueShareReduceFactor correction factor to reduce target proce of the access point to avoid descrepancy
	// Returns percent from 0 to 1 for reducing of the value
	// If there is 10% of price correction, it means that 10% of the final price must be ignored
	RevenueShareReduceFactor() float64

	// Account object
	Account() *Account

	// AccountID of current target
	AccountID() uint64
}

// TargetFromModel convert datavase model specified model
// which implements Target interface
func TargetFromModel(zone *models.Zone, acc *Account) Target {
	if zone.Type.IsSmartlink() {
		return SmartlinkFromModel(zone, acc)
	}
	return ZoneFromModel(zone, acc)
}
