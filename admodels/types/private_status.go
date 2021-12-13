//
// @project geniusrabbit::corelib 2016 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2018
//

package types

import (
	"database/sql/driver"

	"github.com/geniusrabbit/gosql"
)

// PrivateStatus of the object
type PrivateStatus uint

// Status private
const (
	StatusPublic  PrivateStatus = 0
	StatusPrivate PrivateStatus = 1
)

// PrivateNameToStatus converts private status name to status
func PrivateNameToStatus(name string) PrivateStatus {
	switch name {
	case `private`, `true`, `yes`, `1`:
		return StatusPrivate
	default:
		return StatusPublic
	}
}

// Name of the status
func (st PrivateStatus) Name() string {
	if st == StatusPrivate {
		return `private`
	}
	return `public`
}

// IsPrivate status of the object
func (st PrivateStatus) IsPrivate() bool {
	return st == StatusPrivate
}

// Value implements the driver.Valuer interface, json field interface
func (st PrivateStatus) Value() (_ driver.Value, err error) {
	var v []byte
	if v, err := st.MarshalJSON(); err == nil && nil != v {
		return string(v), nil
	}
	return v, err
}

// Scan implements the driver.Valuer interface, json field interface
func (st *PrivateStatus) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		*st = PrivateNameToStatus(v[1 : len(v)-1])
	case []byte:
		*st = PrivateNameToStatus(string(v[1 : len(v)-1]))
	case nil:
		*st = StatusPublic
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (st PrivateStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + st.Name() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (st *PrivateStatus) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errInvalidUnmarshalValue
	}
	*st = PrivateNameToStatus(string(b[1 : len(b)-1]))
	return nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (st *PrivateStatus) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var yamlString string
	if err := unmarshal(&yamlString); err != nil {
		return err
	}
	*st = PrivateNameToStatus(yamlString)
	return nil
}
