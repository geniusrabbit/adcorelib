package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/gosql/v2"
)

type DeviceModelVersion struct {
	Min  types.Version `json:"min"`
	Max  types.Version `json:"max"`
	Name string        `json:"name,omitempty"`
}

type DeviceModel struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`

	MatchExp string `json:"match_exp,omitempty"`

	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	// Link to device maker
	MakerID uint64       `json:"maker_id"`
	Maker   *DeviceMaker `json:"maker,omitempty" gorm:"foreignKey:MakerID;references:ID"`

	// Device type
	TypeID uint64      `json:"type_id"`
	Type   *DeviceType `json:"type,omitempty" gorm:"foreignKey:TypeID;references:ID"`

	Versions gosql.NullableJSONArray[DeviceModelVersion] `json:"versions,omitempty" gorm:"type:jsonb"`

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
