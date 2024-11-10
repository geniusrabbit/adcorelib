//
// @project GeniusRabbit corelib 2016 – 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2018
//

package admodels

import (
	"github.com/geniusrabbit/gosql/v2"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

// AdBid submodel
type AdBid struct {
	// Pricing
	BidPrice  billing.Money `json:"bid_price"`                 // Max price per one View (used in DSP auction)
	Price     billing.Money `json:"price" validate:"notempty"` // Price per one view or click
	LeadPrice billing.Money `json:"lead_price"`                // Price per one lead
	TestPrice billing.Money `json:"test_price"`                // Price for test period per one view (as CPM mode)

	// Targeting
	Tags        gosql.NullableStringArray                `json:"tags,omitempty"`
	Zones       gosql.NullableOrderedNumberArray[uint64] `json:"zones,omitempty"`
	Domains     gosql.NullableStringArray                `json:"domains,omitempty"` // site domains or application bundels
	Sex         gosql.NullableOrderedNumberArray[uint]   `json:"sex,omitempty"`
	Age         uint                                     `json:"age,omitempty"`
	Categories  gosql.NullableOrderedNumberArray[uint64] `json:"categories,omitempty"`
	Countries   gosql.NullableOrderedNumberArray[uint64] `json:"countries,omitempty"`
	Cities      gosql.NullableStringArray                `json:"cities,omitempty"`
	Languages   gosql.NullableOrderedNumberArray[uint64] `json:"languages,omitempty"`
	DeviceTypes gosql.NullableOrderedNumberArray[uint64] `json:"device_types,omitempty"`
	Devices     gosql.NullableOrderedNumberArray[uint64] `json:"devices,omitempty"`
	OS          gosql.NullableOrderedNumberArray[uint64] `json:"os,omitempty"`
	Browsers    gosql.NullableOrderedNumberArray[uint64] `json:"browsers,omitempty"`
	Hours       types.Hours                              `json:"hours,omitempty"`
}

// Test is it suites by pointer
func (bid *AdBid) Test(pointer types.TargetPointer) bool {
	return true &&
		(bid.Tags.Len() < 1 || bid.Tags.OneOf(pointer.Tags())) &&
		(bid.Zones.Len() < 1 || bid.Zones.IndexOf(pointer.TargetID()) != -1) &&
		(bid.Domains.Len() < 1 || bid.Domains.OneOf(pointer.Domain())) &&
		(bid.Sex.Len() < 1 || bid.Sex.IndexOf(pointer.Sex()) != -1) &&
		(bid.Age < 1 && bid.Age >= pointer.Age()) &&
		(bid.Categories.Len() < 1 || bid.Categories.OneOf(pointer.Categories())) &&
		(bid.Countries.Len() < 1 || bid.Countries.IndexOf(pointer.GeoID()) != -1) &&
		(bid.Cities.Len() < 1 || bid.Cities.IndexOf(pointer.City()) != -1) &&
		(bid.Languages.Len() < 1 || bid.Languages.IndexOf(pointer.LanguageID()) != -1) &&
		(bid.DeviceTypes.Len() < 1 || bid.DeviceTypes.IndexOf(pointer.DeviceType()) != -1) &&
		(bid.Devices.Len() < 1 || bid.Devices.IndexOf(pointer.DeviceID()) != -1) &&
		(bid.OS.Len() < 1 || bid.OS.IndexOf(pointer.OSID()) != -1) &&
		(bid.Browsers.Len() < 1 || bid.Browsers.IndexOf(pointer.BrowserID()) != -1) &&
		bid.Hours.TestTime(pointer.Time())
}

// AdBids list
type AdBids []AdBid

// Bid by pointer
func (bids AdBids) Bid(pointer types.TargetPointer) *AdBid {
	for i := 0; i < len(bids); i++ {
		if bids[i].Test(pointer) {
			return &bids[i]
		}
	}
	return nil
}
