package adformat

// NodeKind discriminates which field of Node is populated.
type NodeKind string

// Node kinds.
const (
	NodeAsset NodeKind = "asset"
	NodeField NodeKind = "field"
	NodeGroup NodeKind = "group"
	NodeParam NodeKind = "param" // §3.2.14 — hidden, non-editable value
)

// Node is the single recursive element of the format-structure tree.
// Exactly one of Asset/Field/Group/Param is set, matching Kind. It
// replaces three parallel typed slices (Assets/Fields/Groups) that used
// to live directly on Config — see §3.2.7 for the rationale. Node/Config
// use plain json/yaml tags (no custom (Un)Marshal): the wire format is a
// direct mirror of this Go structure — a "children" list where every
// element carries an explicit "kind" plus a nested "asset"/"field"/
// "group"/"param" object.
//
// Application code is not expected to build Node values by hand — see
// builder.go, which exposes Asset()/StringField()/Group()/Param()
// constructors plus Config.WithAssets/WithFields/WithGroups/WithParams
// that wrap values into Node internally.
type Node struct {
	Kind  NodeKind          `json:"kind" yaml:"kind"`
	Asset *AssetRequirement `json:"asset,omitempty" yaml:"asset,omitempty"`
	Field *Field            `json:"field,omitempty" yaml:"field,omitempty"`
	Group *Config           `json:"group,omitempty" yaml:"group,omitempty"` // recursion
	Param *Param            `json:"param,omitempty" yaml:"param,omitempty"`
}

// AssetNode wraps an AssetRequirement into a Node — a small helper for
// the rare bare Go literal that does not go through builder.go.
func AssetNode(a AssetRequirement) Node {
	return Node{Kind: NodeAsset, Asset: &a}
}

// FieldNode wraps a Field into a Node.
func FieldNode(f Field) Node {
	return Node{Kind: NodeField, Field: &f}
}

// GroupNode wraps a Config into a Node.
func GroupNode(c Config) Node {
	return Node{Kind: NodeGroup, Group: &c}
}

// ParamNode wraps a Param into a Node.
func ParamNode(p Param) Node {
	return Node{Kind: NodeParam, Param: &p}
}

// Name is a read-only proxy that returns the name of whichever
// spec (Asset/Field/Group/Param) is set — a single point of access for
// generic engine code (Validate, Walk, ActiveNodes) that does not need to
// know Kind to read this common attribute. Returns "" for a zero Node.
func (n Node) Name() string {
	switch n.Kind {
	case NodeAsset:
		if n.Asset != nil {
			return n.Asset.GetName()
		}
	case NodeField:
		if n.Field != nil {
			return n.Field.GetName()
		}
	case NodeGroup:
		if n.Group != nil {
			return n.Group.Name
		}
	case NodeParam:
		if n.Param != nil {
			return n.Param.GetName()
		}
	}
	return ""
}

// Condition is a read-only proxy returning the Condition of whichever
// spec is set, or nil if none/not applicable.
func (n Node) Condition() *Condition {
	switch n.Kind {
	case NodeAsset:
		if n.Asset != nil {
			return n.Asset.Condition
		}
	case NodeField:
		if n.Field != nil {
			return n.Field.Condition
		}
	case NodeGroup:
		if n.Group != nil {
			return n.Group.Condition
		}
	case NodeParam:
		if n.Param != nil {
			return n.Param.Condition
		}
	}
	return nil
}

// NavigateTo is a read-only proxy returning the fixed screen-navigation
// target of whichever spec is set (only meaningful for asset/field,
// §3.2.11) — empty string means "no navigation".
func (n Node) NavigateTo() string {
	switch n.Kind {
	case NodeAsset:
		if n.Asset != nil {
			return n.Asset.NavigateTo
		}
	case NodeField:
		if n.Field != nil {
			return n.Field.NavigateTo
		}
	}
	return ""
}

// IsActive reports whether the node is active given the current set of
// field/asset values — true when the node carries no Condition, or its
// Condition evaluates to true against values.
func (n Node) IsActive(values map[string]any) bool {
	c := n.Condition()
	return c == nil || c.Eval(values)
}
