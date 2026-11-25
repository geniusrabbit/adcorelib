package personification

import (
	"context"
	"net/netip"
	"testing"

	"github.com/valyala/fasthttp"

	"github.com/geniusrabbit/adcorelib/personification/dummy"
	"github.com/geniusrabbit/udetect"
)

type lookupCount struct {
	calls int
	code  [2]byte
}

func (l *lookupCount) LookupCountry(netip.Addr) ([2]byte, bool) {
	l.calls++
	return l.code, true
}

func TestWhoisUsesHTTPUserAgentAndCloudflare(t *testing.T) {
	var gotUA string
	sign := &Signature{
		Detector: uaDetector(func(req *udetect.Request) {
			gotUA = req.UA
		}),
		CountryByIP: &lookupCount{code: [2]byte{'U', 'S'}},
	}
	var req fasthttp.RequestCtx
	req.Request.Header.SetUserAgent("curl/8.7.1")
	req.Request.Header.Set("CF-IPCountry", "DE")
	person, err := sign.Whois(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}
	if gotUA != "curl/8.7.1" {
		t.Fatalf("detector UA = %q", gotUA)
	}
	if person.UserInfo().GeoInfo().Country.ISO2() != "DE" {
		t.Fatalf("country = %s, want DE from Cloudflare", person.UserInfo().GeoInfo().Country.ISO2())
	}
}

func TestIdentifySkipsIPWhenCountryIsPassed(t *testing.T) {
	lookup := &lookupCount{code: [2]byte{'U', 'S'}}
	sign := &Signature{Detector: dummy.DummyClient{}, CountryByIP: lookup}
	person, err := sign.Identify(context.Background(), &udetect.Request{IP: "8.8.8.8", UA: "Mozilla/5.0"}, "DE")
	if err != nil {
		t.Fatal(err)
	}
	if lookup.calls != 0 {
		t.Fatalf("lookup calls = %d, want 0", lookup.calls)
	}
	if person.UserInfo().GeoInfo().Country.ISO2() != "DE" {
		t.Fatalf("country = %s, want DE", person.UserInfo().GeoInfo().Country.ISO2())
	}

	person, err = sign.Identify(context.Background(), &udetect.Request{IP: "8.8.8.8"}, "XX")
	if err != nil {
		t.Fatal(err)
	}
	if lookup.calls != 1 {
		t.Fatalf("lookup calls = %d, want 1", lookup.calls)
	}
	if person.UserInfo().GeoInfo().Country.ISO2() != "US" {
		t.Fatalf("country = %s, want US from IP", person.UserInfo().GeoInfo().Country.ISO2())
	}
}

type uaDetector func(*udetect.Request)

func (f uaDetector) Detect(_ context.Context, req *udetect.Request) (*udetect.Response, error) {
	f(req)
	return &udetect.Response{
		Device: &udetect.Device{Browser: &udetect.Browser{UA: req.UA}, OS: &udetect.OS{}},
	}, nil
}
