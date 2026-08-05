package prices

import "github.com/geniusrabbit/adcorelib/billing"

// CPMScope prices the impression action. Both values are stored in CPM units,
// which is the price of 1000 impressions.
//
// The scope can be embedded into any structure which holds the price of the
// impression action:
//
//	type Campaign struct {
//		prices.CPMScope
//	}
type CPMScope struct {
	// MaxBidCPM is the maximum price of 1000 impressions which the advertiser
	// agreed to pay. The zero value means that there is no upper limit.
	MaxBidCPM billing.Money `json:"max_bid_cpm,omitempty"`

	// BidCPM is the current price of 1000 impressions which is used in the auction.
	// Always less than or equal to the MaxBidCPM if the last one is defined.
	BidCPM billing.Money `json:"bid_cpm,omitempty"`
}

// HasCPM returns true if the impression pricing is defined for the scope.
//
//go:inline
func (s *CPMScope) HasCPM() bool { return s.BidCPM > 0 || s.MaxBidCPM > 0 }

// SetBidCPM sets the current price of 1000 impressions.
// The value is clamped by the MaxBidCPM if the last one is defined.
func (s *CPMScope) SetBidCPM(bid billing.Money) error {
	return setBid(&s.BidCPM, bid, s.MaxBidCPM)
}

// PrepareBidCPM returns the price of 1000 impressions limited by the MaxBidCPM.
//
//go:inline
func (s *CPMScope) PrepareBidCPM(bid billing.Money) billing.Money {
	return prepareBid(bid, s.MaxBidCPM)
}
