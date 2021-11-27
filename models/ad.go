//
// @project GeniusRabbit::corelib 2016 – 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2018
//

package models

import (
	"time"

	"github.com/geniusrabbit/gosql"
	"github.com/guregu/null"

	"geniusrabbit.dev/corelib/billing"
)

//
// Context describe information about ads
// {
//	 "type": "direct" | "banner" | "context" | "teaser" | "proxy",
//   "ads":[
//		 {"id": <number>, }
// 	 ],
// }
//
// Files
// [
// 	{"id": 123, "hashid": "dhg321h3ndp43u2hfc", "path": "images/a/c/banner1.jpg"},
// ]
//

// Ad model describesinformation about one paticular advertisement
type Ad struct {
	// ID number of the Advertisement in DB
	ID uint64 `json:"id"`

	// Owner Campaign of the Ad
	Campaign   *Campaign `json:"campaign,omitempty"`
	CampaignID uint64    `json:"campaign_id"`

	// Extended bid information from []AdBid - [{"cc":"GB","bid":1000}]
	Bids gosql.NullableJSON `json:"bids,omitempty"`

	// Status of the approvements
	Status ApproveStatus `json:"status,omitempty"`

	// Is Active advertisement
	Active ActiveStatus `json:"active,omitempty"`

	// Format of the advertisement with structure of allowed items
	Format   *Format `json:"format,omitempty" gorm:"association_autoupdate:false"`
	FormatID uint64  `json:"format_id,omitempty"`

	// If advertisement is streatch format then might be needded to support minimal and maximal sizes
	MinWidth  int `json:"min_width,omitempty"`
	MinHeight int `json:"min_height,omitempty"`
	MaxWidth  int `json:"max_width,omitempty"`
	MaxHeight int `json:"max_height,omitempty"`

	// Pricing model of the Ad (CPM/CPC/CPA/etc.)
	PricingModel PricingModel `json:"pricing_model"`

	// Money limit counters
	BidPrice        billing.Money `json:"bid_price,omitempty"`         // Maximal bid for RTB
	Price           billing.Money `json:"price,omitempty"`             // Price per pricing_model
	LeadPrice       billing.Money `json:"lead_bid,omitempty"`          // Price of lead to calculate effectivity
	DailyBudget     billing.Money `json:"daily_budget,omitempty"`      // Max daily budget spent
	Budget          billing.Money `json:"budget,omitempty"`            // Money budget for whole time
	DailyTestBudget billing.Money `json:"daily_test_budget,omitempty"` // Test money amount a day (it stops automaticaly if not profit for this amount)
	TestBudget      billing.Money `json:"test_budget,omitempty"`       // Total test budget of whole time

	// Link to the target
	Link string `json:"link"`

	// Context contains the most improtant information about Ad
	Context gosql.NullableJSON `json:"context,omitempty"`

	// Weight of the Ad in rotation
	Weight int `json:"weight,omitempty"`

	// Frequency Capping of advertisement display to one user
	FrequencyCapping uint `json:"fc,omitempty"`

	// Hours targetting 168 simbols. Every simbol means hour active or blocked
	// 7 lines [day of week] + 24 hours as '1' or '0'
	// '*' or empty - all is on
	Hours string `json:"hours,omitempty"`

	// Assets related to advertisement
	Assets []*AdFile `gorm:"many2many:m2m_adv_ad_file_ad;association_autoupdate:false"`

	// Time marks
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

// TableName in database
func (a *Ad) TableName() string {
	return "adv_ad"
}
