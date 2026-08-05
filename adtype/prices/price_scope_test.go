package prices

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// allModelsScope defines CPM, CPMV and CPC pricing at the same time.
func allModelsScope() PriceScope {
	return PriceScope{
		CPMScope:  CPMScope{MaxBidCPM: billing.MoneyFloat(2.), BidCPM: billing.MoneyFloat(1.)},
		CPMVScope: CPMVScope{MaxBidCPMV: billing.MoneyFloat(6.), BidCPMV: billing.MoneyFloat(3.)},
		CPCScope:  CPCScope{MaxBidCPC: billing.MoneyFloat(0.5), BidCPC: billing.MoneyFloat(0.2)},
	}
}

func TestPriceScopeDispatch(t *testing.T) {
	scope := allModelsScope()
	scope.CPAScope = CPAScope{MaxLeadPrice: billing.MoneyFloat(20.), LeadPrice: billing.MoneyFloat(10.)}

	tests := []struct {
		action      adtype.Action
		expectHas   bool
		expectPrice billing.Money
		expectMax   billing.Money
	}{
		{
			action:      adtype.ActionImpression,
			expectHas:   true,
			expectPrice: billing.MoneyFloat(0.001),
			expectMax:   billing.MoneyFloat(0.002),
		},
		{
			action:      adtype.ActionView,
			expectHas:   true,
			expectPrice: billing.MoneyFloat(0.003),
			expectMax:   billing.MoneyFloat(0.006),
		},
		{
			action:      adtype.ActionClick,
			expectHas:   true,
			expectPrice: billing.MoneyFloat(0.2),
			expectMax:   billing.MoneyFloat(0.5),
		},
		{
			action:      adtype.ActionLead,
			expectHas:   true,
			expectPrice: billing.MoneyFloat(10.),
			expectMax:   billing.MoneyFloat(20.),
		},
		{
			action:      adtype.Action(0),
			expectHas:   false,
			expectPrice: 0,
			expectMax:   0,
		},
	}
	for _, test := range tests {
		t.Run(test.action.String(), func(t *testing.T) {
			assert.Equal(t, test.expectHas, scope.HasAction(test.action), "HasAction")
			assert.Equal(t, test.expectPrice, scope.PricePerAction(test.action), "PricePerAction")
			assert.Equal(t, test.expectMax, scope.MaxPricePerAction(test.action), "MaxPricePerAction")
		})
	}
}

func TestPriceScopePartialModels(t *testing.T) {
	scope := PriceScope{CPCScope: CPCScope{BidCPC: billing.MoneyFloat(0.2)}}

	assert.False(t, scope.HasAction(adtype.ActionImpression))
	assert.False(t, scope.HasAction(adtype.ActionView))
	assert.True(t, scope.HasAction(adtype.ActionClick))
	assert.Zero(t, scope.PricePerAction(adtype.ActionImpression))
	assert.Equal(t, billing.MoneyFloat(0.2), scope.PricePerAction(adtype.ActionClick))
}

func TestPriceScopeSetBidPerAction(t *testing.T) {
	tests := []struct {
		name        string
		action      adtype.Action
		price       billing.Money
		expectPrice billing.Money
		expectErr   error
	}{
		{
			name:        "impression-in-range",
			action:      adtype.ActionImpression,
			price:       billing.MoneyFloat(0.0015),
			expectPrice: billing.MoneyFloat(0.0015),
		},
		{
			name:        "impression-clamped-by-max",
			action:      adtype.ActionImpression,
			price:       billing.MoneyFloat(0.005),
			expectPrice: billing.MoneyFloat(0.002),
		},
		{
			name:        "view-clamped-by-max",
			action:      adtype.ActionView,
			price:       billing.MoneyFloat(0.01),
			expectPrice: billing.MoneyFloat(0.006),
		},
		{
			name:        "click-in-range",
			action:      adtype.ActionClick,
			price:       billing.MoneyFloat(0.4),
			expectPrice: billing.MoneyFloat(0.4),
		},
		{
			name:        "negative-keeps-the-current-price",
			action:      adtype.ActionClick,
			price:       -1,
			expectPrice: billing.MoneyFloat(0.2),
			expectErr:   ErrNegativeBidPrice,
		},
		{
			name:        "unsupported-action",
			action:      adtype.Action(0),
			price:       billing.MoneyFloat(1.),
			expectPrice: 0,
			expectErr:   ErrUnsupportedAction,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scope := allModelsScope()
			assert.ErrorIs(t, scope.SetBidPerAction(test.action, test.price), test.expectErr)
			assert.Equal(t, test.expectPrice, scope.PricePerAction(test.action))
		})
	}
}

func TestPriceScopeSetBidPrice(t *testing.T) {
	factors := StaticFactors{Commission: 0.2, Source: 0.1, Target: 0.05}

	t.Run("without commission behaves like SetBidPerAction", func(t *testing.T) {
		scope := allModelsScope()
		assert.NoError(t, scope.SetBidPrice(adtype.ActionClick, billing.MoneyFloat(0.4), factors, false))
		assert.Equal(t, billing.MoneyFloat(0.4), scope.PricePerAction(adtype.ActionClick))
	})

	t.Run("with commission grosses the price up before storing", func(t *testing.T) {
		scope := allModelsScope()
		// 0.1 / 0.8 / 0.95 / 0.9 ≈ 0.1462, well below MaxBidCPC=0.5
		assert.NoError(t, scope.SetBidPrice(adtype.ActionClick, billing.MoneyFloat(0.1), factors, true))
		assert.Equal(t, billing.MoneyFloat(0.1/0.8/0.95/0.9), scope.PricePerAction(adtype.ActionClick))
		assert.Equal(t, billing.MoneyFloat(0.1), PublisherPrice(&scope, factors, adtype.ActionClick),
			"round-trip: PublisherPrice recovers the requested net")
	})

	t.Run("with commission still clamps to MaxBid", func(t *testing.T) {
		scope := allModelsScope()
		// Grossed-up value 0.5/0.684 ≈ 0.73 exceeds MaxBidCPC=0.5
		assert.NoError(t, scope.SetBidPrice(adtype.ActionClick, billing.MoneyFloat(0.5), factors, true))
		assert.Equal(t, billing.MoneyFloat(0.5), scope.PricePerAction(adtype.ActionClick))
	})

	t.Run("zero factors is a no-op gross-up", func(t *testing.T) {
		scope := allModelsScope()
		assert.NoError(t, scope.SetBidPrice(adtype.ActionClick, billing.MoneyFloat(0.3), StaticFactors{}, true))
		assert.Equal(t, billing.MoneyFloat(0.3), scope.PricePerAction(adtype.ActionClick))
	})

	t.Run("unsupported action", func(t *testing.T) {
		scope := allModelsScope()
		assert.ErrorIs(t, scope.SetBidPrice(adtype.Action(0), billing.MoneyFloat(1.), factors, false), ErrUnsupportedAction)
	})

	t.Run("negative price", func(t *testing.T) {
		scope := allModelsScope()
		assert.ErrorIs(t, scope.SetBidPrice(adtype.ActionClick, -1, factors, false), ErrNegativeBidPrice)
		assert.Equal(t, billing.MoneyFloat(0.2), scope.PricePerAction(adtype.ActionClick),
			"negative price must leave the previous bid unchanged")
	})
}

func TestPriceScopePrepareBidPerAction(t *testing.T) {
	scope := allModelsScope()
	assert.Equal(t, billing.MoneyFloat(0.002),
		scope.PrepareBidPerAction(adtype.ActionImpression, billing.MoneyFloat(0.01)))
	assert.Equal(t, billing.MoneyFloat(0.001),
		scope.PrepareBidPerAction(adtype.ActionImpression, billing.MoneyFloat(0.001)))
}

func TestPriceScopeEffectiveCPM(t *testing.T) {
	scope := allModelsScope()
	assert.Equal(t, billing.MoneyFloat(1.), scope.EffectiveCPM(), "fallback to the impression bid")

	scope.ECPM = billing.MoneyFloat(4.)
	assert.Equal(t, billing.MoneyFloat(4.), scope.EffectiveCPM(), "predicted ECPM has priority")
}
