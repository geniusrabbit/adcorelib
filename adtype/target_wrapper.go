package adtype

import (
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

// VirtualTarget it's a wrapper which don't have ID
type VirtualTarget struct {
	trg Target
}

// NewVirtualTarget wrapper of exists target
func NewVirtualTarget(trg Target) VirtualTarget {
	return VirtualTarget{trg: trg}
}

// Target object accessor
func (tw VirtualTarget) Target() Target { return tw.trg }

// ID of object (Zone OR SmartLink only)
func (tw VirtualTarget) ID() uint64 { return 0 }

// Codename of the target (equal to tagid)
func (tw VirtualTarget) Codename() string { return "" }

// PricingModel of the target
// Undefined as any priceing model type
func (tw VirtualTarget) PricingModel() types.PricingModel { return tw.trg.PricingModel() }

// PurchasePrice gives the price of view from external resource
func (tw VirtualTarget) PurchasePrice(a Action) billing.Money { return tw.trg.PurchasePrice(a) }

// CommissionShareFactor of current target
func (tw VirtualTarget) CommissionShareFactor() float64 { return tw.trg.CommissionShareFactor() }

// Account object
func (tw VirtualTarget) Account() Account { return tw.trg.Account() }
