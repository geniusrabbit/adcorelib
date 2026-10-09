package types

import (
	"encoding/json"
	"fmt"
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
			{Version{0, 0, 0}, ""},
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

	t.Run("parse", func(t *testing.T) {
		tests := []struct {
			in      string
			want    Version
			wantErr bool
		}{
			{in: "10.0", want: Version{10, 0, 0}},
			{in: "145.0.0.0", want: Version{145, 0, 0}},
			{in: "59.0.3071.115", want: Version{59, 0, 3071}},
			{in: "10.0.19045.3693", want: Version{10, 0, 19045}},
			{in: "10_15_7", want: Version{10, 15, 7}},
			{in: "v1.2.3", want: Version{1, 2, 3}},
			{in: "V10.0", want: Version{10, 0, 0}},
			{in: "1.2.3-rc.1+build", want: Version{1, 2, 3}},
			{in: "1.2.3+build.7", want: Version{1, 2, 3}},
			{in: "  10.0  ", want: Version{10, 0, 0}},
			{in: "", want: Version{}},
			{in: "0", want: Version{}},
			{in: "undefined", want: Version{}},
			{in: "65535", want: Version{65535, 0, 0}},
			{in: "abc", wantErr: true},
			{in: "v", wantErr: true},
			{in: "10..0", wantErr: true},
			{in: "65536", wantErr: true},
			{in: "10.0.65536", wantErr: true},
		}

		for _, tt := range tests {
			t.Run(fmt.Sprintf("%q", tt.in), func(t *testing.T) {
				got, err := ParseVersion(tt.in)
				if tt.wantErr {
					assert.ErrorIs(t, err, ErrInvalidParseVersion)
					return
				}
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			})
		}
	})
}
