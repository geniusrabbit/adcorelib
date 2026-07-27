package adformat

import (
	"fmt"
	"testing"
)

// Ported unchanged from admodels/types/format_test.go (Test_Aspect) —
// compareAspect itself is untouched, only reused by new callers
// (AssetRequirement.SoftEqual, SizeOption.Suits/Format.Suits).
func Test_compareAspect(t *testing.T) {
	var tests = []struct {
		V1, MinV1 int
		V2, MinV2 int
		target    int
	}{
		{V1: 0, MinV1: -1, V2: 1000, MinV2: 400, target: 0},
		{V1: 100, MinV1: 0, V2: 100, MinV2: 0, target: 0},
		{V1: 150, MinV1: 50, V2: 100, MinV2: 50, target: 0},
		{V1: 150, MinV1: 0, V2: 100, MinV2: 50, target: 1},
		{V1: 150, MinV1: 110, V2: 100, MinV2: 50, target: 1},
		{V1: 50, MinV1: 0, V2: 100, MinV2: 0, target: -1},
		{V1: 90, MinV1: 50, V2: 100, MinV2: 0, target: -1},
		{V1: 200, MinV1: 50, V2: 100, MinV2: 20, target: 0},
		{V1: 200, MinV1: 50, V2: 300, MinV2: 150, target: 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d_%d___%d_%d", test.V1, test.MinV1, test.V2, test.MinV2), func(t *testing.T) {
			if compareAspect(test.V1, test.MinV1, test.V2, test.MinV2) != test.target {
				t.Errorf("Fail compare, should be %d", test.target)
			}
		})
	}
}
