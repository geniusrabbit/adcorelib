package models

import (
	"time"

	"github.com/geniusrabbit/gosql/v2"
	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

type OSVersion struct {
	Min  types.Version `json:"min"`
	Max  types.Version `json:"max"`
	Name string        `json:"name,omitempty"`
}

// OS model description
type OS struct {
	ID uint64 `json:"id" gorm:"primaryKey"`

	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	MatchExp string `json:"match_exp,omitempty"`

	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	Versions gosql.NullableJSONArray[OSVersion] `gorm:"type:JSONB" json:"versions,omitempty"`

	// Time marks
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName in database
func (m *OS) TableName() string {
	return `type_os`
}

// RBACResourceName returns the name of the resource for the RBAC
func (m *OS) RBACResourceName() string {
	return "type_os"
}
