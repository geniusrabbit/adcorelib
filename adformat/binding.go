package adformat

import "encoding/json"

// Binding associates a Node (asset/field) or a whole Format with a concrete
// element of some external target representation — a bidding protocol's
// markup (OpenRTB Native asset, OpenRTB Video/VAST object, MRAID feature,
// a specific partner's proprietary JSON), by an opaque target-defined key.
// adformat itself has zero knowledge of what "openrtb_native" or "vast"
// mean — that semantics lives entirely in a small adapter package per
// target (see adformat/openrtbnative, adformat/openrtbvideo,
// adformat/tracking, adformat/vasttracking). One asset/field can carry
// several Bindings for several targets at once — e.g. the same "main"
// image asset is simultaneously the OpenRTB Native "main image" AND, for a
// partner that doesn't speak OpenRTB, a specific field path in their
// proprietary JSON.
type Binding struct {
	Target string         `json:"target" yaml:"target"`             // "openrtb_native", "vast", "mraid", or a partner-specific name
	Params map[string]any `json:"params,omitempty" yaml:"params,omitempty"` // shape defined by the adapter package for Target
}

// Get looks up the first binding matching target and decodes its Params
// into out (a pointer) via encoding/json round-trip. Adapter packages use
// this helper to implement their own typed Get* functions. Returns false
// when no binding for target is present.
func (b Binding) hasTarget(target string) bool {
	return b.Target == target
}

// FindBinding returns the first binding in bindings matching target.
func FindBinding(bindings []Binding, target string) (Binding, bool) {
	for _, b := range bindings {
		if b.hasTarget(target) {
			return b, true
		}
	}
	return Binding{}, false
}

// DecodeBindingParams finds the first binding in bindings matching target
// and decodes its Params into a value of type P via an encoding/json
// round-trip. It is the shared implementation behind every adapter
// package's own typed Get*/GetImage/GetData/... helpers (see
// adformat/openrtbnative, adformat/openrtbvideo, adformat/tracking,
// adformat/vasttracking), so none of them need to touch map[string]any by
// hand. Returns false when no binding for target is present or decoding
// fails.
func DecodeBindingParams[P any](bindings []Binding, target string) (P, bool) {
	var out P
	b, ok := FindBinding(bindings, target)
	if !ok {
		return out, false
	}
	raw, err := json.Marshal(b.Params)
	if err != nil {
		return out, false
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		return out, false
	}
	return out, true
}

// NewBinding builds a Binding for target from an arbitrary params value by
// round-tripping it through encoding/json into map[string]any — used by
// adapter packages' constructors (ImageAsset, TitleField, Video,
// Tracker, ...).
func NewBinding(target string, params any) Binding {
	raw, err := json.Marshal(params)
	if err != nil {
		return Binding{Target: target}
	}
	var m map[string]any
	_ = json.Unmarshal(raw, &m)
	return Binding{Target: target, Params: m}
}
