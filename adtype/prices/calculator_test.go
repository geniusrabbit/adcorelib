package prices

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// fixedFactors extends the static factors by the fixed purchase price of the target.
type fixedFactors struct {
	StaticFactors
	fixed billing.Money
}

func (f fixedFactors) FixedPurchasePrice(action adtype.Action) billing.Money {
	if action.IsImpression() {
		return f.fixed
	}
	return 0
}

var _ FixedPurchasePricer = fixedFactors{}

func TestCalculatePrices(t *testing.T) {
	scope := PriceScope{
		CPMScope:  CPMScope{MaxBidCPM: billing.MoneyFloat(2.), BidCPM: billing.MoneyFloat(1.)},
		CPMVScope: CPMVScope{MaxBidCPMV: billing.MoneyFloat(6.), BidCPMV: billing.MoneyFloat(3.)},
		CPCScope:  CPCScope{MaxBidCPC: billing.MoneyFloat(0.5), BidCPC: billing.MoneyFloat(0.2)},
		CPAScope:  CPAScope{MaxLeadPrice: billing.MoneyFloat(20.), LeadPrice: billing.MoneyFloat(10.)},
	}
	factors := StaticFactors{Commission: 0.2, Source: 0.1, Target: 0.05}

	tests := []struct {
		action          adtype.Action
		expectPotential billing.Money
		expectAdvertier billing.Money
		expectPublisher billing.Money
	}{
		{
			action:          adtype.ActionImpression,
			expectPotential: billing.MoneyFloat(0.002),
			expectAdvertier: billing.MoneyFloat(0.001),
			expectPublisher: billing.MoneyFloat(0.001 * 0.9 * 0.95 * 0.8),
		},
		{
			action:          adtype.ActionView,
			expectPotential: billing.MoneyFloat(0.006),
			expectAdvertier: billing.MoneyFloat(0.003),
			expectPublisher: billing.MoneyFloat(0.003 * 0.9 * 0.95 * 0.8),
		},
		{
			action:          adtype.ActionClick,
			expectPotential: billing.MoneyFloat(0.5),
			expectAdvertier: billing.MoneyFloat(0.2),
			expectPublisher: billing.MoneyFloat(0.2 * 0.9 * 0.95 * 0.8),
		},
		{
			action:          adtype.ActionLead,
			expectPotential: billing.MoneyFloat(20.),
			expectAdvertier: billing.MoneyFloat(10.),
			expectPublisher: billing.MoneyFloat(10. * 0.9 * 0.95 * 0.8),
		},
	}
	for _, test := range tests {
		t.Run(test.action.String(), func(t *testing.T) {
			potential := PotentialPrice(&scope, test.action)
			advertiser := AdvertiserPrice(&scope, factors, test.action)
			publisher := PublisherPrice(&scope, factors, test.action)

			assert.Equal(t, test.expectPotential, potential, "PotentialPrice")
			assert.Equal(t, test.expectAdvertier, advertiser, "AdvertiserPrice")
			assert.Equal(t, test.expectPublisher, publisher, "PublisherPrice")

			assert.LessOrEqual(t, publisher, advertiser, "PublisherPrice <= AdvertiserPrice")
			assert.LessOrEqual(t, advertiser, potential, "AdvertiserPrice <= PotentialPrice")

			assert.Equal(t, potential, scope.PotentialPricePerAction(test.action))
			assert.Equal(t, advertiser, scope.AdvertiserPricePerAction(test.action, factors))
			assert.Equal(t, publisher, scope.PublisherPricePerAction(test.action, factors))
		})
	}
}

func TestCalculatePricesWithoutFactors(t *testing.T) {
	scope := PriceScope{CPCScope: CPCScope{MaxBidCPC: billing.MoneyFloat(1.), BidCPC: billing.MoneyFloat(0.7)}}

	for name, factors := range map[string]Factors{"nil": nil, "zero": StaticFactors{}} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, billing.MoneyFloat(0.7), AdvertiserPrice(&scope, factors, adtype.ActionClick))
			assert.Equal(t, billing.MoneyFloat(0.7), PublisherPrice(&scope, factors, adtype.ActionClick))
			assert.Equal(t, billing.MoneyFloat(1.), PotentialPrice(&scope, adtype.ActionClick))
		})
	}
}

func TestPublisherPriceFixedOverride(t *testing.T) {
	scope := PriceScope{CPMScope: CPMScope{MaxBidCPM: billing.MoneyFloat(2.), BidCPM: billing.MoneyFloat(1.)}}
	factors := fixedFactors{
		StaticFactors: StaticFactors{Commission: 0.2, Source: 0.1, Target: 0.05},
		fixed:         billing.MoneyFloat(0.005),
	}

	assert.Equal(t, billing.MoneyFloat(0.005),
		PublisherPrice(&scope, factors, adtype.ActionImpression), "fixed price of the target")
	assert.Zero(t, PublisherPrice(&scope, factors, adtype.ActionView),
		"the view action has neither a fixed nor a defined price")

	factors.fixed = 0
	assert.Equal(t, billing.MoneyFloat(0.001*0.9*0.95*0.8),
		PublisherPrice(&scope, factors, adtype.ActionImpression), "zero fixed price is ignored")
}

func TestNetworkProfit(t *testing.T) {
	scope := PriceScope{CPMScope: CPMScope{MaxBidCPM: billing.MoneyFloat(2.), BidCPM: billing.MoneyFloat(1.)}}
	factors := StaticFactors{Commission: 0.2, Source: 0.1, Target: 0.05}

	profit := NetworkProfit(&scope, factors, adtype.ActionImpression)
	assert.Equal(t, billing.MoneyFloat(0.001-0.001*0.9*0.95*0.8), profit)
	assert.Equal(t,
		AdvertiserPrice(&scope, factors, adtype.ActionImpression)-PublisherPrice(&scope, factors, adtype.ActionImpression),
		profit)
	assert.Equal(t, profit, scope.NetworkProfitPerAction(adtype.ActionImpression, factors))

	t.Run("zero factors", func(t *testing.T) {
		assert.Zero(t, NetworkProfit(&scope, StaticFactors{}, adtype.ActionImpression),
			"no factors means the publisher gets the same price as the advertiser was charged")
	})

	t.Run("fixed purchase price", func(t *testing.T) {
		fixed := fixedFactors{
			StaticFactors: factors,
			fixed:         billing.MoneyFloat(0.0006),
		}
		assert.Equal(t, billing.MoneyFloat(0.001-0.0006),
			NetworkProfit(&scope, fixed, adtype.ActionImpression),
			"profit is the declared price minus the fixed publisher payout")
	})
}

func TestBidUpPrice(t *testing.T) {
	factors := StaticFactors{Commission: 0.2, Source: 0.1, Target: 0.05}

	// True inverse of reduce: 1 / 0.8 / 0.95 / 0.9
	assert.Equal(t, billing.MoneyFloat(1./0.8/0.95/0.9), BidUpPrice(billing.MoneyFloat(1.), factors))
	assert.Equal(t, billing.MoneyFloat(1.), BidUpPrice(billing.MoneyFloat(1.), StaticFactors{}))
	assert.Equal(t, billing.MoneyFloat(1.), BidUpPrice(billing.MoneyFloat(1.), nil))

	t.Run("round-trip with PublisherPrice", func(t *testing.T) {
		floor := billing.MoneyFloat(0.0015)
		scope := PriceScope{
			CPCScope: CPCScope{MaxBidCPC: billing.MoneyFloat(1.), BidCPC: billing.MoneyFloat(0.2)},
		}
		assert.NoError(t, scope.SetBidPrice(adtype.ActionClick, floor, factors, true))
		assert.Equal(t, floor, PublisherPrice(&scope, factors, adtype.ActionClick),
			"PublisherPrice must recover the floor after SetBidPrice(..., withCommission=true)")
	})

	t.Run("factor at or above 1 yields zero", func(t *testing.T) {
		assert.Zero(t, BidUpPrice(billing.MoneyFloat(1.), StaticFactors{Commission: 1.}))
		assert.Zero(t, BidUpPrice(billing.MoneyFloat(1.), StaticFactors{Commission: 1.5}))
	})
}

func TestReduceFactorOverflow(t *testing.T) {
	assert.Zero(t, reduce(billing.MoneyFloat(1.), 1.5), "factor above 1 must not produce a negative price")
	assert.Equal(t, billing.MoneyFloat(1.), reduce(billing.MoneyFloat(1.)), "no factors keep the value")
}
