//
// @project GeniusRabbit corelib 2016 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2019
//

package models

import (
	"strings"
	"time"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/gosql/v2"
	"github.com/guregu/null"
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
	ID       uint64             `json:"id"`
	Codename string             `json:"codename"`
	Type     string             `json:"type"`
	Title    string             `json:"title"`
	Active   types.ActiveStatus `gorm:"type:ActiveStatus" json:"active,omitempty"`

	Width     int `json:"width,omitempty"`
	Height    int `json:"height,omitempty"`
	MinWidth  int `json:"min_width,omitempty"`
	MinHeight int `json:"min_height,omitempty"`

	Config gosql.NullableJSON[types.FormatConfig] `gorm:"type:JSONB" json:"config,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt null.Time `json:"deleted_at,omitempty"`
}

// TableName in database
func (f *Format) TableName() string {
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
