package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// OS model description
type OS struct {
	ID uint64 `json:"id" gorm:"primaryKey"`

	Name        string        `json:"name"`
	Version     types.Version `json:"version,omitempty"`
	Description string        `json:"description,omitempty"`

	YearRelease    int `json:"year_release,omitempty"`
	YearEndSupport int `json:"year_end_support,omitempty"`

	// Match expressions
	MatchNameExp       string `json:"match_name_exp,omitempty"`
	MatchUserAgentExp  string `json:"match_ua_exp,omitempty"`
	MatchVersionMinExp string `json:"match_ver_min_exp,omitempty"`
	MatchVersionMaxExp string `json:"match_ver_max_exp,omitempty"`

	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	ParentID sql.Null[uint64] `json:"parent_id,omitempty"`
	Parent   *OS              `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID"`
	Versions []*OS            `json:"versions,omitempty" gorm:"foreignKey:ParentID;references:ID"`

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
