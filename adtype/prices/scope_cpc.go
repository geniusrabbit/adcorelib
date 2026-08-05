package prices

import "github.com/geniusrabbit/adcorelib/billing"

// CPCScope prices the click action. Unlike the CPM based scopes both values are
// stored as the price of one single click.
//
// The scope can be embedded into any structure which holds the price of the
// click action:
//
//	type Campaign struct {
//		prices.CPCScope
//	}
type CPCScope struct {
	// MaxBidCPC is the maximum price of one click which the advertiser agreed to
	// pay. The zero value means that there is no upper limit.
	MaxBidCPC billing.Money `json:"max_bid_cpc,omitempty"`

	// BidCPC is the current price of one click which is used in the auction.
	// Always less than or equal to the MaxBidCPC if the last one is defined.
	BidCPC billing.Money `json:"bid_cpc,omitempty"`
}

// HasCPC returns true if the click pricing is defined for the scope.
//
//go:inline
func (s *CPCScope) HasCPC() bool { return s.BidCPC > 0 || s.MaxBidCPC > 0 }

// SetBidCPC sets the current price of one click.
// The value is clamped by the MaxBidCPC if the last one is defined.
func (s *CPCScope) SetBidCPC(bid billing.Money) error {
	return setBid(&s.BidCPC, bid, s.MaxBidCPC)
}

// PrepareBidCPC returns the price of one click limited by the MaxBidCPC.
//
//go:inline
func (s *CPCScope) PrepareBidCPC(bid billing.Money) billing.Money {
	return prepareBid(bid, s.MaxBidCPC)
}
