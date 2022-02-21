package models

import (
	"time"

	"github.com/geniusrabbit/gosql"
	"github.com/guregu/null"

	"geniusrabbit.dev/corelib/admodels/types"
)

// Browser model description
type Browser struct {
	ID uint64 `json:"id"`

	Name   string             `json:"name"`
	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	Versions gosql.NullableStringArray `gorm:"type:TEXT[]" json:"versions,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

// TableName in database
func (m *Browser) TableName() string {
	return `type_browser`
}
