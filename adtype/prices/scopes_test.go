package prices

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/geniusrabbit/adcorelib/billing"
)

func TestCPMUnitConversion(t *testing.T) {
	tests := []struct {
		name        string
		cpm         billing.Money
		expectPrice billing.Money
	}{
		{name: "zero", cpm: 0, expectPrice: 0},
		{name: "one-dollar-cpm", cpm: billing.MoneyFloat(1.), expectPrice: billing.MoneyFloat(0.001)},
		{name: "fractional-cpm", cpm: billing.MoneyFloat(1.5), expectPrice: billing.MoneyFloat(0.0015)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectPrice, PriceFromCPM(test.cpm), "PriceFromCPM")
			assert.Equal(t, test.cpm, CPMFromPrice(test.expectPrice), "CPMFromPrice")
		})
	}
}

func TestScopesDefined(t *testing.T) {
	assert.False(t, (&CPMScope{}).HasCPM())
	assert.True(t, (&CPMScope{BidCPM: 1}).HasCPM())
	assert.True(t, (&CPMScope{MaxBidCPM: 1}).HasCPM(), "the max bid alone defines the pricing")

	assert.False(t, (&CPMVScope{}).HasCPMV())
	assert.True(t, (&CPMVScope{BidCPMV: 1}).HasCPMV())
	assert.True(t, (&CPMVScope{MaxBidCPMV: 1}).HasCPMV())

	assert.False(t, (&CPCScope{}).HasCPC())
	assert.True(t, (&CPCScope{BidCPC: 1}).HasCPC())
	assert.True(t, (&CPCScope{MaxBidCPC: 1}).HasCPC())

	assert.False(t, (&CPAScope{}).HasCPA())
	assert.True(t, (&CPAScope{LeadPrice: 1}).HasCPA())
	assert.True(t, (&CPAScope{MaxLeadPrice: 1}).HasCPA())
}

func TestSetBidClamping(t *testing.T) {
	const initialBid = billing.Money(40)

	tests := []struct {
		name       string
		maxBid     billing.Money
		bid        billing.Money
		expect     billing.Money
		expectPrep billing.Money
		expectErr  error
	}{
		{name: "in-range", maxBid: 100, bid: 50, expect: 50, expectPrep: 50},
		{name: "equal-to-max", maxBid: 100, bid: 100, expect: 100, expectPrep: 100},
		{name: "above-max-clamped", maxBid: 100, bid: 150, expect: 100, expectPrep: 100},
		{name: "raise-up-to-max-allowed", maxBid: 100, bid: 99, expect: 99, expectPrep: 99},
		{name: "no-max-limit", maxBid: 0, bid: 150, expect: 150, expectPrep: 150},
		{name: "zero", maxBid: 100, bid: 0, expect: 0, expectPrep: 0},
		{
			name:       "negative-keeps-the-current-bid",
			maxBid:     100,
			bid:        -1,
			expect:     initialBid,
			expectPrep: 0,
			expectErr:  ErrNegativeBidPrice,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cpm := CPMScope{MaxBidCPM: test.maxBid, BidCPM: initialBid}
			assert.ErrorIs(t, cpm.SetBidCPM(test.bid), test.expectErr, "SetBidCPM")
			assert.Equal(t, test.expect, cpm.BidCPM, "BidCPM")
			assert.Equal(t, test.expectPrep, cpm.PrepareBidCPM(test.bid), "PrepareBidCPM")

			cpmv := CPMVScope{MaxBidCPMV: test.maxBid, BidCPMV: initialBid}
			assert.ErrorIs(t, cpmv.SetBidCPMV(test.bid), test.expectErr, "SetBidCPMV")
			assert.Equal(t, test.expect, cpmv.BidCPMV, "BidCPMV")
			assert.Equal(t, test.expectPrep, cpmv.PrepareBidCPMV(test.bid), "PrepareBidCPMV")

			cpc := CPCScope{MaxBidCPC: test.maxBid, BidCPC: initialBid}
			assert.ErrorIs(t, cpc.SetBidCPC(test.bid), test.expectErr, "SetBidCPC")
			assert.Equal(t, test.expect, cpc.BidCPC, "BidCPC")
			assert.Equal(t, test.expectPrep, cpc.PrepareBidCPC(test.bid), "PrepareBidCPC")

			cpa := CPAScope{MaxLeadPrice: test.maxBid, LeadPrice: initialBid}
			assert.ErrorIs(t, cpa.SetLeadPrice(test.bid), test.expectErr, "SetLeadPrice")
			assert.Equal(t, test.expect, cpa.LeadPrice, "LeadPrice")
			assert.Equal(t, test.expectPrep, cpa.PrepareLeadPrice(test.bid), "PrepareLeadPrice")
		})
	}
}

// TestEmbedding ensures that all the scopes can be embedded into the same structure
// without the selector conflicts.
func TestEmbedding(t *testing.T) {
	type campaign struct {
		CPMScope
		CPMVScope
		CPCScope
		CPAScope
	}

	item := campaign{
		CPMScope:  CPMScope{MaxBidCPM: billing.MoneyFloat(2.)},
		CPMVScope: CPMVScope{MaxBidCPMV: billing.MoneyFloat(3.)},
		CPCScope:  CPCScope{MaxBidCPC: billing.MoneyFloat(0.1)},
		CPAScope:  CPAScope{MaxLeadPrice: billing.MoneyFloat(5.)},
	}

	assert.True(t, item.HasCPM() && item.HasCPMV() && item.HasCPC() && item.HasCPA())
	assert.NoError(t, item.SetBidCPM(billing.MoneyFloat(1.)))
	assert.NoError(t, item.SetBidCPMV(billing.MoneyFloat(1.5)))
	assert.NoError(t, item.SetBidCPC(billing.MoneyFloat(0.05)))
	assert.NoError(t, item.SetLeadPrice(billing.MoneyFloat(4.)))

	assert.Equal(t, billing.MoneyFloat(1.), item.BidCPM)
	assert.Equal(t, billing.MoneyFloat(1.5), item.BidCPMV)
	assert.Equal(t, billing.MoneyFloat(0.05), item.BidCPC)
	assert.Equal(t, billing.MoneyFloat(4.), item.LeadPrice)
}
