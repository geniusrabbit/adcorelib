package adformat

import "strconv"

// Orientation is a supported screen orientation for a creative (§3.2.2).
type Orientation string

// Orientation values.
const (
	OrientationPortrait  Orientation = "portrait"
	OrientationLandscape Orientation = "landscape"
)

// Format describes the creative that one advertiser prepares and uploads
// for one ad (was admodels/types.Format) — see §3.0 for the exact scope
// boundary (one creative of one advertiser, never a placement-level
// composition of several ads).
type Format struct {
	ID       uint64 `json:"id,omitempty" yaml:"id,omitempty"`
	Codename string `json:"codename,omitempty" yaml:"codename,omitempty"`
	Title    string `json:"title,omitempty" yaml:"title,omitempty"`

	// Kind is the primary media/content type; Tags is the open-ended
	// list of secondary classifiers (§3.2, п.1).
	Kind Kind     `json:"kind,omitempty" yaml:"kind,omitempty"`
	Tags []string `json:"tags,omitempty" yaml:"tags,omitempty"`

	// Orientations supported by this creative as a whole — empty means
	// "any" (§3.2.2).
	Orientations []Orientation `json:"orientations,omitempty" yaml:"orientations,omitempty"`

	// Version of the schema this Format was authored against — empty is
	// treated as "1" (§3.2.9). Purely informational: it does not affect
	// current validation, but gives future incompatible changes a place
	// to branch on instead of guessing from field presence.
	Version string `json:"version,omitempty" yaml:"version,omitempty"`

	// Bindings associate the whole format (not any single asset/field)
	// with an external protocol representation — e.g. OpenRTB Video
	// Protocols/APIs/Linearity/PlaybackMethod, which have no
	// protocol-agnostic equivalent on any single AssetRequirement
	// (§3.2.12).
	Bindings []Binding `json:"bindings,omitempty" yaml:"bindings,omitempty"`

	// Sizes is the multiselect of overall display-area sizes, replacing
	// the old single Width/Height/MinWidth/MinHeight pair (§3.2.13).
	Sizes []SizeOption `json:"sizes,omitempty" yaml:"sizes,omitempty"`

	// Config is the recursive structure of assets/fields/groups/params
	// that make up this creative.
	Config *Config `json:"config,omitempty" yaml:"config,omitempty"`
}

// GetConfig of the format object, safe on a nil Format.
func (f *Format) GetConfig() *Config {
	if f == nil {
		return nil
	}
	return f.Config
}

// String representation of the format, preferring its default size when
// one is set.
func (f Format) String() string {
	if size, ok := f.DefaultSize(); ok && (size.Width > 0 || size.Height > 0) {
		return sizeString(size)
	}
	if f.Codename != "" {
		return f.Codename
	}
	return string(f.Kind)
}

func sizeString(s SizeOption) string {
	switch {
	case s.Width > 0 && s.Height > 0:
		return strconv.Itoa(s.Width) + "x" + strconv.Itoa(s.Height)
	case s.Flexible:
		return "flexible"
	default:
		return "stretch"
	}
}

// Suits reports whether width/height matches at least one of f.Sizes. A
// Format with no Sizes at all (e.g. direct/popunder, native "strip" with
// no fixed container) suits any size — it has no display-area constraint
// to check.
func (f Format) Suits(width, height int) bool {
	if len(f.Sizes) == 0 {
		return true
	}
	for _, s := range f.Sizes {
		if s.Suits(width, height) {
			return true
		}
	}
	return false
}

// DefaultSize returns the SizeOption explicitly marked Default; if none
// is marked and there is exactly one option, that one is the default;
// otherwise the first option is returned. ok is false only when Sizes is
// empty.
func (f Format) DefaultSize() (SizeOption, bool) {
	if len(f.Sizes) == 0 {
		return SizeOption{}, false
	}
	if len(f.Sizes) == 1 {
		return f.Sizes[0], true
	}
	for _, s := range f.Sizes {
		if s.Default {
			return s, true
		}
	}
	return f.Sizes[0], true
}

// IsFixedSize reports whether every SizeOption is a plain exact size
// (none Flexible) — helps a UI decide whether a size picker is needed at
// all.
func (f Format) IsFixedSize() bool {
	if len(f.Sizes) == 0 {
		return false
	}
	for _, s := range f.Sizes {
		if !s.IsFixed() {
			return false
		}
	}
	return true
}

// IsDirect format kind (popunder/proxy-click).
func (f Format) IsDirect() bool {
	return f.Kind == KindDirect
}

// IsBanner format kind.
func (f Format) IsBanner() bool {
	return f.Kind == KindBanner
}

// IsVideo format kind.
func (f Format) IsVideo() bool {
	return f.Kind == KindVideo
}

// IsAudio format kind.
func (f Format) IsAudio() bool {
	return f.Kind == KindAudio
}

// IsNative format kind.
func (f Format) IsNative() bool {
	return f.Kind == KindNative
}

// IsPush format kind.
func (f Format) IsPush() bool {
	return f.Kind == KindPush
}

// HasTag reports whether the format carries the given tag.
func (f Format) HasTag(tag string) bool {
	for _, t := range f.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// SupportsOrientation reports whether o is supported — an empty
// Orientations list means "supports any orientation" (§3.2.2).
func (f Format) SupportsOrientation(o Orientation) bool {
	if len(f.Orientations) == 0 {
		return true
	}
	for _, ori := range f.Orientations {
		if ori == o {
			return true
		}
	}
	return false
}

// Validate the whole format: its Config tree (§3.2.6) plus every
// Format-level SizeOption bound.
func (f Format) Validate() error {
	for _, s := range f.Sizes {
		if err := validateSizeOption(s); err != nil {
			return err
		}
	}
	if f.Config != nil {
		return f.Config.Validate()
	}
	return nil
}

func validateSizeOption(s SizeOption) error {
	if !s.Flexible {
		return nil
	}
	if s.MinWidth > 0 && s.MaxWidth > 0 && s.MinWidth > s.MaxWidth {
		return errInvalidSizeOption("min_width is greater than max_width")
	}
	if s.MinHeight > 0 && s.MaxHeight > 0 && s.MinHeight > s.MaxHeight {
		return errInvalidSizeOption("min_height is greater than max_height")
	}
	return nil
}

func errInvalidSizeOption(msg string) error {
	return &sizeOptionError{msg: msg}
}

type sizeOptionError struct{ msg string }

func (e *sizeOptionError) Error() string { return "adformat: invalid size option: " + e.msg }
