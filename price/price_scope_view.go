package price

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// PriceScopeView holds pricing information for view-based auctions.
type PriceScopeView struct {
	// MaxBidViewPrice is the maximum allowed bid price in the auction.
	MaxBidViewPrice billing.Money `json:"max_bid_price,omitempty"`

	// BidViewPrice is the bid price set for the auction. The actual charge is determined by ViewPrice.
	BidViewPrice billing.Money `json:"bid_price,omitempty"`

	// ViewPrice is the price charged for each view action.
	ViewPrice billing.Money `json:"view_price,omitempty"`

	// ECPM is the effective cost per thousand views.
	ECPM billing.Money `json:"ecpm,omitempty"`
}

// PricePerAction returns the price for the specified action type.
// Only returns ViewPrice for ActionView; otherwise returns 0.
func (ps *PriceScopeView) PricePerAction(actionType adtype.Action) billing.Money {
	if actionType == adtype.ActionView {
		return ps.ViewPrice
	}
	return 0
}

// SetBidViewPrice sets the bid price for the auction.
// If price is zero or negative and maxIfZero is true, sets BidViewPrice to MaxBidViewPrice.
// Returns false if price exceeds MaxBidViewPrice; otherwise sets BidViewPrice and returns true.
func (ps *PriceScopeView) SetBidViewPrice(price billing.Money, maxIfZero bool) bool {
	if price <= 0 && maxIfZero {
		ps.BidViewPrice = ps.MaxBidViewPrice
		return true
	}
	if price > ps.MaxBidViewPrice {
		return false
	}
	ps.BidViewPrice = max(price, 0)
	return true
}

// SetViewPrice sets the price for the view action.
// If price is zero or negative and maxIfZero is true, sets ViewPrice to MaxBidViewPrice.
// If price exceeds MaxBidViewPrice, sets ViewPrice to MaxBidViewPrice.
// Otherwise, sets ViewPrice to the given price (minimum zero).
func (ps *PriceScopeView) SetViewPrice(price billing.Money, maxIfZero bool) bool {
	if price <= 0 && maxIfZero {
		ps.ViewPrice = ps.MaxBidViewPrice
		return true
	}
	if price > ps.MaxBidViewPrice {
		ps.ViewPrice = ps.MaxBidViewPrice
		return true
	}
	ps.ViewPrice = max(price, 0)
	return true
}

// PrepareBidViewPrice returns a valid bid price for the auction.
// If price exceeds MaxBidViewPrice and MaxBidViewPrice is positive, returns MaxBidViewPrice.
// Otherwise, returns the given price.
func (ps *PriceScopeView) PrepareBidViewPrice(price billing.Money) billing.Money {
	if price > ps.MaxBidViewPrice && ps.MaxBidViewPrice > 0 {
		price = ps.MaxBidViewPrice
	}
	return price
}
