package adformat

import "testing"

func TestSizeOptionSuits(t *testing.T) {
	tests := []struct {
		name   string
		size   SizeOption
		w, h   int
		expect bool
	}{
		{"exact match", FixedSize("", 300, 250), 300, 250, true},
		{"exact mismatch", FixedSize("", 300, 250), 300, 100, false},
		{"flexible within bounds", FlexibleSize("", 250, 970, 90, 415), 320, 100, true},
		{"flexible below min width", FlexibleSize("", 250, 970, 90, 415), 100, 100, false},
		{"flexible above max width", FlexibleSize("", 250, 970, 90, 415), 1000, 100, false},
		{"flexible unbounded", FlexibleSize("", 0, 0, 0, 0), 5000, 5000, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.size.Suits(tt.w, tt.h); got != tt.expect {
				t.Errorf("Suits(%d,%d) = %v, want %v", tt.w, tt.h, got, tt.expect)
			}
		})
	}
}

func TestSizeOptionIsFixed(t *testing.T) {
	if !FixedSize("", 300, 250).IsFixed() {
		t.Error("expected fixed size to be IsFixed")
	}
	if FlexibleSize("", 0, 0, 0, 0).IsFixed() {
		t.Error("flexible size should not be IsFixed")
	}
}

func TestFormatSuits(t *testing.T) {
	f := Format{Sizes: []SizeOption{
		FixedSize("Medium Rectangle", 300, 250).AsDefault(),
		FixedSize("Square", 250, 250),
		FlexibleSize("Flexible", 250, 970, 90, 415),
	}}

	if !f.Suits(300, 250) {
		t.Error("expected 300x250 to suit")
	}
	if !f.Suits(250, 250) {
		t.Error("expected 250x250 to suit")
	}
	if !f.Suits(970, 415) {
		t.Error("expected 970x415 (flexible bound) to suit")
	}
	if f.Suits(999, 999) {
		t.Error("expected 999x999 not to suit")
	}

	noSizes := Format{Kind: KindDirect}
	if !noSizes.Suits(1, 1) {
		t.Error("a format with no Sizes at all should suit any size")
	}
}

func TestFormatDefaultSize(t *testing.T) {
	f := Format{Sizes: []SizeOption{
		FixedSize("A", 300, 250),
		FixedSize("B", 250, 250).AsDefault(),
	}}
	size, ok := f.DefaultSize()
	if !ok || size.Title != "B" {
		t.Errorf("expected explicit default B, got %+v (ok=%v)", size, ok)
	}

	single := Format{Sizes: []SizeOption{FixedSize("Only", 320, 480)}}
	size, ok = single.DefaultSize()
	if !ok || size.Title != "Only" {
		t.Errorf("expected the only size to be default, got %+v (ok=%v)", size, ok)
	}

	none := Format{}
	if _, ok := none.DefaultSize(); ok {
		t.Error("expected ok=false for a format with no Sizes")
	}
}

func TestFormatIsFixedSize(t *testing.T) {
	fixed := Format{Sizes: []SizeOption{FixedSize("", 300, 250), FixedSize("", 250, 250)}}
	if !fixed.IsFixedSize() {
		t.Error("expected all-fixed Sizes to report IsFixedSize=true")
	}

	mixed := Format{Sizes: []SizeOption{FixedSize("", 300, 250), FlexibleSize("", 0, 0, 0, 0)}}
	if mixed.IsFixedSize() {
		t.Error("expected mixed Sizes to report IsFixedSize=false")
	}

	empty := Format{}
	if empty.IsFixedSize() {
		t.Error("expected a format with no Sizes to report IsFixedSize=false")
	}
}
