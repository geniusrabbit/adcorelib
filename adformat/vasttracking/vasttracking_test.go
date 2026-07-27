package vasttracking

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/adformat"
)

func TestTrackerBindingRoundTrip(t *testing.T) {
	f := adformat.URLField("complete_trackers").
		Count(0, 5).
		Bind(Tracker(EventComplete))

	p, ok := Get(f)
	if !ok || p.Event != EventComplete {
		t.Errorf("Get() = (%+v, %v), want (complete, true)", p, ok)
	}
}

func TestAllVASTEventsAreDistinct(t *testing.T) {
	events := []string{
		EventImpression, EventStart, EventFirstQuartile, EventMidpoint,
		EventThirdQuartile, EventComplete, EventPause, EventResume,
		EventMute, EventUnmute, EventFullscreen, EventSkip,
		EventClickThrough, EventClickTracking,
	}
	seen := map[string]bool{}
	for _, e := range events {
		if seen[e] {
			t.Errorf("duplicate VAST event constant: %s", e)
		}
		seen[e] = true
	}
	if len(seen) != 14 {
		t.Errorf("expected 14 distinct VAST events, got %d", len(seen))
	}
}
