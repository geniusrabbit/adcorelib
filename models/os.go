package models

import (
	"time"

	"github.com/geniusrabbit/gosql/v2"
	"github.com/guregu/null"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// OS model description
type OS struct {
	ID uint64 `json:"id"`

	Name   string             `json:"name"`
	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	Versions gosql.NullableStringArray `gorm:"type:TEXT[]" json:"versions,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at"`
}

// TableName in database
func (m *OS) TableName() string {
	return `type_os`
}
