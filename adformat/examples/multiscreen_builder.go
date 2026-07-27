package examples

import "github.com/geniusrabbit/adcorelib/adformat"

// introConfig, offerConfig, hotspotConfig and multiscreenConfig are the
// Go builder (§3.2.10) equivalent of multiscreen.yaml — same three
// screens/groups, demonstrating .AsEntry(), .WithNavigateTo(...) and the
// ScreenRefField constructor specifically (§3.2.11).
var introConfig = adformat.Group("intro").AsEntry().WithAssets(
	adformat.Asset("background").Require().Types("image"),
).WithFields(
	adformat.StringField("headline").Len(0, 40).WithTitle("Headline"),
	adformat.StringField("cta_label").Len(0, 20).WithTitle("Button label").WithNavigateTo("offer"),
)

var offerConfig = adformat.Group("offer").WithAssets(
	adformat.Asset("background").Require().Types("image"),
).WithFields(
	adformat.StringField("offer_text").Len(0, 80).WithTitle("Offer text"),
	adformat.StringField("back_label").Len(0, 20).WithTitle("Back button label").WithNavigateTo(adformat.NavigateBack),
	adformat.URLField("url").WithTitle("Click-through URL"),
)

var hotspotConfig = adformat.Group("hotspot").Count(0, 5).WithAssets(
	adformat.Asset("icon").Require().Types("image"),
).WithFields(
	adformat.StringField("label").Len(0, 30).WithTitle("Label"),
	adformat.ScreenRefField("target_screen").WithTitle("Target screen").WithOptions(
		adformat.Opt("intro", "Intro"), adformat.Opt("offer", "Offer"),
	),
)

var multiscreenConfig = adformat.Config{}.WithGroups(introConfig, offerConfig, hotspotConfig)

// MultiscreenBanner is the Go var-table equivalent of multiscreen.yaml
// (§4.7/§3.2.10).
var MultiscreenBanner = adformat.Format{
	Codename: "multiscreen_banner",
	Title:    "Multi-screen Interactive Banner",
	Kind:     adformat.KindBanner,
	Tags:     []string{"expandable", "multi-screen"},
	Sizes:    []adformat.SizeOption{adformat.FixedSize("", 320, 480)},
	Config:   &multiscreenConfig,
}
