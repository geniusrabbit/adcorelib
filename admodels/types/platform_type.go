//
// @project GeniusRabbit::corelib 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2022
//

package types

import (
	"bytes"
	"database/sql/driver"

	"github.com/geniusrabbit/gosql/v2"
	"github.com/pkg/errors"
)

// PlatformType type
// CREATE TYPE PlatformType AS ENUM ('web', 'desktop', 'mobile', 'smart.phone', 'tablet', 'smart.tv', 'gamestation', 'smart.watch', 'vr', 'smart.glasses', 'smart.billboard')
type PlatformType uint

// Status approve
const (
	PlatformUndefined          PlatformType = 0
	PlatformWeb                PlatformType = 1
	PlatformDesktop            PlatformType = 2
	PlatformMobile             PlatformType = 3 // Smart.Phone & Tablet
	PlatformSmartPhone         PlatformType = 4
	PlatformTablet             PlatformType = 5
	PlatformSmartTV            PlatformType = 6
	PlatformGameStation        PlatformType = 7
	PlatformSmartWatch         PlatformType = 8
	PlatformVR                 PlatformType = 9
	PlatformSmartGlasses       PlatformType = 10
	PlatformSmartBillboard     PlatformType = 11
	PlatformUndefinedName                   = `undefined`
	PlatformWebName                         = `web`
	PlatformDesktopName                     = `desktop`
	PlatformMobileName                      = `mobile`
	PlatformSmartPhoneName                  = `smart.phone`
	PlatformTabletName                      = `tablet`
	PlatformSmartTVName                     = `smart.tv`
	PlatformGameStationName                 = `gamestation`
	PlatformSmartWatchName                  = `smart.watch`
	PlatformVRName                          = `vr`
	PlatformSmartGlassesName                = `smart.glasses`
	PlatformSmartBillboardName              = `smart.billboard`
)

// PlatformTypeNameList contains list of available platform names
var PlatformTypeNameList = []string{
	PlatformUndefinedName,
	PlatformWebName,
	PlatformDesktopName,
	PlatformMobileName,
	PlatformSmartPhoneName,
	PlatformTabletName,
	PlatformSmartTVName,
	PlatformGameStationName,
	PlatformSmartWatchName,
	PlatformVRName,
	PlatformSmartGlassesName,
	PlatformSmartBillboardName,
}

// Name of the type
func (tp PlatformType) Name() string {
	switch tp {
	case PlatformUndefined:
		return PlatformUndefinedName
	case PlatformWeb:
		return PlatformWebName
	case PlatformDesktop:
		return PlatformDesktopName
	case PlatformMobile:
		return PlatformMobileName
	case PlatformSmartPhone:
		return PlatformSmartPhoneName
	case PlatformTablet:
		return PlatformTabletName
	case PlatformSmartTV:
		return PlatformSmartTVName
	case PlatformGameStation:
		return PlatformGameStationName
	case PlatformSmartWatch:
		return PlatformSmartWatchName
	case PlatformVR:
		return PlatformVRName
	case PlatformSmartGlasses:
		return PlatformSmartGlassesName
	case PlatformSmartBillboard:
		return PlatformSmartBillboardName
	}
	return PlatformUndefinedName
}

// DisplayName of the type
func (tp PlatformType) DisplayName() string {
	return tp.Name()
}

// ApproveNameToStatus name to const
func PlatformTypeNameToType(name string) PlatformType {
	switch name {
	case PlatformWebName, `1`:
		return PlatformWeb
	case PlatformDesktopName, `2`:
		return PlatformDesktop
	case PlatformMobileName, `3`:
		return PlatformMobile
	case PlatformSmartPhoneName, `4`:
		return PlatformSmartPhone
	case PlatformTabletName, `5`:
		return PlatformTablet
	case PlatformSmartTVName, `6`:
		return PlatformSmartTV
	case PlatformGameStationName, `7`:
		return PlatformGameStation
	case PlatformSmartWatchName, `8`:
		return PlatformSmartWatch
	case PlatformVRName, `9`:
		return PlatformVR
	case PlatformSmartGlassesName, `10`:
		return PlatformSmartGlasses
	case PlatformSmartBillboardName, `11`:
		return PlatformSmartBillboard
	}
	return PlatformUndefined
}

// Value implements the driver.Valuer interface, json field interface
func (tp PlatformType) Value() (driver.Value, error) {
	return tp.Name(), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (tp *PlatformType) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		return tp.UnmarshalJSON([]byte(v))
	case []byte:
		return tp.UnmarshalJSON(v)
	case nil:
		*tp = PlatformUndefined
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (tp PlatformType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + tp.Name() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (tp *PlatformType) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.Wrap(errInvalidUnmarshalValue, "`"+string(b)+"`")
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*tp = PlatformTypeNameToType(string(b[1 : len(b)-1]))
	} else {
		*tp = PlatformTypeNameToType(string(b))
	}
	return nil
}
