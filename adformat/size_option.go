package adformat

// SizeOption is one selectable overall display-area size for a Format
// (the ad unit's outer container) — not any individual asset's file
// dimensions, which stay on AssetRequirement (see MinWidth/MaxWidth
// there). Generalizes: a single classic fixed-size banner (one
// exact-size SizeOption), an IAB "multi-size" ad unit (several
// exact-size options), and a stretchable iframe/native container (a
// Flexible option with its own optional bounds) — under one mechanism,
// for any Kind, not only banner (§3.2.13).
type SizeOption struct {
	Title string `json:"title,omitempty" yaml:"title,omitempty"` // optional, for a UI label ("Medium Rectangle")

	// Exact size, when Flexible == false (the common case — a classic
	// fixed-size banner, a standard IAB size, a fixed native strip).
	Width  int `json:"width,omitempty" yaml:"width,omitempty"`
	Height int `json:"height,omitempty" yaml:"height,omitempty"`

	// Flexible: the concrete size is picked not by the format's author
	// but by whoever creates a particular creative from this format
	// (same principle as FocalPoint — the concrete value is data of a
	// particular creative, not part of the format schema). The Min/Max*
	// fields below are optional bounds on that free choice (needed for
	// an iframe banner, which still must match a reasonable range of
	// third-party slots rather than being fully unbounded); unset = no
	// bound.
	Flexible  bool `json:"flexible,omitempty" yaml:"flexible,omitempty"`
	MinWidth  int  `json:"min_width,omitempty" yaml:"min_width,omitempty"`
	MaxWidth  int  `json:"max_width,omitempty" yaml:"max_width,omitempty"`
	MinHeight int  `json:"min_height,omitempty" yaml:"min_height,omitempty"`
	MaxHeight int  `json:"max_height,omitempty" yaml:"max_height,omitempty"`

	// Default marks this option as pre-selected when there is more than
	// one. Not needed when Sizes has exactly one element — the only
	// option is unambiguously the default without an explicit flag.
	Default bool `json:"default,omitempty" yaml:"default,omitempty"`
}

// Suits reports whether the given width/height matches this size option —
// an exact match for a non-flexible option, or falling within
// Min/MaxWidth/Height for a flexible one (unset bounds = unbounded on
// that side).
func (s SizeOption) Suits(width, height int) bool {
	if !s.Flexible {
		return s.Width == width && s.Height == height
	}
	if s.MinWidth > 0 && width < s.MinWidth {
		return false
	}
	if s.MaxWidth > 0 && width > s.MaxWidth {
		return false
	}
	if s.MinHeight > 0 && height < s.MinHeight {
		return false
	}
	if s.MaxHeight > 0 && height > s.MaxHeight {
		return false
	}
	return true
}

// IsFixed reports whether this option describes a single exact size (not
// Flexible, with both Width/Height set).
func (s SizeOption) IsFixed() bool {
	return !s.Flexible && s.Width > 0 && s.Height > 0
}
