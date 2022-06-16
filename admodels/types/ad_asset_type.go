//
// @project GeniusRabbit rotator 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package types

import (
	"bytes"
	"database/sql/driver"

	"github.com/geniusrabbit/gosql/v2"
)

// import (
// 	diskmodels "geniusrabbit.dev/disk/models"
// )

// AdAssetType type
type AdAssetType uint

// AdAssetType values
const (
	AdAssetUndefinedType AdAssetType = 0
	AdAssetImageType     AdAssetType = 1
	AdAssetVideoType     AdAssetType = 2
	AdAssetHTML5Type     AdAssetType = 3
)

// AdAssetTypeByName returns adfile value type
func AdAssetTypeByName(name string) AdAssetType {
	switch name {
	case "image", "img", "1":
		return AdAssetImageType
	case "video", "2":
		return AdAssetVideoType
	case "html5", "3":
		return AdAssetHTML5Type
	}
	return AdAssetUndefinedType
}

// Name of the asset
func (ft AdAssetType) Name() string {
	return ft.Code()
}

// Code of the option
func (ft AdAssetType) Code() string {
	switch ft {
	case AdAssetImageType:
		return "image"
	case AdAssetVideoType:
		return "video"
	case AdAssetHTML5Type:
		return "html5"
	}
	return "undefined"
}

// Num of the option
func (ft AdAssetType) Num() int {
	switch ft {
	case AdAssetImageType:
		return 1
	case AdAssetVideoType:
		return 2
	case AdAssetHTML5Type:
		return 3
	}
	return 0
}

// IsImage file type
func (ft AdAssetType) IsImage() bool {
	return ft == AdAssetImageType
}

// IsVideo file type
func (ft AdAssetType) IsVideo() bool {
	return ft == AdAssetVideoType
}

// Value implements the driver.Valuer interface, json field interface
func (ft AdAssetType) Value() (driver.Value, error) {
	return []byte(ft.Code()), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (ft *AdAssetType) Scan(value any) error {
	switch v := value.(type) {
	case string:
		return ft.UnmarshalJSON([]byte(v))
	case []byte:
		return ft.UnmarshalJSON(v)
	case nil:
		*ft = AdAssetUndefinedType
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (ft AdAssetType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ft.Code() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (ft *AdAssetType) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errInvalidUnmarshalValue
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*ft = AdAssetTypeByName(string(b[1 : len(b)-1]))
	} else {
		*ft = AdAssetTypeByName(string(b))
	}
	return nil
}
