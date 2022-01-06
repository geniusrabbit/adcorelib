//
// @project geniusrabbit::corelib 2016 – 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017, 2019
//

package models

import (
	"strconv"
	"time"

	"geniusrabbit.dev/corelib/admodels/types"
	"github.com/geniusrabbit/gosql"
	"github.com/geniusrabbit/gosql/pgtype"
	"github.com/guregu/null"
)

// RTB price type
const (
	RTBPricePerMille = iota
	RTBPricePerOne
)

// RTBSource for SSP connect
type RTBSource struct {
	ID        uint64   `json:"id"`
	Company   *Company `json:"company,omitempty"`
	CompanyID uint64   `json:"company_id,omitempty"`
	Title     string   `json:"title,omitempty"`

	Status types.ApproveStatus `json:"status,omitempty"`
	Active types.ActiveStatus  `json:"active,omitempty"`
	Flags  pgtype.Hstore       `json:"flags,omitempty"`

	Protocol      string               `json:"protocol"`          // rtb as default
	MinimalWeight float64              `json:"minimal_weight"`    //
	URL           string               `json:"url"`               // RTB client request URL
	Method        string               `json:"method"`            // HTTP method GET, POST, ect; Default POST
	RequestType   types.RTBRequestType `json:"request_type"`      // 1 - json, 2 - xml, 3 - ProtoBUFF, 4 - PLAINTEXT
	Headers       pgtype.Hstore        `json:"headers,omitempty"` //
	RPS           int                  `json:"rps"`               // 0 – unlimit
	Timeout       int                  `json:"timeout"`           // In milliseconds

	// Money configs
	Accuracy              float64           `json:"accuracy,omitempty"`                // Price accuracy for auction in percentages
	PriceCorrectionReduce float64           `json:"price_correction_reduce,omitempty"` // % 100_00, 10000 -> 100%, 6550 -> 65.5%
	AuctionType           types.AuctionType `json:"auction_type,omitempty"`            // default: 0 – first price type, 1 – second price type

	// Targeting filters
	Formats         gosql.StringArray             `json:"formats,omitempty"`                       // => Filters
	DeviceTypes     gosql.NullableOrderedIntArray `json:"device_types,omitempty"`                  //
	Devices         gosql.NullableOrderedIntArray `json:"devices,omitempty"`                       //
	OS              gosql.NullableOrderedIntArray `json:"os,omitempty"`                            //
	Browsers        gosql.NullableOrderedIntArray `json:"browsers,omitempty"`                      //
	Carriers        gosql.NullableOrderedIntArray `json:"carriers,omitempty"`                      //
	Categories      gosql.NullableOrderedIntArray `json:"categories,omitempty"`                    //
	Countries       gosql.StringArray             `json:"countries,omitempty"`                     //
	Languages       gosql.StringArray             `json:"languages,omitempty"`                     //
	Applications    gosql.NullableOrderedIntArray `json:"apps,omitempty" gorm:"column:apps"`       //
	Domains         gosql.StringArray             `json:"domains,omitempty"`                       //
	Zones           gosql.NullableOrderedIntArray `json:"zones,omitempty"`                         //
	ExternalZones   gosql.NullableOrderedIntArray `json:"external_zones,omitempty"`                //
	Secure          int                           `json:"secure,omitempty"`                        // 0 - any, 1 - only, 2 - exclude
	AdBlock         int                           `json:"adblock,omitempty" gorm:"column:adblock"` // 0 - any, 1 - only, 2 - exclude
	PrivateBrowsing int                           `json:"private_browsing,omitempty"`              // 0 - any, 1 - only, 2 - exclude
	IP              int                           `json:"ip,omitempty"`                            // 0 - any, 1 - IPv4, 2 - IPv6

	Config gosql.NullableJSON `json:"config,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"-"`
}

// TableName in database
func (c *RTBSource) TableName() string {
	return "rtb_source"
}

// ProtocolCode name
func (c *RTBSource) ProtocolCode() string {
	if len(c.Protocol) < 1 {
		c.Protocol = "rtb"
	}
	return c.Protocol
}

// Flag get by key
func (c *RTBSource) Flag(flagName string) int {
	if val, ok := c.Flags.Get(flagName); ok {
		i, _ := strconv.Atoi(val)
		return i
	}
	return -1
}

// SetFlag for object
func (c *RTBSource) SetFlag(flagName string, flagValue int) {
	c.Flags.Set(flagName, strconv.Itoa(flagValue))
}

// PriceCorrectionReduceFactor which is a potential
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (c *RTBSource) PriceCorrectionReduceFactor() float64 {
	return c.PriceCorrectionReduce / 100.
}
