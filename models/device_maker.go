package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

type DeviceMaker struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	Codename    string `json:"codename" gorm:"unique"`
	Name        string `json:"name"`
	Description string `json:"description"`

	// "github.com/IGLOU-EU/go-wildcard/v2" package sintax
	MatchExp string `json:"match_exp,omitempty"`

	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	Models []*DeviceModel `json:"models,omitempty" gorm:"foreignKey:MakerCodename;references:Codename"`

	// Time marks
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName in database
func (m *DeviceMaker) TableName() string {
	return `type_device_maker`
}

// RBACResourceName returns the name of the resource for the RBAC
func (m *DeviceMaker) RBACResourceName() string {
	return "device_maker"
}
