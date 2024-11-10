package adresponse

import (
	"reflect"
	"testing"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/stretchr/testify/assert"
)

func Test_ItemPricing(t *testing.T) {
	var (
		acc = &admodels.Account{
			IDval:        1,
			RevenueShare: 0.9,
		}
		imp   = adtype.Impression{Target: &admodels.Smartlink{Acc: acc}}
		items = []adtype.ResponserItem{newAdResponse(acc, imp)}
	)

	for _, item := range items {
		prefix := reflect.TypeOf(item).String()

		t.Run(prefix+"_lead_price", func(t *testing.T) {
			assert.Equal(t, leadPrice(item).Float64(), item.Price(admodels.ActionLead).Float64(), "wrong_lead_price")
		})

		// Check price per one view
		t.Run(prefix+"_bid_price", func(t *testing.T) {
			assert.Equal(t, 10., item.Price(admodels.ActionView).Float64(), "wrong_bid_price")
		})

		// Check revenue per one view
		t.Run(prefix+"_revenue_value", func(t *testing.T) {
			rev := item.RevenueShareFactor() * item.Price(admodels.ActionView).Float64()
			assert.Equal(t, float64(9), rev, "wrong_revenue value")
		})

		// Check comission per one view
		t.Run(prefix+"_comission_value", func(t *testing.T) {
			com := item.ComissionShareFactor() * item.Price(admodels.ActionView).Float64()
			assert.True(t, com >= 0.999 && com <= 1, "wrong_comission value")
		})

		t.Run(prefix+"_cpm_price", func(t *testing.T) {
			assert.Equal(t, 5000., item.CPMPrice().Float64(), "wrong_cpm_price")
		})
	}
}

func newAdResponse(_ *admodels.Account, imp adtype.Impression) *ResponseAdItem {
	return &ResponseAdItem{
		ItemID: "1",
		Src:    &adtype.SourceEmpty{PriceCorrectionReduce: 0},
		Req:    &adtype.BidRequest{ID: "xxx", Imps: []adtype.Impression{imp}},
		Imp:    &imp,
		Ad: &admodels.Ad{
			Format:      &types.Format{Width: 250, Height: 250},
			BidPrice:    billing.MoneyFloat(5.),
			Price:       billing.MoneyFloat(60.),
			LeadPrice:   billing.MoneyFloat(120.),
			DailyBudget: billing.MoneyFloat(1200.),
			Budget:      billing.MoneyFloat(10000.),
			Hours:       nil,
		},
		PriceScope: adtype.PriceScope{
			MaxBidPrice: billing.MoneyFloat(60.),
			BidPrice:    billing.MoneyFloat(5.),
			ViewPrice:   billing.MoneyFloat(10.),
			ClickPrice:  billing.MoneyFloat(0.),
			LeadPrice:   billing.MoneyFloat(120.),
		},
		// BidPrice:    billing.MoneyFloat(10.),
		// CPMBidPrice: billing.MoneyFloat(5.),
		SecondAd: adtype.SecondAd{},
	}
}

func leadPrice(item adtype.ResponserItem) billing.Money {
	switch it := item.(type) {
	case *ResponseAdItem:
		return it.PriceScope.LeadPrice
	}
	return 0
}
