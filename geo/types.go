package geo

import "github.com/geniusrabbit/gogeo"

// Type aliases keep gogeo methods and stay assignable to udetect.Geo fields.
type (
	CountryCode2 = gogeo.CountryCode2
	RegionCode   = gogeo.RegionCode
	Continent    = gogeo.Continent
	Country      = gogeo.Country
	Region       = gogeo.Region
	Coordinates  = gogeo.Coordinates
	TimeZone     = gogeo.TimeZone
)

const (
	UndefinedCountryCodeISO2 = gogeo.UndefinedCountryCodeISO2
	UndefinedCountryCodeISO3 = gogeo.UndefinedCountryCodeISO3
	UndefinedRegionISO3166   = gogeo.UndefinedRegionISO3166
)

var (
	UndefinedCountryCode2        = gogeo.UndefinedCountryCode2
	UndefinedRegionCode          = gogeo.UndefinedRegionCode
	ErrCode2InvalidScanType      = gogeo.ErrCode2InvalidScanType
	ErrCode2InvalidValueSize     = gogeo.ErrCode2InvalidValueSize
	ErrRegionCodeInvalidScanType = gogeo.ErrRegionCodeInvalidScanType

	// Continents is the static continent catalog (IDs 1..7).
	Continents = gogeo.Continents
	// Countries is the static country catalog. Index equals Country.ID.
	Countries = gogeo.Countries
	// Regions is the static ISO 3166-2 catalog. Index equals Region.ID.
	Regions = gogeo.Regions
)
