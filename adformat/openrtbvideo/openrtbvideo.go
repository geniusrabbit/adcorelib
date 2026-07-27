// Package openrtbvideo is a reference adformat.Binding adapter (§3.2.12)
// associating a whole adformat.Format with an OpenRTB Video object's
// properties. It is a second, independent adapter (alongside
// adformat/openrtbnative) demonstrating that adding support for another
// external protocol never requires a change to adcorelib/adformat's core
// — only a new small package like this one.
package openrtbvideo

import "github.com/geniusrabbit/adcorelib/adformat"

// Target identifies this adapter's Binding.Target value.
const Target = "openrtb_video"

// Params carries the OpenRTB Video fields that have no protocol-agnostic
// equivalent on AssetRequirement (duration/allow_sound and similar are
// already covered by the generic model and are intentionally not
// duplicated here).
type Params struct {
	Protocols      []int `json:"protocols,omitempty"`       // VAST versions
	APIs           []int `json:"apis,omitempty"`            // VPAID/MRAID versions
	Linearity      int   `json:"linearity,omitempty"`       // OpenRTB Video.Linearity
	PlaybackMethod []int `json:"playback_method,omitempty"` // OpenRTB Video.PlaybackMethod
}

// Video builds a Binding associating a Format with an OpenRTB Video
// object's properties.
func Video(p Params) adformat.Binding {
	return adformat.NewBinding(Target, p)
}

// Get decodes the openrtb_video binding of f, if any.
func Get(f adformat.Format) (Params, bool) {
	return adformat.DecodeBindingParams[Params](f.Bindings, Target)
}
