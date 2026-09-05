package trafaret

import (
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adquery/bidresponse"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/adtype/prices"
	"github.com/geniusrabbit/adcorelib/billing"
)

func blankAd(id, impID string, eCPM float64) *bidresponse.ResponseItemBlank {
	return &bidresponse.ResponseItemBlank{
		ItemID:          id,
		Imp:             &adtype.Impression{ID: impID},
		Src:             &adtype.SourceEmpty{},
		FormatVal:       &types.Format{},
		PricingModelVal: types.PricingModelCPM,
		PriceScope:      prices.PriceScope{ECPM: billing.MoneyFloat(eCPM)},
	}
}

func blockAd(items ...adtype.ResponseItem) *bidresponse.ResponseItemBlock {
	return &bidresponse.ResponseItemBlock{Items: items}
}

func adIDs(list []adtype.ResponseItemCommon) []string {
	ids := make([]string, 0, len(list))
	for _, ad := range list {
		if ad != nil {
			ids = append(ids, ad.ID())
		}
	}
	return ids
}

func adsListRealSize(list []adtype.ResponseItemCommon) int {
	realSize := 0
	for _, item := range list {
		realSize += adSize(item)
	}
	return realSize
}
