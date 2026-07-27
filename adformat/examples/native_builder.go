package examples

import (
	"github.com/geniusrabbit/adcorelib/adformat"
	"github.com/geniusrabbit/adcorelib/adformat/openrtbnative"
)

// nativeConfig is the Go builder (§3.2.10) equivalent of native.yaml —
// same fields, same Condition (phone active only for adtype=business),
// same openrtb_native Bindings, defined as a plain Go value instead of
// parsed YAML. WithAssets/WithFields are methods, not direct
// Assets:/Fields: struct literals — Config no longer has those, see
// §3.2.7 — so wrapping every asset/field into a Node is entirely hidden
// from this code.
var nativeConfig = adformat.Config{}.
	WithAssets(
		adformat.Asset("main").Require().Size(50, 1500, 50, 1500).Types("image", "video").Focal(0.5, 0.5).
			Bind(openrtbnative.ImageAsset(1, openrtbnative.ImageTypeMain)),
		adformat.Asset("logo").Size(50, 100, 50, 100).Types("image").
			Bind(openrtbnative.ImageAsset(2, openrtbnative.ImageTypeLogo)),
	).
	WithFields(
		adformat.StringField("title").Require().Len(5, 40).WithTitle("Title").
			Bind(openrtbnative.TitleField(3)),
		adformat.StringField("description").Require().Len(5, 80).WithTitle("Description").
			Bind(openrtbnative.DataField(4, openrtbnative.DataTypeDesc)),
		adformat.StringField("adtype").WithTitle("Advertiser type").WithOptions(
			adformat.Opt("business", "Business"), adformat.Opt("game", "Game"), adformat.Opt("app", "App"),
		),
		adformat.PhoneField("phone").WithTitle("Phone").When(adformat.FieldIn("adtype", "business")).
			Bind(openrtbnative.DataField(6, openrtbnative.DataTypePhone)),
		adformat.URLField("url").WithTitle("Promotion URL").
			Bind(openrtbnative.DataField(7, openrtbnative.DataTypeDisplayURL)),
	).
	// Nested group inside a nested group — a genuine demonstration of
	// the Node recursion (§3.2.7), not present in native.yaml itself: a
	// hypothetical multi-step extension of the native format, 1-5
	// repeatable "step" screens, each with an optional "cta" sub-group.
	WithGroups(
		adformat.Group("step").Count(1, 5).WithFields(
			adformat.StringField("headline").Require().Len(0, 40),
		).WithGroups(
			adformat.Group("cta").Count(0, 1).WithFields(
				adformat.URLField("url").Require(),
			),
		),
	)

// NativeAd is the Go var-table equivalent of native.yaml (§3.2.10).
var NativeAd = adformat.Format{
	Codename: "native",
	Title:    "Native Ad",
	Kind:     adformat.KindNative,
	Config:   &nativeConfig,
}
