//
// @project GeniusRabbit::corelib 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package models

import (
	"fmt"
	"time"

	"github.com/geniusrabbit/gosql"
	"github.com/guregu/null"

	"geniusrabbit.dev/corelib/admodels/types"
)

// Zone model
type Zone struct {
	ID        uint64   `json:"id"`
	Title     string   `json:"title"`
	Company   *Company `json:"company,omitempty"`
	CompanyID uint64   `json:"company_id,omitempty"`

	Type   types.ZoneType      `gorm:"type:ZoneType" json:"type,omitempty"`
	Status types.ApproveStatus `gorm:"type:ApproveStatus" json:"status"`
	Active types.ActiveStatus  `gorm:"type:ActiveStatus" json:"active"`

	DefaultCode       gosql.NullableJSON            `gorm:"type:JSONB" json:"default_code,omitempty"`
	Context           gosql.NullableJSON            `gorm:"type:JSONB" json:"context,omitempty"`            //
	MinECPM           float64                       `json:"min_ecpm,omitempty"`                             // Default
	MinECPMByGeo      gosql.NullableJSON            `gorm:"type:JSONB" json:"min_ecpm_by_geo,omitempty"`    // {"CODE": <ecpm>, ...}
	Price             float64                       `json:"price,omitempty"`                                // The cost of single view
	AllowedFormats    gosql.NullableOrderedIntArray `gorm:"type:INT[]" json:"allowed_formats,omitempty"`    //
	AllowedTypes      gosql.NullableOrderedIntArray `gorm:"type:INT[]" json:"allowed_types,omitempty"`      //
	AllowedSources    gosql.NullableOrderedIntArray `gorm:"type:INT[]" json:"allowed_sources,omitempty"`    //
	DisallowedSources gosql.NullableOrderedIntArray `gorm:"type:INT[]" json:"disallowed_sources,omitempty"` //
	Campaigns         gosql.NullableOrderedIntArray `gorm:"type:INT[]" json:"campaigns,omitempty"`          // Strict campaigns targeting (smartlinks only)

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at,omitempty"`
}

// TableName in database
func (z *Zone) TableName() string {
	return "adv_zone"
}

// RevenueShare amount %
func (z *Zone) RevenueShare() float64 {
	return 0
}

// SetCompany object
func (z *Zone) SetCompany(c interface{}) error {
	switch v := c.(type) {
	case *Company:
		z.Company = v
		z.CompanyID = v.ID
	case uint64:
		z.Company = nil
		z.CompanyID = v
	default:
		return fmt.Errorf("undefined value type: %t", c)
	}
	return nil
}
