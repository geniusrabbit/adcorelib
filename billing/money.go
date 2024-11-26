//
// @project GeniusRabbit.com Billing 2015-2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2015-2017
//

package billing

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"golang.org/x/exp/constraints"

	"github.com/geniusrabbit/gosql/v2"
)

const (
	moneyFloatDelimeter = 1000000000.0
	moneyIntDelimeter   = 1000000000
)

// Money value type: first 9 numbers it's float part
// Example:
//
//	1_000000000 => 1    $
//	0_010000000 => 0.01 $
type Money int64

// MoneyFloat object
//
//go:inline
func MoneyFloat[T constraints.Float](amount T) Money {
	return Money(math.Floor(float64(amount)*moneyFloatDelimeter + .5))
}

// MoneyInt object
//
//go:inline
func MoneyInt[T constraints.Integer](amount T) Money {
	return Money(int64(amount) * moneyIntDelimeter)
}

// String implementation of Stringer interface
func (m Money) String() string {
	return fmt.Sprintf("%.09f", m.Float64())
}

// Float64 value from money
//
//go:inline
func (m Money) Float64() float64 {
	return float64(m) / moneyFloatDelimeter
}

// Int64 value from money
//
//go:inline
func (m Money) Int64() int64 {
	return int64(m)
}

// SetFromInt64 value
//
//go:inline
func (m *Money) SetFromInt64(v int64) {
	*m = MoneyInt(v)
}

// SetFromFloat64 value
//
//go:inline
func (m *Money) SetFromFloat64(v float64) {
	*m = MoneyFloat(v)
}

///////////////////////////////////////////////////////////////////////////////
// Encode/Decode
///////////////////////////////////////////////////////////////////////////////

// Value implements the driver.Valuer interface
func (m Money) Value() (driver.Value, error) {
	return m.Int64(), nil
}

// Scan implements the sql.Scanner interface
func (m *Money) Scan(value any) error {
	var data []byte
	switch v := value.(type) {
	case int:
		*m = Money(v)
		return nil
	case int32:
		*m = Money(v)
		return nil
	case int64:
		*m = Money(v)
		return nil
	case uint:
		*m = Money(v)
		return nil
	case uint32:
		*m = Money(v)
		return nil
	case uint64:
		*m = Money(v)
		return nil
	case float32:
		*m = MoneyFloat(v)
		return nil
	case float64:
		*m = MoneyFloat(v)
		return nil
	case string:
		data = []byte(v)
	case []byte:
		data = v
	case nil:
		*m = 0
		return nil
	default:
		return gosql.ErrInvalidScan
	}
	if bytes.ContainsAny(data, ".") {
		if val, err := strconv.ParseFloat(string(data), 64); err == nil {
			*m = MoneyFloat(val)
			return nil
		} else {
			return err
		}
	}
	val, err := strconv.ParseInt(string(data), 10, 64)
	if err == nil {
		*m = Money(val)
	}
	return err
}

// MarshalJSON implements the json.Marshaler
func (m Money) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Float64())
}

// UnmarshalJSON implements the json.Unmarshaller
func (m *Money) UnmarshalJSON(b []byte) error {
	if bytes.ContainsAny(b, ".") {
		if val, err := strconv.ParseFloat(string(b), 64); err == nil {
			*m = MoneyFloat(val)
			return nil
		} else {
			return err
		}
	}

	val, err := strconv.ParseInt(string(b), 10, 64)
	if err == nil {
		*m = MoneyInt(val)
	}
	return err
}

// DecodeValue implements the gocast.Decoder
func (m *Money) DecodeValue(v any) (err error) {
	switch val := v.(type) {
	case []byte:
		err = m.UnmarshalJSON(val)
	case string:
		err = m.UnmarshalJSON([]byte(val))
	case float64:
		*m = MoneyFloat(val)
	case float32:
		*m = MoneyFloat(float64(val))
	case int64:
		*m = Money(val)
	case Money:
		*m = val
	default:
		err = fmt.Errorf("invalid decode value")
	}
	return
}

// Implements the Unmarshaler interface of the yaml pkg.
func (m *Money) UnmarshalYAML(unmarshal func(any) error) error {
	var yamlValue float64
	if err := unmarshal(&yamlValue); err != nil {
		return err
	}
	m.SetFromFloat64(yamlValue)
	return nil
}
