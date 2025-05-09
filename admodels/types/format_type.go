//
// @project GeniusRabbit corelib 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package types

import (
	"bytes"
	"database/sql/driver"

	"github.com/geniusrabbit/gosql/v2"
	"github.com/pkg/errors"
)

var errInvalidUnmarshalValue = errors.New("invalid unmarshal value")

// FormatType value
// CREATE TYPE FormatType AS ENUM ('invalid', 'undefined', 'direct', 'proxy', 'video', 'banner', 'html5', 'native', 'custom')
type FormatType int

// Format types
const (
	FormatInvalidType     FormatType = -1
	FormatUndefinedType   FormatType = 0
	FormatDirectType      FormatType = 1
	FormatProxyType       FormatType = 2
	FormatVideoType       FormatType = 3 // It's kinde of integrated into video player
	FormatBannerType      FormatType = 4
	FormatBannerHTML5Type FormatType = 5
	FormatNativeType      FormatType = 6
	FormatCustomType      FormatType = 31
)

// FormatTypeList of types
var FormatTypeList = []FormatType{
	FormatInvalidType,
	FormatUndefinedType,
	FormatDirectType,
	FormatProxyType, // One of the banner types which works by iframe
	FormatVideoType,
	FormatBannerType,
	FormatBannerHTML5Type,
	FormatNativeType,
	FormatCustomType,
}

// FormatMapping from name to constant
var FormatMapping = map[string]FormatType{
	`invalid`:   FormatInvalidType,
	`undefined`: FormatUndefinedType,
	`direct`:    FormatDirectType,
	`proxy`:     FormatProxyType,
	`video`:     FormatVideoType,
	`banner`:    FormatBannerType,
	`html5`:     FormatBannerHTML5Type,
	`native`:    FormatNativeType,
	`custom`:    FormatCustomType,
}

// FormatTypeByName returns format type by name
func FormatTypeByName(name string) FormatType {
	if t, ok := FormatMapping[name]; ok {
		return t
	}
	return FormatInvalidType
}

// Name by format type
func (t FormatType) Name() string {
	switch t {
	case FormatInvalidType:
		return `invalid`
	case FormatDirectType:
		return `direct`
	case FormatProxyType:
		return `proxy`
	case FormatVideoType:
		return `video`
	case FormatBannerType:
		return `banner`
	case FormatBannerHTML5Type:
		return `html5`
	case FormatNativeType:
		return `native`
	case FormatCustomType:
		return `custom`
	}
	return `undefined`
}

// DisplayName of the format type
func (t FormatType) DisplayName() string {
	return t.Name()
}

// IsInvalid type of format
func (t FormatType) IsInvalid() bool {
	return t == FormatInvalidType
}

// IsDirect format type
func (t FormatType) IsDirect() bool {
	return t == FormatDirectType
}

// IsBanner format type
func (t FormatType) IsBanner() bool {
	return t == FormatBannerType
}

// IsBannerHTML5 format type
func (t FormatType) IsBannerHTML5() bool {
	return t == FormatBannerHTML5Type
}

// IsProxy format type
func (t FormatType) IsProxy() bool {
	return t == FormatProxyType
}

// IsNative format type
func (t FormatType) IsNative() bool {
	return t == FormatNativeType
}

// Value implements the driver.Valuer interface, json field interface
func (t FormatType) Value() (driver.Value, error) {
	return t.Name(), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (t *FormatType) Scan(value any) error {
	switch v := value.(type) {
	case string:
		return t.UnmarshalJSON([]byte(v))
	case []byte:
		return t.UnmarshalJSON(v)
	case nil:
		*t = FormatUndefinedType
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (t FormatType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Name() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (t *FormatType) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.Wrap(errInvalidUnmarshalValue, "`"+string(b)+"`")
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*t = FormatMapping[string(b[1:len(b)-1])]
	} else {
		*t = FormatMapping[string(b)]
	}
	return nil
}
