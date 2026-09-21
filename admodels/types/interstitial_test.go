package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterstitialSupport(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		assert.Equal(t, "none", InterstitialSupportNone.Name())
		assert.True(t, InterstitialSupportNone.IsNone())
		assert.False(t, InterstitialSupportNone.IsOnly())
		assert.False(t, InterstitialSupportNone.IsBoth())
		assert.Equal(t, "only", InterstitialSupportOnly.Name())
		assert.True(t, InterstitialSupportOnly.IsOnly())
		assert.Equal(t, "both", InterstitialSupportBoth.Name())
		assert.True(t, InterstitialSupportBoth.IsBoth())
		assert.Equal(t, "none", InterstitialSupport(0).Name())
	})
	t.Run("decode", func(t *testing.T) {
		tests := []struct {
			name     string
			expected InterstitialSupport
		}{
			{name: "none", expected: InterstitialSupportNone},
			{name: "NONE", expected: InterstitialSupportNone},
			{name: "0", expected: InterstitialSupportNone},
			{name: "only", expected: InterstitialSupportOnly},
			{name: "ONLY", expected: InterstitialSupportOnly},
			{name: "1", expected: InterstitialSupportOnly},
			{name: "both", expected: InterstitialSupportBoth},
			{name: "BOTH", expected: InterstitialSupportBoth},
			{name: "2", expected: InterstitialSupportBoth},
			{name: "other", expected: InterstitialSupportNone},
			{name: "", expected: InterstitialSupportNone},
		}
		for _, test := range tests {
			assert.Equal(t, InterstitialSupportByName(test.name), test.expected)
		}
	})
	t.Run("json", func(t *testing.T) {
		data, err := json.Marshal(InterstitialSupportBoth)
		assert.NoError(t, err)
		assert.Equal(t, `"both"`, string(data))
		var s InterstitialSupport
		err = json.Unmarshal(data, &s)
		assert.NoError(t, err)
		assert.Equal(t, InterstitialSupportBoth, s)
	})
}
