package geo

import "testing"

func TestRegionCodeByPartial(t *testing.T) {
	sk := CountryByCode2("SK").Code2
	us := CountryByCode2("US").Code2

	if got := RegionCodeByPartial("US-CA", us); got.ISO3166() != "US-CA" {
		t.Fatalf("full code = %q, want US-CA", got.ISO3166())
	}
	if got := RegionCodeByPartial("CA", us); got.ISO3166() != "US-CA" {
		t.Fatalf("suffix CA + US = %q, want US-CA", got.ISO3166())
	}
	if got := RegionCodeByPartial("BL", sk); got.ISO3166() != "SK-BL" {
		t.Fatalf("suffix BL + SK = %q, want SK-BL", got.ISO3166())
	}
	if got := RegionCodeByPartial("", sk); got != UndefinedRegionCode {
		t.Fatalf("empty = %q", got.ISO3166())
	}
	if got := RegionCodeByPartial("NOPE", UndefinedCountryCode2); got != UndefinedRegionCode {
		t.Fatalf("unknown without country = %q", got.ISO3166())
	}
}
