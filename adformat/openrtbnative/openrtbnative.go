// Package openrtbnative is a reference adformat.Binding adapter
// (§3.2.12) associating adformat assets/fields with OpenRTB Native 1.2
// asset objects — a concrete, typed replacement for the name-based
// guessing that adsource-openrtb currently does (see
// adsource-openrtb/request/v2/builder.go), and for the always-empty
// FormatFileRequirement.ID/FormatField.ID of the old model.
//
// adformat's core package has zero knowledge of this package or of
// OpenRTB at all — Binding.Target/Params is a generic, opaque
// key/value pair, and everything OpenRTB-specific (constants, typed
// params, getters) lives here instead. Constants below are declared
// locally per the IAB Native Ads 1.2 spec, not reused from
// github.com/bsm/openrtb, so adcorelib/adformat does not gain an
// external dependency just to describe a format.
package openrtbnative

import (
	"encoding/json"
	"fmt"

	"github.com/geniusrabbit/adcorelib/adformat"
)

// Target identifies this adapter's Binding.Target value.
const Target = "openrtb_native"

// Native asset image types (IAB Native Ads 1.2, Image Object "type").
const (
	ImageTypeIcon = 1
	ImageTypeLogo = 2
	ImageTypeMain = 3
)

// Native asset data types (IAB Native Ads 1.2, Data Object "type").
const (
	DataTypeSponsored      = 1
	DataTypeDesc           = 2
	DataTypeRating         = 3
	DataTypeLikes          = 4
	DataTypeDownloads      = 5
	DataTypePrice          = 6
	DataTypeSalePrice      = 7
	DataTypePhone          = 8
	DataTypeAddress        = 9
	DataTypeDescAdditional = 10
	DataTypeDisplayURL     = 11
	DataTypeCTAText        = 12
)

// ImageParams is the Binding.Params shape for an image asset.
type ImageParams struct {
	ID        int `json:"id"`         // stable, format-author-assigned asset id (passed 1:1 into the OpenRTB request/response)
	ImageType int `json:"image_type"` // ImageTypeIcon/Logo/Main
}

// DataParams is the Binding.Params shape for a data asset (a Field bound
// to a non-title OpenRTB Native data object).
type DataParams struct {
	ID       int `json:"id"`
	DataType int `json:"data_type"` // DataTypeDesc/Rating/.../CTAText
}

// TitleParams is the Binding.Params shape for the title asset. The
// title's length limit is intentionally not duplicated here — it is
// already Field.Max (a general, protocol-agnostic property, see
// adformat.Field), and an adapter building an OpenRTB request reads it
// from there instead, so the same constraint cannot drift between the
// two places.
type TitleParams struct {
	ID int `json:"id"`
}

// VideoParams is the Binding.Params shape for a video asset nested
// inside a native ad.
type VideoParams struct {
	ID        int   `json:"id"`
	Protocols []int `json:"protocols,omitempty"` // VAST versions (OpenRTB Video.Protocols)
	APIs      []int `json:"apis,omitempty"`      // MRAID/VPAID versions (OpenRTB Video.API)
}

// ImageAsset builds a Binding associating an AssetRequirement with an
// OpenRTB Native image asset.
func ImageAsset(id, imageType int) adformat.Binding {
	return adformat.NewBinding(Target, ImageParams{ID: id, ImageType: imageType})
}

// DataField builds a Binding associating a Field with an OpenRTB Native
// data asset.
func DataField(id, dataType int) adformat.Binding {
	return adformat.NewBinding(Target, DataParams{ID: id, DataType: dataType})
}

// TitleField builds a Binding associating a Field with the OpenRTB
// Native title asset.
func TitleField(id int) adformat.Binding {
	return adformat.NewBinding(Target, TitleParams{ID: id})
}

// VideoAsset builds a Binding associating an AssetRequirement with an
// OpenRTB Native video asset.
func VideoAsset(id int, protocols, apis []int) adformat.Binding {
	return adformat.NewBinding(Target, VideoParams{ID: id, Protocols: protocols, APIs: apis})
}

// GetImage decodes the openrtb_native image binding of a, if any.
func GetImage(a adformat.AssetRequirement) (ImageParams, bool) {
	return adformat.DecodeBindingParams[ImageParams](a.Bindings, Target)
}

// GetData decodes the openrtb_native data binding of f, if any.
func GetData(f adformat.Field) (DataParams, bool) {
	return adformat.DecodeBindingParams[DataParams](f.Bindings, Target)
}

// GetTitle decodes the openrtb_native title binding of f, if any.
func GetTitle(f adformat.Field) (TitleParams, bool) {
	return adformat.DecodeBindingParams[TitleParams](f.Bindings, Target)
}

// GetVideo decodes the openrtb_native video binding of a, if any.
func GetVideo(a adformat.AssetRequirement) (VideoParams, bool) {
	return adformat.DecodeBindingParams[VideoParams](a.Bindings, Target)
}

type idOnlyParams struct {
	ID int `json:"id"`
}

// Validate checks that every openrtb_native binding in cfg has a
// non-zero, unique (within cfg) id — the one cross-Node invariant this
// specific target actually needs; deliberately not part of the generic
// Config.Validate (§3.2.6), which knows nothing about any target.
func Validate(cfg *adformat.Config) error {
	if cfg == nil {
		return nil
	}
	seenBy := map[int]string{}
	return cfg.Walk(func(path string, n adformat.Node) error {
		var bindings []adformat.Binding
		switch n.Kind {
		case adformat.NodeAsset:
			if n.Asset != nil {
				bindings = n.Asset.Bindings
			}
		case adformat.NodeField:
			if n.Field != nil {
				bindings = n.Field.Bindings
			}
		default:
			return nil
		}
		for _, b := range bindings {
			if b.Target != Target {
				continue
			}
			var p idOnlyParams
			raw, err := json.Marshal(b.Params)
			if err != nil {
				return fmt.Errorf("openrtbnative: %s: invalid binding params: %w", path, err)
			}
			if err := json.Unmarshal(raw, &p); err != nil {
				return fmt.Errorf("openrtbnative: %s: invalid binding params: %w", path, err)
			}
			if p.ID == 0 {
				return fmt.Errorf("openrtbnative: %s: binding has no (or zero) id", path)
			}
			if prev, ok := seenBy[p.ID]; ok {
				return fmt.Errorf("openrtbnative: id %d is used by both %q and %q", p.ID, prev, path)
			}
			seenBy[p.ID] = path
		}
		return nil
	})
}
