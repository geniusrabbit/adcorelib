package price

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

type PriceScopeView struct {
	// MaxBidViewPrice represents the maximum price for the bid on the auction.
	MaxBidViewPrice billing.Money `json:"max_bid_price,omitempty"`

	// BidViewPrice represents the price for the bid on the auction. But charged will be by ViewPrice
	BidViewPrice billing.Money `json:"bid_price,omitempty"`

	// ViewPrice represents the price for the view action.
	ViewPrice billing.Money `json:"view_price,omitempty"`

	// ECPM represents the price for the 1000 views.
	ECPM billing.Money `json:"ecpm,omitempty"`
}

// PricePerAction returns the price for the action type.
func (ps *PriceScopeView) PricePerAction(actionType adtype.Action) billing.Money {
	if actionType == adtype.ActionView {
		return ps.ViewPrice
	}
	return 0
}

// SetBidViewPrice sets the price for the bid on the auction.
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

// PrepareBidViewPrice prepares the bid view price for the auction.
func (ps *PriceScopeView) PrepareBidViewPrice(price billing.Money) billing.Money {
	if price > ps.MaxBidViewPrice {
		price = ps.MaxBidViewPrice
	}
	return price
}
