package adformat

import "testing"

func TestFormatString(t *testing.T) {
	sized := Format{Codename: "banner_300x250", Sizes: []SizeOption{FixedSize("", 300, 250)}}
	if got, want := sized.String(), "300x250"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	direct := Format{Codename: "direct", Kind: KindDirect}
	if got, want := direct.String(), "direct"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestFormatIsKindHelpers(t *testing.T) {
	tests := []struct {
		f    Format
		want string
	}{
		{Format{Kind: KindDirect}, "direct"},
		{Format{Kind: KindBanner}, "banner"},
		{Format{Kind: KindVideo}, "video"},
		{Format{Kind: KindAudio}, "audio"},
		{Format{Kind: KindNative}, "native"},
		{Format{Kind: KindPush}, "push"},
	}
	for _, tt := range tests {
		switch tt.want {
		case "direct":
			if !tt.f.IsDirect() {
				t.Errorf("expected IsDirect for %+v", tt.f)
			}
		case "banner":
			if !tt.f.IsBanner() {
				t.Errorf("expected IsBanner for %+v", tt.f)
			}
		case "video":
			if !tt.f.IsVideo() {
				t.Errorf("expected IsVideo for %+v", tt.f)
			}
		case "audio":
			if !tt.f.IsAudio() {
				t.Errorf("expected IsAudio for %+v", tt.f)
			}
		case "native":
			if !tt.f.IsNative() {
				t.Errorf("expected IsNative for %+v", tt.f)
			}
		case "push":
			if !tt.f.IsPush() {
				t.Errorf("expected IsPush for %+v", tt.f)
			}
		}
	}
}

func TestFormatHasTag(t *testing.T) {
	f := Format{Tags: []string{"playable", "mraid"}}
	if !f.HasTag("playable") {
		t.Error("expected HasTag(playable) to be true")
	}
	if f.HasTag("rewarded") {
		t.Error("expected HasTag(rewarded) to be false")
	}
}

func TestFormatSupportsOrientation(t *testing.T) {
	unrestricted := Format{}
	if !unrestricted.SupportsOrientation(OrientationPortrait) {
		t.Error("expected empty Orientations to support any orientation")
	}

	portraitOnly := Format{Orientations: []Orientation{OrientationPortrait}}
	if !portraitOnly.SupportsOrientation(OrientationPortrait) {
		t.Error("expected portrait to be supported")
	}
	if portraitOnly.SupportsOrientation(OrientationLandscape) {
		t.Error("expected landscape not to be supported")
	}
}

func TestFormatValidateCatchesInvalidFlexibleSize(t *testing.T) {
	f := Format{Sizes: []SizeOption{{Flexible: true, MinWidth: 500, MaxWidth: 100}}}
	if err := f.Validate(); err == nil {
		t.Error("expected an error for min_width > max_width")
	}
}

func TestFormatGetConfigSafeOnNil(t *testing.T) {
	var f *Format
	if f.GetConfig() != nil {
		t.Error("expected nil GetConfig on a nil Format")
	}
}
