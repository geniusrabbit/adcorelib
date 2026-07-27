package openrtbvideo

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/adformat"
)

func TestVideoBindingRoundTrip(t *testing.T) {
	f := adformat.Format{
		Kind: adformat.KindVideo,
		Bindings: []adformat.Binding{
			Video(Params{Protocols: []int{2, 3, 5, 6}, APIs: []int{1, 2}, Linearity: 1}),
		},
	}
	p, ok := Get(f)
	if !ok {
		t.Fatal("expected to find the openrtb_video binding")
	}
	if len(p.Protocols) != 4 || p.Linearity != 1 || len(p.APIs) != 2 {
		t.Errorf("unexpected params: %+v", p)
	}
}

func TestGetReturnsFalseWhenNoBinding(t *testing.T) {
	if _, ok := Get(adformat.Format{}); ok {
		t.Error("expected Get to report false with no binding present")
	}
}
