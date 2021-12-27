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

// ActiveStatus of the object
type ActiveStatus uint

// Status active
const (
	StatusPause  ActiveStatus = 0
	StatusActive ActiveStatus = 1
)

// ActiveNameToStatus converts avtivity status name to status
func ActiveNameToStatus(name string) ActiveStatus {
	switch name {
	case `active`, `true`, `yes`, `1`:
		return StatusActive
	default:
		return StatusPause
	}
}

// Name of the status
func (st ActiveStatus) Name() string {
	if st == StatusActive {
		return `active`
	}
	return `pause`
}

// IsActive status of the object
func (st ActiveStatus) IsActive() bool {
	return st == StatusActive
}

// Value implements the driver.Valuer interface, json field interface
func (st ActiveStatus) Value() (_ driver.Value, err error) {
	var v []byte
	if v, err := st.MarshalJSON(); err == nil && nil != v {
		return string(v), nil
	}
	return v, err
}

// Scan implements the driver.Valuer interface, json field interface
func (st *ActiveStatus) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		*st = ActiveNameToStatus(v[1 : len(v)-1])
	case []byte:
		*st = ActiveNameToStatus(string(v[1 : len(v)-1]))
	case nil:
		*st = StatusPause
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (st ActiveStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + st.Name() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (st *ActiveStatus) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.Wrap(errInvalidUnmarshalValue, "`"+string(b)+"`")
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*st = ActiveNameToStatus(string(b[1 : len(b)-1]))
	} else {
		*st = ActiveNameToStatus(string(b))
	}
	return nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (st *ActiveStatus) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var yamlString string
	if err := unmarshal(&yamlString); err != nil {
		return err
	}
	*st = ActiveNameToStatus(yamlString)
	return nil
}