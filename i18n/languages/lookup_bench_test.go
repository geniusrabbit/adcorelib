package languages

import "testing"

var (
	benchLanguage *Language
	benchCode     Code
	benchString   string
	benchID       uint
)

func BenchmarkGetLanguageByCodeString(b *testing.B) {
	b.ReportAllocs()
	var lang *Language
	for i := 0; i < b.N; i++ {
		lang = GetLanguageByCodeString("en")
	}
	benchLanguage = lang
}

func BenchmarkGetLanguageByCode(b *testing.B) {
	code := Code{'e', 'n'}
	b.ReportAllocs()
	var lang *Language
	for i := 0; i < b.N; i++ {
		lang = GetLanguageByCode(code)
	}
	benchLanguage = lang
}

func BenchmarkCodeLanguage(b *testing.B) {
	code := Code{'e', 'n'}
	b.ReportAllocs()
	var lang *Language
	for i := 0; i < b.N; i++ {
		lang = code.Language()
	}
	benchLanguage = lang
}

func BenchmarkCodeISO2(b *testing.B) {
	code := Code{'e', 'n'}
	b.ReportAllocs()
	var s string
	for i := 0; i < b.N; i++ {
		s = code.ISO2()
	}
	benchString = s
}

func BenchmarkCodeFromStringBCP47(b *testing.B) {
	b.ReportAllocs()
	var c Code
	for i := 0; i < b.N; i++ {
		c = CodeFromString("en-US")
	}
	benchCode = c
}

func BenchmarkGetLanguageIdByCodeString(b *testing.B) {
	b.ReportAllocs()
	var id uint
	for i := 0; i < b.N; i++ {
		id = GetLanguageIdByCodeString("en")
	}
	benchID = id
}

func BenchmarkGetLanguageByID(b *testing.B) {
	b.ReportAllocs()
	var lang *Language
	for i := 0; i < b.N; i++ {
		lang = GetLanguageByID(40)
	}
	benchLanguage = lang
}
