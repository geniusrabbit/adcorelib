package adtype

import (
	"testing"

	"geniusrabbit.dev/corelib/billing"
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
	price -= PriceSourceFactors(price, &SourceEmpty{PriceCorrectionReduce: 10}, true)
	price -= PriceSystemComission(price, item, true)
	price -= PriceRevenueShareReduceFactors(price, item, true)
	assert.True(t, price > 0 && price < billing.MoneyFloat(1.123))
}
