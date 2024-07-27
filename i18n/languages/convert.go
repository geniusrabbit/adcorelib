package languages

import "github.com/geniusrabbit/gosql/v2"

// LangCodes2IDs convert language codes to language IDs
func LangCodes2IDs(codes []string) gosql.NullableOrderedNumberArray[uint64] {
	if len(codes) == 0 {
		return nil
	}
	result := make(gosql.NullableOrderedNumberArray[uint64], 0, len(codes))
	for _, lg := range codes {
		result = append(result, uint64(GetLanguageIdByCodeString(lg)))
	}
	return result.Sort()
}
