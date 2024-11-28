//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package admodels

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

func Test_AdModel(t *testing.T) {
	ad := Ad{
		ID:       1,
		Campaign: &Campaign{},
		Format:   &types.Format{},
		Hours:    nil,

		Weight:           10,
		FrequencyCapping: 10,
		Flags:            AdFlagActive,

		PricingModel: types.PricingModelCPM,

		Bids:        nil,
		BidPrice:    billing.MoneyFloat(0.5),
		Price:       billing.MoneyFloat(1.5),
		LeadPrice:   billing.MoneyFloat(32.5),
		DailyBudget: billing.MoneyFloat(100.),
		Budget:      billing.MoneyFloat(1000.),

		Content: nil,
	}

	t.Run("PricingModel", func(t *testing.T) {
		if ad.PricingModel != types.PricingModelCPM {
			t.Error("Wrong pricing model")
		}
	})

	t.Run("Weight", func(t *testing.T) {
		if ad.Weight != 10 {
			t.Error("Wrong weight of model")
		}
	})

	t.Run("FrequencyCapping", func(t *testing.T) {
		if ad.FrequencyCapping != 10 {
			t.Error("Wrong frequency capping of model")
		}
	})
}
