package adformat

import (
	"errors"
	"strings"
	"testing"
)

func TestFieldPrepareString(t *testing.T) {
	f := StringField("title").Require().Len(5, 10)

	if _, err := f.Prepare("hi"); err == nil || !strings.Contains(err.Error(), "min length") {
		t.Errorf("expected min length error, got %v", err)
	}
	if _, err := f.Prepare("way too long a title"); err == nil || !strings.Contains(err.Error(), "max length") {
		t.Errorf("expected max length error, got %v", err)
	}
	v, err := f.Prepare("just right")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "just right" {
		t.Errorf("got %v, want %q", v, "just right")
	}
	if _, err := f.Prepare(""); !errors.Is(err, ErrFieldIsRequired) {
		t.Errorf("expected ErrFieldIsRequired, got %v", err)
	}
}

func TestFieldPrepareFloat_MaxBoundActuallyChecked(t *testing.T) {
	// Regression test for the old admodels/types bug (§2): the float
	// branch compared the upper bound against Min instead of Max, so the
	// maximum was never actually enforced.
	f := FloatField("price").Range(1, 10)

	if _, err := f.Prepare(5.0); err != nil {
		t.Errorf("5.0 should be within [1,10], got error: %v", err)
	}
	if _, err := f.Prepare(0.5); err == nil || !strings.Contains(err.Error(), "min value") {
		t.Errorf("expected min value error for 0.5, got %v", err)
	}
	if _, err := f.Prepare(15.0); err == nil || !strings.Contains(err.Error(), "max value") {
		t.Errorf("expected max value error for 15.0 (regression check), got %v", err)
	}
	// Exact error formatting bug fix: the old model produced
	// "max value is %f.3" (wrong verb/precision order) instead of
	// "%.3f".
	_, err := f.Prepare(15.0)
	if err == nil || !strings.Contains(err.Error(), "10.000") {
		t.Errorf("expected properly formatted %%.3f value in error, got %v", err)
	}
}

func TestFieldPrepareBoolPhoneURL_ValueNoLongerDropped(t *testing.T) {
	// Regression test for the old admodels/types bug (§2): bool/phone (and
	// any unrecognized/url type) branches never wrote to the named
	// `result`, so a valid value was silently discarded.
	b := BoolField("active")
	v, err := b.Prepare(true)
	if err != nil || v != true {
		t.Errorf("Prepare(true) = (%v, %v), want (true, nil)", v, err)
	}

	p := PhoneField("phone")
	v, err = p.Prepare("+1 415 555 0100")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "+1 415 555 0100" {
		t.Errorf("phone value lost, got %v", v)
	}

	u := URLField("url")
	v, err = u.Prepare("https://example.com/landing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "https://example.com/landing" {
		t.Errorf("url value lost, got %v", v)
	}
}

func TestFieldPrepareURLType(t *testing.T) {
	u := URLField("url").Require()
	if _, err := u.Prepare("not a url"); err == nil {
		t.Error("expected error for invalid url")
	}
	if _, err := u.Prepare("https://example.com"); err != nil {
		t.Errorf("unexpected error for valid url: %v", err)
	}
}

func TestFieldPrepareEmailType(t *testing.T) {
	e := EmailField("email")
	if _, err := e.Prepare("not-an-email"); err == nil {
		t.Error("expected error for invalid email")
	}
	if _, err := e.Prepare("user@example.com"); err != nil {
		t.Errorf("unexpected error for valid email: %v", err)
	}
}

func TestFieldPrepareGeoType(t *testing.T) {
	g := GeoField("location")
	if _, err := g.Prepare(map[string]any{"lat": 91.0, "lng": 0.0}); err == nil {
		t.Error("expected error for out-of-range latitude")
	}
	v, err := g.Prepare(map[string]any{"lat": 40.7, "lng": -74.0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	geo, ok := v.(GeoValue)
	if !ok || geo.Lat != 40.7 || geo.Lng != -74.0 {
		t.Errorf("unexpected geo value: %+v (ok=%v)", v, ok)
	}
}

func TestFieldPrepareRegExp(t *testing.T) {
	f := StringField("zip").WithRegExp(`^\d{5}$`)
	if _, err := f.Prepare("abc"); err == nil {
		t.Error("expected regexp validation error")
	}
	if _, err := f.Prepare("12345"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFieldPrepareOptions(t *testing.T) {
	f := StringField("adtype").WithOptions(Opt("business", "Business"), Opt("game", "Game"))
	if _, err := f.Prepare("business"); err != nil {
		t.Errorf("unexpected error for valid option: %v", err)
	}
	if _, err := f.Prepare("not-an-option"); err == nil {
		t.Error("expected error for value not in Options")
	}
}

func TestFieldIsValidOption(t *testing.T) {
	free := StringField("free")
	if !free.IsValidOption("anything") {
		t.Error("a field without Options should accept any value")
	}

	restricted := StringField("restricted").WithOptions(Opt("a", ""), Opt("b", ""))
	if !restricted.IsValidOption("a") {
		t.Error("expected 'a' to be a valid option")
	}
	if restricted.IsValidOption("z") {
		t.Error("expected 'z' to be an invalid option")
	}
}

func TestFieldSoftEqual(t *testing.T) {
	a := Field{Name: "title", Type: FieldStringType}
	b := Field{Name: "title", Type: FieldStringType}
	if !a.SoftEqual(b) {
		t.Error("expected identical string fields to be SoftEqual")
	}
	c := Field{Name: "other", Type: FieldStringType}
	if a.SoftEqual(c) {
		t.Error("expected different names not to be SoftEqual")
	}
}
