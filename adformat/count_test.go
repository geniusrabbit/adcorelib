package adformat

import "testing"

func TestEffectiveCount(t *testing.T) {
	tests := []struct {
		name               string
		minCount, maxCount int
		required           bool
		wantMin, wantMax   int
	}{
		{"unset optional", 0, 0, false, 0, 1},
		{"unset required", 0, 0, true, 1, 1},
		{"explicit range", 1, 3, false, 1, 3},
		{"unlimited", 1, Unlimited, false, 1, Unlimited},
		{"min only, required inferred by min", 2, 0, false, 2, 2},
		{"zero min but required with max set", 0, 5, true, 1, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mn, mx := effectiveCount(tt.minCount, tt.maxCount, tt.required)
			if mn != tt.wantMin || mx != tt.wantMax {
				t.Errorf("effectiveCount(%d,%d,%v) = (%d,%d), want (%d,%d)",
					tt.minCount, tt.maxCount, tt.required, mn, mx, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestAssetRequirementCount(t *testing.T) {
	a := AssetRequirement{Required: true}
	if got := a.GetMinCount(); got != 1 {
		t.Errorf("GetMinCount() = %d, want 1", got)
	}
	if got := a.GetMaxCount(); got != 1 {
		t.Errorf("GetMaxCount() = %d, want 1", got)
	}
	if a.IsMultiple() {
		t.Error("single required asset should not be IsMultiple")
	}

	logo := AssetRequirement{MinCount: 1, MaxCount: 3}
	if !logo.IsMultiple() {
		t.Error("expected logo with max_count=3 to be IsMultiple")
	}
	mn, mx := logo.RangeCount()
	if mn != 1 || mx != 3 {
		t.Errorf("RangeCount() = (%d,%d), want (1,3)", mn, mx)
	}

	unlimited := AssetRequirement{MinCount: 1, MaxCount: Unlimited}
	if !unlimited.IsMultiple() {
		t.Error("expected unlimited asset to be IsMultiple")
	}
}

func TestFieldCount(t *testing.T) {
	categories := StringField("categories").WithOptions(Opt("a", ""), Opt("b", ""), Opt("c", "")).Count(1, 3)
	if !categories.IsMultiple() {
		t.Error("expected categories field with max_count=3 to be IsMultiple")
	}
	mn, mx := categories.RangeCount()
	if mn != 1 || mx != 3 {
		t.Errorf("RangeCount() = (%d,%d), want (1,3)", mn, mx)
	}

	title := StringField("title").Require()
	if title.IsMultiple() {
		t.Error("plain required field should not be IsMultiple")
	}
}
