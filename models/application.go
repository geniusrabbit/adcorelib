package models

import (
	"time"

	"github.com/geniusrabbit/gosql/v2"
	"github.com/guregu/null"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// Application model describes site or mobile/desktop application
type Application struct {
	ID        uint64 `json:"id"` // Authoincrement key
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
	Private types.PrivateStatus `gorm:"type:PrivateStatus" json:"private"`

	Categories gosql.NullableNumberArray[uint] `gorm:"type:INT[]" json:"categories,omitempty"`

	// RevenueShare it's amount of percent of the raw incode which will be shared with the publisher company
	// For example:
	//   Displayed ads for 100$
	//   Company revenue share 60%
	//   In such case the ad network have 40$
	//   The publisher have 60$
	// Optional
	RevenueShare float64 `json:"revenue_share,omitempty"` // % 100_00, 10000 -> 100%, 6550 -> 65.5%

	// Time marks
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at,omitempty"`
}

// TableName in database
func (app *Application) TableName() string {
	return "adv_application"
}

// RBACResourceName returns the name of the resource for the RBAC
func (app *Application) RBACResourceName() string {
	return "adv_application"
}
