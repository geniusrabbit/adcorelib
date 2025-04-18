package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/gosql/v2"
)

// TrafficRouter model represents the traffic router for the traffic distribution
// between the RTB sources from the Applications and Zones
type TrafficRouter struct {
	ID        uint64 `json:"id"`
	AccountID uint64 `json:"account_id"`

	Title   string  `json:"title"`
	Percent float64 `json:"percent"`

	Description string              `json:"description"`
	Status      types.ApproveStatus `gorm:"type:ApproveStatus" json:"status,omitempty"`
	Active      types.ActiveStatus  `json:"active"`

	// Target RTB`s and sources of the Advertisement
	RTBSourceIDs gosql.NumberArray[uint64] `gorm:"type:BIGINT[]" json:"rtb_source_ids"`

	// Targeting filters
	Formats         gosql.NullableStringArray                `gorm:"type:TEXT[]" json:"formats,omitempty"`            // => Filters
	DeviceTypes     gosql.NullableOrderedNumberArray[uint64] `gorm:"type:BIGINT[]" json:"device_types,omitempty"`     //
	Devices         gosql.NullableOrderedNumberArray[uint64] `gorm:"type:BIGINT[]" json:"devices,omitempty"`          //
	OS              gosql.NullableOrderedNumberArray[uint64] `gorm:"type:BIGINT[]" json:"os,omitempty"`               //
	Browsers        gosql.NullableOrderedNumberArray[uint64] `gorm:"type:BIGINT[]" json:"browsers,omitempty"`         //
	Carriers        gosql.NullableOrderedNumberArray[uint64] `gorm:"type:BIGINT[]" json:"carriers,omitempty"`         //
	Categories      gosql.NullableOrderedNumberArray[uint64] `gorm:"type:BIGINT[]" json:"categories,omitempty"`       //
	Countries       gosql.NullableStringArray                `gorm:"type:TEXT[]" json:"countries,omitempty"`          //
	Languages       gosql.NullableStringArray                `gorm:"type:TEXT[]" json:"languages,omitempty"`          //
	Domains         gosql.NullableStringArray                `gorm:"type:TEXT[]" json:"domains,omitempty"`            //
	Applications    gosql.NullableOrderedNumberArray[uint64] `gorm:"column:apps;type:BIGINT[]" json:"apps,omitempty"` //
	Zones           gosql.NullableOrderedNumberArray[uint64] `gorm:"type:BIGINT[]" json:"zones,omitempty"`            //
	Secure          int                                      `json:"secure,omitempty"`                                // 0 - any, 1 - only, 2 - exclude
	AdBlock         int                                      `json:"adblock,omitempty" gorm:"column:adblock"`         // 0 - any, 1 - only, 2 - exclude
	PrivateBrowsing int                                      `json:"private_browsing,omitempty"`                      // 0 - any, 1 - only, 2 - exclude
	IP              int                                      `json:"ip,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// TableName in database
func (tr *TrafficRouter) TableName() string {
	return "adv_traffic_router"
}

// RBACResourceName returns the name of the resource for the RBAC
func (tr *TrafficRouter) RBACResourceName() string {
	return "traffic_router"
}
