package adtype

import (
	"fmt"
	"testing"

	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/stretchr/testify/assert"
)

type revenueShareReduceTest struct {
	ComissionShare     float64
	RevenueShareReduce float64
}

func (r *revenueShareReduceTest) ComissionShareFactor() float64 {
	return r.ComissionShare
}

func (r *revenueShareReduceTest) RevenueShareReduceFactor() float64 {
	return r.RevenueShareReduce
}

func TestPriceCorrection(t *testing.T) {
	item := &revenueShareReduceTest{
		ComissionShare:     0.05, // System comission
		RevenueShareReduce: 0.15, // Potential descrepancy
	}
	const testNumber = 1.123
	price := billing.MoneyFloat(testNumber)

	// Price correction for source descrepancy factor
	price += PriceSourceFactors(price, &SourceEmpty{PriceCorrectionReduce: 0.1}, true)

	if assert.Equal(t, billing.MoneyFloat(testNumber/1.1), price, "source price factor") {

		// Price correction for system comission
		price += PriceSystemComission(price, item, true)

		if assert.Equal(t, billing.MoneyFloat(testNumber/1.1/1.05), price, "system comission") {
			// Price correction for revenue share reduce
			price += PriceRevenueShareReduceFactors(price, item, true)
			assert.Equal(t, billing.MoneyFloat(testNumber/1.1/1.05/1.15), price, "revenue share reduce")
		}
	}
}

func TestPriceConvertions(t *testing.T) {
	tests := []struct {
		Price    billing.Money
		Factor   PriceFactor
		Expected float64
		Item     ResponserItem
	}{
		{0, NonePriceFactor, 0, nil},
		{1, NonePriceFactor, 0, nil},
		{1, SourcePriceFactor, 0, &ResponseItemEmpty{}},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			assert.Equal(t, test.Expected, test.Factor.AddComission(test.Price, test.Item).Float64(), "add comission")
			assert.Equal(t, test.Expected, test.Factor.RemoveComission(test.Price, test.Item).Float64(), "remove comission")
		})
	}
}

func TestAdjustPrice(t *testing.T) {
	const testNumber = 1.123
	originalPrice := billing.MoneyFloat(testNumber)

	// Adjust adds to price for 10% factor comission
	price := originalPrice + AdjustPrice(originalPrice, 0.1, false)

	if assert.Equal(t, billing.MoneyFloat(testNumber*1.1), price, "adjust price") {
		// Adjust removes from price for 10% factor comission
		price = price + AdjustPrice(price, 0.1, true)
		assert.Equal(t, originalPrice, price, "adjust price")
	}
}
