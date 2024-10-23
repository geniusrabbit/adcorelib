package models

import (
	"time"

	"github.com/geniusrabbit/gosql/v2"
	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// Application model describes site or mobile/desktop application
type Application struct {
	ID uint64 `json:"id" gorm:"primaryKey"` // Authoincrement key

	Title       string `json:"title"`
	Description string `json:"description"`

	AccountID uint64 `json:"account_id"`
	CreatorID uint64 `json:"creator_id"`

	// Unical application identificator like:
	//   - site domain -> domain.com
	//   - mobile/desktop application bundle -> com.application.game
	URI      string                `json:"uri"`
	Type     types.ApplicationType `gorm:"type:ApplicationType" json:"type"`
	Platform types.PlatformType    `gorm:"type:PlatformType" json:"platform"`
	Premium  bool                  `json:"premium"`

	// Status of the application
	Status types.ApproveStatus `gorm:"type:ApproveStatus" json:"status"`

	// Is Active application
	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active"`

	// Is private campaign type
	Private types.PrivateStatus `gorm:"type:PrivateStatus;not null" json:"private"`

	Categories gosql.NullableNumberArray[uint] `gorm:"type:BIGINT[]" json:"categories,omitempty"`

	// RevenueShare it's amount of percent of the raw incode which will be shared with the publisher company
	// For example:
	//   Displayed ads for 100$
	//   Company revenue share 60%
	//   In such case the ad network have 40$
	//   The publisher have 60$
	// Optional
	RevenueShare float64 `json:"revenue_share,omitempty"` // % 1.0 -> 100%, 0.655 -> 65.5%

	// Advertisement sources
	AllowedSources    gosql.NullableOrderedNumberArray[int64] `gorm:"type:BIGINT[]" json:"allowed_sources,omitempty"`
	DisallowedSources gosql.NullableOrderedNumberArray[int64] `gorm:"type:BIGINT[]" json:"disallowed_sources,omitempty"`

	// Time marks
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName in database
func (app *Application) TableName() string {
	return "adv_application"
}

// RBACResourceName returns the name of the resource for the RBAC
func (app *Application) RBACResourceName() string {
	return "adv_application"
}
