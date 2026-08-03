//
// @project GeniusRabbit corelib 2018 - 2019, 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019, 2022
//

package types

import (
	"bytes"
	"database/sql/driver"
	"strings"

	"github.com/geniusrabbit/gosql/v2"
	"github.com/pkg/errors"
)

// PricingModel value
// CREATE TYPE PricingModel AS ENUM ('undefined', 'CPM', 'CPMV', 'CPC', 'CPA')
type PricingModel uint8

// PricingModel consts
const (
	PricingModelUndefined PricingModel = iota
	PricingModelCPM
	PricingModelCPMV
	PricingModelCPC
	PricingModelCPA
)

// PricingModelByName string
func PricingModelByName(model string) PricingModel {
	switch strings.ToUpper(model) {
	case `CPM`, `1`:
		return PricingModelCPM
	case `CPMV`, `2`:
		return PricingModelCPMV
	case `CPC`, `3`:
		return PricingModelCPC
	case `CPA`, `4`:
		return PricingModelCPA
	}
	return PricingModelUndefined
}

func (pm PricingModel) String() string {
	return pm.Name()
}

// Or returns current value if not undefined or alternative value
func (pm PricingModel) Or(npm PricingModel) PricingModel {
	if pm == PricingModelUndefined {
		return npm
	}
	return pm
}

// Name value
func (pm PricingModel) Name() string {
	switch pm {
	case PricingModelCPM:
		return `CPM`
	case PricingModelCPMV:
		return `CPMV`
	case PricingModelCPC:
		return `CPC`
	case PricingModelCPA:
		return `CPA`
	}
	return `undefined`
}

// IsCPM model
//
//go:inline
func (pm PricingModel) IsCPM() bool {
	return pm == PricingModelCPM
}

// IsCPMV model
//
//go:inline
func (pm PricingModel) IsCPMV() bool {
	return pm == PricingModelCPMV
}

// IsCPC model
//
//go:inline
func (pm PricingModel) IsCPC() bool {
	return pm == PricingModelCPC
}

// IsCPA model
//
//go:inline
func (pm PricingModel) IsCPA() bool {
	return pm == PricingModelCPA
}

// UInt value
//
//go:inline
func (pm PricingModel) UInt() uint {
	return uint(pm)
}

// Value implements the driver.Valuer interface, json field interface
func (pm PricingModel) Value() (driver.Value, error) {
	return pm.Name(), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (pm *PricingModel) Scan(value any) error {
	switch v := value.(type) {
	case string:
		*pm = PricingModelByName(v)
	case []byte:
		*pm = PricingModelByName(string(v))
	case nil:
		*pm = PricingModelUndefined
	default:
		return gosql.ErrInvalidScan
	}
	return nil
}

// MarshalJSON implements the json.Marshaler
func (pm PricingModel) MarshalJSON() ([]byte, error) {
	return []byte(`"` + pm.Name() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaller
func (pm *PricingModel) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return errors.Wrap(errInvalidUnmarshalValue, "`"+string(b)+"`")
	}
	if bytes.HasPrefix(b, []byte(`"`)) {
		*pm = PricingModelByName(string(b[1 : len(b)-1]))
	} else {
		*pm = PricingModelByName(string(b))
	}
	return nil
}
