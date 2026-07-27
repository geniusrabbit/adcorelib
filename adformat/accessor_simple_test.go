package adformat

import "testing"

func TestSimpleAccessor(t *testing.T) {
	direct := &Format{ID: 1, Codename: "direct", Kind: KindDirect}
	square := &Format{ID: 2, Codename: "banner_200x200", Kind: KindBanner,
		Sizes: []SizeOption{FixedSize("", 200, 200)}}
	medrect := &Format{ID: 3, Codename: "banner_300x250", Kind: KindBanner,
		Sizes: []SizeOption{FixedSize("", 300, 250)}}

	acc := NewSimpleAccessor(medrect, direct, square) // intentionally unsorted input

	if got := len(acc.Formats()); got != 3 {
		t.Fatalf("Formats() length = %d, want 3", got)
	}

	if f := acc.FormatByID(2); f == nil || f.Codename != "banner_200x200" {
		t.Errorf("FormatByID(2) = %+v, want banner_200x200", f)
	}
	if f := acc.FormatByID(999); f != nil {
		t.Errorf("FormatByID(999) = %+v, want nil", f)
	}

	if f := acc.FormatByCode("direct"); f == nil {
		t.Error("expected to find direct format by code")
	}
	if f := acc.FormatByCode("missing"); f != nil {
		t.Error("expected nil for unknown code")
	}

	bySize := acc.FormatsBySize(200, 200)
	foundSquare, foundDirect := false, false
	for _, f := range bySize {
		if f.Codename == "banner_200x200" {
			foundSquare = true
		}
		if f.Codename == "direct" {
			foundDirect = true
		}
		if f.Codename == "banner_300x250" {
			t.Error("banner_300x250 should not match a 200x200 lookup")
		}
	}
	if !foundSquare {
		t.Error("expected banner_200x200 to match a 200x200 lookup")
	}
	if !foundDirect {
		t.Error("expected direct format to always be included regardless of size")
	}

	directs := acc.DirectFormatSet()
	if !directs.Has(1) {
		t.Error("expected direct format id 1 to be in DirectFormatSet")
	}
	if directs.Has(2) {
		t.Error("did not expect banner id 2 in DirectFormatSet")
	}
}
