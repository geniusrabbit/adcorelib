package geo

import (
	"strings"

	"github.com/geniusrabbit/gosql/v2"
)

// continentCodes are the 2-letter continent codes used by Country.Continent().
//
// EU is both a continent code and a GeoIP extra country in the catalog.
// This helper treats EU/AS/AF/OC/SA/NA/AN as continents and expands them to
// every matching country ID. Live targeting via CountryFilter does not use
// this expansion and resolves EU as the GeoIP country instead.
var continentCodes = map[string]struct{}{
	"EU": {},
	"AS": {},
	"AF": {},
	"OC": {},
	"SA": {},
	"NA": {},
	"AN": {},
}

// CountryCodes2IDs convert country codes to country IDs.
// Continent codes EU, AS, AF, OC, SA, NA, AN expand to all countries on that continent.
func CountryCodes2IDs(codes []string) gosql.NullableOrderedNumberArray[uint64] {
	if len(codes) == 0 {
		return nil
	}
	result := make(gosql.NullableOrderedNumberArray[uint64], 0, len(codes))
	for _, cc := range codes {
		cc = strings.ToUpper(cc)
		if _, ok := continentCodes[cc]; ok {
			for _, country := range CountriesByContinent(cc) {
				result = append(result, uint64(country.ID))
			}
			continue
		}
		result = append(result, uint64(CountryByCode2(cc).ID))
	}
	return result.Sort()
}

// RegionCodes2IDs convert ISO 3166-2 region codes to region IDs.
// Unrecognised codes resolve to ID 0 (undefined).
func RegionCodes2IDs(codes []string) gosql.NullableOrderedNumberArray[uint64] {
	if len(codes) == 0 {
		return nil
	}
	result := make(gosql.NullableOrderedNumberArray[uint64], 0, len(codes))
	for _, code := range codes {
		result = append(result, uint64(RegionByCode(strings.ToUpper(code)).ID))
	}
	return result.Sort()
}
