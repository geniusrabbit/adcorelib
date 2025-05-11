package adtype

import (
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

// TargetEmpty model represents empty target
// which is used for the case when target is not specified
type TargetEmpty struct {
	Acc Account
}

// ID of object
func (l *TargetEmpty) ID() uint64 { return 0 }

// Codename of the target (equal to tagid)
func (l *TargetEmpty) Codename() string { return "" }

// ObjectKey of the target
func (l *TargetEmpty) ObjectKey() string { return l.Codename() }

// PricingModel of the target
func (l *TargetEmpty) PricingModel() types.PricingModel {
	return types.PricingModelUndefined
}

// AlternativeAdCode returns URL or any code (HTML, XML, etc)
func (l *TargetEmpty) AlternativeAdCode(key string) string { return "" }

// PurchasePrice gives the price of view from external resource
func (l *TargetEmpty) PurchasePrice(action Action) billing.Money { return 0 }

// Account object
func (l *TargetEmpty) Account() Account { return l.Acc }

// AccountID of current target
func (l *TargetEmpty) AccountID() uint64 { return l.Acc.ID() }

// SetAccount for target
func (l *TargetEmpty) SetAccount(acc Account) { l.Acc = acc }

// CommissionShareFactor which system get from publisher
func (l *TargetEmpty) CommissionShareFactor() float64 {
	return l.Acc.CommissionShareFactor()
}

// RevenueShareReduceFactor correction factor to reduce target proce of the access point to avoid descrepancy
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (l *TargetEmpty) RevenueShareReduceFactor() float64 { return 0 }

// IsAllowedSource for targeting
func (l *TargetEmpty) IsAllowedSource(id uint64, types []int) bool {
	return true
}
