package adtype

import (
	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/billing"
)

type PriceScope struct {
	MaxBidPrice billing.Money `json:"max_bid_price,omitempty"`
	ViewPrice   billing.Money `json:"view_price,omitempty"`
	ClickPrice  billing.Money `json:"click_price,omitempty"`
	LeadPrice   billing.Money `json:"lead_price,omitempty"`
	ECPM        billing.Money `json:"ecpm,omitempty"`
}

func (ps *PriceScope) PricePerAction(actionType admodels.Action) billing.Money {
	switch actionType {
	case admodels.ActionView:
		return ps.ViewPrice
	case admodels.ActionClick:
		return ps.ClickPrice
	case admodels.ActionLead:
		return ps.LeadPrice
	default:
		return ps.MaxBidPrice
	}
}
