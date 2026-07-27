package adformat

import (
	"strings"
	"testing"
)

func TestValidateCatchesBadRegExp(t *testing.T) {
	cfg := Config{}.WithFields(StringField("zip").WithRegExp(`(`))
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "invalid regexp") {
		t.Errorf("expected invalid regexp error, got %v", err)
	}
}

func TestValidateCatchesUnknownConditionField(t *testing.T) {
	cfg := Config{}.WithFields(
		PhoneField("phone").When(FieldIn("does_not_exist", "business")),
	)
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "unknown field/asset") {
		t.Errorf("expected unknown condition field error, got %v", err)
	}
}

func TestValidateCatchesConditionCycle(t *testing.T) {
	cfg := Config{}.WithFields(
		StringField("a").When(FieldIn("b", "x")),
		StringField("b").When(FieldIn("a", "y")),
	)
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "cyclic condition dependency") {
		t.Errorf("expected cyclic condition dependency error, got %v", err)
	}
}

func TestValidateCatchesUnknownNavigateTo(t *testing.T) {
	cfg := Config{}.WithGroups(
		Group("intro").AsEntry().WithFields(
			StringField("cta").WithNavigateTo("nonexistent"),
		),
		Group("offer"),
	)
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "navigate_to references unknown screen") {
		t.Errorf("expected unknown navigate_to error, got %v", err)
	}
}

func TestValidateAllowsReservedNavigateCommands(t *testing.T) {
	cfg := Config{}.WithGroups(
		Group("intro").AsEntry().WithFields(
			StringField("back").WithNavigateTo(NavigateBack),
			StringField("close").WithNavigateTo(NavigateClose),
		),
	)
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error for reserved navigate_to commands: %v", err)
	}
}

func TestValidateCatchesUnknownScreenRefOption(t *testing.T) {
	cfg := Config{}.WithGroups(
		Group("intro").AsEntry(),
		Group("hotspot").WithFields(
			ScreenRefField("target").WithOptions(Opt("nonexistent", "")),
		),
	)
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "does not reference an existing screen") {
		t.Errorf("expected unknown screen_ref option error, got %v", err)
	}
}

func TestValidateNavigateToScopeIsInheritedFromEnclosingScreens(t *testing.T) {
	// Regression test: a field nested one level inside a named screen
	// must resolve navigate_to against the screen's own siblings, not
	// against something inside the screen itself (§3.2.11).
	cfg := Config{}.WithGroups(
		Group("intro").AsEntry().WithFields(
			StringField("cta").WithNavigateTo("offer"),
		),
		Group("offer"),
	)
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateCatchesBadAspectRatio(t *testing.T) {
	cfg := Config{}.WithAssets(Asset("main").WithAspectRatio("not-a-ratio"))
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "invalid aspect_ratio") {
		t.Errorf("expected invalid aspect_ratio error, got %v", err)
	}
}

func TestValidateCatchesInvertedBitrateRange(t *testing.T) {
	cfg := Config{}.WithAssets(Asset("main").Bitrate(8000, 500))
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "bitrate_min") {
		t.Errorf("expected an inverted bitrate range error, got %v", err)
	}
}

func TestValidateAllowsGoodAspectRatioAndBitrateRange(t *testing.T) {
	cfg := Config{}.WithAssets(
		Asset("main").WithAspectRatio("1.91:1").Bitrate(500, 8000),
	)
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateOKForCleanConfig(t *testing.T) {
	cfg := Config{}.
		WithAssets(Asset("main").Require()).
		WithFields(
			StringField("adtype").WithOptions(Opt("business", ""), Opt("game", "")),
			PhoneField("phone").When(FieldIn("adtype", "business")),
		)
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error for a clean config: %v", err)
	}
}
