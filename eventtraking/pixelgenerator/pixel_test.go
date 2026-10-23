package pixelgenerator

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/geniusrabbit/adcorelib/eventtraking/events"
)

type stubEvent struct {
	code events.Code
}

func (e stubEvent) Pack() events.Code { return e.code }

type stubLead string

func (s stubLead) String() string { return string(s) }

func TestPixelGenerator_MatchesURLValues(t *testing.T) {
	const host = "https://t.example"
	g := NewPixelGenerator[stubEvent, stubLead](host)
	ev := stubEvent{code: events.CodeObj([]byte("impression payload=1&x"), nil)}

	gif, err := g.Event(ev, false)
	if err != nil {
		t.Fatal(err)
	}
	if gif != referenceEvent(host, ev.code, false) {
		t.Fatalf("gif\n got %s\nwant %s", gif, referenceEvent(host, ev.code, false))
	}

	js, err := g.Event(ev, true)
	if err != nil {
		t.Fatal(err)
	}
	if js != referenceEvent(host, ev.code, true) {
		t.Fatalf("js\n got %s\nwant %s", js, referenceEvent(host, ev.code, true))
	}

	directs := []string{
		"",
		"https://landing.example/a",
		"a b=c&d",
		"привет ~ok",
		string([]byte{0, 1, ' ', '=', '&', '~', '*'}),
	}
	for _, direct := range directs {
		got, err := g.EventDirect(ev, direct)
		if err != nil {
			t.Fatal(err)
		}
		want := referenceDirect(host, ev.code, direct)
		if got != want {
			t.Fatalf("direct %q\n got %s\nwant %s", direct, got, want)
		}
	}

	leads := []stubLead{"", "a b=c&d", "привет", stubLead(string([]byte{0xff, '~', ' '}))}
	for _, lead := range leads {
		got, err := g.Lead(lead)
		if err != nil {
			t.Fatal(err)
		}
		want := fmt.Sprintf("%s/lead?l=%s", host, url.QueryEscape(lead.String()))
		if got != want {
			t.Fatalf("lead %q\n got %s\nwant %s", lead, got, want)
		}
	}
}

func TestPixelGenerator_EmptyCode(t *testing.T) {
	g := NewPixelGenerator[stubEvent, stubLead]("https://t.example")
	ev := stubEvent{code: events.CodeObj(nil, events.ErrEmptyData)}

	gif, err := g.Event(ev, false)
	if gif != "" || !errors.Is(err, events.ErrEmptyData) {
		t.Fatalf("Event() = %q, %v", gif, err)
	}
	direct, err := g.EventDirect(ev, "https://landing.example")
	if direct != "" || !errors.Is(err, events.ErrEmptyData) {
		t.Fatalf("EventDirect() = %q, %v", direct, err)
	}
}

func TestQueryEscapeMatchesNetURL(t *testing.T) {
	samples := []string{
		"",
		"abcXYZ09-_.~",
		" ",
		"=",
		"&",
		"привет",
		"a b=c&d~*",
	}
	var all [256]byte
	for i := range all {
		all[i] = byte(i)
	}
	samples = append(samples, string(all[:]))

	for _, sample := range samples {
		var b strings.Builder
		b.Grow(queryEscapedLen(sample))
		appendQueryEscape(&b, sample)
		if got, want := b.String(), url.QueryEscape(sample); got != want {
			t.Fatalf("escape %q\n got %s\nwant %s", sample, got, want)
		}
	}
}

func referenceEvent(host string, code events.Code, js bool) string {
	encoded := code.Compress().URLEncode()
	u := url.Values{"i": []string{encoded.String()}}
	path := "/t/px.gif?"
	if js {
		path = "/t/px.js?"
	}
	return host + path + u.Encode()
}

func referenceDirect(host string, code events.Code, direct string) string {
	encoded := code.Compress().URLEncode()
	u := url.Values{
		"i": []string{encoded.String()},
		"u": []string{direct},
	}
	return host + "/go/m?" + u.Encode()
}

var sink string

func BenchmarkPixelGenerator(b *testing.B) {
	g := NewPixelGenerator[stubEvent, stubLead]("https://t.example")
	ev := stubEvent{code: events.CodeObj([]byte("impression-payload-1234567890"), nil)}
	direct := "https://landing.example/path?a=1&b=2 hello"
	lead := stubLead("order=42&name=привет")

	b.Run("Event", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			s, err := g.Event(ev, false)
			if err != nil {
				b.Fatal(err)
			}
			sink = s
		}
	})
	b.Run("EventDirect", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			s, err := g.EventDirect(ev, direct)
			if err != nil {
				b.Fatal(err)
			}
			sink = s
		}
	})
	b.Run("Lead", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			s, err := g.Lead(lead)
			if err != nil {
				b.Fatal(err)
			}
			sink = s
		}
	})
}
