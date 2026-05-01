//
// @project GeniusRabbit corelib 2018, 2025
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018, 2025
//

package types

import (
	"bytes"
	"database/sql/driver"

	"github.com/geniusrabbit/gosql/v2"
)

// AdFileAssetType represents the type of the asset
type AdFileAssetType uint

// AdFileAssetType values
const (
	AdFileAssetUndefinedType AdFileAssetType = 0
	AdFileAssetImageType     AdFileAssetType = 1
	AdFileAssetVideoType     AdFileAssetType = 2
	AdFileAssetHTML5Type     AdFileAssetType = 3
	AdFileAssetVASTTagType   AdFileAssetType = 4
)

// AdFileAssetTypeByName returns adfile value type
func AdFileAssetTypeByName(name string) AdFileAssetType {
	switch name {
	case "image", "img", "1":
		return AdFileAssetImageType
	case "video", "2":
		return AdFileAssetVideoType
	case "html5", "3":
		return AdFileAssetHTML5Type
	case "vast_tag", "vast", "4":
		return AdFileAssetVASTTagType
	}
	return AdFileAssetUndefinedType
}

// Name of the asset
func (ft AdFileAssetType) Name() string {
	return ft.Code()
}

// Code of the option
func (ft AdFileAssetType) Code() string {
	switch ft {
	case AdFileAssetImageType:
		return "image"
	case AdFileAssetVideoType:
		return "video"
	case AdFileAssetHTML5Type:
		return "html5"
	case AdFileAssetVASTTagType:
		return "vast_tag"
	}
	return "undefined"
}

// Num of the option
func (ft AdFileAssetType) Num() int {
	switch ft {
	case AdFileAssetImageType:
		return 1
	case AdFileAssetVideoType:
		return 2
	case AdFileAssetHTML5Type:
		return 3
	case AdFileAssetVASTTagType:
		return 4
	}
	return 0
}

// IsImage file type
func (ft AdFileAssetType) IsImage() bool {
	return ft == AdFileAssetImageType
}

// IsVideo file type
func (ft AdFileAssetType) IsVideo() bool {
	return ft == AdFileAssetVideoType
}

// IsHTML5 file type
func (ft AdFileAssetType) IsHTML5() bool {
	return ft == AdFileAssetHTML5Type
}

// IsVASTTag file type
func (ft AdFileAssetType) IsVASTTag() bool {
	return ft == AdFileAssetVASTTagType
}

// IsUndefined file type
func (ft AdFileAssetType) IsUndefined() bool {
	return ft == AdFileAssetUndefinedType
}

// Value implements the driver.Valuer interface, json field interface
func (ft AdFileAssetType) Value() (driver.Value, error) {
	return []byte(ft.Code()), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (ft *AdFileAssetType) Scan(value any) error {
	switch v := value.(type) {
	case string:
		return ft.UnmarshalJSON([]byte(v))
	case []byte:
		return ft.UnmarshalJSON(v)
	case nil:
		*ft = AdFileAssetUndefinedType
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (ft AdFileAssetType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ft.Code() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (ft *AdFileAssetType) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errInvalidUnmarshalValue
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*ft = AdFileAssetTypeByName(string(b[1 : len(b)-1]))
	} else {
		*ft = AdFileAssetTypeByName(string(b))
	}
	return nil
}
