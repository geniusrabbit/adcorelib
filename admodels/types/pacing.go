//
// @project GeniusRabbit corelib
//

package types

import (
	"bytes"
	"database/sql/driver"
	"strings"

	"github.com/geniusrabbit/gosql/v2"
	"github.com/pkg/errors"
)

// AdPacing is the AdItem budget-delivery policy.
// CREATE TYPE AdPacing AS ENUM ('asap', 'evenly')
type AdPacing uint8

// AdPacing consts. Zero value is ASAP (default).
const (
	AdPacingASAP AdPacing = iota
	AdPacingEvenly
)

// AdPacingByName string
func AdPacingByName(name string) AdPacing {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case `evenly`, `even`, `1`:
		return AdPacingEvenly
	}
	return AdPacingASAP
}

func (p AdPacing) String() string {
	return p.Name()
}

// Name value
func (p AdPacing) Name() string {
	if p == AdPacingEvenly {
		return `evenly`
	}
	return `asap`
}

// IsEvenly pacing
//
//go:inline
func (p AdPacing) IsEvenly() bool {
	return p == AdPacingEvenly
}

// Value implements the driver.Valuer interface, json field interface
func (p AdPacing) Value() (driver.Value, error) {
	return p.Name(), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (p *AdPacing) Scan(value any) error {
	switch v := value.(type) {
	case string:
		*p = AdPacingByName(v)
	case []byte:
		*p = AdPacingByName(string(v))
	case nil:
		*p = AdPacingASAP
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (p AdPacing) MarshalJSON() ([]byte, error) {
	return []byte(`"` + p.Name() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (p *AdPacing) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.Wrap(errInvalidUnmarshalValue, "`"+string(b)+"`")
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*p = AdPacingByName(string(b[1 : len(b)-1]))
	} else {
		*p = AdPacingByName(string(b))
	}
	return nil
}
