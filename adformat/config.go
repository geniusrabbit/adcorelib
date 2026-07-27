package adformat

// Config is the recursive structure describing an ad format's creative
// (was admodels/types.FormatConfig). It is built on top of a single
// Children []Node list (assets + fields + groups + params, in the order
// the author declared them) instead of three parallel typed slices — see
// node.go and §3.2.7 for the rationale. There is no separate "Group"
// type: a nested group is simply another Config referenced by a Node of
// Kind NodeGroup.
//
// Name/MinCount/MaxCount/Condition/Entry are meaningless on the root
// Config hanging off Format.Config — they only matter for a Config used
// as a nested group/screen (a Node's Group field).
type Config struct {
	// Name of this group (empty for the root Config).
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Title of this group for the UI.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`

	// MinCount/MaxCount: how many instances of this group one ad
	// creative carries (§3.4/§3.5) — same effective semantics as
	// AssetRequirement/Field (see GetMinCount/GetMaxCount/RangeCount).
	// Only legitimate when every instance is guaranteed to be part of
	// one advertiser's single submission — see §3.0.
	MinCount int `json:"min_count,omitempty" yaml:"min_count,omitempty"`
	MaxCount int `json:"max_count,omitempty" yaml:"max_count,omitempty"`

	// Entry marks this group as the starting screen of a multi-screen
	// creative (§3.2.11). If no nested group is marked Entry, the first
	// one in Children is assumed to be the start screen.
	Entry bool `json:"entry,omitempty" yaml:"entry,omitempty"`

	// Condition controls when this group is active (§3.2.5) — e.g. a
	// nested "panel" that only applies for a certain layout choice.
	Condition *Condition `json:"condition,omitempty" yaml:"condition,omitempty"`

	// Children holds assets, fields, groups and params together, in
	// authoring order.
	Children []Node `json:"children,omitempty" yaml:"children,omitempty"`
}

// IsEmpty config.
func (c Config) IsEmpty() bool {
	return len(c.Children) == 0
}

// GetMinCount effective minimal number of instances of this group.
func (c Config) GetMinCount() int {
	mn, _ := effectiveCount(c.MinCount, c.MaxCount, false)
	return mn
}

// GetMaxCount effective maximal number of instances (-1 = unlimited).
func (c Config) GetMaxCount() int {
	_, mx := effectiveCount(c.MinCount, c.MaxCount, false)
	return mx
}

// RangeCount returns the effective (min, max) instance-count range.
func (c Config) RangeCount() (min, max int) {
	return effectiveCount(c.MinCount, c.MaxCount, false)
}

// IsMultiple reports whether this group may have more than one instance.
func (c Config) IsMultiple() bool {
	return isMultipleCount(c.GetMaxCount())
}

// Assets returns the top-level asset requirements of this Config (does
// not recurse into nested groups).
func (c Config) Assets() []AssetRequirement {
	var list []AssetRequirement
	for _, n := range c.Children {
		if n.Kind == NodeAsset && n.Asset != nil {
			list = append(list, *n.Asset)
		}
	}
	return list
}

// Fields returns the top-level fields of this Config (does not recurse
// into nested groups).
func (c Config) Fields() []Field {
	var list []Field
	for _, n := range c.Children {
		if n.Kind == NodeField && n.Field != nil {
			list = append(list, *n.Field)
		}
	}
	return list
}

// Groups returns the top-level nested groups of this Config.
func (c Config) Groups() []Config {
	var list []Config
	for _, n := range c.Children {
		if n.Kind == NodeGroup && n.Group != nil {
			list = append(list, *n.Group)
		}
	}
	return list
}

// Params returns the top-level params of this Config.
func (c Config) Params() []Param {
	var list []Param
	for _, n := range c.Children {
		if n.Kind == NodeParam && n.Param != nil {
			list = append(list, *n.Param)
		}
	}
	return list
}

// AssetByName from config; an empty name matches the main asset.
func (c Config) AssetByName(name string) *AssetRequirement {
	for _, n := range c.Children {
		if n.Kind != NodeAsset || n.Asset == nil {
			continue
		}
		if n.Asset.Name == name || (name == "" && n.Asset.IsMain()) {
			a := *n.Asset
			return &a
		}
	}
	return nil
}

// MainAsset returns the main asset if it exists.
func (c Config) MainAsset() *AssetRequirement {
	for _, n := range c.Children {
		if n.Kind == NodeAsset && n.Asset != nil && n.Asset.IsMain() {
			a := *n.Asset
			return &a
		}
	}
	return nil
}

// SimpleAsset returns the main asset in case it is the only required
// asset in the config.
func (c Config) SimpleAsset() *AssetRequirement {
	var as *AssetRequirement
	for _, n := range c.Children {
		if n.Kind != NodeAsset || n.Asset == nil {
			continue
		}
		a := *n.Asset
		if a.IsMain() {
			as = &a
		} else if a.IsRequired() {
			return nil
		}
	}
	return as
}

// GetField by name.
func (c Config) GetField(name string) *Field {
	for _, n := range c.Children {
		if n.Kind == NodeField && n.Field != nil && n.Field.Name == name {
			f := *n.Field
			return &f
		}
	}
	return nil
}

// GetParam by name (§3.2.14).
func (c Config) GetParam(name string) (Param, bool) {
	for _, n := range c.Children {
		if n.Kind == NodeParam && n.Param != nil && n.Param.Name == name {
			return *n.Param, true
		}
	}
	return Param{}, false
}

// RequiredField returns a required field, optionally restricted to one of
// the given names — same lookup semantics as the old model.
func (c Config) RequiredField(fields ...string) *Field {
	allFields := c.Fields()
	if len(fields) == 0 {
		for i := range allFields {
			if allFields[i].Required {
				return &allFields[i]
			}
		}
		return nil
	}

	haveRequired := false
	for _, fl := range fields {
		for i := range allFields {
			field := allFields[i]
			if !haveRequired {
				haveRequired = field.Required
			}
			if field.Required && fl == field.Name {
				return &field
			}
		}
		if !haveRequired {
			return nil
		}
	}
	return nil
}

// RequiredFieldExcept returns any required field not present in the
// given exception list.
func (c Config) RequiredFieldExcept(fields ...string) *Field {
	allFields := c.Fields()
	for i := range allFields {
		field := allFields[i]
		if !field.Required {
			continue
		}
		excepted := false
		for _, fl := range fields {
			if fl == field.Name {
				excepted = true
				break
			}
		}
		if !excepted {
			return &field
		}
	}
	return nil
}

// ContainsAsset in the list of top-level assets.
func (c Config) ContainsAsset(asset AssetRequirement, revers ...bool) bool {
	revCheck := len(revers) > 0 && revers[0]
	for _, a := range c.Assets() {
		if revCheck {
			if asset.SoftEqual(a) {
				return true
			}
		} else if a.SoftEqual(asset) {
			return true
		}
	}
	return false
}

// SimilarField in the list of top-level fields.
func (c Config) SimilarField(field Field, revers ...bool) *Field {
	revCheck := len(revers) > 0 && revers[0]
	for _, f := range c.Fields() {
		if revCheck {
			if field.SoftEqual(f) {
				ff := f
				return &ff
			}
		} else if f.SoftEqual(field) {
			ff := f
			return &ff
		}
	}
	return nil
}

// Intersec reports whether c and other are compatible: every required
// top-level asset/field of one config is satisfiable by the other. Only
// considers top-level assets/fields, same as the old flat model — nested
// groups are a new structural feature and are not part of this
// comparison.
func (c Config) Intersec(other Config) bool {
	if c.IsEmpty() && other.IsEmpty() {
		return true
	}

	assets, otherAssets := c.Assets(), other.Assets()
	if len(assets) > 0 && len(otherAssets) > 0 {
		for _, asset := range otherAssets {
			if asset.IsRequired() && !c.ContainsAsset(asset) {
				return false
			}
		}
		for _, asset := range assets {
			if asset.IsRequired() && !other.ContainsAsset(asset, true) {
				return false
			}
		}
	} else if len(assets) > 0 {
		for _, asset := range assets {
			if asset.IsRequired() {
				return false
			}
		}
	}

	fields, otherFields := c.Fields(), other.Fields()
	if len(fields) > 0 && len(otherFields) > 0 {
		for _, field := range otherFields {
			if field.Required && c.SimilarField(field) == nil {
				return false
			}
		}
		for _, field := range fields {
			if field.Required && other.SimilarField(field, true) == nil {
				return false
			}
		}
	} else if len(fields) > 0 {
		for _, field := range fields {
			if field.Required {
				return false
			}
		}
	}
	return true
}

// ActiveAssets returns the top-level assets whose Condition (if any)
// evaluates to true against values (§3.2.5).
func (c Config) ActiveAssets(values map[string]any) []AssetRequirement {
	var list []AssetRequirement
	for _, n := range c.Children {
		if n.Kind == NodeAsset && n.Asset != nil && n.IsActive(values) {
			list = append(list, *n.Asset)
		}
	}
	return list
}

// ActiveFields returns the top-level fields whose Condition (if any)
// evaluates to true against values (§3.2.5).
func (c Config) ActiveFields(values map[string]any) []Field {
	var list []Field
	for _, n := range c.Children {
		if n.Kind == NodeField && n.Field != nil && n.IsActive(values) {
			list = append(list, *n.Field)
		}
	}
	return list
}

// ActiveGroups returns the top-level nested groups whose Condition (if
// any) evaluates to true against values.
func (c Config) ActiveGroups(values map[string]any) []Config {
	var list []Config
	for _, n := range c.Children {
		if n.Kind == NodeGroup && n.Group != nil && n.IsActive(values) {
			list = append(list, *n.Group)
		}
	}
	return list
}

// ActiveParams returns the top-level params whose Condition (if any)
// evaluates to true against values.
func (c Config) ActiveParams(values map[string]any) []Param {
	var list []Param
	for _, n := range c.Children {
		if n.Kind == NodeParam && n.Param != nil && n.IsActive(values) {
			list = append(list, *n.Param)
		}
	}
	return list
}

// EntryGroup returns the nested group explicitly marked Entry, or,
// failing that, the first nested group in Children — the default
// "convention over configuration" start-screen rule from §3.2.11. The
// second return value is false only when there are no nested groups at
// all.
func (c Config) EntryGroup() (Config, bool) {
	groups := c.Groups()
	if len(groups) == 0 {
		return Config{}, false
	}
	for _, g := range groups {
		if g.Entry {
			return g, true
		}
	}
	return groups[0], true
}
