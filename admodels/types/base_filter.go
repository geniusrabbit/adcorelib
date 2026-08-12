//
// @project GeniusRabbit corelib 2017 - 2018, 2022, 2026
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018, 2022, 2026
//

package types

import (
	"github.com/geniusrabbit/adcorelib/errtype"
	"github.com/geniusrabbit/gosql/v2"
)

var (
	ErrFormatNotAllowed              = errtype.Error("format not allowed")
	ErrSecureNotAllowed              = errtype.Error("secure not allowed")
	ErrSecureOnlyNotAllowed          = errtype.Error("secure only not allowed")
	ErrAdBlockNotAllowed             = errtype.Error("ad block not allowed")
	ErrAdBlockOnlyNotAllowed         = errtype.Error("ad block only not allowed")
	ErrPrivateBrowsingNotAllowed     = errtype.Error("private browsing not allowed")
	ErrPrivateBrowsingOnlyNotAllowed = errtype.Error("private browsing only not allowed")
	ErrIPv6NotAllowed                = errtype.Error("IPv6 not allowed")
	ErrIPv4NotAllowed                = errtype.Error("IPv4 not allowed")
	ErrTrafficSourceNotAllowed       = errtype.Error("traffic source not allowed")
	ErrTargetNotAllowed              = errtype.Error("target not allowed")
	ErrAppNotAllowed                 = errtype.Error("app not allowed")
	ErrDomainNotAllowed              = errtype.Error("domain not allowed")
	ErrDeviceTypeNotAllowed          = errtype.Error("device type not allowed")
	ErrDeviceIDNotAllowed            = errtype.Error("device ID not allowed")
	ErrOSIDNotAllowed                = errtype.Error("OS ID not allowed")
	ErrBrowserIDNotAllowed           = errtype.Error("browser ID not allowed")
	ErrCategoriesNotAllowed          = errtype.Error("categories not allowed")
	ErrCountryIDNotAllowed           = errtype.Error("country ID not allowed")
	ErrLanguageIDNotAllowed          = errtype.Error("language ID not allowed")
)

// FilterField identifies a filter dimension in [BaseFilter].
// Prefer the typed setter methods (SetFormats, SetDeviceTypes, etc.) over the
// legacy Set() dispatcher — they provide compile-time type safety.
type FilterField = uint64

// Base filter fields.
const (
	FieldFormat              FilterField = iota // ad format codenames (Formats / InterstitialFormats)
	FieldDeviceTypes                            // device category (phone, tablet, desktop…)
	FieldDevices                                // specific device model IDs
	FieldOS                                     // operating system IDs
	FieldBrowsers                               // browser IDs
	FieldCategories                             // IAB content category IDs
	FieldCountries                              // country IDs (or ISO 3166-1 alpha-2 codes)
	FieldLanguages                              // language IDs (or BCP-47 codes)
	FieldTrafficSources                         // traffic-source IDs
	FieldDomains                                // domain / bundle-name allowlist or blocklist
	FieldApps                                   // app IDs
	FieldZones                                  // zone IDs
	FieldInterstitialFormats                    // ad format codenames for interstitial requests
)

// Secure request filter values.
const (
	SecureAny     int8 = iota // do not filter by HTTPS
	SecureOnly                // HTTPS only
	SecureExclude             // HTTP only
)

// AdBlock filter values.
const (
	AdBlockAny     int8 = iota // do not filter by ad-block presence
	AdBlockOnly                // ad-block traffic only
	AdBlockExclude             // exclude ad-block traffic
)

// PrivateBrowsing filter values.
const (
	PrivateBrowsingAny     int8 = iota // do not filter by private-browsing mode
	PrivateBrowsingOnly                // private-browsing traffic only
	PrivateBrowsingExclude             // exclude private-browsing traffic
)

// IP version filter values.
const (
	IPAny    int8 = iota // do not filter by IP version
	IPv4Only             // IPv4 traffic only
	IPv6Only             // IPv6 traffic only
)

// BaseFilter holds the targeting criteria that a [TargetPointer] must satisfy.
//
// Array-type fields can act as either an include list (allow only matching
// values) or an exclude list (reject matching values). The polarity for each
// field is recorded in the internal excludeMask bitset:
//
//   - bit CLEAR → include list: request passes when the value IS found
//   - bit SET   → exclude list: request passes when the value is NOT found
//   - empty list → no constraint (field is ignored)
//
// For signed-integer source arrays (int64 convention): positive values build
// an include list, negative values build an exclude list (absolute values are
// stored; see IDArrayFilter). Use [BaseFilter.SetPositive] to set the polarity
// directly when building from uint64 or pre-processed data.
//
// Format selection is context-aware:
//   - Non-interstitial requests are matched against Formats.
//   - Interstitial requests are matched against InterstitialFormats when it is
//     non-empty; otherwise Formats is used as the fallback.
type BaseFilter struct {
	// excludeMask is a per-field polarity bitset.
	// Bit n CLEAR = field n is an include list; bit n SET = exclude list.
	excludeMask         uint64
	Formats             gosql.StringArray // format codenames for non-interstitial requests
	InterstitialFormats gosql.StringArray // format codenames for interstitial requests (overrides Formats when set)
	DeviceTypes         gosql.NullableOrderedNumberArray[uint64]
	Devices             gosql.NullableOrderedNumberArray[uint64]
	OS                  gosql.NullableOrderedNumberArray[uint64]
	Browsers            gosql.NullableOrderedNumberArray[uint64]
	Categories          gosql.NullableOrderedNumberArray[uint64]
	Countries           gosql.NullableOrderedNumberArray[uint64]
	Languages           gosql.NullableOrderedNumberArray[uint64]
	TrafficSources      gosql.NullableOrderedNumberArray[uint64]
	Domains             gosql.StringArray
	Apps                gosql.NullableOrderedNumberArray[uint64]
	Zones               gosql.NullableOrderedNumberArray[uint64]
	Secure              int8 // SecureAny | SecureOnly | SecureExclude
	AdBlock             int8 // AdBlockAny | AdBlockOnly | AdBlockExclude
	PrivateBrowsing     int8 // PrivateBrowsingAny | PrivateBrowsingOnly | PrivateBrowsingExclude
	IP                  int8 // IPAny | IPv4Only | IPv6Only
}

// ---------------------------------------------------------------------------
// Typed setters — prefer these over the generic Set() dispatcher.
// ---------------------------------------------------------------------------

// SetFormats sets the format codename allowlist used for non-interstitial
// requests. An empty slice removes the format constraint.
func (fl *BaseFilter) SetFormats(arr []string) {
	fl.Formats = arr
}

// SetInterstitialFormats sets the format codename allowlist used when the
// request is interstitial (IsInterstitial() == true). When non-empty it
// overrides Formats for such requests; an empty slice falls back to Formats.
func (fl *BaseFilter) SetInterstitialFormats(arr []string) {
	fl.InterstitialFormats = arr
}

// SetDeviceTypes sets the device-type filter from a signed ([]int64) or
// unsigned ([]uint64) slice. Negative int64 values denote exclusion.
func (fl *BaseFilter) SetDeviceTypes(data any) {
	var positive bool
	fl.DeviceTypes, positive = IDArrayFilterAny(data, "invalid type for DeviceTypes")
	fl.SetPositive(FieldDeviceTypes, positive)
}

// SetDevices sets the device-model filter from a signed ([]int64) or
// unsigned ([]uint64) slice. Negative int64 values denote exclusion.
func (fl *BaseFilter) SetDevices(data any) {
	var positive bool
	fl.Devices, positive = IDArrayFilterAny(data, "invalid type for Devices")
	fl.SetPositive(FieldDevices, positive)
}

// SetOS sets the operating-system filter from a signed ([]int64) or
// unsigned ([]uint64) slice. Negative int64 values denote exclusion.
func (fl *BaseFilter) SetOS(data any) {
	var positive bool
	fl.OS, positive = IDArrayFilterAny(data, "invalid type for OS")
	fl.SetPositive(FieldOS, positive)
}

// SetBrowsers sets the browser filter from a signed ([]int64) or unsigned
// ([]uint64) slice. Negative int64 values denote exclusion.
func (fl *BaseFilter) SetBrowsers(data any) {
	var positive bool
	fl.Browsers, positive = IDArrayFilterAny(data, "invalid type for Browsers")
	fl.SetPositive(FieldBrowsers, positive)
}

// SetCategories sets the IAB content-category filter from a signed ([]int64)
// or unsigned ([]uint64) slice. Negative int64 values denote exclusion.
func (fl *BaseFilter) SetCategories(data any) {
	var positive bool
	fl.Categories, positive = IDArrayFilterAny(data, "invalid type for Categories")
	fl.SetPositive(FieldCategories, positive)
}

// SetCountries sets the country filter. Accepts numeric ID slices ([]int64,
// []uint64) or ISO 3166-1 alpha-2 code arrays (gosql.StringArray /
// gosql.NullableStringArray). Prefix a code with '-' to exclude that country.
func (fl *BaseFilter) SetCountries(data any) {
	var positive bool
	switch vl := data.(type) {
	case []int64, []uint64:
		fl.Countries, positive = IDArrayFilterAny(vl, "")
	case gosql.StringArray:
		fl.Countries, positive = CountryFilter(gosql.NullableStringArray(vl))
	case gosql.NullableStringArray:
		fl.Countries, positive = CountryFilter(vl)
	}
	fl.SetPositive(FieldCountries, positive)
}

// SetLanguages sets the language filter. Accepts numeric ID slices ([]int64,
// []uint64) or BCP-47 code arrays (gosql.StringArray / gosql.NullableStringArray).
// Prefix a code with '-' to exclude that language.
func (fl *BaseFilter) SetLanguages(data any) {
	var positive bool
	switch vl := data.(type) {
	case []int64, []uint64:
		fl.Languages, positive = IDArrayFilterAny(vl, "")
	case gosql.StringArray:
		fl.Languages, positive = LanguageFilter(gosql.NullableStringArray(vl))
	case gosql.NullableStringArray:
		fl.Languages, positive = LanguageFilter(vl)
	}
	fl.SetPositive(FieldLanguages, positive)
}

// SetTrafficSources sets the traffic-source filter from a signed ([]int64) or
// unsigned ([]uint64) slice. Negative int64 values denote exclusion.
func (fl *BaseFilter) SetTrafficSources(data any) {
	var positive bool
	fl.TrafficSources, positive = IDArrayFilterAny(data, "invalid type for TrafficSources")
	fl.SetPositive(FieldTrafficSources, positive)
}

// SetDomains sets the domain / bundle-name filter. Values with a leading '-'
// are treated as an exclusion list; all others form an inclusion list.
func (fl *BaseFilter) SetDomains(arr gosql.NullableStringArray) {
	var positive bool
	fl.Domains, positive = StringArrayFilter(arr)
	fl.SetPositive(FieldDomains, positive)
}

// SetApps sets the app filter from a signed integer slice.
// Positive values form an include list; negative values (absolute value stored)
// form an exclude list.
func (fl *BaseFilter) SetApps(arr []int64) {
	var positive bool
	fl.Apps, positive = IDArrayFilter(gosql.NullableOrderedNumberArray[int64](arr))
	fl.SetPositive(FieldApps, positive)
}

// SetAppIDs sets the app filter as an explicit include list of unsigned IDs.
func (fl *BaseFilter) SetAppIDs(arr []uint64) {
	fl.Apps = arr
	fl.SetPositive(FieldApps, false) // false = include mode (bit clear)
}

// SetZones sets the zone filter from a signed integer slice.
// Positive values form an include list; negative values (absolute value stored)
// form an exclude list.
func (fl *BaseFilter) SetZones(arr []int64) {
	var positive bool
	fl.Zones, positive = IDArrayFilter(gosql.NullableOrderedNumberArray[int64](arr))
	fl.SetPositive(FieldZones, positive)
}

// SetZoneIDs sets the zone filter as an explicit include list of unsigned IDs.
func (fl *BaseFilter) SetZoneIDs(arr []uint64) {
	fl.Zones = arr
	fl.SetPositive(FieldZones, false) // false = include mode (bit clear)
}

// SetPositive records the include/exclude polarity for the given field in
// excludeMask. Despite the name, positive=true activates exclude mode
// (sets the corresponding bit); positive=false activates include mode
// (clears the bit). The parameter mirrors the "executed" return value of
// [IDArrayFilter]: true means the exclude path was taken.
func (fl *BaseFilter) SetPositive(field uint64, positive bool) {
	if positive {
		fl.excludeMask |= 1 << field // exclude mode: pass when NOT found
	} else {
		fl.excludeMask &= ^(1 << field) // include mode: pass when found
	}
}

// Test evaluates whether the target satisfies every configured filter.
// Checks are applied in order: format → tristate flags (secure, adblock,
// private browsing, IP version) → source identifiers (traffic source, zone,
// app, domain) → device / OS / browser / geo / language.
// Any single failing check short-circuits and returns a preallocated sentinel error.
func (fl *BaseFilter) Test(t TargetPointer) error {
	formatList := t.Formats().List()
	// Select the active format allowlist: InterstitialFormats takes precedence
	// when the request is interstitial and the list is configured.
	activeFormats := fl.Formats
	if t.IsInterstitial() && fl.InterstitialFormats.Len() > 0 {
		activeFormats = fl.InterstitialFormats
	}

	found := len(formatList) < 1 || activeFormats.Len() <= 0
	if !found {
		for _, f := range formatList {
			if found = activeFormats.IndexOf(f.Codename) >= 0; found {
				break
			}
		}
	}

	if !found {
		return ErrFormatNotAllowed
	}

	// ===========================================================================
	// Basic quick checks
	// ===========================================================================

	if fl.Secure != SecureAny && (fl.Secure == SecureOnly) != t.IsSecure() {
		if fl.Secure == SecureOnly {
			return ErrSecureOnlyNotAllowed
		}
		return ErrSecureNotAllowed
	}

	if fl.AdBlock != AdBlockAny && (fl.AdBlock == AdBlockOnly) != t.IsAdBlock() {
		if fl.AdBlock == AdBlockOnly {
			return ErrAdBlockOnlyNotAllowed
		}
		return ErrAdBlockNotAllowed
	}

	if fl.PrivateBrowsing != PrivateBrowsingAny && (fl.PrivateBrowsing == PrivateBrowsingOnly) != t.IsPrivateBrowsing() {
		if fl.PrivateBrowsing == PrivateBrowsingOnly {
			return ErrPrivateBrowsingOnlyNotAllowed
		}
		return ErrPrivateBrowsingNotAllowed
	}

	if fl.IP != IPAny && (fl.IP == IPv6Only) != t.IsIPv6() {
		if fl.IP == IPv6Only {
			return ErrIPv4NotAllowed
		}
		return ErrIPv6NotAllowed
	}

	// ===========================================================================
	// Sources filter
	// ===========================================================================

	if !fl.checkUintArr(t.TrafficSourceID(), FieldTrafficSources, fl.TrafficSources) {
		return ErrTrafficSourceNotAllowed
	}

	if !fl.checkUintArr(t.TargetID(), FieldZones, fl.Zones) {
		return ErrTargetNotAllowed
	}

	if !fl.checkUintArr(t.AppID(), FieldApps, fl.Apps) {
		return ErrAppNotAllowed
	}

	if !fl.checkStringArr(t.Domain(), FieldDomains, fl.Domains) {
		return ErrDomainNotAllowed
	}

	// ===========================================================================
	// General filters
	// ===========================================================================

	if !fl.checkUintArr(uint64(t.DeviceInfo().DeviceType), FieldDeviceTypes, fl.DeviceTypes) {
		return ErrDeviceTypeNotAllowed
	}

	if !fl.checkUintArr(uint64(t.DeviceInfo().ID), FieldDevices, fl.Devices) {
		return ErrDeviceIDNotAllowed
	}

	if !fl.checkUintArr(uint64(t.OSInfo().ID), FieldOS, fl.OS) {
		return ErrOSIDNotAllowed
	}

	if !fl.checkUintArr(t.BrowserInfo().ID, FieldBrowsers, fl.Browsers) {
		return ErrBrowserIDNotAllowed
	}

	if !fl.multyCheckUintArr(t.Categories(), FieldCategories, fl.Categories) {
		return ErrCategoriesNotAllowed
	}

	if !fl.checkUintArr(uint64(t.GeoInfo().ID), FieldCountries, fl.Countries) {
		return ErrCountryIDNotAllowed
	}

	if !fl.checkUintArr(t.LanguageID(), FieldLanguages, fl.Languages) {
		return ErrLanguageIDNotAllowed
	}

	return nil
}

// TestFormat reports whether format f is permitted by the Formats allowlist.
// An empty Formats slice permits every format.
// Note: does not consider InterstitialFormats; call [BaseFilter.Test] for
// context-aware format matching that respects [TargetPointer.IsInterstitial].
//
//go:inline
func (fl *BaseFilter) TestFormat(f *Format) bool {
	return len(fl.Formats) == 0 || fl.Formats.IndexOf(f.Codename) >= 0
}

// TestInterstitialFormat reports whether format f is permitted by the InterstitialFormats allowlist.
// An empty InterstitialFormats slice permits every format.
// Note: does not consider Formats; call [BaseFilter.Test] for context-aware
// format matching that respects [TargetPointer.IsInterstitial].
//
//go:inline
func (fl *BaseFilter) TestInterstitialFormat(f *Format) bool {
	return len(fl.InterstitialFormats) == 0 || fl.InterstitialFormats.IndexOf(f.Codename) >= 0
}

// checkUintArr reports whether the single value v satisfies the filter arr
// under the include/exclude polarity stored at bit off of excludeMask.
// An empty arr always passes.
//
//go:inline
func (fl *BaseFilter) checkUintArr(v uint64, off uint64, arr gosql.NullableOrderedNumberArray[uint64]) bool {
	return arr.Len() < 1 || (arr.IndexOf(v) >= 0) == (fl.excludeMask&(1<<off) == 0)
}

// multyCheckUintArr is like checkUintArr but accepts a slice of values and
// passes when at least one of them satisfies the filter (OneOf semantics).
//
//go:inline
func (fl *BaseFilter) multyCheckUintArr(v []uint64, off uint64, arr gosql.NullableOrderedNumberArray[uint64]) bool {
	return arr.Len() < 1 || arr.OneOf(v) == (fl.excludeMask&(1<<off) == 0)
}

// checkStringArr is the string-slice equivalent of multyCheckUintArr.
//
//go:inline
func (fl *BaseFilter) checkStringArr(v []string, off uint64, arr gosql.StringArray) bool {
	return arr.Len() < 1 || arr.OneOf(v) == (fl.excludeMask&(1<<off) == 0)
}

// Reset clears all filter fields and resets excludeMask and tristate flags to
// their zero / Any defaults. Underlying array memory is reused where possible.
func (fl *BaseFilter) Reset() {
	fl.excludeMask = 0
	fl.Formats = fl.Formats[:0]
	fl.InterstitialFormats = fl.InterstitialFormats[:0]
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
