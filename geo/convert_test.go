package geo

import (
	"slices"
	"testing"
)

func TestCountryCodes2IDs(t *testing.T) {
	ids := CountryCodes2IDs([]string{"us", "JP"})
	us := uint64(CountryByCode2("US").ID)
	jp := uint64(CountryByCode2("JP").ID)
	if !slices.Contains(ids, us) || !slices.Contains(ids, jp) {
		t.Fatalf("got %v, want US=%d JP=%d", ids, us, jp)
	}
}

func TestCountryCodes2IDsExpandsContinent(t *testing.T) {
	ids := CountryCodes2IDs([]string{"EU"})
	if len(ids) < 10 {
		t.Fatalf("EU expansion too small: %d", len(ids))
	}
	de := uint64(CountryByCode2("DE").ID)
	if !slices.Contains(ids, de) {
		t.Fatal("EU expansion must include DE")
	}
}

func TestCountryCodes2IDsEmpty(t *testing.T) {
	if CountryCodes2IDs(nil) != nil {
		t.Fatal("empty input must be nil")
	}
}

func TestRegionCodes2IDs(t *testing.T) {
	ids := RegionCodes2IDs([]string{"us-ca", "GB-ENG"})
	ca := uint64(RegionByCode("US-CA").ID)
	eng := uint64(RegionByCode("UK-ENG").ID)
	if !slices.Contains(ids, ca) || !slices.Contains(ids, eng) {
		t.Fatalf("got %v, want US-CA=%d GB-ENG=%d", ids, ca, eng)
	}
}

func TestRegionCodes2IDsUnknown(t *testing.T) {
	ids := RegionCodes2IDs([]string{"ZZ-ZZ"})
	if len(ids) != 1 || ids[0] != 0 {
		t.Fatalf("unknown region = %v, want [0]", ids)
	}
}
