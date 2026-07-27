package adformat

// Unlimited is the sentinel MaxCount/GetMaxCount value meaning "no upper
// bound on the number of instances".
const Unlimited = -1

// effectiveCount derives the actual (min, max) instance-count range from
// the raw MinCount/MaxCount fields declared on AssetRequirement or Field,
// plus the Required flag, per the shared semantics (§3.2, п.2/§3.2.15):
//   - MinCount == 0 → effective min is 1 when required, else 0.
//   - MaxCount == 0 → effective max is "exactly one" (max(effective min, 1)),
//     matching the pre-MinCount/MaxCount behavior of the old model.
//   - MaxCount == Unlimited (-1) → no upper bound.
func effectiveCount(minCount, maxCount int, required bool) (mn, mx int) {
	mn = minCount
	if mn <= 0 {
		if required {
			mn = 1
		} else {
			mn = 0
		}
	}
	switch {
	case maxCount == Unlimited:
		mx = Unlimited
	case maxCount == 0:
		mx = mn
		if mx < 1 {
			mx = 1
		}
	default:
		mx = maxCount
	}
	return mn, mx
}

// isMultipleCount reports whether the effective max count allows more than
// one instance.
func isMultipleCount(maxCount int) bool {
	return maxCount == Unlimited || maxCount > 1
}
