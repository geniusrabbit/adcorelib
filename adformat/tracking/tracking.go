// Package tracking is a reference adformat.Binding adapter (§3.2.16) for
// protocol-agnostic fire-and-forget impression/click tracking pixels,
// applicable to any Format.Kind (native/banner/push/direct). It carries
// no new type in adformat's core: a tracking pixel is just a
// Field{Type: FieldURLType}, optionally repeated via
// Field.MinCount/MaxCount (§3.2.15) for several simultaneous
// verification vendors, bound to a specific event via this package's
// Tracker Binding. See adformat/vasttracking for the VAST-specific
// event dictionary used by in-stream video creatives.
package tracking

import "github.com/geniusrabbit/adcorelib/adformat"

// Target identifies this adapter's Binding.Target value.
const Target = "tracking"

// Tracking events.
const (
	EventImpression = "impression"
	EventClick      = "click"
)

// Params is the Binding.Params shape for a tracking pixel Field.
type Params struct {
	Event string `json:"event"`
}

// Tracker builds a Binding associating a Field with a tracking event.
func Tracker(event string) adformat.Binding {
	return adformat.NewBinding(Target, Params{Event: event})
}

// Get decodes the tracking binding of f, if any.
func Get(f adformat.Field) (Params, bool) {
	return adformat.DecodeBindingParams[Params](f.Bindings, Target)
}
