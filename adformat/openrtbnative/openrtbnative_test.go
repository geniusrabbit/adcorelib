package openrtbnative

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/adformat"
)

func TestBindingsRoundTrip(t *testing.T) {
	asset := adformat.Asset("main").Bind(ImageAsset(1, ImageTypeMain))
	img, ok := GetImage(asset)
	if !ok || img.ID != 1 || img.ImageType != ImageTypeMain {
		t.Errorf("GetImage() = (%+v, %v), want ({1 %d}, true)", img, ok, ImageTypeMain)
	}

	titleField := adformat.StringField("title").Bind(TitleField(3))
	title, ok := GetTitle(titleField)
	if !ok || title.ID != 3 {
		t.Errorf("GetTitle() = (%+v, %v), want ({3}, true)", title, ok)
	}

	descField := adformat.StringField("description").Bind(DataField(4, DataTypeDesc))
	data, ok := GetData(descField)
	if !ok || data.ID != 4 || data.DataType != DataTypeDesc {
		t.Errorf("GetData() = (%+v, %v)", data, ok)
	}

	videoAsset := adformat.Asset("video").Bind(VideoAsset(9, []int{2, 3}, []int{1, 2}))
	video, ok := GetVideo(videoAsset)
	if !ok || video.ID != 9 || len(video.Protocols) != 2 || len(video.APIs) != 2 {
		t.Errorf("GetVideo() = (%+v, %v)", video, ok)
	}
}

func TestGetReturnsFalseWhenNoBinding(t *testing.T) {
	asset := adformat.Asset("main")
	if _, ok := GetImage(asset); ok {
		t.Error("expected GetImage to report false with no binding present")
	}
}

func TestValidateDetectsMissingAndDuplicateIDs(t *testing.T) {
	cfg := adformat.Config{}.WithAssets(
		adformat.Asset("main").Bind(ImageAsset(1, ImageTypeMain)),
		adformat.Asset("logo").Bind(ImageAsset(1, ImageTypeLogo)), // duplicate id
	)
	if err := Validate(&cfg); err == nil {
		t.Error("expected a duplicate id error")
	}

	zeroID := adformat.Config{}.WithAssets(
		adformat.Asset("main").Bind(ImageAsset(0, ImageTypeMain)),
	)
	if err := Validate(&zeroID); err == nil {
		t.Error("expected a zero id error")
	}
}

func TestValidateOKForUniqueIDs(t *testing.T) {
	cfg := adformat.Config{}.
		WithAssets(adformat.Asset("main").Bind(ImageAsset(1, ImageTypeMain))).
		WithFields(adformat.StringField("title").Bind(TitleField(2)))
	if err := Validate(&cfg); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateNilConfig(t *testing.T) {
	if err := Validate(nil); err != nil {
		t.Errorf("expected nil error for a nil config, got %v", err)
	}
}
