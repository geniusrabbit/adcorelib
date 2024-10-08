//
// @project GeniusRabbit corelib 2016 – 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017
//

package models

import (
	"time"

	"github.com/geniusrabbit/gosql/v2"
	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

type RTBAccessPointFlags struct {
	Trace        int8 `json:"trace,omitempty"`
	ErrorsIgnore int8 `json:"errors_ignore,omitempty"`
}

// RTBAccessPoint for DSP connect.
// It means that this is entry point which contains
// information for access and search data
type RTBAccessPoint struct {
	ID        uint64 `json:"id" gorm:"primaryKey"`
	AccountID uint64 `json:"account_id"`
	Codename  string `json:"codename,omitempty" gorm:"unique_index"`

	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	// RevenueShareReduce represents extra reduce factor to nevilate AdExchange and SSP discrepancy.
	// It means that the final bid respose will be decresed by RevenueShareReduce %
	// Example:
	//   1. Found advertisement with `bid=1.0$`
	//   2. Final `bid = bid - $AdSourceComission{%} - $AdExchangeComission{%} - $RevenueShareReduce{%}`
	RevenueShareReduce float64           `json:"revenue_share_reduce,omitempty"`                 // % 100_00, 10000 -> 100%, 6550 -> 65.5%
	AuctionType        types.AuctionType `gorm:"type:AuctionType" json:"auction_type,omitempty"` // default: 1 – first price type, 2 – second price type

	Status types.ApproveStatus                     `gorm:"type:ApproveStatus" json:"status,omitempty"`
	Active types.ActiveStatus                      `gorm:"type:ActiveStatus" json:"active,omitempty"`
	Flags  gosql.NullableJSON[RTBAccessPointFlags] `gorm:"type:JSONB" json:"flags,omitempty"`

	// Protocol configs
	Protocol      string                                `json:"protocol,omitempty"`
	Timeout       int                                   `json:"timeout,omitempty"`
	RPS           int                                   `json:"rps,omitempty"`
	DomainDefault string                                `json:"domain_default,omitempty"`
	RequestType   types.RTBRequestType                  `gorm:"type:RTBRequestType" json:"request_type"` // 1 - json, 2 - xml, 3 - ProtoBUFF, 4 - PLAINTEXT
	Headers       gosql.NullableJSON[map[string]string] `gorm:"type:JSONB" json:"headers,omitempty"`

	// Price limits
	MaxBid             float64 `json:"max_bid,omitempty"`
	FixedPurchasePrice float64 `json:"fixed_purchase_price,omitempty"`

	// Targeting filters
	Formats         gosql.NullableStringArray               `gorm:"type:TEXT[]" json:"formats,omitempty"`                       // => Filters [direct,banner_250x300]
	DeviceTypes     gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"device_types,omitempty"`                   //
	Devices         gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"devices,omitempty"`                        //
	OS              gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"os,omitempty"`                             //
	Browsers        gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"browsers,omitempty"`                       //
	Categories      gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"categories,omitempty"`                     //
	Carriers        gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"carriers,omitempty"`                       //
	Countries       gosql.NullableStringArray               `gorm:"type:TEXT[]" json:"countries,omitempty"`                     //
	Languages       gosql.NullableStringArray               `gorm:"type:TEXT[]" json:"languages,omitempty"`                     //
	Applications    gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[];column:apps" json:"apps,omitempty"`               //
	Zones           gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[];column:zones" json:"zones,omitempty"`             //
	Domains         gosql.NullableStringArray               `gorm:"type:TEXT[]" json:"domains,omitempty"`                       //
	Sources         gosql.NullableOrderedNumberArray[int64] `gorm:"column:rtb_sources;type:INT[]" json:"rtb_sources,omitempty"` //
	Secure          int                                     `gorm:"notnull" json:"secure,omitempty"`                            // 0 - any, 1 - only, 2 - exclude
	AdBlock         int                                     `gorm:"column:adblock;notnull" json:"adblock,omitempty"`            // 0 - any, 1 - only, 2 - exclude
	PrivateBrowsing int                                     `gorm:"notnull" json:"private_browsing,omitempty"`                  // 0 - any, 1 - only, 2 - exclude
	IP              int                                     `gorm:"notnull" json:"ip,omitempty"`                                // 0 - any, 1 - IPv4, 2 - IPv6

	// Time marks
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName in database
func (s *RTBAccessPoint) TableName() string {
	return "rtb_access_point"
}

// RBACResourceName returns the name of the resource for the RBAC
func (c *RTBAccessPoint) RBACResourceName() string {
	return "rtb_access_point"
}
