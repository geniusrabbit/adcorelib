package personification

import (
	"testing"
	"time"

	"github.com/geniusrabbit/udetect"
	"github.com/valyala/fasthttp"

	adgeo "github.com/geniusrabbit/adcorelib/geo"
)

func TestApplyCloudflareGeoFillMissing(t *testing.T) {
	h := cfHeaders(map[string]string{
		"CF-IPCountry":   "SK",
		"CF-Region-Code": "BL", // ISO 3166-2 suffix; Bratislava region is SK-BL
		"CF-IPCity":      "Bratislava",
		"CF-Postal-Code": "81101",
		"CF-Metro-Code":  "123",
		"CF-IPLatitude":  "48.1486",
		"CF-IPLongitude": "17.1077",
		"CF-Timezone":    "Europe/Bratislava",
		"CF-IPContinent": "EU",         // no Geo field — ignored
		"CF-Region":      "Bratislava", // name — ignored
	})
	geo := &udetect.Geo{}
	applyCloudflareGeo(geo, h)

	if geo.Country.ISO2() != "SK" {
		t.Fatalf("Country = %q, want SK", geo.Country.ISO2())
	}
	if geo.ID == 0 {
		t.Fatal("Country ID was not filled")
	}
	if geo.Region.ISO3166() != "SK-BL" {
		t.Fatalf("Region = %q, want SK-BL", geo.Region.ISO3166())
	}
	if geo.City != "Bratislava" {
		t.Fatalf("City = %q, want Bratislava", geo.City)
	}
	if geo.ZIP != "81101" {
		t.Fatalf("ZIP = %q, want 81101", geo.ZIP)
	}
	if geo.Metro != "123" {
		t.Fatalf("Metro = %q, want 123", geo.Metro)
	}
	if geo.Lat != 48.1486 || geo.Lon != 17.1077 {
		t.Fatalf("LatLon = (%v, %v)", geo.Lat, geo.Lon)
	}
	wantOff := timezoneOffsetHoursForTest(t, "Europe/Bratislava")
	if geo.UTCOffset != wantOff {
		t.Fatalf("UTCOffset = %d, want %d hours", geo.UTCOffset, wantOff)
	}
}

func TestApplyCloudflareGeoRegionSuffix(t *testing.T) {
	sk := adgeo.CountryByCode2("SK").Code2

	if got := adgeo.RegionCodeByPartial("BL", sk); got.ISO3166() != "SK-BL" {
		t.Fatalf("BL + SK = %q, want SK-BL", got.ISO3166())
	}
	if got := adgeo.RegionCodeByPartial("SK-BL", sk); got.ISO3166() != "SK-BL" {
		t.Fatalf("full SK-BL = %q", got.ISO3166())
	}

	h := cfHeaders(map[string]string{
		"CF-IPCountry":   "US",
		"CF-Region-Code": "CA",
	})
	geo := &udetect.Geo{}
	applyCloudflareGeo(geo, h)
	if geo.Country.ISO2() != "US" || geo.Region.ISO3166() != "US-CA" {
		t.Fatalf("US + CA → country %q region %q, want US / US-CA", geo.Country.ISO2(), geo.Region.ISO3166())
	}
}

func TestApplyCloudflareGeoKeepDetectorCountry(t *testing.T) {
	sk := adgeo.CountryByCode2("SK")
	geo := &udetect.Geo{
		ID:        uint(sk.ID),
		Country:   sk.Code2,
		City:      "Bratislava",
		ZIP:       "81101",
		Lat:       48.1,
		Lon:       17.1,
		UTCOffset: 1,
	}
	h := cfHeaders(map[string]string{
		"CF-IPCountry":   "US",
		"CF-Region-Code": "BL",
		"CF-IPCity":      "Los Angeles",
		"CF-Postal-Code": "90001",
		"CF-IPLatitude":  "34.05",
		"CF-IPLongitude": "-118.25",
		"CF-Timezone":    "America/Los_Angeles",
	})
	applyCloudflareGeo(geo, h)

	if geo.Country.ISO2() != "SK" {
		t.Fatalf("Country overwritten: %q", geo.Country.ISO2())
	}
	if geo.City != "Bratislava" || geo.ZIP != "81101" {
		t.Fatalf("City/ZIP overwritten: %q / %q", geo.City, geo.ZIP)
	}
	if geo.Lat != 48.1 || geo.Lon != 17.1 || geo.UTCOffset != 1 {
		t.Fatal("Lat/Lon/UTCOffset overwritten")
	}
	if geo.Region.ISO3166() != "SK-BL" {
		t.Fatalf("empty Region should still fill from CF using detector country: got %q", geo.Region.ISO3166())
	}
}

func TestApplyCloudflareGeoHeaderCaseInsensitive(t *testing.T) {
	geo := &udetect.Geo{}
	applyCloudflareGeo(geo, cfHeaders(map[string]string{
		"cf-ipcountry":   "SK",
		"cf-region-code": "BL",
		"cf-ipcity":      "Bratislava",
		"cf-iplatitude":  "48.1486",
		"cf-iplongitude": "17.1077",
	}))
	if geo.Country.ISO2() != "SK" || geo.City != "Bratislava" || geo.Region.ISO3166() != "SK-BL" {
		t.Fatalf("lowercase CF headers: country=%q city=%q region=%q", geo.Country.ISO2(), geo.City, geo.Region.ISO3166())
	}
	if geo.Lat != 48.1486 || geo.Lon != 17.1077 {
		t.Fatalf("lowercase lat/lon: (%v, %v)", geo.Lat, geo.Lon)
	}
}

func TestApplyCloudflareGeoSkipXX(t *testing.T) {
	geo := &udetect.Geo{}
	applyCloudflareGeo(geo, cfHeaders(map[string]string{
		"CF-IPCountry": "XX",
		"CF-IPCity":    "Unknown",
	}))
	if geo.Country.ISO2() != adgeo.UndefinedCountryCodeISO2 {
		t.Fatalf("XX must not set Country, got %q", geo.Country.ISO2())
	}
	if geo.ID != 0 {
		t.Fatalf("XX must not set Country ID, got %d", geo.ID)
	}
	if geo.City != "Unknown" {
		t.Fatalf("City should still fill when country is skipped, got %q", geo.City)
	}

	applyCloudflareGeo(geo, cfHeaders(map[string]string{"CF-IPCountry": "T1"}))
	if geo.Country.ISO2() != adgeo.UndefinedCountryCodeISO2 {
		t.Fatalf("T1 must not set Country, got %q", geo.Country.ISO2())
	}
}

func cfHeaders(values map[string]string) *fasthttp.RequestHeader {
	h := &fasthttp.RequestHeader{}
	for k, v := range values {
		h.Set(k, v)
	}
	return h
}

func timezoneOffsetHoursForTest(t *testing.T, name string) int {
	t.Helper()
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Fatal(err)
	}
	_, seconds := time.Now().In(loc).Zone()
	return seconds / 3600
}
