package languages

// Language is a catalog entry. Slice index equals ID.
type Language struct {
	ID         uint   `json:"ID"`
	Code       Code   `json:"code"`
	Name       string `json:"name"`
	NativeName string `json:"native_name"`
}

// IntCode for ISO packed bytes.
func (l Language) IntCode() uint {
	return l.Code.Int()
}

// ISO2 returns the interned two-letter code (never allocates after init).
func (l *Language) ISO2() string {
	if l == nil {
		return UndefinedISO2
	}
	return languageISO2[l.ID]
}
