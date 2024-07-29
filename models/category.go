package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// Category model description
type Category struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	IABCode string `json:"iab_code"` // IAB category code of OpenRTB

	ParentID uint64    `json:"parent_id"`
	Parent   *Category `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Position uint64    `json:"position"`

	// Is Active advertisement
	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active"`

	// Time marks
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName in database
func (c *Category) TableName() string {
	return "adv_category"
}

// RBACResourceName returns the name of the resource for the RBAC
func (c *Category) RBACResourceName() string {
	return "adv_category"
}
