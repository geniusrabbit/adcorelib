package languages

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// UndefinedISO2 is the sentinel ISO-639-1 stand-in for an unknown language.
const UndefinedISO2 = "**"

// UndefinedCode is the two-byte form of [UndefinedISO2].
var UndefinedCode = Code{UndefinedISO2[0], UndefinedISO2[1]}

// Error list...
var (
	ErrCodeInvalidScanType  = errors.New("[languages.code] invalid scan type, supports only bytes and string")
	ErrCodeInvalidValueSize = errors.New("[languages.code] invalid value size, supports only two 2 chars")
)

// Code is a 2-letter ISO-639-1 language code.
type Code [2]byte

// CodeFromString returns the canonical catalog code for s.
// BCP-47 tags use the first two letters (`en-US` → en). Unknown input is [UndefinedCode].
func CodeFromString(s string) Code {
	return GetLanguageByCodeString(s).Code
}

func (c Code) String() string {
	return c.ISO2()
}

// ISO2 returns the interned two-letter catalog code, or [UndefinedISO2].
func (c Code) ISO2() string {
	return c.Language().ISO2()
}

// Language returns the catalog entry for c. Unknown codes resolve to Languages[0] (never nil).
func (c Code) Language() *Language {
	return lookup(c[0], c[1])
}

// ID is the catalog language ID (0 for unknown / undefined).
func (c Code) ID() uint {
	return c.Language().ID
}

// IsUndefined reports whether c maps to Languages[0].
func (c Code) IsUndefined() bool {
	return c.Language().ID == 0
}

// Int is the packed little-endian numeric form of the two bytes.
func (c Code) Int() uint {
	return uint(c[0]) | uint(c[1])<<8
}

// Value implementation of sql driver.Valuer interface.
func (c Code) Value() (driver.Value, error) {
	return c.ISO2(), nil
}

// Scan implementation of sql.Scanner interface.
func (c *Code) Scan(data any) error {
	switch v := data.(type) {
	case []byte:
		if len(v) != 2 {
			return ErrCodeInvalidValueSize
		}
		*c = Code{asciiLower(v[0]), asciiLower(v[1])}
	case string:
		if len(v) != 2 {
			return ErrCodeInvalidValueSize
		}
		*c = Code{asciiLower(v[0]), asciiLower(v[1])}
	default:
		return ErrCodeInvalidScanType
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
func (c Code) MarshalJSON() ([]byte, error) {
	return []byte{'"', asciiLower(c[0]), asciiLower(c[1]), '"'}, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Code) UnmarshalJSON(data []byte) error {
	if n := len(data); n >= 2 && data[0] == '"' && data[n-1] == '"' {
		s := data[1 : n-1]
		if len(s) != 2 {
			*c = UndefinedCode
			return nil
		}
		*c = Code{asciiLower(s[0]), asciiLower(s[1])}
		return nil
	}
	var code string
	if err := json.Unmarshal(data, &code); err != nil {
		return err
	}
	if code == UndefinedISO2 || len(code) != 2 {
		*c = UndefinedCode
	} else {
		*c = Code{asciiLower(code[0]), asciiLower(code[1])}
	}
	return nil
}

var (
	_ driver.Valuer    = UndefinedCode
	_ sql.Scanner      = &UndefinedCode
	_ json.Marshaler   = UndefinedCode
	_ json.Unmarshaler = &UndefinedCode
)
