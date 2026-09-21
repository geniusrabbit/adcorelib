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

// InterstitialSupport is the AdItem interstitial inventory policy.
// CREATE TYPE InterstitialSupport AS ENUM ('none', 'only', 'both')
type InterstitialSupport uint8

// InterstitialSupport consts. Zero value is None (default).
const (
	InterstitialSupportNone InterstitialSupport = iota
	InterstitialSupportOnly
	InterstitialSupportBoth
)

// InterstitialSupportByName string. Unknown/empty names map to None.
func InterstitialSupportByName(name string) InterstitialSupport {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case `only`, `1`:
		return InterstitialSupportOnly
	case `both`, `2`:
		return InterstitialSupportBoth
	}
	return InterstitialSupportNone
}

func (s InterstitialSupport) String() string {
	return s.Name()
}

// Name value
func (s InterstitialSupport) Name() string {
	switch s {
	case InterstitialSupportOnly:
		return `only`
	case InterstitialSupportBoth:
		return `both`
	default:
		return `none`
	}
}

// IsNone interstitial support (non-interstitial only)
//
//go:inline
func (s InterstitialSupport) IsNone() bool {
	return s == InterstitialSupportNone
}

// IsOnly interstitial support (interstitial only)
//
//go:inline
func (s InterstitialSupport) IsOnly() bool {
	return s == InterstitialSupportOnly
}

// IsBoth interstitial support (interstitial and non-interstitial)
//
//go:inline
func (s InterstitialSupport) IsBoth() bool {
	return s == InterstitialSupportBoth
}

// Value implements the driver.Valuer interface, json field interface
func (s InterstitialSupport) Value() (driver.Value, error) {
	return s.Name(), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (s *InterstitialSupport) Scan(value any) error {
	switch v := value.(type) {
	case string:
		*s = InterstitialSupportByName(v)
	case []byte:
		*s = InterstitialSupportByName(string(v))
	case nil:
		*s = InterstitialSupportNone
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (s InterstitialSupport) MarshalJSON() ([]byte, error) {
	return []byte(`"` + s.Name() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (s *InterstitialSupport) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.Wrap(errInvalidUnmarshalValue, "`"+string(b)+"`")
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*s = InterstitialSupportByName(string(b[1 : len(b)-1]))
	} else {
		*s = InterstitialSupportByName(string(b))
	}
	return nil
}
