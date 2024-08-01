package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	t.Run("encode/decode", func(t *testing.T) {
		tests := []struct {
			ver Version
			str string
		}{
			{Version{1, 2, 3}, "1.2.3"},
			{Version{1, 2, 0}, "1.2"},
			{Version{1, 0, 0}, "1"},
			{Version{0, 0, 0}, "0"},
		}

		for _, tt := range tests {
			if assert.Equal(t, tt.str, tt.ver.String(), "Version.String()") {
				v := Version{}
				if assert.NoError(t, v.SetFromStr(tt.str)) {
					assert.Equal(t, tt.ver, v)
				}
			}
		}
	})

	t.Run("json", func(t *testing.T) {
		v := Version{}
		err := json.Unmarshal([]byte(`"1.2.3"`), &v)
		if assert.NoError(t, err) {
			assert.Equal(t, Version{1, 2, 3}, v)
		}
	})
}
