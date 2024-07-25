//
// @project GeniusRabbit corelib 2016 – 2017, 2019, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017, 2019, 2024
//

package models

import (
	"time"

	"github.com/geniusrabbit/gosql/v2"
	"github.com/guregu/null"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

// RTB price type
const (
	RTBPricePerMille = iota
	RTBPricePerOne
)

type RTBSourceFlags struct {
	Trace        int8 `json:"trace,omitempty"`
	ErrorsIgnore int8 `json:"errors_ignore,omitempty"`
}

// RTBSource for SSP connect
type RTBSource struct {
	ID        uint64 `json:"id"`
	AccountID uint64 `json:"account_id,omitempty"`

	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`

	Status types.ApproveStatus                `gorm:"type:ApproveStatus" json:"status,omitempty"`
	Active types.ActiveStatus                 `gorm:"type:ActiveStatus" json:"active,omitempty"`
	Flags  gosql.NullableJSON[RTBSourceFlags] `gorm:"type:JSONB" json:"flags,omitempty"`

	Protocol      string                                `json:"protocol"`                                // rtb as default
	MinimalWeight float64                               `json:"minimal_weight"`                          //
	URL           string                                `json:"url"`                                     // RTB client request URL
	Method        string                                `json:"method"`                                  // HTTP method GET, POST, ect; Default POST
	RequestType   types.RTBRequestType                  `gorm:"type:RTBRequestType" json:"request_type"` // 1 - json, 2 - xml, 3 - ProtoBUFF, 4 - PLAINTEXT
	Headers       gosql.NullableJSON[map[string]string] `gorm:"type:JSONB" json:"headers,omitempty"`     //
	RPS           int                                   `json:"rps"`                                     // 0 – unlimit
	Timeout       int                                   `json:"timeout"`                                 // In milliseconds

	// Money configs
	Accuracy              float64           `json:"accuracy,omitempty"`                             // Price accuracy for auction in percentages
	PriceCorrectionReduce float64           `json:"price_correction_reduce,omitempty"`              // % 100_00, 10000 -> 100%, 6550 -> 65.5%
	AuctionType           types.AuctionType `gorm:"type:AuctionType" json:"auction_type,omitempty"` // default: 0 – first price type, 1 – second price type

	// Price limits
	MinBid billing.Money `json:"min_bid,omitempty"` // Minimal bid value
	MaxBid billing.Money `json:"max_bid,omitempty"` // Maximal bid value

	// Targeting filters
	Formats         gosql.NullableStringArray               `gorm:"type:TEXT[]" json:"formats,omitempty"`         // => Filters
	DeviceTypes     gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"device_types,omitempty"`     //
	Devices         gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"devices,omitempty"`          //
	OS              gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"os,omitempty"`               //
	Browsers        gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"browsers,omitempty"`         //
	Carriers        gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"carriers,omitempty"`         //
	Categories      gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"categories,omitempty"`       //
	Countries       gosql.NullableStringArray               `gorm:"type:TEXT[]" json:"countries,omitempty"`       //
	Languages       gosql.NullableStringArray               `gorm:"type:TEXT[]" json:"languages,omitempty"`       //
	Applications    gosql.NullableOrderedNumberArray[int64] `gorm:"column:apps;type:INT[]" json:"apps,omitempty"` //
	Domains         gosql.NullableStringArray               `gorm:"type:TEXT[]" json:"domains,omitempty"`         //
	Zones           gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"zones,omitempty"`            //
	ExternalZones   gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"external_zones,omitempty"`   //
	Secure          int                                     `json:"secure,omitempty"`                             // 0 - any, 1 - only, 2 - exclude
	AdBlock         int                                     `json:"adblock,omitempty" gorm:"column:adblock"`      // 0 - any, 1 - only, 2 - exclude
	PrivateBrowsing int                                     `json:"private_browsing,omitempty"`                   // 0 - any, 1 - only, 2 - exclude
	IP              int                                     `json:"ip,omitempty"`                                 // 0 - any, 1 - IPv4, 2 - IPv6

	Config gosql.NullableJSON[any] `gorm:"type:JSONB" json:"config,omitempty"`

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

// PriceCorrectionReduceFactor which is a potential
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (c *RTBSource) PriceCorrectionReduceFactor() float64 {
	return c.PriceCorrectionReduce / 100.
}

// RBACResourceName returns the name of the resource for the RBAC
func (c *RTBSource) RBACResourceName() string {
	return "rtb_source"
}
