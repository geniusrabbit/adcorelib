//
// @project geniusrabbit::corelib 2016 – 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017
//

package models

import (
	"time"

	"github.com/demdxx/gocast"
	"github.com/geniusrabbit/gosql"
	"github.com/guregu/null"

	"geniusrabbit.dev/corelib/admodels/types"
	"geniusrabbit.dev/corelib/billing"
)

// RTBAccessPoint for DSP connect.
// It means that this is entry point which contains
// information for access and search data
type RTBAccessPoint struct {
	ID        uint64   `json:"id"`
	Company   *Company `json:"company,omitempty"`
	CompanyID uint64   `json:"company_id,omitempty"`
	Title     string   `json:"title,omitempty"`
	Codename  string   `json:"codename,omitempty"`

	// RevenueShareReduce represents extra reduce factor to nevilate AdExchange and SSP discrepancy.
	// It means that the final bid respose will be decresed by RevenueShareReduce %
	// Example:
	//   1. Found advertisement with `bid=1.0$`
	//   2. Final `bid = bid - $AdSourceComission{%} - $AdExchangeComission{%} - $RevenueShareReduce{%}`
	RevenueShareReduce float64           `json:"revenue_share_reduce,omitempty"`                 // % 100_00, 10000 -> 100%, 6550 -> 65.5%
	AuctionType        types.AuctionType `gorm:"type:AuctionType" json:"auction_type,omitempty"` // default: 0 – first price type, 1 – second price type

	Status types.ApproveStatus `gorm:"type:ApproveStatus" json:"status,omitempty"`
	Active types.ActiveStatus  `gorm:"type:ActiveStatus" json:"active,omitempty"`
	Flags  gosql.NullableJSON  `gorm:"type:JSONB" json:"flags,omitempty"`

	// Protocol configs
	Protocol      string             `json:"protocol,omitempty"`
	Timeout       int                `json:"timeout,omitempty"`
	RPS           int                `json:"rps,omitempty"`
	MaxBid        billing.Money      `json:"max_bid,omitempty"`
	DomainDefault string             `json:"domain_default,omitempty"`
	Headers       gosql.NullableJSON `gorm:"type:JSONB" json:"headers,omitempty"`

	// Targeting filters
	Formats         gosql.NullableStringArray     `gorm:"type:TEXT[]" json:"formats,omitempty"`            // => Filters [direct,banner_250x300]
	DeviceTypes     gosql.NullableOrderedIntArray `gorm:"type:INT[]" json:"device_types,omitempty"`        //
	Devices         gosql.NullableOrderedIntArray `gorm:"type:INT[]" json:"devices,omitempty"`             //
	OS              gosql.NullableOrderedIntArray `gorm:"type:INT[]" json:"os,omitempty"`                  //
	Browsers        gosql.NullableOrderedIntArray `gorm:"type:INT[]" json:"browsers,omitempty"`            //
	Categories      gosql.NullableOrderedIntArray `gorm:"type:INT[]" json:"categories,omitempty"`          //
	Carriers        gosql.NullableOrderedIntArray `gorm:"type:INT[]" json:"carriers,omitempty"`            //
	Countries       gosql.NullableStringArray     `gorm:"type:TEXT[]" json:"countries,omitempty"`          //
	Languages       gosql.NullableStringArray     `gorm:"type:TEXT[]" json:"languages,omitempty"`          //
	Applications    gosql.NullableOrderedIntArray `gorm:"type:INT[];column:apps" json:"apps,omitempty"`    //
	Zones           gosql.NullableOrderedIntArray `gorm:"type:INT[];column:zones" json:"zones,omitempty"`  //
	Domains         gosql.NullableStringArray     `gorm:"type:TEXT[]" json:"domains,omitempty"`            //
	Sources         gosql.NullableOrderedIntArray `gorm:"type:INT[]" json:"rtb_sources,omitempty"`         //
	Secure          int                           `gorm:"notnull" json:"secure,omitempty"`                 // 0 - any, 1 - only, 2 - exclude
	AdBlock         int                           `gorm:"column:adblock;notnull" json:"adblock,omitempty"` // 0 - any, 1 - only, 2 - exclude
	PrivateBrowsing int                           `gorm:"notnull" json:"private_browsing,omitempty"`       // 0 - any, 1 - only, 2 - exclude
	IP              int                           `gorm:"notnull" json:"ip,omitempty"`                     // 0 - any, 1 - IPv4, 2 - IPv6

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"-"`
}

// TableName in database
func (s *RTBAccessPoint) TableName() string {
	return "rtb_access_point"
}

// Flag get by key
func (s *RTBAccessPoint) Flag(flagName string) int {
	var m map[string]int
	if err := s.Flags.UnmarshalTo(&m); err == nil {
		return gocast.ToInt(m[flagName])
	}
	return -1
}

// SetFlag for object
func (s *RTBAccessPoint) SetFlag(flagName string, flagValue int) {
	var m map[string]int
	_ = s.Flags.UnmarshalTo(&m)
	if m == nil {
		m = map[string]int{}
	}
	m[flagName] = flagValue
	_ = s.Flags.SetValue(m)
}
