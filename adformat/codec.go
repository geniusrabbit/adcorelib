package adformat

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// LoadYAML decodes a single Format from YAML data and validates it
// (§3.2.6). Since Node/Config use plain yaml tags with no custom
// (Un)Marshal code, the YAML shape is a direct mirror of the Go
// structures — see node.go.
func LoadYAML(data []byte) (*Format, error) {
	var f Format
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("adformat: yaml decode: %w", err)
	}
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return &f, nil
}

// LoadJSON decodes a single Format from JSON data and validates it.
func LoadJSON(data []byte) (*Format, error) {
	var f Format
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("adformat: json decode: %w", err)
	}
	if err := f.Validate(); err != nil {
		return nil, err
	}
	return &f, nil
}

// DecodeAny decodes a Format from data, choosing the JSON or YAML decoder
// based on the file extension of name (".json" vs anything else, treated
// as YAML — a superset of JSON). Useful together with embed.FS, where the
// file name/extension is already known.
func DecodeAny(name string, data []byte) (*Format, error) {
	if strings.ToLower(filepath.Ext(name)) == ".json" {
		return LoadJSON(data)
	}
	return LoadYAML(data)
}

// ReadYAML reads and decodes a single Format from r.
func ReadYAML(r io.Reader) (*Format, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("adformat: read yaml: %w", err)
	}
	return LoadYAML(data)
}

// ReadJSON reads and decodes a single Format from r.
func ReadJSON(r io.Reader) (*Format, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("adformat: read json: %w", err)
	}
	return LoadJSON(data)
}
