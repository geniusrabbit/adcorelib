//
// @project GeniusRabbit corelib 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package models

import (
	"time"

	"github.com/geniusrabbit/gosql/v2"
	"github.com/guregu/null"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

// Zone model
type Zone struct {
	ID        uint64 `json:"id"`
	Title     string `json:"title"`
	AccountID uint64 `json:"account_id,omitempty"`

	Type   types.ZoneType      `gorm:"type:ZoneType" json:"type,omitempty"`
	Status types.ApproveStatus `gorm:"type:ApproveStatus" json:"status"`
	Active types.ActiveStatus  `gorm:"type:ActiveStatus" json:"active"`

	DefaultCode  gosql.NullableJSON[map[string]string]  `gorm:"type:JSONB" json:"default_code,omitempty"`
	Context      gosql.NullableJSON[map[string]any]     `gorm:"type:JSONB" json:"context,omitempty"`         //
	MinECPM      float64                                `json:"min_ecpm,omitempty"`                          // Default
	MinECPMByGeo gosql.NullableJSON[map[string]float64] `gorm:"type:JSONB" json:"min_ecpm_by_geo,omitempty"` // {"CODE": <ecpm>, ...}

	// The cost of the traffic acceptance
	FixedPurchasePrice billing.Money `json:"fixed_purchase_price,omitempty"`

	// Filtering
	AllowedFormats    gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"allowed_formats,omitempty"`    //
	AllowedTypes      gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"allowed_types,omitempty"`      //
	AllowedSources    gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"allowed_sources,omitempty"`    //
	DisallowedSources gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"disallowed_sources,omitempty"` //
	Campaigns         gosql.NullableOrderedNumberArray[int64] `gorm:"type:INT[]" json:"campaigns,omitempty"`          // Strict campaigns targeting (smartlinks only)

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at,omitempty"`
}

// TableName in database
func (z *Zone) TableName() string {
	return "adv_zone"
}

// RBACResourceName returns the name of the resource for the RBAC
func (z *Zone) RBACResourceName() string {
	return "adv_zone"
}

// RevenueShare amount %
func (z *Zone) RevenueShare() float64 {
	return 0
}

func (z *Zone) GetAccountID() uint64 {
	return z.AccountID
}
