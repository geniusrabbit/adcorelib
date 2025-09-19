package price

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

type PriceScope struct {
	// TestMode represents the flag for the test budget usage for the view price.
	TestMode bool `json:"test_mode,omitempty"`

	// MaxBidImpPrice represents the maximum price for the bid on the auction.
	MaxBidImpPrice billing.Money `json:"max_bid_imp_price,omitempty"`

	// BidImpPrice represents the price for the bid on the auction. But charged will be by ImpPrice
	BidImpPrice billing.Money `json:"bid_imp_price,omitempty"`

	// ImpPrice represents the price for the impression action.
	ImpPrice billing.Money `json:"imp_price,omitempty"`

	// ViewPrice represents the price for the view action.
	ViewPrice billing.Money `json:"view_price,omitempty"`

	// ClickPrice represents the price for the click action.
	ClickPrice billing.Money `json:"click_price,omitempty"`

	// LeadPrice represents the price for the lead action.
	LeadPrice billing.Money `json:"lead_price,omitempty"`

	// ECPM represents the price for the 1000 views.
	ECPM billing.Money `json:"ecpm,omitempty"`
}

// PricePerAction returns the price for the action type.
func (ps *PriceScope) PricePerAction(actionType adtype.Action) billing.Money {
	switch actionType {
	case adtype.ActionImpression:
		return ps.ImpPrice
	case adtype.ActionView:
		return ps.ViewPrice
	case adtype.ActionClick:
		return ps.ClickPrice
	case adtype.ActionLead:
		return ps.LeadPrice
	default:
		return 0
	}
}

// SetBidImpressionPrice sets the price for the bid on the auction.
func (ps *PriceScope) SetBidImpressionPrice(price billing.Money, maxIfZero bool) bool {
	if price <= 0 && maxIfZero {
		ps.BidImpPrice = ps.MaxBidImpPrice
		return true
	}
	if price > ps.BidImpPrice {
		return false
	}
	ps.BidImpPrice = max(price, 0)
	return true
}

// SetImpressionPrice sets the price for the impression action.
func (ps *PriceScope) SetImpressionPrice(price billing.Money, maxIfZero bool) bool {
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

// PrepareBidImpressionPrice prepares the bid impression price for the auction.
func (ps *PriceScope) PrepareBidImpressionPrice(price billing.Money) billing.Money {
	if price > ps.MaxBidImpPrice {
		price = ps.MaxBidImpPrice
	}
	return price
}
