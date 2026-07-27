package examples

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/adformat/openrtbnative"
)

func TestAllExamplesDecodeAndValidate(t *testing.T) {
	formats, err := All()
	if err != nil {
		t.Fatalf("failed to load examples: %v", err)
	}
	if len(formats) != len(Names) {
		t.Fatalf("expected %d formats, got %d", len(Names), len(formats))
	}
	for i, f := range formats {
		if f.Codename == "" {
			t.Errorf("%s: empty codename", Names[i])
		}
		if err := f.Validate(); err != nil {
			t.Errorf("%s: Validate: %v", Names[i], err)
		}
	}
}

func TestNativeExampleOpenRTBBindingsAreValid(t *testing.T) {
	f, err := Load("native.yaml")
	if err != nil {
		t.Fatalf("load native.yaml: %v", err)
	}
	if err := openrtbnative.Validate(f.Config); err != nil {
		t.Errorf("openrtbnative.Validate: %v", err)
	}
}

func TestGoBuilderExamplesValidate(t *testing.T) {
	if err := NativeAd.Validate(); err != nil {
		t.Errorf("NativeAd.Validate: %v", err)
	}
	if err := openrtbnative.Validate(NativeAd.Config); err != nil {
		t.Errorf("openrtbnative.Validate(NativeAd): %v", err)
	}
	if err := MultiscreenBanner.Validate(); err != nil {
		t.Errorf("MultiscreenBanner.Validate: %v", err)
	}
}

func TestGoBuilderNativeMatchesYAMLShape(t *testing.T) {
	yamlNative, err := Load("native.yaml")
	if err != nil {
		t.Fatalf("load native.yaml: %v", err)
	}
	if got, want := len(NativeAd.Config.Assets()), len(yamlNative.Config.Assets()); got != want {
		t.Errorf("asset count = %d, want %d", got, want)
	}
	if NativeAd.Config.GetField("phone") == nil {
		t.Error("expected a phone field on the Go builder NativeAd")
	}
	if len(NativeAd.Config.Groups()) != 1 {
		t.Errorf("expected 1 nested group (step), got %d", len(NativeAd.Config.Groups()))
	}
}

func TestGoBuilderMultiscreenMatchesYAMLShape(t *testing.T) {
	yamlMultiscreen, err := Load("multiscreen.yaml")
	if err != nil {
		t.Fatalf("load multiscreen.yaml: %v", err)
	}
	if got, want := len(MultiscreenBanner.Config.Groups()), len(yamlMultiscreen.Config.Groups()); got != want {
		t.Errorf("group count = %d, want %d", got, want)
	}
	entry, ok := MultiscreenBanner.Config.EntryGroup()
	if !ok || entry.Name != "intro" {
		t.Errorf("expected entry group %q, got %q (ok=%v)", "intro", entry.Name, ok)
	}
}
