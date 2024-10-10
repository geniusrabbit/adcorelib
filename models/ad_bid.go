//
// @project GeniusRabbit corelib 2016 – 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2018
//s

package models

import (
	"github.com/geniusrabbit/gosql/v2"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

// AdBid submodel
type AdBid struct {
	BidPrice  billing.Money `json:"bid_price"`
	Price     billing.Money `json:"price" validate:"notempty"`
	LeadPrice billing.Money `json:"lead_price"`

	// Targeting options
	Applications gosql.NullableNumberArray[uint64] `json:"apps,omitempty"`
	Tags         gosql.NullableStringArray         `json:"tags,omitempty"`
	Zones        gosql.NullableNumberArray[uint64] `json:"zones,omitempty"`
	Domains      gosql.NullableStringArray         `json:"domains,omitempty"` // site domains or application bundels
	Sex          gosql.NullableNumberArray[uint]   `json:"sex,omitempty"`
	Age          gosql.NullableNumberArray[uint]   `json:"age,omitempty"`
	Categories   gosql.NullableNumberArray[uint64] `json:"categories,omitempty"`
	Countries    gosql.NullableStringArray         `json:"countries,omitempty"`
	Cities       gosql.NullableStringArray         `json:"cities,omitempty"`
	Languages    gosql.NullableStringArray         `json:"languages,omitempty"`
	DeviceTypes  gosql.NullableNumberArray[uint64] `json:"device_types,omitempty"`
	Devices      gosql.NullableNumberArray[uint64] `json:"devices,omitempty"`
	Os           gosql.NullableNumberArray[uint64] `json:"os,omitempty"`
	Browsers     gosql.NullableNumberArray[uint64] `json:"browsers,omitempty"`
	Hours        types.Hours                       `json:"hours,omitempty"`
}
