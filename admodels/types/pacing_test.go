package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdPacing(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		assert.Equal(t, "asap", AdPacingASAP.Name())
		assert.False(t, AdPacingASAP.IsEvenly())
		assert.Equal(t, "evenly", AdPacingEvenly.Name())
		assert.True(t, AdPacingEvenly.IsEvenly())
		assert.Equal(t, "asap", AdPacing(0).Name())
	})
	t.Run("decode", func(t *testing.T) {
		tests := []struct {
			name     string
			expected AdPacing
		}{
			{name: "asap", expected: AdPacingASAP},
			{name: "ASAP", expected: AdPacingASAP},
			{name: "evenly", expected: AdPacingEvenly},
			{name: "EVENLY", expected: AdPacingEvenly},
			{name: "even", expected: AdPacingEvenly},
			{name: "other", expected: AdPacingASAP},
		}
		for _, test := range tests {
			assert.Equal(t, AdPacingByName(test.name), test.expected)
		}
	})
	t.Run("json", func(t *testing.T) {
		data, err := json.Marshal(AdPacingEvenly)
		assert.NoError(t, err)
		assert.Equal(t, `"evenly"`, string(data))
		var p AdPacing
		err = json.Unmarshal(data, &p)
		assert.NoError(t, err)
		assert.Equal(t, AdPacingEvenly, p)
	})
}
