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

// IntArrayToUint array type
func IntArrayToUint64(arr []int) (res gosql.NullableOrderedNumberArray[uint64]) {
	if len(arr) < 1 {
		return
	}
	for _, v := range arr {
		res = append(res, uint64(v))
	}
	res.Sort()
	return
}

// IDArrayFilter array which could be or positive (include) or negative (exclude)
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

// IDArrayFilterAny array which could be or positive (include) or negative (exclude)
func IDArrayFilterAny(v any, panicMsg string) (gosql.NullableOrderedNumberArray[uint64], bool) {
	switch vl := v.(type) {
	case gosql.NullableOrderedNumberArray[int64]:
		return IDArrayFilter(gosql.NullableOrderedNumberArray[int64](vl))
	case gosql.NullableOrderedNumberArray[uint64]:
		return gosql.NullableOrderedNumberArray[uint64](vl), true
	case []int:
		return IntArrayToUint64(vl), true
	case []int64:
		return IDArrayFilter(gosql.NullableOrderedNumberArray[int64](vl))
	case []uint64:
		return gosql.NullableOrderedNumberArray[uint64](vl), true
	default:
		if panicMsg != "" {
			panic(panicMsg)
		}
	}
	return nil, false
}

// StringArrayFilter array
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

// CountryFilter array which could be or positive (include) or negative (exclude)
func CountryFilter(arr gosql.NullableStringArray) (narr gosql.NullableOrderedNumberArray[uint64], executed bool) {
	var sarr gosql.StringArray
	if sarr, executed = StringArrayFilter(arr); sarr.Len() < 1 {
		return narr, executed
	}
	// Countries filter
	for _, cc := range sarr {
		narr = append(narr, uint64(gogeo.CountryByCode2(cc).ID))
	}
	narr.Sort()

	return narr, executed
}

// LanguageFilter array which could be or positive (include) or negative (exclude)
func LanguageFilter(arr gosql.NullableStringArray) (narr gosql.NullableOrderedNumberArray[uint64], executed bool) {
	var sarr gosql.StringArray
	if sarr, executed = StringArrayFilter(arr); sarr.Len() < 1 {
		return narr, executed
	}
	// Countries filter
	for _, lg := range sarr {
		narr = append(narr, uint64(languages.GetLanguageIdByCodeString(lg)))
	}
	narr.Sort()

	return narr, executed
}
