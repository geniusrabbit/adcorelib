package price

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

type PriceScopeView struct {
	// MaxBidPrice represents the maximum price for the bid on the auction.
	MaxBidPrice billing.Money `json:"max_bid_price,omitempty"`

	// BidPrice represents the price for the bid on the auction. But charged will be by ViewPrice
	BidPrice billing.Money `json:"bid_price,omitempty"`

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

// SetBidPrice sets the price for the bid on the auction.
func (ps *PriceScopeView) SetBidPrice(price billing.Money, maxIfZero bool) bool {
	if price <= 0 && maxIfZero {
		ps.BidPrice = ps.MaxBidPrice
		return true
	}
	if price > ps.MaxBidPrice {
		return false
	}
	ps.BidPrice = max(price, 0)
	return true
}

// SetViewPrice sets the price for the view action.
func (ps *PriceScopeView) SetViewPrice(price billing.Money, maxIfZero bool) bool {
	if price <= 0 && maxIfZero {
		ps.ViewPrice = ps.MaxBidPrice
		return true
	}
	if price > ps.MaxBidPrice {
		ps.ViewPrice = ps.MaxBidPrice
		return true
	}
	ps.ViewPrice = max(price, 0)
	return true
}
