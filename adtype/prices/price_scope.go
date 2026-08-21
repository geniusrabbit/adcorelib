package prices

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// PriceScope combines all the pricing models of the advertisement. Any number of
// them can be defined and used at the same time – the scope does not restrict the
// combination and does not decide which of them has to be charged.
//
// The scope is designed to be embedded into any campaign or advertisement type:
//
//	type Campaign struct {
//		prices.PriceScope
//	}
//
// All the dispatch methods operate with the price of one single action, so the
// CPM units of the impression and the view scopes stay an implementation detail.
type PriceScope struct {
	CPMScope  // Impression pricing
	CPMVScope // View pricing
	CPCScope  // Click pricing
	CPAScope  // Lead pricing

	// ECPM of the advertisement which is used for the internal auction ranking.
	// For the CPM model it is equal to the BidCPM, for all the other models it has
	// to be predicted from the statistics and assigned by the optimizer module.
	ECPM billing.Money `json:"ecpm,omitempty"`
}

// HasAction returns true if the pricing of the action is defined for the scope.
func (ps *PriceScope) HasAction(action adtype.Action) bool {
	switch action {
	case adtype.ActionImpression, adtype.ActionDirect:
		return ps.HasCPM()
	case adtype.ActionView:
		return ps.HasCPMV()
	case adtype.ActionClick:
		return ps.HasCPC()
	case adtype.ActionLead:
		return ps.HasCPA()
	}
	return false
}

// PricePerAction returns the current price of one action.
func (ps *PriceScope) PricePerAction(action adtype.Action) billing.Money {
	switch action {
	case adtype.ActionImpression, adtype.ActionDirect:
		return PriceFromCPM(ps.BidCPM)
	case adtype.ActionView:
		return PriceFromCPM(ps.BidCPMV)
	case adtype.ActionClick:
		return ps.BidCPC
	case adtype.ActionLead:
		return ps.LeadPrice
	}
	return 0
}

// MaxPricePerAction returns the maximal possible price of one action.
func (ps *PriceScope) MaxPricePerAction(action adtype.Action) billing.Money {
	switch action {
	case adtype.ActionImpression, adtype.ActionDirect:
		return PriceFromCPM(ps.MaxBidCPM)
	case adtype.ActionView:
		return PriceFromCPM(ps.MaxBidCPMV)
	case adtype.ActionClick:
		return ps.MaxBidCPC
	case adtype.ActionLead:
		return ps.MaxLeadPrice
	}
	return 0
}

// SetBidPerAction sets the current price of one action.
// The value is clamped by the maximal price of the action if the last one is defined.
func (ps *PriceScope) SetBidPerAction(action adtype.Action, price billing.Money) error {
	switch action {
	case adtype.ActionImpression, adtype.ActionDirect:
		return ps.SetBidCPM(CPMFromPrice(price))
	case adtype.ActionView:
		return ps.SetBidCPMV(CPMFromPrice(price))
	case adtype.ActionClick:
		return ps.SetBidCPC(price)
	case adtype.ActionLead:
		return ps.SetLeadPrice(price)
	}
	return ErrUnsupportedAction
}

// SetBidPrice sets the current bid price for the given action. If withCommission
// is true, the price is treated as the net amount that must remain after the
// discrepancy corrections and the commission share are deducted downstream (see
// [PublisherPrice]) and is therefore grossed up via [BidUpPrice] before being
// stored. If withCommission is false, the price is stored as given. The stored
// value is then clamped by the maximal price of the action if the last one is
// defined.
func (ps *PriceScope) SetBidPrice(action adtype.Action, price billing.Money, factors Factors, withCommission bool) error {
	if withCommission {
		price = BidUpPrice(price, factors)
	}
	return ps.SetBidPerAction(action, price)
}

// PrepareBidPerAction returns the price of one action limited by the maximal price
// of the action.
func (ps *PriceScope) PrepareBidPerAction(action adtype.Action, price billing.Money) billing.Money {
	return prepareBid(price, ps.MaxPricePerAction(action))
}

// EffectiveCPM returns the CPM value which is used for the internal auction ranking.
// Fallbacks to the impression bid if the ECPM is not predicted.
//
//go:inline
func (ps *PriceScope) EffectiveCPM() billing.Money {
	if ps.ECPM > 0 {
		return ps.ECPM
	}
	return ps.BidCPM
}

// PotentialPricePerAction returns the maximal price which the advertiser could have
// paid for the action. See [PotentialPrice].
//
//go:inline
func (ps *PriceScope) PotentialPricePerAction(action adtype.Action) billing.Money {
	return PotentialPrice(ps, action)
}

// AdvertiserPricePerAction returns the price which will be charged from the
// advertiser for the action. See [AdvertiserPrice].
//
//go:inline
func (ps *PriceScope) AdvertiserPricePerAction(action adtype.Action, factors Factors) billing.Money {
	return AdvertiserPrice(ps, factors, action)
}

// PublisherPricePerAction returns the price which the system has to pay to the
// publisher for the action. See [PublisherPrice].
//
//go:inline
func (ps *PriceScope) PublisherPricePerAction(action adtype.Action, factors Factors) billing.Money {
	return PublisherPrice(ps, factors, action)
}

// NetworkProfitPerAction returns the profit of the network for the action. See
// [NetworkProfit].
//
//go:inline
func (ps *PriceScope) NetworkProfitPerAction(action adtype.Action, factors Factors) billing.Money {
	return NetworkProfit(ps, factors, action)
}

var _ PriceProvider = (*PriceScope)(nil)
