package types

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/geo"
	"github.com/geniusrabbit/gosql/v2"
)

func TestRegionFilterInclude(t *testing.T) {
	ids, exclude := RegionFilter(gosql.NullableStringArray{"US-CA", "gb-eng"})
	if exclude {
		t.Fatal("want include mode")
	}
	wantCA := uint64(geo.RegionByCode("US-CA").ID)
	wantENG := uint64(geo.RegionByCode("GB-ENG").ID)
	foundCA, foundENG := false, false
	for _, id := range ids {
		switch id {
		case wantCA:
			foundCA = true
		case wantENG:
			foundENG = true
		}
	}
	if !foundCA || !foundENG {
		t.Fatalf("ids = %v, want CA=%d ENG=%d", ids, wantCA, wantENG)
	}
}

func TestRegionFilterExclude(t *testing.T) {
	ids, exclude := RegionFilter(gosql.NullableStringArray{"-US-CA"})
	if !exclude {
		t.Fatal("want exclude mode")
	}
	if len(ids) != 1 || ids[0] != uint64(geo.RegionByCode("US-CA").ID) {
		t.Fatalf("ids = %v", ids)
	}
}
