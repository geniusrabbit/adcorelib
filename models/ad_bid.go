//
// @project GeniusRabbit::corelib 2016 – 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2018
//s

package models

import (
	"github.com/geniusrabbit/gosql/v2"

	"geniusrabbit.dev/adcorelib/admodels/types"
	"geniusrabbit.dev/adcorelib/billing"
)

// AdBid submodel
type AdBid struct {
	BidPrice     billing.Money                   `json:"bid_price"`
	Price        billing.Money                   `json:"price" validate:"notempty"`
	LeadPrice    billing.Money                   `json:"lead_price"`
	Applications gosql.NullableNumberArray[uint] `json:"apps,omitempty"`
	Tags         gosql.NullableStringArray       `json:"tags,omitempty"`
	Zones        gosql.NullableNumberArray[uint] `json:"zones,omitempty"`
	Domains      gosql.NullableStringArray       `json:"domains,omitempty"` // site domains or application bundels
	Sex          gosql.NullableNumberArray[uint] `json:"sex,omitempty"`
	Age          gosql.NullableNumberArray[uint] `json:"age,omitempty"`
	Categories   gosql.NullableNumberArray[uint] `json:"categories,omitempty"`
	Countries    gosql.NullableStringArray       `json:"countries,omitempty"`
	Cities       gosql.NullableStringArray       `json:"cities,omitempty"`
	Languages    gosql.NullableStringArray       `json:"languages,omitempty"`
	DeviceTypes  gosql.NullableNumberArray[uint] `json:"device_types,omitempty"`
	Devices      gosql.NullableNumberArray[uint] `json:"devices,omitempty"`
	Os           gosql.NullableNumberArray[uint] `json:"os,omitempty"`
	Browsers     gosql.NullableNumberArray[uint] `json:"browsers,omitempty"`
	Hours        types.Hours                     `json:"hours,omitempty"`
}
