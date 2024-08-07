//
// @project GeniusRabbit corelib 2016 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2019
//

package models

import (
	"strings"
	"time"

	"github.com/geniusrabbit/gosql/v2"
	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// Format types...
// TODO: switch format-type check on type
const (
	FormatTypeDirect = `direct`
	FormatTypeVideo  = `video`
	FormatTypeProxy  = `proxy`
	FormatTypeNative = `native`
)

// Format model description
type Format struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Codename string `json:"codename"`
	Type     string `json:"type"`

	Title       string `json:"title"`
	Description string `json:"description"`

	Active types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	Width     int `json:"width,omitempty"`
	Height    int `json:"height,omitempty"`
	MinWidth  int `json:"min_width,omitempty"`
	MinHeight int `json:"min_height,omitempty"`

	Config gosql.NullableJSON[types.FormatConfig] `gorm:"type:JSONB" json:"config,omitempty"`

	// Time marks
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// TableName in database
func (f *Format) TableName() string {
	return "adv_format"
}

// RBACResourceName returns the name of the resource for the RBAC
func (f *Format) RBACResourceName() string {
	return "adv_format"
}

// IsStretch format
func (f *Format) IsStretch() bool {
	return f != nil && !f.IsFixed()
}

// IsFixed format
func (f *Format) IsFixed() bool {
	return f != nil &&
		f.Width > 0 && (f.MinWidth <= 0 || f.Width == f.MinWidth) &&
		f.Height > 0 && (f.MinHeight <= 0 || f.Height == f.MinHeight)
}

// IsDirect type (popunder or proxy or etc.)
func (f Format) IsDirect() bool {
	return f.Type == FormatTypeDirect || f.Codename == FormatTypeDirect
}

// IsVideo type of advertisement
func (f Format) IsVideo() bool {
	return f.Type == FormatTypeVideo || strings.HasPrefix(f.Codename, FormatTypeVideo)
}

// IsProxy type of format
func (f Format) IsProxy() bool {
	return f.Type == FormatTypeProxy || strings.HasPrefix(f.Codename, FormatTypeProxy)
}

// IsNative type of format
func (f Format) IsNative() bool {
	return f.Type == FormatTypeNative || strings.HasPrefix(f.Codename, FormatTypeNative)
}
