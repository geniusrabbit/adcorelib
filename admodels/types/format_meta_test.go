package types

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMeta(t *testing.T) {
	data, err := fileData("../assets/format.native.json")
	if err != nil {
		t.Error(err)
	}

	var config FormatConfig
	if err = json.Unmarshal(data, &config); err != nil {
		t.Error(err)
	}

	if len(config.Assets) != 2 {
		t.Errorf("Invalid count of assets: %d", len(config.Assets))
	}

	if len(config.Fields) != 5 {
		t.Errorf("Invalid count of fields: %d", len(config.Fields))
	}
}

func TestConfigIntersec(t *testing.T) {
	var tests = []struct {
		name   string
		c1     *FormatConfig
		c2     *FormatConfig
		result bool
	}{
		{
			name:   "empty",
			c1:     &FormatConfig{},
			c2:     &FormatConfig{},
			result: true,
		},
		{
			name: "basic_stretch",
			c1: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Required: true},
				},
			},
			c2: &FormatConfig{
				Assets: []FormatFileRequirement{{}},
			},
			result: true,
		},
		{
			name: "basic_assets_similar_sizes",
			c1: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Required: true, Width: 1000, Height: 1000, MinHeight: 100, MinWidth: 100},
					{Required: false, Name: "icon", Width: 80, Height: 80},
				},
			},
			c2: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Required: true, Width: 1100, Height: 900, MinHeight: 200, MinWidth: 150},
					{Required: true, Name: "icon", Width: 100, Height: 100, MinHeight: 50, MinWidth: 50},
				},
			},
			result: true,
		},
		{
			name: "basic_assets_similar_sizes_negative",
			c1: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Required: true, Width: 1000, Height: 1000, MinHeight: 100, MinWidth: 100},
					{Required: true, Name: "icon", Width: 100, Height: 100, MinHeight: 50, MinWidth: 50},
				},
			},
			c2: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Required: true, Width: 1100, Height: 900, MinHeight: 200, MinWidth: 150},
					{Required: false, Name: "icon", Width: 80, Height: 80},
				},
			},
			result: false,
		},
		{
			name: "basic_ext",
			c1: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Required: true, Name: "main", Width: 300, Height: 100},
				},
				Fields: []FormatField{
					{Required: true, Name: "title"},
				},
			},
			c2: &FormatConfig{
				Assets: []FormatFileRequirement{
					{Width: 300, Height: 100},
					{Required: false, Name: "icon", Width: 30, Height: 30},
				},
				Fields: []FormatField{
					{Required: false, Name: "title"},
				},
			},
			result: true,
		},
		{
			name: "basic_field_negative",
			c1: &FormatConfig{
				Fields: []FormatField{
					{Required: true, Name: "title"},
				},
			},
			c2:     &FormatConfig{},
			result: false,
		},
		{
			name: "basic_field_negative2",
			c1: &FormatConfig{
				Fields: []FormatField{
					{Required: true, Name: "title"},
				},
			},
			c2: &FormatConfig{
				Fields: []FormatField{
					{Required: true, Name: "title"},
					{Required: true, Name: "icon"},
				},
			},
			result: false,
		},
		{
			name: "basic_field_positive",
			c1: &FormatConfig{
				Fields: []FormatField{
					{Required: true, Name: "title"},
					{Required: false, Name: "icon"},
				},
			},
			c2: &FormatConfig{
				Fields: []FormatField{
					{Required: true, Name: "title"},
					{Required: true, Name: "icon"},
				},
			},
			result: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.c1.Intersec(test.c2) != test.result {
				t.Errorf("Intersec must be %t", test.result)
			}
		})
	}
}

func TestFormatFieldUIMetadataJSON(t *testing.T) {
	t.Run("omitted editable defaults to shown", func(t *testing.T) {
		var f FormatField
		if err := json.Unmarshal([]byte(`{"name":"title","title":"Title"}`), &f); err != nil {
			t.Fatal(err)
		}
		if !f.IsEditable() {
			t.Error("omitted editable must default to true")
		}
		if f.IsMultilang() {
			t.Error("omitted multilang must default to false")
		}
		if f.MultilineRows() != 0 {
			t.Errorf("omitted multiline must be 0, got %d", f.MultilineRows())
		}
	})

	t.Run("editable false hides the field", func(t *testing.T) {
		var f FormatField
		if err := json.Unmarshal([]byte(`{"name":"internal","editable":false}`), &f); err != nil {
			t.Fatal(err)
		}
		if f.IsEditable() {
			t.Error("editable:false must hide the field")
		}
	})

	t.Run("round-trip multiline description multilang", func(t *testing.T) {
		in := FormatField{
			Name:        "description",
			Title:       "Description",
			Description: "Body text shown with the ad",
			Multiline:   3,
			Multilang:   true,
		}
		data, err := json.Marshal(in)
		if err != nil {
			t.Fatal(err)
		}
		var out FormatField
		if err := json.Unmarshal(data, &out); err != nil {
			t.Fatal(err)
		}
		if out.Description != in.Description || out.Multiline != 3 || !out.IsMultilang() || !out.IsEditable() {
			t.Errorf("round-trip mismatch: %+v", out)
		}
		if strings.Contains(string(data), `"editable"`) {
			t.Errorf("unset editable should be omitted, got %s", data)
		}
		var fromNull FormatField
		if err := json.Unmarshal([]byte(`{"name":"title","editable":null}`), &fromNull); err != nil {
			t.Fatal(err)
		}
		if !fromNull.IsEditable() {
			t.Error("editable:null must default to shown")
		}
	})
}
