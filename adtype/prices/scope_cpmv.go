package prices

import "github.com/geniusrabbit/adcorelib/billing"

// CPMVScope prices the view action. Both values are stored in CPM units,
// which is the price of 1000 views.
//
// The scope can be embedded into any structure which holds the price of the
// view action:
//
//	type Campaign struct {
//		prices.CPMVScope
//	}
type CPMVScope struct {
	// MaxBidCPMV is the maximum price of 1000 views which the advertiser agreed
	// to pay. The zero value means that there is no upper limit.
	MaxBidCPMV billing.Money `json:"max_bid_cpmv,omitempty"`

	// BidCPMV is the current price of 1000 views which is used in the auction.
	// Always less than or equal to the MaxBidCPMV if the last one is defined.
	BidCPMV billing.Money `json:"bid_cpmv,omitempty"`
}

// HasCPMV returns true if the view pricing is defined for the scope.
//
//go:inline
func (s *CPMVScope) HasCPMV() bool { return s.BidCPMV > 0 || s.MaxBidCPMV > 0 }

// SetBidCPMV sets the current price of 1000 views.
// The value is clamped by the MaxBidCPMV if the last one is defined.
func (s *CPMVScope) SetBidCPMV(bid billing.Money) error {
	return setBid(&s.BidCPMV, bid, s.MaxBidCPMV)
}

// PrepareBidCPMV returns the price of 1000 views limited by the MaxBidCPMV.
//
//go:inline
func (s *CPMVScope) PrepareBidCPMV(bid billing.Money) billing.Money {
	return prepareBid(bid, s.MaxBidCPMV)
}
