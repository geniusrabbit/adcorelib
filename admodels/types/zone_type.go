//
// @project geniusrabbit::corelib 2016 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2018
//

package types

import (
	"bytes"
	"database/sql/driver"

	"github.com/geniusrabbit/gosql"
	"github.com/pkg/errors"
)

// PrivateStatus of the object
type ZoneType uint

// Zone Types enum
const (
	ZoneTypeDefault   ZoneType = iota // 0
	ZoneTypeSmartlink                 // 1
)

// ZoneNameToType converts zone type name to const
func ZoneNameToType(name string) ZoneType {
	switch name {
	case `smartlink`, `smart`, `1`:
		return ZoneTypeSmartlink
	default:
		return ZoneTypeDefault
	}
}

// Name of the status
func (tp ZoneType) Name() string {
	if tp == ZoneTypeSmartlink {
		return `smartlink`
	}
	return `zone`
}

// IsSmartlink type of the object
func (tp ZoneType) IsSmartlink() bool {
	return tp == ZoneTypeSmartlink
}

// Value implements the driver.Valuer interface, json field interface
func (st ZoneType) Value() (_ driver.Value, err error) {
	var v []byte
	if v, err := st.MarshalJSON(); err == nil && nil != v {
		return string(v), nil
	}
	return v, err
}

// Scan implements the driver.Valuer interface, json field interface
func (st *ZoneType) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		*st = ZoneNameToType(v[1 : len(v)-1])
	case []byte:
		*st = ZoneNameToType(string(v[1 : len(v)-1]))
	case nil:
		*st = ZoneTypeDefault
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (st ZoneType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + st.Name() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (st *ZoneType) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.Wrap(errInvalidUnmarshalValue, "`"+string(b)+"`")
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*st = ZoneNameToType(string(b[1 : len(b)-1]))
	} else {
		*st = ZoneNameToType(string(b))
	}
	return nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (st *ZoneType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var yamlString string
	if err := unmarshal(&yamlString); err != nil {
		return err
	}
	*st = ZoneNameToType(yamlString)
	return nil
}
