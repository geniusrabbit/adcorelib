//
// @project geniusrabbit::corelib 2016 – 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017
//

package models

import (
	"strconv"
	"time"

	"github.com/geniusrabbit/gosql"
	"github.com/geniusrabbit/gosql/pgtype"
	"github.com/guregu/null"
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

	RevenueShareReduce float64 `json:"revenue_share_reduce,omitempty"` // % 100_00, 10000 -> 100%, 6550 -> 65.5%
	AuctionType        int     `json:"auction_type,omitempty"`         // default: 0 – first price type, 1 – second price type

	Status int           `json:"status,omitempty"`
	Active int           `json:"active,omitempty"`
	Flags  pgtype.Hstore `json:"flags,omitempty"`

	// Money configs
	Protocol      string        `json:"protocol,omitempty"`
	Timeout       int           `json:"timeout,omitempty"`
	RPS           int           `json:"rps,omitempty"`
	DomainDefault string        `json:"domain_default,omitempty"`
	Headers       pgtype.Hstore `json:"headers,omitempty"`

	// Targeting filters
	Formats         gosql.StringArray             `json:"formats,omitempty"`                       // => Filters [pop,250x300]
	DeviceTypes     gosql.NullableOrderedIntArray `json:"device_types,omitempty"`                  //
	Devices         gosql.NullableOrderedIntArray `json:"devices,omitempty"`                       //
	OS              gosql.NullableOrderedIntArray `json:"os,omitempty"`                            //
	Browsers        gosql.NullableOrderedIntArray `json:"browsers,omitempty"`                      //
	Categories      gosql.NullableOrderedIntArray `json:"categories,omitempty"`                    //
	Carriers        gosql.NullableOrderedIntArray `json:"carriers,omitempty"`                      //
	Countries       gosql.StringArray             `json:"countries,omitempty"`                     //
	Languages       gosql.StringArray             `json:"languages,omitempty"`                     //
	Applications    gosql.NullableOrderedIntArray `json:"apps,omitempty" gorm:"column:apps"`       //
	Zones           gosql.NullableOrderedIntArray `json:"zones,omitempty" gorm:"column:zones"`     //
	Domains         gosql.StringArray             `json:"domains,omitempty"`                       //
	Sources         gosql.NullableOrderedIntArray `json:"rtb_sources,omitempty"`                   //
	Secure          int                           `json:"secure,omitempty"`                        // 0 - any, 1 - only, 2 - exclude
	AdBlock         int                           `json:"adblock,omitempty" gorm:"column:adblock"` // 0 - any, 1 - only, 2 - exclude
	PrivateBrowsing int                           `json:"private_browsing,omitempty"`              // 0 - any, 1 - only, 2 - exclude
	IP              int                           `json:"ip,omitempty"`                            // 0 - any, 1 - IPv4, 2 - IPv6

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
	if val, ok := s.Flags.Get(flagName); ok {
		i, _ := strconv.Atoi(val)
		return i
	}
	return -1
}

// SetFlag for object
func (s *RTBAccessPoint) SetFlag(flagName string, flagValue int) {
	s.Flags.Set(flagName, strconv.Itoa(flagValue))
}
