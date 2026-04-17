package types

import (
	"testing"
)

func TestSimpleFormatAccessor_Basics(t *testing.T) {
	formats := MockFormats()
	fa := NewSimpleFormatAccessor(formats)

	if len(fa.Formats()) != len(formats) {
		t.Errorf("invalid formats count: %d != %d", len(fa.Formats()), len(formats))
	}

	f := fa.FormatByID(1)
	if f == nil || f.Codename != "direct" {
		t.Errorf("FormatByID(1) expected direct, got %#v", f)
	}

	if fa.FormatByID(999) != nil {
		t.Errorf("FormatByID(999) must be nil")
	}

	p := fa.FormatByCode("proxy")
	if p == nil || p.ID == 0 {
		t.Errorf("FormatByCode(proxy) expected non-nil format")
	}
}

func TestSimpleFormatAccessor_PrepareSort(t *testing.T) {
	unsorted := []*Format{
		{ID: 5, Codename: "f5"},
		{ID: 2, Codename: "f2"},
		{ID: 3, Codename: "f3"},
	}
	fa := NewSimpleFormatAccessor(unsorted)
	fa.Prepare()

	prev := uint64(0)
	for i, f := range fa.Formats() {
		if i == 0 {
			prev = f.ID
			continue
		}
		if prev > f.ID {
			t.Fatalf("formats not sorted after Prepare: %d > %d", prev, f.ID)
		}
		prev = f.ID
	}
}

func TestSimpleFormatAccessor_FormatsBySizeAndDirectSet(t *testing.T) {
	fa := NewSimpleFormatAccessor(MockFormats())

	// default search should exclude direct/video types
	list := fa.FormatsBySize(200, 200, 10, 10)
	if len(list) == 0 {
		t.Fatalf("expected formats for 200x200, got none")
	}

	foundProxy200 := false
	foundBanner200 := false
	for _, f := range list {
		if f.Codename == "proxy_200x200" {
			foundProxy200 = true
		}
		if f.Codename == "banner_200x200" {
			foundBanner200 = true
		}
	}
	if !foundProxy200 || !foundBanner200 {
		t.Errorf("expected proxy_200x200 and banner_200x200 in results, got: %+v", list)
	}

	// include video explicitly
	videoList := fa.FormatsBySize(0, 0, 0, 0, *NewFormatTypeBitset(FormatVideoType))
	hasVideo := false
	for _, f := range videoList {
		if f.IsVideo() {
			hasVideo = true
			break
		}
	}
	if !hasVideo {
		t.Errorf("expected video formats when filtering by video type")
	}

	// DirectFormatSet contains direct format id
	ds := fa.DirectFormatSet()
	if !ds.Has(uint(1)) {
		t.Errorf("DirectFormatSet must contain id 1 for direct format")
	}
}
