//
// @project GeniusRabbit rotator 2016 – 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2018
//

package admodels

import (
	"github.com/geniusrabbit/gosql/v2"

	"geniusrabbit.dev/adcorelib/admodels/types"
	"geniusrabbit.dev/adcorelib/billing"
)

// TargetBid object
type TargetBid struct {
	ECPM      billing.Money
	BidPrice  billing.Money // Max price per one View (used in DSP auction)
	Price     billing.Money // Price per one view or click
	LeadPrice billing.Money // Price per one lead
	Ad        *Ad
}

// Less then other target
func (t TargetBid) Less(tb TargetBid) bool {
	return t.ECPM < tb.ECPM
}

// AdBid submodel
type AdBid struct {
	BidPrice    billing.Money                   `json:"bid_price"`
	Price       billing.Money                   `json:"price" validate:"notempty"`
	LeadPrice   billing.Money                   `json:"lead_price"`
	Tags        gosql.NullableStringArray       `json:"tags,omitempty"`
	Zones       gosql.NullableNumberArray[uint] `json:"zones,omitempty"`
	Domains     gosql.NullableStringArray       `json:"domains,omitempty"` // site domains or application bundels
	Sex         gosql.NullableNumberArray[uint] `json:"sex,omitempty"`
	Age         uint                            `json:"age,omitempty"`
	Categories  gosql.NullableNumberArray[uint] `json:"categories,omitempty"`
	Countries   gosql.NullableNumberArray[uint] `json:"countries,omitempty"`
	Cities      gosql.NullableStringArray       `json:"cities,omitempty"`
	Languages   gosql.NullableNumberArray[uint] `json:"languages,omitempty"`
	DeviceTypes gosql.NullableNumberArray[uint] `json:"device_types,omitempty"`
	Devices     gosql.NullableNumberArray[uint] `json:"devices,omitempty"`
	Os          gosql.NullableNumberArray[uint] `json:"os,omitempty"`
	Browsers    gosql.NullableNumberArray[uint] `json:"browsers,omitempty"`
	Hours       types.Hours                     `json:"hours,omitempty"`
}

// Test is it suites by pointer
func (bid *AdBid) Test(pointer types.TargetPointer) bool {
	return true &&
		(bid.Tags.Len() < 1 || bid.Tags.OneOf(pointer.Tags())) &&
		(bid.Zones.Len() < 1 || bid.Zones.IndexOf(uint(pointer.TargetID())) != -1) &&
		(bid.Domains.Len() < 1 || bid.Domains.OneOf(pointer.Domain())) &&
		(bid.Sex.Len() < 1 || bid.Sex.IndexOf(pointer.Sex()) != -1) &&
		(bid.Age < 1 && bid.Age >= pointer.Age()) &&
		(bid.Categories.Len() < 1 || bid.Categories.OneOf(pointer.Categories())) &&
		(bid.Countries.Len() < 1 || bid.Countries.IndexOf(pointer.GeoID()) != -1) &&
		(bid.Cities.Len() < 1 || bid.Cities.IndexOf(pointer.City()) != -1) &&
		(bid.Languages.Len() < 1 || bid.Languages.IndexOf(pointer.LanguageID()) != -1) &&
		(bid.DeviceTypes.Len() < 1 || bid.DeviceTypes.IndexOf(pointer.DeviceType()) != -1) &&
		(bid.Devices.Len() < 1 || bid.Devices.IndexOf(pointer.DeviceID()) != -1) &&
		(bid.Os.Len() < 1 || bid.Os.IndexOf(pointer.OSID()) != -1) &&
		(bid.Browsers.Len() < 1 || bid.Browsers.IndexOf(pointer.BrowserID()) != -1) &&
		bid.Hours.TestTime(pointer.Time())
}

// AdBids list
type AdBids []AdBid

// Bid by pointer
func (bids AdBids) Bid(pointer types.TargetPointer) *AdBid {
	for i, bid := range bids {
		if bid.Test(pointer) {
			return &bids[i]
		}
	}
	return nil
}
