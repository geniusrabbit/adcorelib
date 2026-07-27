package adformat

// compareAspect compares one (value, minValue) size range against another
// the same way the old admodels/types package did — kept unchanged for
// AssetRequirement.SoftEqual (asset-level range comparison) and reused by
// SizeOption.Suits/Format.Suits (format-level display-area comparison,
// §3.2.13) so both layers share one, already-tested notion of "range A is
// compatible with range B".
func compareAspect(v, mv, targetV, targetMV int) int {
	if v == 0 && mv <= 0 {
		return 0
	}

	if mv == -1 {
		mv = v / 2
	}

	if targetMV == -1 {
		targetMV = targetV / 2
	}

	if mv > targetV || (mv == 0 && v > targetV) {
		return 1
	}

	if v != 0 && ((targetMV > 0 && v < targetMV) || (targetMV == 0 && v < targetV)) {
		return -1
	}

	return 0
}