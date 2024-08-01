package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

type DeviceModel struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	MakerID uint64       `json:"maker_id"`
	Maker   *DeviceMaker `json:"maker,omitempty" gorm:"foreignKey:MakerID;references:ID"`

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
