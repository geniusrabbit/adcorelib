package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

type DeviceModel struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	Codename    string `json:"codename" gorm:"unique"`
	Name        string `json:"name"`
	Description string `json:"description"`

	ParentID uint64       `json:"parent_id,omitempty"`
	Parent   *DeviceModel `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID"`

	YearRelease int `json:"year_release,omitempty"`

	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	MatchExp string `json:"match_exp,omitempty"`

	// Link to device maker
	MakerCodename string       `json:"maker_codename,omitempty"`
	Maker         *DeviceMaker `json:"maker,omitempty" gorm:"foreignKey:MakerCodename;references:Codename"`

	// Device type
	TypeCodename string      `json:"type_codename,omitempty"`
	Type         *DeviceType `json:"type,omitempty" gorm:"foreignKey:TypeCodename;references:Codename"`

	// Versions of the model
	Versions []*DeviceModel `json:"versions,omitempty" gorm:"foreignKey:ParentID;references:ID"`

	// Time marks
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName in database
func (m *DeviceModel) TableName() string {
	return `type_device_model`
}

// RBACResourceName returns the name of the resource for the RBAC
func (m *DeviceModel) RBACResourceName() string {
	return "device_model"
}
