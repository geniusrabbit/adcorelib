package adtype

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/stretchr/testify/assert"
)

type revenueShareReduceTest struct {
	ComissionShare     float64
	RevenueShareReduce float64
}

func (r *revenueShareReduceTest) ComissionShareFactor() float64 {
	return r.ComissionShare / 100.
}

func (r *revenueShareReduceTest) RevenueShareReduceFactor() float64 {
	return r.RevenueShareReduce / 100.
}

func TestPriceCorrection(t *testing.T) {
	item := &revenueShareReduceTest{
		ComissionShare:     5,  // System comission
		RevenueShareReduce: 15, // Potential descrepancy
	}
	price := billing.MoneyFloat(1.123)
	price += PriceSourceFactors(price, &SourceEmpty{PriceCorrectionReduce: 10}, true)
	if assert.Equal(t, billing.MoneyFloat(1.123*0.9), price, "source price factor") {
		price += PriceSystemComission(price, item, true)
		if assert.Equal(t, billing.MoneyFloat(1.123*0.9*0.95), price, "system comission") {
			price += PriceRevenueShareReduceFactors(price, item, true)
			assert.Equal(t, billing.MoneyFloat(1.123*0.9*0.95*0.85), price, "revenue share reduce")
		}
	}
}

func TestPriceCorrection2(t *testing.T) {
	var (
		acc = &admodels.Account{
			IDval:        1,
			RevenueShare: 85,
		}
		imp  = Impression{Target: &admodels.Smartlink{Acc: acc}}
		item = newRTBResponse(acc, imp)
	)
	price := billing.MoneyFloat(1.123)
	price += PriceFactorFromList(SourcePriceFactor, SystemComissionPriceFactor, TargetReducePriceFactor).
		Remove(price, item)
	assert.True(t, price > 0 && price < billing.MoneyFloat(1.123))
	assert.Equal(t, billing.MoneyFloat(1.123*0.85).Float64(), price.Float64())
}
