package adformat

// Kind is the primary media/content type of an ad format (banner, video,
// native, ...). It is a plain, open-ended string instead of a bitmask (the
// old admodels/types.FormatTypeBitset was capped at 32 flags and mixed two
// independent axes — media type and display surface/behavior) so new
// values can be added without any change to this package. The secondary,
// open-ended axis (display surface/behavior: interstitial, dooh, tv,
// ingame, vr, preroll, push-notification, product-card, mraid, playable,
// interactive-end-card, expandable, rewarded, ...) lives in Format.Tags.
type Kind string

// Kind values covering all media types of the old model.
const (
	KindDirect Kind = "direct" // popunder / proxy-click
	KindBanner Kind = "banner"
	KindVideo  Kind = "video"
	KindAudio  Kind = "audio"
	KindNative Kind = "native"
	KindPush   Kind = "push"
	KindCustom Kind = "custom"
)

// String implements fmt.Stringer.
func (k Kind) String() string {
	return string(k)
}
