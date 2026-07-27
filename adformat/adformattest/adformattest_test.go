package adformattest

import "testing"

func TestMockFormats(t *testing.T) {
	formats := MockFormats()
	if len(formats) != 8 {
		t.Fatalf("expected 8 mock formats, got %d", len(formats))
	}
	seen := map[uint64]bool{}
	for _, f := range formats {
		if seen[f.ID] {
			t.Errorf("duplicate mock format id %d", f.ID)
		}
		seen[f.ID] = true
		if err := f.Validate(); err != nil {
			t.Errorf("%s: Validate: %v", f.Codename, err)
		}
	}
}

func TestNewAccessor(t *testing.T) {
	acc := NewAccessor()
	if len(acc.Formats()) != 8 {
		t.Fatalf("expected 8 formats in accessor, got %d", len(acc.Formats()))
	}
	if f := acc.FormatByCode("native"); f == nil {
		t.Error("expected to find native format by code")
	}
	if f := acc.FormatByID(1); f == nil {
		t.Error("expected to find format by id 1")
	}
}

func TestFormatByCodename(t *testing.T) {
	if f := FormatByCodename("push_ad"); f == nil {
		t.Error("expected to find push_ad format")
	}
	if f := FormatByCodename("does-not-exist"); f != nil {
		t.Error("expected nil for unknown codename")
	}
}
