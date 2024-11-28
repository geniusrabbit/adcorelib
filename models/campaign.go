//
// @project GeniusRabbit corelib 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package models

import (
	"time"

	"github.com/geniusrabbit/gosql/v2"
	"github.com/guregu/null"
	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// Campaign model
type Campaign struct {
	// ID number of the Advertisement in DB
	ID    uint64 `json:"id"`
	Title string `json:"title"`

	// Owner/moderator Company of the Campaign
	AccountID   uint64 `json:"account_id"`
	CreatorID   uint64 `json:"creator_id"`
	ModeratorID uint64 `json:"moderator_id"`

	// Status of the campaign
	Status types.ApproveStatus `gorm:"type:ApproveStatus" json:"status"`

	// Is Active campaign
	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active"`

	// Is private campaign type
	Private types.PrivateStatus `gorm:"type:PrivateStatus" json:"private"`

	// Money limit counters
	DailyBudget     float64 `json:"daily_budget,omitempty"`      // Max daily budget spent
	Budget          float64 `json:"budget,omitempty"`            // Money budget for whole time
	DailyTestBudget float64 `json:"daily_test_budget,omitempty"` // Test money amount a day (it stops automaticaly if not profit for this amount)
	TestBudget      float64 `json:"test_budget,omitempty"`       // Total test budget of whole time

	Context gosql.NullableJSON[map[string]any] `gorm:"type:JSONB" json:"context,omitempty"`

	// Targeting scope incofrmation
	Zones       gosql.NullableNumberArray[uint64] `gorm:"type:BIGINT[]" json:"zones,omitempty"`
	Domains     gosql.NullableStringArray         `gorm:"type:TEXT[]" json:"domains,omitempty"` // site domains or application bundels
	Categories  gosql.NullableNumberArray[uint64] `gorm:"type:BIGINT[]" json:"categories,omitempty"`
	Geos        gosql.NullableStringArray         `gorm:"type:TEXT[]" json:"geos,omitempty"`
	Languages   gosql.NullableStringArray         `gorm:"type:TEXT[]" json:"languages,omitempty"`
	Browsers    gosql.NullableNumberArray[uint64] `gorm:"type:BIGINT[]" json:"browsers,omitempty"`
	OS          gosql.NullableNumberArray[uint64] `gorm:"type:BIGINT[]" json:"os,omitempty"`
	DeviceTypes gosql.NullableNumberArray[uint64] `gorm:"type:BIGINT[]" json:"device_types,omitempty"`
	Devices     gosql.NullableNumberArray[uint64] `gorm:"type:BIGINT[]" json:"devices,omitempty"`
	DateStart   null.Time                         `json:"date_start,omitempty"`
	DateEnd     null.Time                         `json:"date_end,omitempty"`
	Hours       null.String                       `json:"hours,omitempty"`
	Sex         gosql.NullableNumberArray[uint]   `gorm:"type:INT[]" json:"sex,omitempty"`
	Age         gosql.NullableNumberArray[uint]   `gorm:"type:INT[]" json:"age,omitempty"`

	// Advertisement list
	Ads   []*Ad     `json:"ads,omitempty" gorm:"ForeignKey:CampaignID"`
	Links []*AdLink `json:"links,omitempty" gorm:"ForeignKey:CampaignID"`

	Trace        gosql.NullableStringArray `gorm:"type:TEXT[]" json:"trace,omitempty"`
	TracePercent int                       `json:"trace_percent,omitempty"`

	// Time marks
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName in database
func (c *Campaign) TableName() string {
	return "adv_campaign"
}

// RBACResourceName returns the name of the resource for the RBAC
func (c *Campaign) RBACResourceName() string {
	return "adv_campaign"
}
