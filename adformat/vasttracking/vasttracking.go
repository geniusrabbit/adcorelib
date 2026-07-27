// Package vasttracking is a reference adformat.Binding adapter (§3.2.16)
// for the full VAST 4.x tracking-event dictionary, used by in-stream
// video creatives (Format.Kind == KindVideo). Like adformat/tracking, it
// carries no new type in adformat's core — a VAST tracking pixel is a
// Field{Type: FieldURLType}, typically repeated via
// Field.MinCount/MaxCount (§3.2.15) so several verification vendors
// (IAS, DoubleVerify, Moat, ...) can fire on the same event
// independently, bound to a specific VAST event via this package's
// Tracker Binding.
package vasttracking

import "github.com/geniusrabbit/adcorelib/adformat"

// Target identifies this adapter's Binding.Target value.
const Target = "vast_tracking"

// VAST 4.x tracking events.
const (
	EventImpression    = "impression"
	EventStart         = "start"
	EventFirstQuartile = "firstQuartile"
	EventMidpoint      = "midpoint"
	EventThirdQuartile = "thirdQuartile"
	EventComplete      = "complete"
	EventPause         = "pause"
	EventResume        = "resume"
	EventMute          = "mute"
	EventUnmute        = "unmute"
	EventFullscreen    = "fullscreen"
	EventSkip          = "skip"
	EventClickThrough  = "clickThrough"
	EventClickTracking = "clickTracking"
)

// Params is the Binding.Params shape for a VAST tracking pixel Field.
type Params struct {
	Event string `json:"event"`
}

// Tracker builds a Binding associating a Field with a VAST tracking
// event.
func Tracker(event string) adformat.Binding {
	return adformat.NewBinding(Target, Params{Event: event})
}

// Get decodes the vast_tracking binding of f, if any.
func Get(f adformat.Field) (Params, bool) {
	return adformat.DecodeBindingParams[Params](f.Bindings, Target)
}
