package adformat

// Param is a fixed, author-supplied value embedded directly in the
// Config tree — unlike Field, never shown to or editable by whoever
// fills in the ad creative. Used to carry hidden/internal context
// alongside the visible structure: a fixed template/layout id, a
// partner-specific tag, a hardcoded default consumed by the rendering
// engine, or any other constant that needs to travel WITH a specific
// group/screen (so nested Param values inside different named
// groups/screens can differ) rather than living as a single
// Format-level constant (§3.2.14).
type Param struct {
	Name      string     `json:"name" yaml:"name"`
	Value     any        `json:"value" yaml:"value"`
	Condition *Condition `json:"condition,omitempty" yaml:"condition,omitempty"`
}

// GetName of the param.
func (p Param) GetName() string {
	return p.Name
}
