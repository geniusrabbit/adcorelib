//
// @project GeniusRabbit corelib 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package models

import "geniusrabbit.dev/adcorelib/admodels/types"

// PricingModel value
type PricingModel = types.PricingModel

// PricingModel consts
const (
	PricingModelUndefined = types.PricingModelUndefined
	PricingModelCPM       = types.PricingModelCPM
	PricingModelCPC       = types.PricingModelCPC
	PricingModelCPA       = types.PricingModelCPA
)

// PricingModelByName string
func PricingModelByName(model string) PricingModel {
	return types.PricingModelByName(model)
}
