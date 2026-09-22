package geo

import "testing"

func TestCountryByCode2(t *testing.T) {
	us := CountryByCode2("US")
	if us.ISO2() != "US" || us.Name == "" {
		t.Fatalf("US lookup: ISO2=%q name=%q", us.ISO2(), us.Name)
	}
	if CountryByCode2("XX").ISO2() != UndefinedCountryCodeISO2 {
		t.Fatal("unknown country must be undefined")
	}
}

func TestRegionByCodeUSCA(t *testing.T) {
	ca := RegionByCode("US-CA")
	if ca.Code() != "US-CA" {
		t.Fatalf("code = %q, want US-CA", ca.Code())
	}
	if ca.Country().ISO2() != "US" {
		t.Fatalf("country = %q, want US", ca.Country().ISO2())
	}
	if ca.ID == 0 {
		t.Fatal("US-CA must have a non-zero ID")
	}
}

func TestRegionUKGBAlias(t *testing.T) {
	engUK := RegionByCode("UK-ENG")
	engGB := RegionByCode("GB-ENG")
	if engUK.ID == 0 || engUK.ID != engGB.ID {
		t.Fatalf("UK-ENG ID=%d GB-ENG ID=%d, want same non-zero", engUK.ID, engGB.ID)
	}
	if engUK.Code() != "GB-ENG" {
		t.Fatalf("canonical code = %q, want GB-ENG", engUK.Code())
	}
	if engUK.Country().ISO2() != "GB" {
		t.Fatalf("country ISO2 = %q, want GB", engUK.Country().ISO2())
	}
}

func TestRegionByLatLng(t *testing.T) {
	near := RegionByLatLng(42.51, 1.52)
	if near.ID == 0 {
		t.Fatal("expected a nearest region with coordinates")
	}
}

func TestCountriesByContinent(t *testing.T) {
	eu := CountriesByContinent("EU")
	if len(eu) == 0 {
		t.Fatal("EU must have countries")
	}
	for _, c := range eu {
		if c.Continent() != "EU" {
			t.Fatalf("%s continent = %q", c.ISO2(), c.Continent())
		}
	}
	if len(CountriesByContinent("ZZ")) != 0 {
		t.Fatal("unknown continent must be empty")
	}
}

func TestContinentByCode2(t *testing.T) {
	eu := ContinentByCode2("EU")
	if eu == nil || eu.Code2 != "EU" {
		t.Fatalf("EU continent = %+v", eu)
	}
	if ContinentByCode2("ZZ") != nil {
		t.Fatal("unknown continent must be nil")
	}
}
