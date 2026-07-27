package tracking

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/adformat"
)

func TestTrackerBindingRoundTrip(t *testing.T) {
	f := adformat.URLField("impression_trackers").
		Count(0, 5).
		Bind(Tracker(EventImpression))

	p, ok := Get(f)
	if !ok || p.Event != EventImpression {
		t.Errorf("Get() = (%+v, %v), want (impression, true)", p, ok)
	}
	if !f.IsMultiple() {
		t.Error("expected the tracker field to be repeatable (several vendors at once)")
	}
}

func TestGetReturnsFalseWhenNoBinding(t *testing.T) {
	if _, ok := Get(adformat.URLField("url")); ok {
		t.Error("expected Get to report false with no binding present")
	}
}
