package models

import (
	"time"

	"github.com/geniusrabbit/gosql/v2"
	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

type BrowserVersion struct {
	Min  types.Version `json:"min"`
	Max  types.Version `json:"max"`
	Name string        `json:"name"`
}

// Browser model description
type Browser struct {
	ID uint64 `json:"id"`

	Name   string             `json:"name"`
	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	Versions gosql.NullableJSONArray[BrowserVersion] `gorm:"type:JSONB" json:"versions,omitempty"`

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
