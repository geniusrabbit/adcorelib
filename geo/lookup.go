package geo

import "github.com/geniusrabbit/gogeo"

// CountryByID returns the country with the given ID, or the undefined country.
func CountryByID(id uint8) *Country {
	return gogeo.CountryByID(id)
}

// CountryByCode2 returns the country for an ISO-2 code.
// Unknown codes resolve to the undefined country (never nil).
func CountryByCode2(code string) *Country {
	return gogeo.CountryByCode2(code)
}

// CountryByCode2Bytes is the allocation-free ISO-2 lookup.
// Unknown codes resolve to the undefined country (never nil).
func CountryByCode2Bytes(code Code2) *Country {
	return gogeo.CountryByCode2Bytes(code)
}

// CountryByCode3 returns the country for an ISO-3 code.
// Empty and unknown codes resolve to the undefined country (never nil).
func CountryByCode3(code string) *Country {
	return gogeo.CountryByCode3(code)
}

// CountryCode2ByString returns a Code2 from an ISO-2 or ISO-3 string.
func CountryCode2ByString(code string) Code2 {
	return gogeo.CountryCode2ByString(code)
}

// ContinentByCode2 returns a continent by 2-letter code, or nil.
func ContinentByCode2(code string) *Continent {
	return gogeo.ContinentByCode2(code)
}

// RegionByID returns the region with the given ID, or the undefined region.
func RegionByID(id uint16) *Region {
	return gogeo.RegionByID(id)
}

// RegionByCode returns the region for an ISO 3166-2 code (for example "US-CA").
// "GB-*" and "UK-*" resolve to the same United Kingdom subdivisions.
// Unknown codes resolve to the undefined region (never nil).
func RegionByCode(code string) *Region {
	return gogeo.RegionByCode(code)
}

// RegionCodeByString parses an ISO 3166-2 code. Unknown values become undefined.
func RegionCodeByString(code string) RegionCode {
	return gogeo.RegionCodeByString(code)
}

// RegionByLatLng returns the nearest region with coordinates, or undefined.
func RegionByLatLng(lat, lng float32) *Region {
	return gogeo.RegionByLatLng(lat, lng)
}

// CountriesByContinent returns countries whose continent code matches.
// Unknown or empty continent codes yield nil.
func CountriesByContinent(code string) []Country {
	if code == "" {
		return nil
	}
	out := make([]Country, 0, 64)
	for _, country := range Countries {
		if country.ID == 0 {
			continue
		}
		if country.Continent() == code {
			out = append(out, country)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
