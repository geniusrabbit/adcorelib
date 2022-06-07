//
// @project GeniusRabbit rotator 2017
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017
//

package types

import (
	"github.com/geniusrabbit/gogeo"
	"github.com/geniusrabbit/gosql/v2"

	"geniusrabbit.dev/corelib/i18n/languages"
)

// IntArrayToUint array type
func IntArrayToUint(arr []int) (res gosql.NullableOrderedNumberArray[uint]) {
	if len(arr) < 1 {
		return
	}
	for _, v := range arr {
		res = append(res, uint(v))
	}
	res.Sort()
	return
}

// IDArrayFilter array which could be or positive (include) or negative (exclude)
func IDArrayFilter(arr gosql.NullableOrderedNumberArray[int]) (narr gosql.NullableOrderedNumberArray[uint], executed bool) {
	if arr.Len() < 1 {
		return
	}

	subarr := arr.Map(func(v int) (int, bool) { return v, v > 0 })
	if subarr.Len() < 1 {
		subarr = arr.Map(func(v int) (int, bool) { return -v, v < 0 })
		executed = true
	}

	narr = IntArrayToUint(subarr)
	return
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
func CountryFilter(arr gosql.NullableStringArray) (narr gosql.NullableOrderedNumberArray[uint], executed bool) {
	var sarr gosql.StringArray
	if sarr, executed = StringArrayFilter(arr); sarr.Len() < 1 {
		return narr, executed
	}
	// Countries filter
	for _, cc := range sarr {
		narr = append(narr, uint(gogeo.CountryByCode2(cc).ID))
	}
	narr.Sort()

	return narr, executed
}

// LanguageFilter array which could be or positive (include) or negative (exclude)
func LanguageFilter(arr gosql.NullableStringArray) (narr gosql.NullableOrderedNumberArray[uint], executed bool) {
	var sarr gosql.StringArray
	if sarr, executed = StringArrayFilter(arr); sarr.Len() < 1 {
		return narr, executed
	}
	// Countries filter
	for _, lg := range sarr {
		narr = append(narr, languages.GetLanguageIdByCodeString(lg))
	}
	narr.Sort()

	return narr, executed
}
