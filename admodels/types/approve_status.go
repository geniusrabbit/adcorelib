//
// @project GeniusRabbit corelib 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2022
//

package types

import (
	"bytes"
	"database/sql/driver"

	"github.com/geniusrabbit/gosql/v2"
	"github.com/pkg/errors"
)

// ApproveStatus type
// CREATE TYPE ApproveStatus AS ENUM ('pending', 'approved', 'rejected')
type ApproveStatus uint

// Status approve
const (
	StatusPending      ApproveStatus = 0
	StatusApproved     ApproveStatus = 1
	StatusRejected     ApproveStatus = 2
	StatusPendingName                = `pending`
	StatusApprovedName               = `approved`
	StatusRejectedName               = `rejected`
)

// ApproveStatusNameList contains available names list
var ApproveStatusNameList = []string{
	StatusPendingName,
	StatusApprovedName,
	StatusRejectedName,
}

// Name of the status
func (st ApproveStatus) Name() string {
	switch st {
	case StatusApproved:
		return StatusApprovedName
	case StatusRejected:
		return StatusRejectedName
	}
	return StatusPendingName
}

// DisplayName of the status
func (st ApproveStatus) DisplayName() string {
	return st.Name()
}

// IsApproved status of the object
func (st ApproveStatus) IsApproved() bool {
	return st == StatusApproved
}

// ApproveNameToStatus name to const
func ApproveNameToStatus(name string) ApproveStatus {
	switch name {
	case StatusApprovedName, `1`:
		return StatusApproved
	case StatusRejectedName, `2`:
		return StatusRejected
	}
	return StatusPending
}

// Value implements the driver.Valuer interface, json field interface
func (st ApproveStatus) Value() (driver.Value, error) {
	return st.Name(), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (st *ApproveStatus) Scan(value any) error {
	switch v := value.(type) {
	case string:
		return st.UnmarshalJSON([]byte(v))
	case []byte:
		return st.UnmarshalJSON(v)
	case nil:
		*st = StatusPending
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (st ApproveStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + st.Name() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (st *ApproveStatus) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.Wrap(errInvalidUnmarshalValue, "`"+string(b)+"`")
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*st = ApproveNameToStatus(string(b[1 : len(b)-1]))
	} else {
		*st = ApproveNameToStatus(string(b))
	}
	return nil
}

// Implements the Unmarshaler interface of the yaml pkg.
func (st *ApproveStatus) UnmarshalYAML(unmarshal func(any) error) error {
	var yamlString string
	if err := unmarshal(&yamlString); err != nil {
		return err
	}
	*st = ApproveNameToStatus(yamlString)
	return nil
}
