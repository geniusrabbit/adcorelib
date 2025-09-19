package price

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// PriceScopeImpression holds pricing information for impression-based auctions.
type PriceScopeImpression struct {
	// MaxBidImpPrice is the maximum allowed bid price in the auction.
	MaxBidImpPrice billing.Money `json:"max_bid_price,omitempty"`

	// BidImpPrice is the bid price set for the auction. The actual charge is determined by ImpPrice.
	BidImpPrice billing.Money `json:"bid_price,omitempty"`

	// ImpPrice is the price charged for each impression action.
	ImpPrice billing.Money `json:"imp_price,omitempty"`

	// ECPM is the effective cost per thousand impressions.
	ECPM billing.Money `json:"ecpm,omitempty"`
}

// PricePerAction returns the price for the specified action type.
// Only returns ImpPrice for ActionImpression; otherwise returns 0.
func (ps *PriceScopeImpression) PricePerAction(actionType adtype.Action) billing.Money {
	if actionType == adtype.ActionImpression {
		return ps.ImpPrice
	}
	return 0
}

// SetBidImpressionPrice sets the bid price for the auction.
// If price is zero or negative and maxIfZero is true, sets BidImpPrice to MaxBidImpPrice.
// Returns false if price exceeds MaxBidImpPrice; otherwise sets BidImpPrice and returns true.
func (ps *PriceScopeImpression) SetBidImpressionPrice(price billing.Money, maxIfZero bool) bool {
	if price <= 0 && maxIfZero {
		ps.BidImpPrice = ps.MaxBidImpPrice
		return true
	}
	if price > ps.MaxBidImpPrice {
		return false
	}
	ps.BidImpPrice = max(price, 0)
	return true
}

// SetImpressionPrice sets the price for the impression action.
// If price is zero or negative and maxIfZero is true, sets ImpPrice to MaxBidImpPrice.
// If price exceeds MaxBidImpPrice, sets ImpPrice to MaxBidImpPrice.
// Otherwise, sets ImpPrice to the given price (minimum zero).
func (ps *PriceScopeImpression) SetImpressionPrice(price billing.Money, maxIfZero bool) bool {
	if price <= 0 && maxIfZero {
		ps.ImpPrice = ps.MaxBidImpPrice
		return true
	}
	if price > ps.MaxBidImpPrice {
		ps.ImpPrice = ps.MaxBidImpPrice
		return true
	}
	ps.ImpPrice = max(price, 0)
	return true
}

// PrepareBidImpressionPrice returns a valid bid price for the auction.
// If price exceeds MaxBidImpPrice and MaxBidImpPrice is positive, returns MaxBidImpPrice.
// Otherwise, returns the given price.
func (ps *PriceScopeImpression) PrepareBidImpressionPrice(price billing.Money) billing.Money {
	if price > ps.MaxBidImpPrice && ps.MaxBidImpPrice > 0 {
		price = ps.MaxBidImpPrice
	}
	return price
}
