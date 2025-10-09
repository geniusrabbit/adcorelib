package trafaret

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adquery/bidresponse"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/price"
	"github.com/stretchr/testify/assert"
)

func TestFiller(t *testing.T) {
	filler := &Filler{}

	src1 := adtype.SourceEmpty{}

	imp1 := adtype.Impression{
		ID: "imp1",
	}

	format1 := types.Format{}

	ad1 := &bidresponse.ResponseItemBlank{
		ItemID:          "ad1",
		Imp:             &imp1,
		Src:             &src1,
		FormatVal:       &format1,
		PricingModelVal: types.PricingModelCPM,
		PriceScope:      price.PriceScope{ECPM: billing.MoneyFloat(1.0)},
	}

	ad2 := &bidresponse.ResponseItemBlank{
		ItemID:          "ad2",
		Imp:             &imp1,
		Src:             &src1,
		FormatVal:       &format1,
		PricingModelVal: types.PricingModelCPM,
		PriceScope:      price.PriceScope{ECPM: billing.MoneyFloat(0.5)},
	}

	ad3 := &bidresponse.ResponseItemBlank{
		ItemID:          "ad3",
		Imp:             &imp1,
		Src:             &src1,
		FormatVal:       &format1,
		PricingModelVal: types.PricingModelCPM,
		PriceScope:      price.PriceScope{ECPM: billing.MoneyFloat(2.0)},
	}

	adx1 := &bidresponse.ResponseItemBlock{
		Items: []adtype.ResponseItem{ad1, ad3},
	}

	filler.Push(0.3, ad1, ad2)
	filler.Push(0.5, ad1, ad2, ad3)
	filler.Push(0.7, adx1, ad3)

	for _, size := range []int{1, 2, 3} {
		ls := filler.Copy().Fill("imp1", size)
		resSize := adsListRealSize(ls)
		assert.Equal(t, size, resSize, "Expected %d items, got %d", size, resSize)
	}
}

func adsListRealSize(list []adtype.ResponseItemCommon) int {
	realSize := 0
	for _, item := range list {
		realSize += adSize(item)
	}
	return realSize
}
