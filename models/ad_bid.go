//
// @project GeniusRabbit::corelib 2016 – 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2018
//s

package models

import (
	"github.com/geniusrabbit/gosql"

	"geniusrabbit.dev/corelib/billing"
)

// AdBid submodel
type AdBid struct {
	Bid          billing.Money             `json:"bid" validate:"notempty"`
	Applications gosql.NullableUintArray   `json:"apps,omitempty"`
	Tags         gosql.NullableStringArray `json:"tags,omitempty"`
	Zones        gosql.NullableUintArray   `json:"zones,omitempty"`
	Domains      gosql.NullableStringArray `json:"domains,omitempty"` // site domains or application bundels
	Sex          gosql.NullableUintArray   `json:"sex,omitempty"`
	Age          gosql.NullableUintArray   `json:"age,omitempty"`
	Categories   gosql.NullableUintArray   `json:"categories,omitempty"`
	Countries    gosql.NullableStringArray `json:"countries,omitempty"`
	Cities       gosql.NullableStringArray `json:"cities,omitempty"`
	Languages    gosql.NullableStringArray `json:"languages,omitempty"`
	DeviceTypes  gosql.NullableUintArray   `json:"device_types,omitempty"`
	Devices      gosql.NullableUintArray   `json:"devices,omitempty"`
	Os           gosql.NullableUintArray   `json:"os,omitempty"`
	Browsers     gosql.NullableUintArray   `json:"browsers,omitempty"`
	Hours        string                    `json:"hours,omitempty"`
}
