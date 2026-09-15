package languages

import (
	"encoding/json"
	"errors"
	"testing"
)

func Test_CodeLanguage(t *testing.T) {
	en := GetLanguageByCodeString("EN")
	if en == nil || en.Name != "English" || en.ID == 0 {
		t.Fatalf("EN: %+v", en)
	}

	tests := []struct {
		in   string
		want uint
		name string
	}{
		{"EN", en.ID, "English"},
		{"en", en.ID, "English"},
		{"en-US", en.ID, "English"},
		{"en_GB", en.ID, "English"},
		{"**", 0, "Undefined"},
		{"xx", 0, "Undefined"},
		{"*", 0, "Undefined"},
		{"", 0, "Undefined"},
		{"X", 0, "Undefined"},
		{"zh-Hans", GetLanguageByCodeString("ZH").ID, "Chinese"},
	}

	for _, test := range tests {
		got := GetLanguageByCodeString(test.in)
		if got == nil {
			t.Fatalf("%q: nil language", test.in)
		}
		if got.ID != test.want || got.Name != test.name {
			t.Errorf("%q: id=%d name=%q, want id=%d name=%q", test.in, got.ID, got.Name, test.want, test.name)
		}
		if CodeFromString(test.in).ID() != test.want {
			t.Errorf("CodeFromString(%q).ID() = %d, want %d", test.in, CodeFromString(test.in).ID(), test.want)
		}
	}

	if !GetLanguageByCode(Code{'x', 'x'}).Code.IsUndefined() {
		t.Fatal("unknown Code should be undefined language")
	}
	if !UndefinedCode.IsUndefined() || !UndefinedCode.Language().Code.IsUndefined() {
		t.Fatal("UndefinedCode should map to Languages[0]")
	}
	if CodeFromString("EN") != (Code{'e', 'n'}) {
		t.Fatalf("canonical code should be lowercase, got %q", CodeFromString("EN"))
	}
	if GetLanguageByCodeString("en") != &languages[en.ID] {
		t.Fatal("lookup should return pointer into catalog")
	}
}

func Test_CodeISO2(t *testing.T) {
	tests := []struct {
		cc   Code
		want string
	}{
		{Code{'E', 'N'}, "en"},
		{Code{'e', 'n'}, "en"},
		{UndefinedCode, UndefinedISO2},
		{Code{'-', '-'}, UndefinedISO2},
		{Code{'X', 'Y'}, UndefinedISO2},
	}
	for _, test := range tests {
		if got := test.cc.ISO2(); got != test.want {
			t.Errorf("Code(%q).ISO2() = %q, want %q", string(test.cc[:]), got, test.want)
		}
		if got := test.cc.String(); got != test.want {
			t.Errorf("Code(%q).String() = %q, want %q", string(test.cc[:]), got, test.want)
		}
	}
}

func Test_CodeJSON(t *testing.T) {
	type object struct {
		Code Code `json:"code"`
	}

	tests := []struct {
		source string
		result string
	}{
		{`{"code":"en"}`, `{"code":"en"}`},
		{`{"code":"EN"}`, `{"code":"en"}`},
		{`{"code":"**"}`, `{"code":"**"}`},
		{`{"code":"***"}`, `{"code":"**"}`},
		{`{"code":"ABC"}`, `{"code":"**"}`},
	}

	for _, test := range tests {
		var obj object
		if err := json.Unmarshal([]byte(test.source), &obj); err != nil {
			t.Errorf("unmarshal %s: %v", test.source, err)
			continue
		}
		data, err := json.Marshal(&obj)
		if err != nil {
			t.Errorf("marshal %s: %v", test.source, err)
			continue
		}
		if string(data) != test.result {
			t.Errorf("round-trip [%s] got [%s] want [%s]", test.source, data, test.result)
		}
	}
}

func Test_CodeScan(t *testing.T) {
	var cc Code
	if err := cc.Scan([]byte("en")); err != nil {
		t.Fatalf("scan bytes: %v", err)
	}
	if cc != (Code{'e', 'n'}) {
		t.Fatalf("scan bytes got %q", cc)
	}
	if err := cc.Scan("DE"); err != nil {
		t.Fatalf("scan string: %v", err)
	}
	if cc != (Code{'d', 'e'}) {
		t.Fatalf("scan string got %q", cc)
	}
	if err := cc.Scan([]byte("ENG")); !errors.Is(err, ErrCodeInvalidValueSize) {
		t.Fatalf("scan long bytes err = %v", err)
	}
	if err := cc.Scan(""); !errors.Is(err, ErrCodeInvalidValueSize) {
		t.Fatalf("scan empty err = %v", err)
	}
	if err := cc.Scan(123); !errors.Is(err, ErrCodeInvalidScanType) {
		t.Fatalf("scan int err = %v", err)
	}
}

func Test_NoAllocLookup(t *testing.T) {
	code := Code{'e', 'n'}
	if n := testing.AllocsPerRun(1000, func() { _ = code.Language() }); n != 0 {
		t.Fatalf("Code.Language allocs = %v", n)
	}
	if n := testing.AllocsPerRun(1000, func() { _ = GetLanguageByCodeString("en") }); n != 0 {
		t.Fatalf("GetLanguageByCodeString allocs = %v", n)
	}
	if n := testing.AllocsPerRun(1000, func() { _ = code.ISO2() }); n != 0 {
		t.Fatalf("Code.ISO2 allocs = %v", n)
	}
	if n := testing.AllocsPerRun(1000, func() { _ = CodeFromString("en-US") }); n != 0 {
		t.Fatalf("CodeFromString allocs = %v", n)
	}
}
