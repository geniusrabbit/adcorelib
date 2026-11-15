package adtype

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

// zoneAccountTarget stores the publisher id separately from a missing account object.
type zoneAccountTarget struct {
	id  uint64
	acc Account
}

func (z zoneAccountTarget) ID() uint64                         { return 1 }
func (z zoneAccountTarget) Codename() string                   { return "z" }
func (z zoneAccountTarget) ObjectKey() string                  { return z.Codename() }
func (z zoneAccountTarget) PricingModel() types.PricingModel   { return types.PricingModelUndefined }
func (z zoneAccountTarget) AlternativeAdCode(string) string    { return "" }
func (z zoneAccountTarget) PurchasePrice(Action) billing.Money { return 0 }
func (z zoneAccountTarget) CommissionShareFactor() float64     { return 0 }
func (z zoneAccountTarget) RevenueShareReduceFactor() float64  { return 0 }
func (z zoneAccountTarget) Account() Account                   { return z.acc }
func (z zoneAccountTarget) AccountID() uint64                  { return z.id }

func TestImpressionAccountIDNilAccount(t *testing.T) {
	var boxed *admodels.Account
	imp := &Impression{Target: zoneAccountTarget{id: 42, acc: boxed}}
	if got := imp.AccountID(); got != 42 {
		t.Fatalf("AccountID = %d, want 42", got)
	}
}

type nilAccountTarget struct{}

func (nilAccountTarget) ID() uint64                         { return 0 }
func (nilAccountTarget) Codename() string                   { return "" }
func (nilAccountTarget) ObjectKey() string                  { return "" }
func (nilAccountTarget) PricingModel() types.PricingModel   { return types.PricingModelUndefined }
func (nilAccountTarget) AlternativeAdCode(string) string    { return "" }
func (nilAccountTarget) PurchasePrice(Action) billing.Money { return 0 }
func (nilAccountTarget) CommissionShareFactor() float64     { return 0 }
func (nilAccountTarget) RevenueShareReduceFactor() float64  { return 0 }
func (nilAccountTarget) Account() Account {
	var acc *admodels.Account
	return acc
}

func TestImpressionAccountIDBoxedNil(t *testing.T) {
	imp := &Impression{Target: nilAccountTarget{}}
	if got := imp.AccountID(); got != 0 {
		t.Fatalf("AccountID = %d, want 0", got)
	}
}
