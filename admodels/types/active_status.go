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
// CREATE TYPE ActiveStatus AS ENUM ('active', 'pause')
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

// DisplayName of the status
func (st ActiveStatus) DisplayName() string {
	return st.Name()
}

// IsActive status of the object
func (st ActiveStatus) IsActive() bool {
	return st == StatusActive
}

// Value implements the driver.Valuer interface, json field interface
func (st ActiveStatus) Value() (driver.Value, error) {
	return st.Name(), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (st *ActiveStatus) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		return st.UnmarshalJSON([]byte(v))
	case []byte:
		return st.UnmarshalJSON(v)
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
