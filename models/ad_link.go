package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// AdLink to the advertisement target
type AdLink struct {
	// ID number of the Advertisement Link in DB
	ID uint64 `json:"id" gorm:"primaryKey"`

	// Link to the target
	Link string `json:"link"`

	// Target campaign
	Campaign   *Campaign `json:"campaign,omitempty"`
	CampaignID uint64    `json:"campaign_id"`

	// Status of the approvements
	Status types.ApproveStatus `gorm:"type:ApproveStatus" json:"status,omitempty"`

	// Is Active advertisement
	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	// Time marks
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName in database
func (link *AdLink) TableName() string {
	return "adv_link"
}

// RBACResourceName returns the name of the resource for the RBAC
func (link *AdLink) RBACResourceName() string {
	return "adv_link"
}
