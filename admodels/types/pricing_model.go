//
// @project GeniusRabbit corelib 2018 - 2019, 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019, 2022
//

package types

import (
	"bytes"
	"database/sql/driver"
	"strings"

	"github.com/geniusrabbit/gosql"
	"github.com/pkg/errors"
)

// PricingModel value
// CREATE TYPE PricingModel AS ENUM ('undefined', 'CPM', 'CPC', 'CPA')
type PricingModel uint8

// PricingModel consts
const (
	PricingModelUndefined PricingModel = iota
	PricingModelCPM
	PricingModelCPC
	PricingModelCPA
)

// PricingModelByName string
func PricingModelByName(model string) PricingModel {
	switch strings.ToUpper(model) {
	case `CPM`, `1`:
		return PricingModelCPM
	case `CPC`, `2`:
		return PricingModelCPC
	case `CPA`, `3`:
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
	case PricingModelCPC:
		return `CPC`
	case PricingModelCPA:
		return `CPA`
	}
	return `undefined`
}

// IsCPM model
func (pm PricingModel) IsCPM() bool {
	return pm == PricingModelCPM
}

// IsCPC model
func (pm PricingModel) IsCPC() bool {
	return pm == PricingModelCPC
}

// IsCPA model
func (pm PricingModel) IsCPA() bool {
	return pm == PricingModelCPA
}

// UInt value
func (pm PricingModel) UInt() uint {
	return uint(pm)
}

// Value implements the driver.Valuer interface, json field interface
func (pm PricingModel) Value() (driver.Value, error) {
	return pm.Name(), nil
}

// Scan implements the driver.Valuer interface, json field interface
func (pm *PricingModel) Scan(value interface{}) error {
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
