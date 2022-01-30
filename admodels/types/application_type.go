//
// @project GeniusRabbit::corelib 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2022
//

package types

import (
	"bytes"
	"database/sql/driver"

	"github.com/geniusrabbit/gosql"
	"github.com/pkg/errors"
)

// ApplicationType type
// CREATE TYPE ApplicationType AS ENUM ('site', 'application', 'game')
type ApplicationType uint

// Status approve
const (
	ApplicationUndefined     ApplicationType = 0
	ApplicationSite          ApplicationType = 1
	ApplicationApp           ApplicationType = 2
	ApplicationGame          ApplicationType = 3
	ApplicationUndefinedName                 = `undefined`
	ApplicationSiteName                      = `site`
	ApplicationAppName                       = `application`
	ApplicationGameName                      = `game`
)

// ApplicationTypeNameList contains list of posible platform name
var ApplicationTypeNameList = []string{
	ApplicationUndefinedName,
	ApplicationSiteName,
	ApplicationAppName,
	ApplicationGameName,
}

// Name of the type
func (tp ApplicationType) Name() string {
	switch tp {
	case ApplicationUndefined:
		return ApplicationUndefinedName
	case ApplicationSite:
		return ApplicationSiteName
	case ApplicationApp:
		return ApplicationAppName
	case ApplicationGame:
		return ApplicationGameName
	}
	return ApplicationUndefinedName
}

// DisplayName of the type
func (tp ApplicationType) DisplayName() string {
	return tp.Name()
}

// ApproveNameToStatus name to const
func ApplicationTypeNameToType(name string) ApplicationType {
	switch name {
	case ApplicationSiteName, `1`:
		return ApplicationSite
	case ApplicationAppName, `2`:
		return ApplicationApp
	case ApplicationGameName, `3`:
		return ApplicationGame
	}
	return ApplicationUndefined
}

// Value implements the driver.Valuer interface, json field interface
func (tp ApplicationType) Value() (driver.Value, error) {
	return tp.Name(), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (tp *ApplicationType) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		return tp.UnmarshalJSON([]byte(v))
	case []byte:
		return tp.UnmarshalJSON(v)
	case nil:
		*tp = ApplicationUndefined
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (tp ApplicationType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + tp.Name() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (tp *ApplicationType) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.Wrap(errInvalidUnmarshalValue, "`"+string(b)+"`")
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*tp = ApplicationTypeNameToType(string(b[1 : len(b)-1]))
	} else {
		*tp = ApplicationTypeNameToType(string(b))
	}
	return nil
}
