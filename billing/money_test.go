package billing

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoney(t *testing.T) {
	t.Run("general", func(t *testing.T) {
		vl := MoneyInt(1)

		assert.Equal(t, int64(1*moneyFloatDelimeter), vl.I64())
		assert.Equal(t, int64(1*moneyFloatDelimeter), vl.Int64())
		assert.Equal(t, float64(1), vl.F64())
		assert.Equal(t, float64(1), vl.Float64())

		vl.SetFromInt64(2)
		assert.Equal(t, int64(2*moneyFloatDelimeter), vl.I64())
		assert.Equal(t, int64(2*moneyFloatDelimeter), vl.Int64())
		assert.Equal(t, float64(2), vl.F64())
		assert.Equal(t, float64(2), vl.Float64())

		vl.SetFromFloat64(3.5)
		assert.Equal(t, int64(3.5*moneyFloatDelimeter), vl.I64())
		assert.Equal(t, int64(3.5*moneyFloatDelimeter), vl.Int64())
		assert.Equal(t, float64(3.5), vl.F64())
		assert.Equal(t, float64(3.5), vl.Float64())
	})

	t.Run("json", func(t *testing.T) {
		vl := MoneyFloat(3.5)

		b, err := vl.MarshalJSON()
		assert.NoError(t, err)
		assert.Equal(t, "3.5", string(b))

		var vl2 Money
		err = json.Unmarshal([]byte("5"), &vl2)

		assert.NoError(t, err)
		assert.Equal(t, int64(5*moneyFloatDelimeter), vl2.I64())
		assert.Equal(t, int64(5*moneyFloatDelimeter), vl2.Int64())
		assert.Equal(t, float64(5), vl2.F64())
		assert.Equal(t, float64(5), vl2.Float64())

		err = json.Unmarshal([]byte("3.5"), &vl2)

		assert.NoError(t, err)
		assert.Equal(t, int64(3.5*moneyFloatDelimeter), vl2.I64())
		assert.Equal(t, int64(3.5*moneyFloatDelimeter), vl2.Int64())
		assert.Equal(t, float64(3.5), vl2.F64())
		assert.Equal(t, float64(3.5), vl2.Float64())
	})

	t.Run("decode", func(t *testing.T) {
		var vl Money
		tests := []struct {
			input     any
			expected  float64
			expectErr bool
		}{
			{input: "3.5", expected: 3.5, expectErr: false},
			{input: []byte("3.5"), expected: 3.5, expectErr: false},
			{input: 3.5, expected: 3.5, expectErr: false},
			{input: 5, expected: 5, expectErr: false},
			{input: int64(5), expected: 5, expectErr: false},
			{input: float32(5), expected: 5, expectErr: false},
			{input: "invalid", expected: 0, expectErr: true},
			{input: nil, expected: 0, expectErr: true},
			{input: "3.5abc", expected: 0, expectErr: true},
			{input: MoneyFloat(3.5), expected: 3.5, expectErr: false},
		}

		for _, test := range tests {
			err := vl.DecodeValue(test.input)
			if test.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, vl.F64())
			}
		}
	})
}
