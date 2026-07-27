package adformat

import (
	"strings"
	"testing"
)

const yamlSample = `
id: 1
codename: native
title: "Native Ad"
kind: native
config:
  children:
    - kind: asset
      asset:
        name: main
        required: true
    - kind: field
      field:
        name: title
        type: string
        required: true
        max: 40
`

const jsonSample = `{
  "id": 1,
  "codename": "native",
  "kind": "native",
  "config": {
    "children": [
      {"kind": "asset", "asset": {"name": "main", "required": true}},
      {"kind": "field", "field": {"name": "title", "type": "string", "required": true, "max": 40}}
    ]
  }
}`

func TestLoadYAML(t *testing.T) {
	f, err := LoadYAML([]byte(yamlSample))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Codename != "native" || f.Kind != KindNative {
		t.Errorf("unexpected format: %+v", f)
	}
	if len(f.Config.Assets()) != 1 || len(f.Config.Fields()) != 1 {
		t.Errorf("unexpected config shape: %+v", f.Config)
	}
}

func TestLoadJSON(t *testing.T) {
	f, err := LoadJSON([]byte(jsonSample))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Codename != "native" || f.Kind != KindNative {
		t.Errorf("unexpected format: %+v", f)
	}
	if len(f.Config.Assets()) != 1 || len(f.Config.Fields()) != 1 {
		t.Errorf("unexpected config shape: %+v", f.Config)
	}
}

func TestDecodeAny(t *testing.T) {
	f, err := DecodeAny("native.json", []byte(jsonSample))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Codename != "native" {
		t.Errorf("unexpected codename: %s", f.Codename)
	}

	f, err = DecodeAny("native.yaml", []byte(yamlSample))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Codename != "native" {
		t.Errorf("unexpected codename: %s", f.Codename)
	}
}

func TestLoadYAMLPropagatesValidationErrors(t *testing.T) {
	bad := `
codename: broken
config:
  children:
    - kind: field
      field:
        name: zip
        regexp: "("
`
	if _, err := LoadYAML([]byte(bad)); err == nil || !strings.Contains(err.Error(), "invalid regexp") {
		t.Errorf("expected a validation error to propagate, got %v", err)
	}
}
