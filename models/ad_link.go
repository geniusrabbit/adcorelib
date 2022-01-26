package models

import (
	"time"

	"github.com/guregu/null"

	"geniusrabbit.dev/corelib/admodels/types"
)

// AdLink to the advertisement target
type AdLink struct {
	// ID number of the Advertisement Link in DB
	ID uint64 `json:"id"`

	// Link to the target
	Link string `json:"link"`

	// Target campaign
	CampaignID uint64 `json:"campaign_id"`

	// Status of the approvements
	Status types.ApproveStatus `gorm:"type:INT" json:"status,omitempty"`

	// Is Active advertisement
	Active types.ActiveStatus `gorm:"type:INT" json:"active,omitempty"`

	// Time marks
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

// TableName in database
func (link *AdLink) TableName() string {
	return "adv_link"
}
