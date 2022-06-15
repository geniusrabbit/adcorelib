package types

import (
	"bytes"
	"database/sql/driver"

	"github.com/geniusrabbit/gosql/v2"
)

// ProcessingStatus of any item type
type ProcessingStatus int

// ProcessingStatus values
const (
	ProcessingUndefined ProcessingStatus = 0
	ProcessingProgress  ProcessingStatus = 1
	ProcessingProcessed ProcessingStatus = 2
	ProcessingError     ProcessingStatus = 3
	ProcessingDeleted   ProcessingStatus = 4
)

func ProcessingStatusByName(name string) ProcessingStatus {
	switch name {
	case "progress", "1":
		return ProcessingProgress
	case "processed", "2":
		return ProcessingProcessed
	case "error", "3":
		return ProcessingError
	case "deleted", "4":
		return ProcessingDeleted
	}
	return ProcessingUndefined
}

// Name returns the name of status
func (st ProcessingStatus) Name() string {
	switch st {
	case ProcessingProgress:
		return "Progress"
	case ProcessingProcessed:
		return "Processed"
	case ProcessingError:
		return "Error"
	case ProcessingDeleted:
		return "Deleted"
	}
	return "Undefined"
}

// Code returns value of status
func (st ProcessingStatus) Code() string {
	switch st {
	case ProcessingProgress:
		return "progress"
	case ProcessingProcessed:
		return "processed"
	case ProcessingError:
		return "error"
	case ProcessingDeleted:
		return "deleted"
	}
	return "undefined"
}

// Value implements the driver.Valuer interface, json field interface
func (st ProcessingStatus) Value() (driver.Value, error) {
	return []byte(st.Code()), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (st *ProcessingStatus) Scan(value any) error {
	switch v := value.(type) {
	case string:
		return st.UnmarshalJSON([]byte(v))
	case []byte:
		return st.UnmarshalJSON(v)
	case nil:
		*st = ProcessingUndefined
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (st ProcessingStatus) MarshalJSON() ([]byte, error) {
	return []byte(`"` + st.Code() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (st *ProcessingStatus) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errInvalidUnmarshalValue
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*st = ProcessingStatusByName(string(b[1 : len(b)-1]))
	} else {
		*st = ProcessingStatusByName(string(b))
	}
	return nil
}
