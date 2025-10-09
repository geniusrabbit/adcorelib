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
	FieldTrafficSources
	FieldDomains
	FieldApps
	FieldZones
)

const (
	SecureAny int8 = iota
	SecureOnly
	SecureExclude
)

const (
	AdBlockAny int8 = iota
	AdBlockOnly
	AdBlockExclude
)

const (
	PrivateBrowsingAny int8 = iota
	PrivateBrowsingOnly
	PrivateBrowsingExclude
)

const (
	IPAny int8 = iota
	IPv4Only
	IPv6Only
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
	TrafficSources  gosql.NullableOrderedNumberArray[uint64]
	Domains         gosql.StringArray
	Apps            gosql.NullableOrderedNumberArray[uint64]
	Zones           gosql.NullableOrderedNumberArray[uint64]
	Secure          int8 // 0 - any, 1 - only, 2 - exclude
	AdBlock         int8 // 0 - any, 1 - only, 2 - exclude
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
		fl.DeviceTypes, positive = IDArrayFilterAny(data, "Invalid type for DeviceTypes")
	case FieldDevices:
		fl.Devices, positive = IDArrayFilterAny(data, "Invalid type for Devices")
	case FieldOS:
		fl.OS, positive = IDArrayFilterAny(data, "Invalid type for OS")
	case FieldBrowsers:
		fl.Browsers, positive = IDArrayFilterAny(data, "Invalid type for Browsers")
	case FieldCategories:
		fl.Categories, positive = IDArrayFilterAny(data, "Invalid type for Categories")
	case FieldCountries:
		switch vl := data.(type) {
		case []int64, []uint64:
			fl.Countries, positive = IDArrayFilterAny(vl, "")
		case gosql.StringArray:
			fl.Countries, positive = CountryFilter(gosql.NullableStringArray(vl))
		case gosql.NullableStringArray:
			fl.Countries, positive = CountryFilter(vl)
		}
	case FieldLanguages:
		switch vl := data.(type) {
		case []int64, []uint64:
			fl.Languages, positive = IDArrayFilterAny(vl, "")
		case gosql.StringArray:
			fl.Languages, positive = LanguageFilter(gosql.NullableStringArray(vl))
		case gosql.NullableStringArray:
			fl.Languages, positive = LanguageFilter(vl)
		}
	case FieldTrafficSources:
		fl.TrafficSources, positive = IDArrayFilterAny(data, "Invalid type for TrafficSources")
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
	formatList := t.Formats().List()
	found := len(formatList) < 1
	for _, f := range formatList {
		if found = fl.TestFormat(f); found {
			break
		}
	}

	if !found {
		return false
	}

	// ===========================================================================
	// Basic quick checks
	// ===========================================================================

	if fl.Secure != SecureAny && (fl.Secure == SecureOnly) != t.IsSecure() {
		return false
	}

	if fl.AdBlock != AdBlockAny && (fl.AdBlock == AdBlockOnly) != t.IsAdBlock() {
		return false
	}

	if fl.PrivateBrowsing != PrivateBrowsingAny && (fl.PrivateBrowsing == PrivateBrowsingOnly) != t.IsPrivateBrowsing() {
		return false
	}

	if fl.IP != IPAny && (fl.IP == IPv6Only) != t.IsIPv6() {
		return false
	}

	// ===========================================================================
	// Sources filter
	// ===========================================================================

	if !fl.checkUintArr(t.TrafficSourceID(), FieldTrafficSources, fl.TrafficSources) {
		return false
	}

	if !fl.checkUintArr(t.TargetID(), FieldZones, fl.Zones) {
		return false
	}

	if !fl.checkUintArr(t.AppID(), FieldApps, fl.Apps) {
		return false
	}

	if !fl.checkStringArr(t.Domain(), FieldDomains, fl.Domains) {
		return false
	}

	// ===========================================================================
	// General filters
	// ===========================================================================

	if !fl.checkUintArr(uint64(t.DeviceInfo().DeviceType), FieldDeviceTypes, fl.DeviceTypes) {
		return false
	}

	if !fl.checkUintArr(uint64(t.DeviceInfo().ID), FieldDevices, fl.Devices) {
		return false
	}

	if !fl.checkUintArr(uint64(t.OSInfo().ID), FieldOS, fl.OS) {
		return false
	}

	if !fl.checkUintArr(t.BrowserInfo().ID, FieldBrowsers, fl.Browsers) {
		return false
	}

	if !fl.multyCheckUintArr(t.Categories(), FieldCategories, fl.Categories) {
		return false
	}

	if !fl.checkUintArr(uint64(t.GeoInfo().ID), FieldCountries, fl.Countries) {
		return false
	}

	if !fl.checkUintArr(t.LanguageID(), FieldLanguages, fl.Languages) {
		return false
	}

	return true
}

// TestFormat available in filter
//
//go:inline
func (fl *BaseFilter) TestFormat(f *Format) bool {
	return len(fl.Formats) < 1 || fl.Formats.IndexOf(f.Codename) >= 0
}

//go:inline
func (fl *BaseFilter) checkUintArr(v uint64, off uint64, arr gosql.NullableOrderedNumberArray[uint64]) bool {
	return arr.Len() < 1 || (arr.IndexOf(v) >= 0) == (fl.excludeMask&(1<<off) == 0)
}

//go:inline
func (fl *BaseFilter) multyCheckUintArr(v []uint64, off uint64, arr gosql.NullableOrderedNumberArray[uint64]) bool {
	return arr.Len() < 1 || arr.OneOf(v) == (fl.excludeMask&(1<<off) == 0)
}

//go:inline
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
	fl.TrafficSources = fl.TrafficSources[:0]
	fl.Domains = fl.Domains[:0]
	fl.Apps = fl.Apps[:0]
	fl.Zones = fl.Zones[:0]
	fl.Secure = SecureAny
	fl.AdBlock = AdBlockAny
	fl.PrivateBrowsing = PrivateBrowsingAny
	fl.IP = IPAny
}
