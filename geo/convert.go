package geo

import (
	"strings"

	"github.com/geniusrabbit/gogeo"
	"github.com/geniusrabbit/gosql/v2"
)

// CountryCodes2IDs convert country codes to country IDs
func CountryCodes2IDs(codes []string) gosql.NullableOrderedNumberArray[uint64] {
	if len(codes) == 0 {
		return nil
	}
	result := make(gosql.NullableOrderedNumberArray[uint64], 0, len(codes))
	for _, cc := range codes {
		switch cc = strings.ToUpper(cc); cc {
		case "EU", "AS", "AF", "OC", "SA", "NA", "AN":
			for _, country := range gogeo.Countries {
				if country.Continent == cc {
					result = append(result, uint64(country.ID))
				}
			}
		default: // ** - as undefined
			result = append(result, uint64(gogeo.CountryByCode2(cc).ID))
		}
	}
	return result.Sort()
}
