package adresponse

import (
	"reflect"
	"testing"

	"github.com/bsm/openrtb"
	"github.com/stretchr/testify/assert"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

func TestItemPricing(t *testing.T) {
	var (
		acc = &admodels.Account{
			IDval:        1,
			RevenueShare: 0.9,
		}
		imp   = adtype.Impression{Target: &admodels.Smartlink{Acc: acc}}
		items = []adtype.ResponserItem{newRTBResponse(acc, imp)}
	)

	for _, item := range items {
		prefix := reflect.TypeOf(item).String()

		t.Run(prefix+"_empty_lead_price", func(t *testing.T) {
			if item.Price(admodels.ActionLead) != 0 {
				t.Error("lead_price should be empty")
			}
		})

		t.Run(prefix+"_bid_price", func(t *testing.T) {
			if item.Price(admodels.ActionImpression) != billing.MoneyFloat(10.) {
				t.Errorf("target price must be 10, not %.3f", item.Price(admodels.ActionImpression).Float64())
			}
		})

		t.Run(prefix+"_revenue_value", func(t *testing.T) {
			rev := item.RevenueShareFactor() * item.Price(admodels.ActionImpression).Float64()
			assert.Equal(t, float64(9), rev, "wrong_revenue value")
		})

		t.Run(prefix+"_comission_value", func(t *testing.T) {
			com := item.ComissionShareFactor() * item.Price(admodels.ActionImpression).Float64()
			assert.True(t, com >= 0.999 && com <= 1, "wrong_comission value")
		})

		t.Run(prefix+"_cpm_price", func(t *testing.T) {
			if item.CPMPrice() != billing.MoneyFloat(5.) {
				t.Errorf("cpm_price value: 5 != %.3f", item.CPMPrice().Float64())
			}
		})
	}
}

func TestPriceCorrection(t *testing.T) {
	var (
		acc = &admodels.Account{
			IDval:        1,
			RevenueShare: 0.85,
		}
		imp  = adtype.Impression{Target: &admodels.Smartlink{Acc: acc}}
		item = newRTBResponse(acc, imp)
	)
	price := billing.MoneyFloat(1.123)
	price += adtype.PriceFactorFromList(adtype.SourcePriceFactor, adtype.SystemComissionPriceFactor, adtype.TargetReducePriceFactor).
		Remove(price, item)
	assert.True(t, price > 0 && price < billing.MoneyFloat(1.123))
	assert.Equal(t, billing.MoneyFloat(1.123*0.85).Float64(), price.Float64())
}

func newRTBResponse(_ *admodels.Account, imp adtype.Impression) *ResponseBidItem {
	return &ResponseBidItem{
		ItemID:      "1",
		Src:         &adtype.SourceEmpty{PriceCorrectionReduce: 0},
		Req:         &adtype.BidRequest{ID: "xxx", Imps: []adtype.Impression{imp}},
		Imp:         &imp,
		Bid:         &openrtb.Bid{Price: 60},
		BidPrice:    billing.MoneyFloat(10.),
		CPMBidPrice: billing.MoneyFloat(5.),
		SecondAd:    adtype.SecondAd{},
	}
}
