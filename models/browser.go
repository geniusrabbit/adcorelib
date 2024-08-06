package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/gosql/v2"
)

// BrowserVersion model
type BrowserVersion struct {
	Min  types.Version `json:"min"`
	Max  types.Version `json:"max"`
	Name string        `json:"name,omitempty"`
}

// Browser model description
type Browser struct {
	ID uint64 `json:"id"`

	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	MatchExp string `json:"match_exp,omitempty"`

	Versions gosql.NullableJSONArray[BrowserVersion] `json:"versions,omitempty" gorm:"type:jsonb"`

	// Time marks
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName in database
func (m *Browser) TableName() string {
	return `type_browser`
}

// RBACResourceName returns the name of the resource for the RBAC
func (m *Browser) RBACResourceName() string {
	return `type_browser`
}
