//
// @project GeniusRabbit corelib 2017 - 2018, 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018, 2022
//

package types

import (
	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/gosql/v2"
)

// Base filter fields
const (
	FieldFormat = iota
	FieldDeviceTypes
	FieldDevices
	FieldOS
	FieldBrowsers
	FieldCategories
	FieldCountries
	FieldLanguages
	FieldDomains
	FieldApps
	FieldZones
)

// BaseFilter object
type BaseFilter struct {
	excludeMask     uint64
	Formats         gosql.StringArray
	DeviceTypes     gosql.NullableOrderedNumberArray[uint64]
	Devices         gosql.NullableOrderedNumberArray[uint64]
	OS              gosql.NullableOrderedNumberArray[uint64]
	Browsers        gosql.NullableOrderedNumberArray[uint64]
	Categories      gosql.NullableOrderedNumberArray[uint64]
	Countries       gosql.NullableOrderedNumberArray[uint64]
	Languages       gosql.NullableOrderedNumberArray[uint64]
	Domains         gosql.StringArray
	Apps            gosql.NullableOrderedNumberArray[uint64]
	Zones           gosql.NullableOrderedNumberArray[uint64]
	Secure          int8 // 0 - any, 1 - only, 2 - exclude
	Adblock         int8 // 0 - any, 1 - only, 2 - exclude
	PrivateBrowsing int8 // 0 - any, 1 - only, 2 - exclude
	IP              int8 // 0 - any, 1 - IPv4, 2 - IPv6
}

// Set filter item
func (fl *BaseFilter) Set(field uint64, data any) {
	var positive bool
	switch field {
	case FieldFormat:
		fl.Formats, _ = data.(gosql.StringArray)
	case FieldDeviceTypes:
		fl.DeviceTypes, positive = IDArrayFilter(data.(gosql.NullableOrderedNumberArray[int64]))
	case FieldDevices:
		fl.Devices, positive = IDArrayFilter(data.(gosql.NullableOrderedNumberArray[int64]))
	case FieldOS:
		fl.OS, positive = IDArrayFilter(data.(gosql.NullableOrderedNumberArray[int64]))
	case FieldBrowsers:
		fl.Browsers, positive = IDArrayFilter(data.(gosql.NullableOrderedNumberArray[int64]))
	case FieldCategories:
		fl.Categories, positive = IDArrayFilter(data.(gosql.NullableOrderedNumberArray[int64]))
	case FieldCountries:
		switch vl := data.(type) {
		case gosql.NullableOrderedNumberArray[int64]:
			fl.Countries, positive = IDArrayFilter(vl)
		case gosql.StringArray:
			fl.Countries, positive = CountryFilter(gosql.NullableStringArray(vl))
		case gosql.NullableStringArray:
			fl.Countries, positive = CountryFilter(vl)
		}
	case FieldLanguages:
		switch vl := data.(type) {
		case gosql.NullableOrderedNumberArray[int64]:
			fl.Languages, positive = IDArrayFilter(vl)
		case gosql.StringArray:
			fl.Languages, positive = LanguageFilter(gosql.NullableStringArray(vl))
		case gosql.NullableStringArray:
			fl.Languages, positive = LanguageFilter(vl)
		}
	case FieldApps:
		fl.Apps, positive = IDArrayFilter(gocast.AnySlice[int64](data))
	case FieldZones:
		fl.Zones, positive = IDArrayFilter(gocast.AnySlice[int64](data))
	case FieldDomains:
		switch arr := data.(type) {
		case gosql.StringArray:
			fl.Domains, positive = StringArrayFilter(gosql.NullableStringArray(arr))
		case gosql.NullableStringArray:
			fl.Domains, positive = StringArrayFilter(arr)
		}
	}
	fl.SetPositive(field, positive)
}

// SetPositive field state
func (fl *BaseFilter) SetPositive(field uint64, positive bool) {
	if positive {
		fl.excludeMask |= 1 << field
	} else {
		fl.excludeMask &= ^(1 << field)
	}
}

// Test filter items
func (fl *BaseFilter) Test(t TargetPointer) bool {
	found := len(t.Formats()) < 1
	for _, f := range t.Formats() {
		if found = fl.TestFormat(f); found {
			break
		}
	}
	return found &&
		(fl.Secure == 0 || (fl.Secure == 1) == t.IsSecure()) &&
		(fl.Adblock == 0 || (fl.Adblock == 1) == t.IsAdblock()) &&
		(fl.PrivateBrowsing == 0 || (fl.PrivateBrowsing == 1) == t.IsPrivateBrowsing()) &&
		(fl.IP == 0 || (fl.IP == 2) == t.IsIPv6()) &&
		fl.checkUintArr(t.DeviceType(), FieldDeviceTypes, fl.DeviceTypes) &&
		fl.checkUintArr(t.DeviceID(), FieldDevices, fl.Devices) &&
		fl.checkUintArr(t.OSID(), FieldOS, fl.OS) &&
		fl.checkUintArr(t.BrowserID(), FieldBrowsers, fl.Browsers) &&
		fl.multyCheckUintArr(t.Categories(), FieldCategories, fl.Categories) &&
		fl.checkUintArr(t.GeoID(), FieldCountries, fl.Countries) &&
		fl.checkUintArr(t.LanguageID(), FieldLanguages, fl.Languages) &&
		fl.checkUintArr(t.TargetID(), FieldZones, fl.Zones) &&
		fl.checkUintArr(t.AppID(), FieldApps, fl.Apps) &&
		fl.checkStringArr(t.Domain(), FieldDomains, fl.Domains)
}

// TestFormat available in filter
func (fl *BaseFilter) TestFormat(f *Format) bool {
	return len(fl.Formats) < 1 || fl.Formats.IndexOf(f.Codename) >= 0
}

func (fl *BaseFilter) checkUintArr(v uint64, off uint64, arr gosql.NullableOrderedNumberArray[uint64]) bool {
	return arr.Len() < 1 || (arr.IndexOf(v) >= 0) == (fl.excludeMask&(1<<off) == 0)
}

func (fl *BaseFilter) multyCheckUintArr(v []uint64, off uint64, arr gosql.NullableOrderedNumberArray[uint64]) bool {
	return arr.Len() < 1 || arr.OneOf(v) == (fl.excludeMask&(1<<off) == 0)
}

func (fl *BaseFilter) checkStringArr(v []string, off uint64, arr gosql.StringArray) bool {
	return arr.Len() < 1 || arr.OneOf(v) == (fl.excludeMask&(1<<off) == 0)
}

// Reset filter object
func (fl *BaseFilter) Reset() {
	fl.excludeMask = 0
	fl.Formats = fl.Formats[:0]
	fl.DeviceTypes = fl.DeviceTypes[:0]
	fl.Devices = fl.Devices[:0]
	fl.OS = fl.OS[:0]
	fl.Browsers = fl.Browsers[:0]
	fl.Categories = fl.Categories[:0]
	fl.Countries = fl.Countries[:0]
	fl.Languages = fl.Languages[:0]
	fl.Domains = fl.Domains[:0]
	fl.Apps = fl.Apps[:0]
	fl.Zones = fl.Zones[:0]
	fl.Secure = 0
	fl.Adblock = 0
	fl.PrivateBrowsing = 0
	fl.IP = 0
}
