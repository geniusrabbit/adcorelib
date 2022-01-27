//
// @project GeniusRabbit::corelib 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package models

import (
	"fmt"
	"time"

	"github.com/bsm/openrtb"
	"github.com/geniusrabbit/gosql"
	"github.com/guregu/null"

	"geniusrabbit.dev/corelib/admodels/types"
)

// Device set of types
const (
	DeviceTypeUnknown   = openrtb.DeviceTypeUnknown
	DeviceTypeMobile    = openrtb.DeviceTypeMobile
	DeviceTypePC        = openrtb.DeviceTypePC
	DeviceTypeTV        = openrtb.DeviceTypeTV
	DeviceTypePhone     = openrtb.DeviceTypePhone
	DeviceTypeTablet    = openrtb.DeviceTypeTablet
	DeviceTypeConnected = openrtb.DeviceTypeConnected
	DeviceTypeSetTopBox = openrtb.DeviceTypeSetTopBox
)

// Campaign model
type Campaign struct {
	// ID number of the Advertisement in DB
	ID    uint64 `json:"id"`
	Title string `json:"title"`

	// Owner/moderator Company of the Campaign
	Company     *Company `json:"company,omitempty"` // Owner Project
	CompanyID   uint64   `json:"company_id"`
	Creator     *User    `json:"creator,omitempty"` // User who created the object
	CreatorID   uint64   `json:"creator_id"`
	Moderator   *User    `json:"moderator,omitempty"`
	ModeratorID uint64   `json:"moderator_id"`

	// Status of the campaign
	Status types.ApproveStatus `gorm:"type:ApproveStatus" json:"status"`

	// Is Active campaign
	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active"`

	// Is private campaign type
	Private types.PrivateStatus `gorm:"type:PrivateStatus" json:"private"`

	// Money limit counters
	DailyBudget     float64 `gorm:"type:NUMERIC(10,10)" json:"daily_budget,omitempty"`      // Max daily budget spent
	Budget          float64 `gorm:"type:NUMERIC(10,10)" json:"budget,omitempty"`            // Money budget for whole time
	DailyTestBudget float64 `gorm:"type:NUMERIC(10,10)" json:"daily_test_budget,omitempty"` // Test money amount a day (it stops automaticaly if not profit for this amount)
	TestBudget      float64 `gorm:"type:NUMERIC(10,10)" json:"test_budget,omitempty"`       // Total test budget of whole time

	Context gosql.NullableJSON `gorm:"type:JSONB" json:"context,omitempty"`

	// Targeting scope incofrmation
	Zones       gosql.NullableUintArray   `gorm:"type:INT[]" json:"zones,omitempty"`
	Domains     gosql.NullableStringArray `gorm:"type:TEXT[]" json:"domains,omitempty"` // site domains or application bundels
	Categories  gosql.NullableUintArray   `gorm:"type:INT[]" json:"categories,omitempty"`
	Geos        gosql.NullableStringArray `gorm:"type:TEXT[]" json:"geos,omitempty"`
	Languages   gosql.NullableStringArray `gorm:"type:TEXT[]" json:"languages,omitempty"`
	Browsers    gosql.NullableUintArray   `gorm:"type:INT[]" json:"browsers,omitempty"`
	Os          gosql.NullableUintArray   `gorm:"type:INT[]" json:"os,omitempty"`
	DeviceTypes gosql.NullableUintArray   `gorm:"type:INT[]" json:"device_types,omitempty"`
	Devices     gosql.NullableUintArray   `gorm:"type:INT[]" json:"devices,omitempty"`
	DateStart   null.Time                 `json:"date_start,omitempty"`
	DateEnd     null.Time                 `json:"date_end,omitempty"`
	Hours       null.String               `json:"hours,omitempty"`
	Sex         gosql.NullableUintArray   `gorm:"type:INT[]" json:"sex,omitempty"`
	Age         gosql.NullableUintArray   `gorm:"type:INT[]" json:"age,omitempty"`

	// Advertisement list
	Ads   []*Ad     `json:"ads,omitempty" gorm:"ForeignKey:CampaignID"`
	Links []*AdLink `json:"links,omitempty" gorm:"ForeignKey:CampaignID"`

	Trace        gosql.NullableStringArray `gorm:"type:TEXT[]" json:"trace,omitempty"`
	TracePercent int                       `json:"trace_percent,omitempty"`

	// Time marks
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at,omitempty"`
}

// TableName in database
func (c *Campaign) TableName() string {
	return "adv_campaign"
}

// SetCompany campaign owner
func (c *Campaign) SetCompany(com interface{}) error {
	switch v := com.(type) {
	case *Company:
		c.Company = v
		c.CompanyID = v.ID
	case uint64:
		if c.CompanyID != v {
			c.Company = nil
			c.CompanyID = v
		}
	default:
		return fmt.Errorf("[models.Campaign] undefined value type: %t", com)
	}
	return nil
}
