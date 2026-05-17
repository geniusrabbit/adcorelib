//
// @project GeniusRabbit corelib 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package types

import (
	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/gogeo"
	"github.com/geniusrabbit/gosql/v2"

	"github.com/geniusrabbit/adcorelib/i18n/languages"
)

// IntArrayToUint64 converts a []int slice to a sorted
// [gosql.NullableOrderedNumberArray][uint64]. Only non-negative values produce
// meaningful results; negative ints are cast directly (two's complement) and
// will produce incorrect large IDs.
func IntArrayToUint64(arr []int) (res gosql.NullableOrderedNumberArray[uint64]) {
	if len(arr) < 1 {
		return res
	}
	for _, v := range arr {
		res = append(res, uint64(v))
	}
	res.Sort()
	return res
}

// IntArrayToInt64 converts a []int slice to a sorted
// [gosql.NullableOrderedNumberArray][int64]. Negative values are preserved.
func IntArrayToInt64(arr []int) (res gosql.NullableOrderedNumberArray[int64]) {
	if len(arr) < 1 {
		return res
	}
	for _, v := range arr {
		res = append(res, int64(v))
	}
	res.Sort()
	return res
}

// IDArrayFilter splits a signed int64 slice into an include or exclude list.
//
// If any element is positive, those elements form an include list and
// executed=false is returned. If all elements are negative (or zero), their
// absolute values form an exclude list and executed=true is returned.
// An empty input returns (nil, false) — no constraint.
//
// The returned bool mirrors the "executed" flag used by [BaseFilter.SetPositive]:
// false = include mode, true = exclude mode.
func IDArrayFilter(arr gosql.NullableOrderedNumberArray[int64]) (narr gosql.NullableOrderedNumberArray[uint64], executed bool) {
	if arr.Len() < 1 {
		return narr, false
	}

	subarr := arr.Map(func(v int64) (int64, bool) { return v, v > 0 })
	if subarr.Len() < 1 {
		subarr = arr.Map(func(v int64) (int64, bool) { return -v, v < 0 })
		executed = true
	}

	narr = gosql.NullableOrderedNumberArray[uint64](
		gocast.Slice[uint64](subarr),
	)
	narr.Sort()
	return narr, executed
}

// IDArrayFilterAny is a type-dispatching wrapper around [IDArrayFilter].
//
// Accepted types and their include/exclude semantics:
//
//   - gosql.NullableOrderedNumberArray[int64] / []int64 / []int — delegated to
//     [IDArrayFilter]: positive values = include (false), negative = exclude (true).
//     []int is converted via [IntArrayToInt64] first.
//   - gosql.NullableOrderedNumberArray[uint64] / []uint64 — raw unsigned IDs
//     with no sign convention; always returns false (include mode). For a
//     pre-processed exclude list, call [BaseFilter.SetPositive] manually after.
//
// If the value does not match any case and panicMsg is non-empty, it panics.
func IDArrayFilterAny(v any, panicMsg string) (gosql.NullableOrderedNumberArray[uint64], bool) {
	switch vl := v.(type) {
	case gosql.NullableOrderedNumberArray[int64]:
		return IDArrayFilter(vl)
	case gosql.NullableOrderedNumberArray[uint64]:
		return vl, false
	case []int:
		return IDArrayFilter(IntArrayToInt64(vl))
	case []int64:
		return IDArrayFilter(gosql.NullableOrderedNumberArray[int64](vl))
	case []uint64:
		// Raw unsigned IDs carry no exclude signal — always include mode.
		return gosql.NullableOrderedNumberArray[uint64](vl), false
	default:
		if panicMsg != "" {
			panic(panicMsg)
		}
	}
	return nil, false
}

// StringArrayFilter splits a string array into an include or exclude list
// based on a leading '-' prefix convention.
//
// Elements without a '-' prefix form an include list (executed=false).
// If no such elements exist, elements with a '-' prefix have the prefix
// stripped and form an exclude list (executed=true).
// An empty input returns (nil, false) — no constraint.
func StringArrayFilter(arr gosql.NullableStringArray) (gosql.StringArray, bool) {
	if arr.Len() < 1 {
		return nil, false
	}
	executed := false
	narr := make(gosql.StringArray, 0, len(arr))
	for _, v := range arr {
		if len(v) > 0 && v[0] != '-' {
			narr = append(narr, v)
		}
	}
	if narr.Len() < 1 {
		narr = nil
		for _, v := range arr {
			if len(v) > 0 && v[0] == '-' {
				narr = append(narr, v[1:])
			}
		}
		executed = true
	}
	return narr, executed
}

// CountryFilter resolves ISO 3166-1 alpha-2 country codes to geo IDs and
// returns a sorted ID array with include/exclude polarity (see [StringArrayFilter]).
// Codes with a leading '-' contribute to an exclude list; all others form an
// include list. An unrecognised code resolves to ID 0.
func CountryFilter(arr gosql.NullableStringArray) (narr gosql.NullableOrderedNumberArray[uint64], executed bool) {
	var sarr gosql.StringArray
	if sarr, executed = StringArrayFilter(arr); sarr.Len() < 1 {
		return narr, executed
	}
	for _, cc := range sarr {
		narr = append(narr, uint64(gogeo.CountryByCode2(cc).ID))
	}
	narr.Sort()
	return narr, executed
}

// LanguageFilter resolves BCP-47 language codes to language IDs and returns a
// sorted ID array with include/exclude polarity (see [StringArrayFilter]).
// Codes with a leading '-' contribute to an exclude list; all others form an
// include list. An unrecognised code resolves to ID 0.
func LanguageFilter(arr gosql.NullableStringArray) (narr gosql.NullableOrderedNumberArray[uint64], executed bool) {
	var sarr gosql.StringArray
	if sarr, executed = StringArrayFilter(arr); sarr.Len() < 1 {
		return narr, executed
	}
	for _, lg := range sarr {
		narr = append(narr, uint64(languages.GetLanguageIdByCodeString(lg)))
	}
	narr.Sort()
	return narr, executed
}
