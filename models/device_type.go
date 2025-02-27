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
	DeviceTypeMobile    = udetect.DeviceTypeMobile    // Mobile/Tablet
	DeviceTypePC        = udetect.DeviceTypePC        // Desktop
	DeviceTypeTV        = udetect.DeviceTypeTV        // TV
	DeviceTypePhone     = udetect.DeviceTypePhone     // SmartPhone, SmallScreen
	DeviceTypeTablet    = udetect.DeviceTypeTablet    // Tablet
	DeviceTypeConnected = udetect.DeviceTypeConnected // Console, EReader, Watch
	DeviceTypeSetTopBox = udetect.DeviceTypeSetTopBox // MediaHub
	DeviceTypeWatch     = udetect.DeviceTypeWatch     // SmartWatch
	DeviceTypeGlasses   = udetect.DeviceTypeGlasses   // Glasses
	DeviceTypeOOH       = udetect.DeviceTypeOOH       // Out of Home - Billboards, Kiosks, etc.
)

type DeviceType struct {
	ID          uint64 `json:"id" gorm:"primaryKey"`
	Codename    string `json:"codename" gorm:"unique"`
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
		Codename:    "unknown",
		Name:        "Unknown",
		Description: "Unknown device type",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeMobile),
		Codename:    "mobile",
		Name:        "Mobile",
		Description: "Mobile device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypePC),
		Codename:    "pc",
		Name:        "PC",
		Description: "Personal Computer",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeTV),
		Codename:    "tv",
		Name:        "TV",
		Description: "TV device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypePhone),
		Codename:    "phone",
		Name:        "Phone",
		Description: "Phone device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeTablet),
		Codename:    "tablet",
		Name:        "Tablet",
		Description: "Tablet device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeConnected),
		Codename:    "connected",
		Name:        "Connected",
		Description: "Connected device (Console, EReader)",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeSetTopBox),
		Codename:    "settopbox",
		Name:        "SetTopBox",
		Description: "SetTopBox device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeWatch),
		Codename:    "watch",
		Name:        "Watch",
		Description: "Watch device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeGlasses),
		Codename:    "glasses",
		Name:        "Glasses",
		Description: "Glasses device",
		Active:      types.StatusActive,
	},
	{
		ID:          uint64(DeviceTypeOOH),
		Codename:    "ooh",
		Name:        "OOH",
		Description: "Out of Home device",
		Active:      types.StatusActive,
	},
}
