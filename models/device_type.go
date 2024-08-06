package models

import (
	"time"

	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/udetect"
)

// Device set of types
const (
	DeviceTypeUnknown   = udetect.DeviceTypeUnknown
	DeviceTypeMobile    = udetect.DeviceTypeMobile
	DeviceTypePC        = udetect.DeviceTypePC
	DeviceTypeTV        = udetect.DeviceTypeTV
	DeviceTypePhone     = udetect.DeviceTypePhone
	DeviceTypeTablet    = udetect.DeviceTypeTablet
	DeviceTypeConnected = udetect.DeviceTypeConnected
	DeviceTypeSetTopBox = udetect.DeviceTypeSetTopBox
	DeviceTypeWatch     = udetect.DeviceTypeWatch
	DeviceTypeGlasses   = udetect.DeviceTypeGlasses
	DeviceTypeOOH       = udetect.DeviceTypeOOH
)

type DeviceType struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	// Time marks
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName in database
func (m *DeviceType) TableName() string {
	return `type_device_type`
}

// RBACResourceName returns the name of the resource for the RBAC
func (m *DeviceType) RBACResourceName() string {
	return "device_type"
}

// DeviceTypeList is a list of DeviceType
var DeviceTypeList = []*DeviceType{
	{
		ID:          uint64(DeviceTypeUnknown),
		Name:        "Unknown",
		Description: "Unknown device type",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeMobile),
		Name:        "Mobile",
		Description: "Mobile device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypePC),
		Name:        "PC",
		Description: "Personal Computer",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeTV),
		Name:        "TV",
		Description: "TV device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypePhone),
		Name:        "Phone",
		Description: "Phone device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeTablet),
		Name:        "Tablet",
		Description: "Tablet device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeConnected),
		Name:        "Connected",
		Description: "Connected device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeSetTopBox),
		Name:        "SetTopBox",
		Description: "SetTopBox device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeWatch),
		Name:        "Watch",
		Description: "Watch device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeGlasses),
		Name:        "Glasses",
		Description: "Glasses device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeOOH),
		Name:        "OOH",
		Description: "Out of Home device",
		Active:      types.StatusActive,
	},
}
