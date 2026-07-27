package adformat

import "testing"

func TestParseAspectRatio(t *testing.T) {
	tests := []struct {
		in      string
		wantW   float64
		wantH   float64
		wantErr bool
	}{
		{"1.91:1", 1.91, 1, false},
		{"16:9", 16, 9, false},
		{" 4 : 3 ", 4, 3, false},
		{"bad", 0, 0, true},
		{"1:0", 0, 0, true},
		{"-1:1", 0, 0, true},
		{"", 0, 0, true},
	}
	for _, tt := range tests {
		w, h, err := ParseAspectRatio(tt.in)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseAspectRatio(%q) error = %v, wantErr %v", tt.in, err, tt.wantErr)
			continue
		}
		if err == nil && (w != tt.wantW || h != tt.wantH) {
			t.Errorf("ParseAspectRatio(%q) = (%v, %v), want (%v, %v)", tt.in, w, h, tt.wantW, tt.wantH)
		}
	}
}

func TestAssetRequirementAspectRatioValue(t *testing.T) {
	a := AssetRequirement{}
	if _, _, ok := a.AspectRatioValue(); ok {
		t.Error("expected ok=false for an empty AspectRatio")
	}

	a.AspectRatio = "1.91:1"
	w, h, ok := a.AspectRatioValue()
	if !ok || w != 1.91 || h != 1 {
		t.Errorf("AspectRatioValue() = (%v, %v, %v), want (1.91, 1, true)", w, h, ok)
	}

	a.AspectRatio = "nope"
	if _, _, ok := a.AspectRatioValue(); ok {
		t.Error("expected ok=false for a malformed AspectRatio")
	}
}

func TestAssetRequirementMatchesAspectRatio(t *testing.T) {
	a := AssetRequirement{AspectRatio: "1.91:1"}
	if !a.MatchesAspectRatio(1200, 628) {
		t.Error("expected 1200x628 to satisfy a 1.91:1 aspect ratio within tolerance")
	}
	if a.MatchesAspectRatio(1000, 1000) {
		t.Error("expected a square image not to satisfy a 1.91:1 aspect ratio")
	}
	if a.MatchesAspectRatio(0, 100) {
		t.Error("expected zero width to never match")
	}

	unrestricted := AssetRequirement{}
	if !unrestricted.MatchesAspectRatio(100, 100) {
		t.Error("expected no AspectRatio constraint to match anything")
	}
}

func TestAssetRequirementSupportsCodec(t *testing.T) {
	a := AssetRequirement{AllowedCodecs: []string{"h264", "vp9"}}
	if !a.SupportsCodec("H264") {
		t.Error("expected case-insensitive codec match")
	}
	if a.SupportsCodec("av1") {
		t.Error("expected av1 not to be supported")
	}

	unrestricted := AssetRequirement{}
	if !unrestricted.SupportsCodec("anything") {
		t.Error("expected an empty AllowedCodecs to accept any codec")
	}
}

func TestAssetRequirementSupportsFramerate(t *testing.T) {
	a := AssetRequirement{AllowedFramerates: []float64{24, 25, 29.97, 30}}
	if !a.SupportsFramerate(29.97) {
		t.Error("expected 29.97 fps to be supported")
	}
	if a.SupportsFramerate(60) {
		t.Error("expected 60 fps not to be supported")
	}

	unrestricted := AssetRequirement{}
	if !unrestricted.SupportsFramerate(123) {
		t.Error("expected an empty AllowedFramerates to accept any frame rate")
	}
}

func TestBuilderGapModifiers(t *testing.T) {
	a := Asset("main").
		WithAspectRatio("16:9").
		WithSafeZone(0.1, 0.2, 0.15, 0.05).
		RequireAltText(100).
		Bitrate(500, 8000).
		Codecs("h264", "vp9").
		Framerates(24, 25, 30)

	if a.AspectRatio != "16:9" {
		t.Errorf("AspectRatio = %q, want 16:9", a.AspectRatio)
	}
	if a.SafeZone == nil || a.SafeZone.Top != 0.1 || a.SafeZone.Right != 0.2 || a.SafeZone.Bottom != 0.15 || a.SafeZone.Left != 0.05 {
		t.Errorf("SafeZone = %+v, unexpected", a.SafeZone)
	}
	if a.AltText == nil || !a.AltText.Required || a.AltText.MaxLength != 100 {
		t.Errorf("AltText = %+v, unexpected", a.AltText)
	}
	if a.BitrateMin != 500 || a.BitrateMax != 8000 {
		t.Errorf("Bitrate = (%d, %d), want (500, 8000)", a.BitrateMin, a.BitrateMax)
	}
	if len(a.AllowedCodecs) != 2 || len(a.AllowedFramerates) != 3 {
		t.Errorf("Codecs/Framerates = %+v/%+v, unexpected", a.AllowedCodecs, a.AllowedFramerates)
	}
}
