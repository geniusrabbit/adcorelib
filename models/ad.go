//
// @project GeniusRabbit::corelib 2016 – 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2018
//

package models

import (
	"time"

	"github.com/geniusrabbit/gosql"
	"github.com/guregu/null"

	"geniusrabbit.dev/corelib/admodels/types"
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
	Bids gosql.NullableJSON `gorm:"type:JSONB" json:"bids,omitempty"`

	// Status of the approvements
	Status types.ApproveStatus `gorm:"type:INT" json:"status,omitempty"`

	// Is Active advertisement
	Active types.ActiveStatus `gorm:"type:INT" json:"active,omitempty"`

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
	BidPrice        billing.Money `gorm:"type:NUMERIC(10,10)" json:"bid_price,omitempty"`         // Maximal bid for RTB
	Price           billing.Money `gorm:"type:NUMERIC(10,10)" json:"price,omitempty"`             // Price per pricing_model
	LeadPrice       billing.Money `gorm:"type:NUMERIC(10,10)" json:"lead_bid,omitempty"`          // Price of lead to calculate effectivity
	DailyBudget     billing.Money `gorm:"type:NUMERIC(10,10)" json:"daily_budget,omitempty"`      // Max daily budget spent
	Budget          billing.Money `gorm:"type:NUMERIC(10,10)" json:"budget,omitempty"`            // Money budget for whole time
	DailyTestBudget billing.Money `gorm:"type:NUMERIC(10,10)" json:"daily_test_budget,omitempty"` // Test money amount a day (it stops automaticaly if not profit for this amount)
	TestBudget      billing.Money `gorm:"type:NUMERIC(10,10)" json:"test_budget,omitempty"`       // Total test budget of whole time

	// Context contains the most improtant information about Ad
	Context gosql.NullableJSON `gorm:"type:JSONB" json:"context,omitempty"`

	// Weight of the Ad in rotation
	Weight int `json:"weight,omitempty"`

	// Frequency Capping of advertisement display to one user
	FrequencyCapping uint `json:"frequency_capping,omitempty"`

	// Hours targetting 168 simbols. Every simbol means hour active or blocked
	// 7 lines [day of week] + 24 hours as '1' or '0'
	// '*' or empty - all is on
	Hours string `json:"hours,omitempty"`

	// Assets related to advertisement
	Assets []*AdFile `json:"assets,omitempty" gorm:"many2many:m2m_adv_ad_file_ad;association_autoupdate:false"`

	// Time marks
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

// TableName in database
func (a *Ad) TableName() string {
	return "adv_ad"
}
