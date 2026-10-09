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
	ErrFilterError                   = errtype.Error("filter")
	ErrFormatNotAllowed              = ErrFilterError.WithMessage("format not allowed")
	ErrSecureNotAllowed              = ErrFilterError.WithMessage("secure not allowed")
	ErrSecureOnlyNotAllowed          = ErrFilterError.WithMessage("secure only not allowed")
	ErrAdBlockNotAllowed             = ErrFilterError.WithMessage("ad block not allowed")
	ErrAdBlockOnlyNotAllowed         = ErrFilterError.WithMessage("ad block only not allowed")
	ErrPrivateBrowsingNotAllowed     = ErrFilterError.WithMessage("private browsing not allowed")
	ErrPrivateBrowsingOnlyNotAllowed = ErrFilterError.WithMessage("private browsing only allowed")
	ErrIPv6NotAllowed                = ErrFilterError.WithMessage("IPv6 not allowed")
	ErrIPv4NotAllowed                = ErrFilterError.WithMessage("IPv4 not allowed")
	ErrTrafficSourceNotAllowed       = ErrFilterError.WithMessage("traffic source not allowed")
	ErrTargetNotAllowed              = ErrFilterError.WithMessage("target not allowed")
	ErrAppNotAllowed                 = ErrFilterError.WithMessage("app not allowed")
	ErrDomainNotAllowed              = ErrFilterError.WithMessage("domain not allowed")
	ErrDeviceTypeNotAllowed          = ErrFilterError.WithMessage("device type not allowed")
	ErrDeviceIDNotAllowed            = ErrFilterError.WithMessage("device ID not allowed")
	ErrOSIDNotAllowed                = ErrFilterError.WithMessage("OS ID not allowed")
	ErrBrowserIDNotAllowed           = ErrFilterError.WithMessage("browser ID not allowed")
	ErrCategoriesNotAllowed          = ErrFilterError.WithMessage("categories not allowed")
	ErrCountryIDNotAllowed           = ErrFilterError.WithMessage("country ID not allowed")
	ErrLanguageIDNotAllowed          = ErrFilterError.WithMessage("language ID not allowed")
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
	FieldExtApps                                // external app/site ExtID strings (ext_app_id)
	FieldExtZones                               // external zone/tagid strings (ext_zone_id)
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
// List setters take the values and an include flag. include=true clears the
// bit; include=false sets it. Numeric lists are sorted so IndexOf/OneOf can
// find them. String values are stored as given (no '-' prefix).
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
	ExtApps             gosql.StringArray // publisher ext_app_id / site ExtID
	ExtZones            gosql.StringArray // publisher ext_zone_id / tagid
	Secure              int8              // SecureAny | SecureOnly | SecureExclude
	AdBlock             int8              // AdBlockAny | AdBlockOnly | AdBlockExclude
	PrivateBrowsing     int8              // PrivateBrowsingAny | PrivateBrowsingOnly | PrivateBrowsingExclude
	IP                  int8              // IPAny | IPv4Only | IPv6Only
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

// SetDeviceTypes sets the device-type filter. include=false marks the list as exclude.
func (fl *BaseFilter) SetDeviceTypes(arr []uint64, include bool) {
	fl.DeviceTypes = gosql.NullableOrderedNumberArray[uint64](arr).Sort()
	fl.setExcluded(FieldDeviceTypes, !include)
}

// SetDevices sets the device-model filter. include=false marks the list as exclude.
func (fl *BaseFilter) SetDevices(arr []uint64, include bool) {
	fl.Devices = gosql.NullableOrderedNumberArray[uint64](arr).Sort()
	fl.setExcluded(FieldDevices, !include)
}

// SetOS sets the operating-system filter. include=false marks the list as exclude.
func (fl *BaseFilter) SetOS(arr []uint64, include bool) {
	fl.OS = gosql.NullableOrderedNumberArray[uint64](arr).Sort()
	fl.setExcluded(FieldOS, !include)
}

// SetBrowsers sets the browser filter. include=false marks the list as exclude.
func (fl *BaseFilter) SetBrowsers(arr []uint64, include bool) {
	fl.Browsers = gosql.NullableOrderedNumberArray[uint64](arr).Sort()
	fl.setExcluded(FieldBrowsers, !include)
}

// SetCategories sets the IAB content-category filter. include=false marks the list as exclude.
func (fl *BaseFilter) SetCategories(arr []uint64, include bool) {
	fl.Categories = gosql.NullableOrderedNumberArray[uint64](arr).Sort()
	fl.setExcluded(FieldCategories, !include)
}

// SetCountries sets the country filter. Accepts numeric ID slices ([]uint64 /
// NullableOrderedNumberArray) or ISO 3166-1 alpha-2 code arrays. Codes are
// resolved to geo IDs. include=false marks the list as exclude.
func (fl *BaseFilter) SetCountries(data any, include bool) {
	switch vl := data.(type) {
	case []uint64:
		fl.Countries = gosql.NullableOrderedNumberArray[uint64](vl).Sort()
	case gosql.NullableOrderedNumberArray[uint64]:
		fl.Countries = vl.Sort()
	case gosql.StringArray:
		fl.Countries = countryCodesToIDs(gosql.NullableStringArray(vl))
	case gosql.NullableStringArray:
		fl.Countries = countryCodesToIDs(vl)
	}
	fl.setExcluded(FieldCountries, !include)
}

// SetLanguages sets the language filter. Accepts numeric ID slices ([]uint64 /
// NullableOrderedNumberArray) or BCP-47 code arrays. Codes are resolved to
// language IDs. include=false marks the list as exclude.
func (fl *BaseFilter) SetLanguages(data any, include bool) {
	switch vl := data.(type) {
	case []uint64:
		fl.Languages = gosql.NullableOrderedNumberArray[uint64](vl).Sort()
	case gosql.NullableOrderedNumberArray[uint64]:
		fl.Languages = vl.Sort()
	case gosql.StringArray:
		fl.Languages = languageCodesToIDs(gosql.NullableStringArray(vl))
	case gosql.NullableStringArray:
		fl.Languages = languageCodesToIDs(vl)
	}
	fl.setExcluded(FieldLanguages, !include)
}

// SetTrafficSources sets the traffic-source filter. include=false marks the list as exclude.
func (fl *BaseFilter) SetTrafficSources(arr []uint64, include bool) {
	fl.TrafficSources = gosql.NullableOrderedNumberArray[uint64](arr).Sort()
	fl.setExcluded(FieldTrafficSources, !include)
}

// SetDomains sets the domain / bundle-name filter. Values are stored as given.
// include=false marks the list as exclude.
func (fl *BaseFilter) SetDomains(arr gosql.NullableStringArray, include bool) {
	fl.Domains = gosql.StringArray(arr)
	fl.setExcluded(FieldDomains, !include)
}

// SetAppIDs sets the app filter. include=false marks the list as exclude.
func (fl *BaseFilter) SetAppIDs(arr []uint64, include bool) {
	fl.Apps = gosql.NullableOrderedNumberArray[uint64](arr).Sort()
	fl.setExcluded(FieldApps, !include)
}

// SetZoneIDs sets the zone filter. include=false marks the list as exclude.
func (fl *BaseFilter) SetZoneIDs(arr []uint64, include bool) {
	fl.Zones = gosql.NullableOrderedNumberArray[uint64](arr).Sort()
	fl.setExcluded(FieldZones, !include)
}

// SetExtApps sets the external app/site ExtID filter. Values are stored as given.
// include=false marks the list as exclude.
func (fl *BaseFilter) SetExtApps(arr gosql.NullableStringArray, include bool) {
	fl.ExtApps = gosql.StringArray(arr)
	fl.setExcluded(FieldExtApps, !include)
}

// SetExtZones sets the external zone/tagid filter. Values are stored as given.
// include=false marks the list as exclude.
func (fl *BaseFilter) SetExtZones(arr gosql.NullableStringArray, include bool) {
	fl.ExtZones = gosql.StringArray(arr)
	fl.setExcluded(FieldExtZones, !include)
}

// setExcluded records whether field is an exclude list.
// excluded=true sets the bit (pass when the value is not found);
// excluded=false clears it (pass when the value is found).
func (fl *BaseFilter) setExcluded(field uint64, excluded bool) {
	if excluded {
		fl.excludeMask |= 1 << field
	} else {
		fl.excludeMask &= ^(1 << field)
	}
}

// IsExcluded reports whether field is an exclude list. An empty list is not a
// constraint either way; this only reads the polarity bit.
func (fl BaseFilter) IsExcluded(field FilterField) bool {
	return fl.excludeMask&(1<<field) != 0
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

	if !fl.checkStringArr(extAppIDs(t), FieldExtApps, fl.ExtApps) {
		return ErrAppNotAllowed
	}

	if !fl.checkStringArr(extZoneIDs(t), FieldExtZones, fl.ExtZones) {
		return ErrTargetNotAllowed
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

	if fl.OS.Len() > 0 {
		os := t.OSInfo()
		if !fl.checkUintArr(uint64(os.ID), FieldOS, fl.OS) {
			if os.VersionID == 0 || !fl.checkUintArr(uint64(os.VersionID), FieldOS, fl.OS) {
				return ErrOSIDNotAllowed
			}
		}
	}

	if fl.Browsers.Len() > 0 {
		brw := t.BrowserInfo()
		if !fl.checkUintArr(brw.ID, FieldBrowsers, fl.Browsers) {
			if brw.VersionID == 0 || !fl.checkUintArr(brw.VersionID, FieldBrowsers, fl.Browsers) {
				return ErrBrowserIDNotAllowed
			}
		}
	}

	if !fl.multyCheckUintArr(t.CategoryIDs(), FieldCategories, fl.Categories) {
		return ErrCategoriesNotAllowed
	}

	if !fl.checkUintArr(t.GeoInfo().CountryID(), FieldCountries, fl.Countries) {
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
	return arr.Len() < 1 || (arr.IndexOf(v) >= 0) == !fl.IsExcluded(off)
}

// multyCheckUintArr is like checkUintArr but accepts a slice of values and
// passes when at least one of them satisfies the filter (OneOf semantics).
//
//go:inline
func (fl *BaseFilter) multyCheckUintArr(v []uint64, off uint64, arr gosql.NullableOrderedNumberArray[uint64]) bool {
	return arr.Len() < 1 || arr.OneOf(v) == !fl.IsExcluded(off)
}

// checkStringArr is the string-slice equivalent of multyCheckUintArr.
//
//go:inline
func (fl *BaseFilter) checkStringArr(v []string, off uint64, arr gosql.StringArray) bool {
	return arr.Len() < 1 || arr.OneOf(v) == !fl.IsExcluded(off)
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
	fl.ExtApps = fl.ExtApps[:0]
	fl.ExtZones = fl.ExtZones[:0]
	fl.Secure = SecureAny
	fl.AdBlock = AdBlockAny
	fl.PrivateBrowsing = PrivateBrowsingAny
	fl.IP = IPAny
}

func extAppIDs(t TargetPointer) []string {
	var ids []string
	if app := t.AppInfo(); app != nil && app.ExtID != "" {
		ids = append(ids, app.ExtID)
	}
	if site := t.SiteInfo(); site != nil && site.ExtID != "" {
		ids = append(ids, site.ExtID)
	}
	return ids
}

func extZoneIDs(t TargetPointer) []string {
	if id := t.ExtarnalTargetID(); id != "" {
		return []string{id}
	}
	return nil
}
