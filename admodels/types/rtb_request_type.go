package types

import (
	"bytes"
	"database/sql/driver"

	"github.com/geniusrabbit/gosql/v2"
	"github.com/pkg/errors"
)

// RTBRequestType contains type of representation of request information
// CREATE TYPE RTBRequestType AS ENUM ('undefined','json','xml','protobuff','postformencoded','plain')
type RTBRequestType int

// Request types
const (
	RTBRequestTypeUndefined       RTBRequestType = 0
	RTBRequestTypeJSON            RTBRequestType = 1
	RTBRequestTypeXML             RTBRequestType = 2
	RTBRequestTypeProtoBUFF       RTBRequestType = 3
	RTBRequestTypePOSTFormEncoded RTBRequestType = 4 // application/x-www-form-urlencoded
	RTBRequestTypePLAINTEXT       RTBRequestType = 5
)

var (
	rtbRequestTypeMapping = map[string]RTBRequestType{
		"undefined":       RTBRequestTypeUndefined,
		"json":            RTBRequestTypeJSON,
		"xml":             RTBRequestTypeXML,
		"protobuff":       RTBRequestTypeProtoBUFF,
		"postformencoded": RTBRequestTypePOSTFormEncoded,
		"plain":           RTBRequestTypePLAINTEXT,
	}
)

// Name of the constant
func (rt RTBRequestType) Name() string {
	switch rt {
	case RTBRequestTypeJSON:
		return "json"
	case RTBRequestTypeXML:
		return "xml"
	case RTBRequestTypeProtoBUFF:
		return "protobuff"
	case RTBRequestTypePOSTFormEncoded:
		return "postformencoded"
	case RTBRequestTypePLAINTEXT:
		return "plain"
	}
	return "undefined"
}

// Value implements the driver.Valuer interface, json field interface
func (rt RTBRequestType) Value() (driver.Value, error) {
	return rt.Name(), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (rt *RTBRequestType) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		*rt = rtbRequestTypeMapping[v[1:len(v)-1]]
	case []byte:
		*rt = rtbRequestTypeMapping[string(v[1:len(v)-1])]
	case nil:
		*rt = RTBRequestTypeUndefined
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (rt RTBRequestType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + rt.Name() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (rt *RTBRequestType) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.Wrap(errInvalidUnmarshalValue, "`"+string(b)+"`")
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*rt = rtbRequestTypeMapping[string(b[1:len(b)-1])]
	} else {
		*rt = rtbRequestTypeMapping[string(b)]
	}
	return nil
}
