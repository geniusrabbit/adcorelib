package adtype

import (
	"reflect"
	"testing"

	"github.com/bsm/openrtb"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

func Test_ItemPricing(t *testing.T) {
	var (
		comp = &admodels.Company{
			ID:           1,
			RevenueShare: 90,
		}
		imp   = Impression{Target: &admodels.Smartlink{Comp: comp}}
		items = []ResponserItem{newRTBResponse(comp, imp), newAdResponse(comp, imp)}
	)

	for _, item := range items {
		prefix := reflect.TypeOf(item).String()

		t.Run(prefix+"_empty_lead_price", func(t *testing.T) {
			if item.Price(admodels.ActionLead) != leadPrice(item) {
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
			if rev != 9 {
				t.Errorf("wrong_revenue value: %.9f", rev)
			}
		})

		t.Run(prefix+"_comission_value", func(t *testing.T) {
			com := item.ComissionShareFactor() * item.Price(admodels.ActionImpression).Float64()
			if com != 1 {
				t.Errorf("wrong_comission value: %.3f", com)
			}
		})

		t.Run(prefix+"_cpm_price", func(t *testing.T) {
			if item.CPMPrice() != billing.MoneyFloat(5.) {
				t.Errorf("cpm_price value: 5 != %.3f", item.CPMPrice().Float64())
			}
		})
	}
}

func newRTBResponse(_ *admodels.Company, imp Impression) *ResponseBidItem {
	return &ResponseBidItem{
		ItemID:      "1",
		Src:         &SourceEmpty{PriceCorrectionReduce: 0},
		Req:         &BidRequest{ID: "xxx", Imps: []Impression{imp}},
		Imp:         &imp,
		Bid:         &openrtb.Bid{Price: 60},
		BidPrice:    billing.MoneyFloat(10.),
		CPMBidPrice: billing.MoneyFloat(5.),
		SecondAd:    SecondAd{},
	}
}

func newAdResponse(_ *admodels.Company, imp Impression) *ResponseAdItem {
	return &ResponseAdItem{
		ItemID: "1",
		Src:    &SourceEmpty{PriceCorrectionReduce: 0},
		Req:    &BidRequest{ID: "xxx", Imps: []Impression{imp}},
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
		BidPrice:    billing.MoneyFloat(10.),
		CPMBidPrice: billing.MoneyFloat(5.),
		SecondAd:    SecondAd{},
	}
}

func leadPrice(item ResponserItem) billing.Money {
	switch it := item.(type) {
	case *ResponseAdItem:
		return it.Ad.LeadPrice
	}
	return 0
}
