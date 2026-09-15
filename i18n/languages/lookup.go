package languages

var (
	codeIndex    [26][26]uint8
	languageISO2 []string
)

func init() {
	languageISO2 = make([]string, len(languages))
	languageISO2[0] = UndefinedISO2
	for i := 1; i < len(languages); i++ {
		c := languages[i].Code
		languageISO2[i] = string([]byte{c[0], c[1]})
		b0, b1 := asciiLower(c[0]), asciiLower(c[1])
		if b0 >= 'a' && b0 <= 'z' && b1 >= 'a' && b1 <= 'z' {
			codeIndex[b0-'a'][b1-'a'] = uint8(i)
		}
	}
}

func asciiLower(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		return b + ('a' - 'A')
	}
	return b
}

func lookup(b0, b1 byte) *Language {
	b0, b1 = asciiLower(b0), asciiLower(b1)
	if b0 >= 'a' && b0 <= 'z' && b1 >= 'a' && b1 <= 'z' {
		return &languages[codeIndex[b0-'a'][b1-'a']]
	}
	return &languages[0]
}

// GetLanguageByCodeString resolves a BCP-47 / ISO-639-1 string (first two letters).
// Unknown and short values return Languages[0] (never nil).
func GetLanguageByCodeString(code string) *Language {
	if len(code) < 2 {
		return &languages[0]
	}
	return lookup(code[0], code[1])
}

// GetLanguageByCode resolves a two-byte ISO code. Unknown values return Languages[0].
func GetLanguageByCode(code Code) *Language {
	return lookup(code[0], code[1])
}

// GetLanguageIdByCode returns the catalog ID (0 if unknown).
func GetLanguageIdByCode(code Code) uint {
	return lookup(code[0], code[1]).ID
}

// GetLanguageIdByCodeString returns the catalog ID (0 if unknown / short).
func GetLanguageIdByCodeString(code string) uint {
	return GetLanguageByCodeString(code).ID
}

// GetLanguageByID returns the language with the given ID, or Languages[0].
func GetLanguageByID(id uint) *Language {
	if id >= uint(len(languages)) {
		return &languages[0]
	}
	return &languages[id]
}

// IntCode is the packed little-endian numeric form of a two-byte code.
func IntCode(geoCode [2]byte) uint {
	return Code(geoCode).Int()
}
