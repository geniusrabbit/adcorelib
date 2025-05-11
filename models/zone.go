//
// @project GeniusRabbit corelib 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package models

import (
	"time"

	"github.com/geniusrabbit/gosql/v2"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// Zone model
type Zone struct {
	ID       uint64 `json:"id"`
	Codename string `json:"codename" gorm:"unique_index"`

	Title       string `json:"title"`
	Description string `json:"description"`

	AccountID uint64 `json:"account_id,omitempty"`

	Type   types.ZoneType      `gorm:"type:ZoneType" json:"type,omitempty"`
	Status types.ApproveStatus `gorm:"type:ApproveStatus" json:"status"`
	Active types.ActiveStatus  `gorm:"type:ActiveStatus" json:"active"`

	DefaultCode  gosql.NullableJSON[map[string]string]  `gorm:"type:JSONB" json:"default_code,omitempty"`
	Context      gosql.NullableJSON[map[string]any]     `gorm:"type:JSONB" json:"context,omitempty"`         //
	MinECPM      float64                                `json:"min_ecpm,omitempty"`                          // Default
	MinECPMByGeo gosql.NullableJSON[map[string]float64] `gorm:"type:JSONB" json:"min_ecpm_by_geo,omitempty"` // {"CODE": <ecpm>, ...}

	// The cost of the traffic acceptance
	FixedPurchasePrice float64 `json:"fixed_purchase_price,omitempty"`

	// Filtering
	AllowedFormats    gosql.NullableStringArray                `gorm:"type:TEXT[]" json:"allowed_formats,omitempty"`
	AllowedTypes      gosql.NullableOrderedNumberArray[uint64] `gorm:"type:BIGINT[]" json:"allowed_types,omitempty"`
	AllowedSources    gosql.NullableOrderedNumberArray[uint64] `gorm:"type:BIGINT[]" json:"allowed_sources,omitempty"`
	DisallowedSources gosql.NullableOrderedNumberArray[uint64] `gorm:"type:BIGINT[]" json:"disallowed_sources,omitempty"`

	// Strict campaigns targeting (smartlinks only)
	Campaigns gosql.NullableOrderedNumberArray[uint64] `gorm:"type:BIGINT[]" json:"campaigns,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
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

// BeforeCreate hook
func (z *Zone) BeforeCreate(tx *gorm.DB) error {
	if z.Codename == "" {
		id, err := ksuid.NewRandom()
		if err != nil {
			return err
		}
		z.Codename = id.String()
	}
	z.CreatedAt = time.Now()
	z.UpdatedAt = z.CreatedAt
	return nil
}
