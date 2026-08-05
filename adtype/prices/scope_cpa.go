package prices

import "github.com/geniusrabbit/adcorelib/billing"

// CPAScope prices the lead action. Unlike the CPM based scopes both values are
// stored as the price of one single lead.
//
// The scope can be embedded into any structure which holds the price of the
// lead action:
//
//	type Campaign struct {
//		prices.CPAScope
//	}
type CPAScope struct {
	// MaxLeadPrice is the maximum price of one lead which the advertiser agreed to
	// pay. The zero value means that there is no upper limit.
	MaxLeadPrice billing.Money `json:"max_lead_price,omitempty"`

	// LeadPrice is the current price of one lead.
	// Always less than or equal to the MaxLeadPrice if the last one is defined.
	LeadPrice billing.Money `json:"lead_price,omitempty"`
}

// HasCPA returns true if the lead pricing is defined for the scope.
//
//go:inline
func (s *CPAScope) HasCPA() bool { return s.LeadPrice > 0 || s.MaxLeadPrice > 0 }

// SetLeadPrice sets the current price of one lead.
// The value is clamped by the MaxLeadPrice if the last one is defined.
func (s *CPAScope) SetLeadPrice(price billing.Money) error {
	return setBid(&s.LeadPrice, price, s.MaxLeadPrice)
}

// PrepareLeadPrice returns the price of one lead limited by the MaxLeadPrice.
//
//go:inline
func (s *CPAScope) PrepareLeadPrice(price billing.Money) billing.Money {
	return prepareBid(price, s.MaxLeadPrice)
}
