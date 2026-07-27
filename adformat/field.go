package adformat

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"sync"

	"github.com/demdxx/gocast/v2"
)

// Errors returned by Field.Prepare.
var (
	ErrFieldIsRequired  = errors.New("adformat: field is required")
	ErrFieldInvalidType = errors.New("adformat: field value has invalid type")
)

// FieldType is the data type of a Field value. It is a plain string, open
// for extension the same way Kind is.
type FieldType string

// Field type values.
const (
	FieldStringType    FieldType = "string"
	FieldIntType       FieldType = "int"
	FieldFloatType     FieldType = "float"
	FieldBoolType      FieldType = "bool"
	FieldPhoneType     FieldType = "phone"
	FieldEmailType     FieldType = "email"
	FieldURLType       FieldType = "url"
	FieldGeoType       FieldType = "geo"
	FieldScreenRefType FieldType = "screen_ref" // §3.2.11 — advertiser picks a target screen from Options
	FieldDefaultType             = FieldStringType
)

// Well-known field names, kept for parity with the old model
// (admodels/types.FormatField* constants).
const (
	FieldNameTitle       = "title"
	FieldNameDescription = "description"
	FieldNameBrandname   = "brandname"
	FieldNamePhone       = "phone"
	FieldNameURL         = "url"
	FieldNameContent     = "content"
	FieldNameRating      = "rating"
	FieldNameLikes       = "likes"
	FieldNameAddress     = "address"
	FieldNameSponsored   = "sponsored"
)

// GeoValue is the value shape of a FieldGeoType field.
type GeoValue struct {
	Lat float64 `json:"lat" yaml:"lat"`
	Lng float64 `json:"lng" yaml:"lng"`
}

// Validate reports whether the coordinates are within valid ranges.
func (g GeoValue) Validate() error {
	if g.Lat < -90 || g.Lat > 90 {
		return fmt.Errorf("adformat: geo latitude %f out of range [-90, 90]", g.Lat)
	}
	if g.Lng < -180 || g.Lng > 180 {
		return fmt.Errorf("adformat: geo longitude %f out of range [-180, 180]", g.Lng)
	}
	return nil
}

// Option is a single selectable choice for a Field, rendered as a
// dropdown/checkbox/select control by a UI builder. Works for both string
// and numeric field types — Value's concrete type should match Field.Type
// (e.g. a string for FieldStringType, a number for
// FieldIntType/FieldFloatType). Replaces the old, inconsistently-shaped
// FormatField.Select []any (§3.2.3).
type Option struct {
	Value any    `json:"value" yaml:"value"`
	Title string `json:"title,omitempty" yaml:"title,omitempty"` // if empty, a UI shows Value as-is
}

// Field describes a requirement for a text/numeric/boolean value that is
// part of one ad creative (was admodels/types.FormatField).
type Field struct {
	// Required marks the field as mandatory when it is active (see
	// Condition below) — a Condition never overrides Required, it only
	// controls whether the field participates at all.
	Required bool `json:"required,omitempty" yaml:"required,omitempty"`

	// Title of the field for the UI.
	Title string `json:"title,omitempty" yaml:"title,omitempty"`

	// Name of the field (required, used for addressing/Condition/Binding).
	Name string `json:"name" yaml:"name"`

	// Type of the field data.
	Type FieldType `json:"type,omitempty" yaml:"type,omitempty"`

	// Options is the list of selectable values, used both for a
	// single-select field and, combined with MaxCount > 1, for a
	// multi-select field (§3.2.3/§3.2.15). Empty = free-form input.
	Options []Option `json:"options,omitempty" yaml:"options,omitempty"`

	// Min/Max: length bounds for string-like types, value bounds for
	// int/float types.
	Min float64 `json:"min,omitempty" yaml:"min,omitempty"`
	Max float64 `json:"max,omitempty" yaml:"max,omitempty"`

	// MinCount/MaxCount: how many instances of this field one ad creative
	// carries — repeatable text variants (e.g. Google Responsive Ads
	// style headlines) or a multi-select out of Options (§3.2.15). See
	// GetMinCount/GetMaxCount/IsMultiple/RangeCount for the effective
	// semantics. MaxCount == -1 means unlimited.
	MinCount int `json:"min_count,omitempty" yaml:"min_count,omitempty"`
	MaxCount int `json:"max_count,omitempty" yaml:"max_count,omitempty"`

	// Mask is a pure UI hint for rendering a masked input (e.g.
	// "+1 (999) 999-9999") — it is never used for server-side validation,
	// unlike RegExp below.
	Mask string `json:"mask,omitempty" yaml:"mask,omitempty"`

	// RegExp is an additional/overriding validation pattern (e.g. a
	// stricter phone format for a specific exchange). Unlike the old
	// model, this is actually applied by Prepare (§3.2.4).
	RegExp string `json:"regexp,omitempty" yaml:"regexp,omitempty"`

	// Condition controls when this field is active (shown + eligible to
	// be required) — replaces the dead FormatField.Exclude (§3.2.5).
	Condition *Condition `json:"condition,omitempty" yaml:"condition,omitempty"`

	// NavigateTo names the screen (sibling Config/group) this field
	// switches the creative to when activated (e.g. a button) — for
	// structural multi-screen creatives (§3.2.11). "$back"/"$close" are
	// reserved commands, not screen names.
	NavigateTo string `json:"navigate_to,omitempty" yaml:"navigate_to,omitempty"`

	// Bindings associate this field with external protocol
	// representations (OpenRTB Native data/title asset, tracking pixel
	// event, ...) — §3.2.12.
	Bindings []Binding `json:"bindings,omitempty" yaml:"bindings,omitempty"`
}

// GetType of the field, defaulting to FieldDefaultType.
func (f Field) GetType() FieldType {
	if f.Type == "" {
		return FieldDefaultType
	}
	return f.Type
}

// GetName of the field.
func (f Field) GetName() string {
	return f.Name
}

// IsRequired field value. Bool fields are never "required" in the
// validation sense (false is a perfectly valid value), same as the old
// model.
func (f Field) IsRequired() bool {
	return f.GetType() != FieldBoolType && f.Required
}

// MaxLength of the field limit, for length-bounded types only.
func (f Field) MaxLength() int {
	switch f.GetType() {
	case FieldStringType, FieldPhoneType, FieldEmailType, FieldURLType:
		return int(f.Max)
	}
	return 0
}

// GetMinCount effective minimal number of instances of this field within
// one ad creative (§3.2.15).
func (f Field) GetMinCount() int {
	mn, _ := effectiveCount(f.MinCount, f.MaxCount, f.Required)
	return mn
}

// GetMaxCount effective maximal number of instances (-1 = unlimited).
func (f Field) GetMaxCount() int {
	_, mx := effectiveCount(f.MinCount, f.MaxCount, f.Required)
	return mx
}

// RangeCount returns the effective (min, max) instance-count range.
func (f Field) RangeCount() (min, max int) {
	return effectiveCount(f.MinCount, f.MaxCount, f.Required)
}

// IsMultiple reports whether this field may have more than one instance
// (repeatable field or multi-select out of Options).
func (f Field) IsMultiple() bool {
	return isMultipleCount(f.GetMaxCount())
}

// HasOptions reports whether the field restricts values to a fixed list.
func (f Field) HasOptions() bool {
	return len(f.Options) > 0
}

// IsValidOption reports whether v equals one of Options[i].Value. Always
// true when the field has no Options (free-form input).
func (f Field) IsValidOption(v any) bool {
	if !f.HasOptions() {
		return true
	}
	for _, opt := range f.Options {
		if opt.Value == v {
			return true
		}
	}
	return false
}

// SoftEqual compares two fields for compatibility purposes (used by
// Config.Intersec), same relaxed semantics as the old model: same name and
// either both are string-typed or both share the exact same type.
func (f Field) SoftEqual(field Field) bool {
	return f.Name == field.Name && (f.GetType() == FieldStringType || f.GetType() == field.GetType())
}

var (
	regexpCacheMu sync.RWMutex
	regexpCache   = map[string]*regexp.Regexp{}
)

func compileRegexp(pattern string) (*regexp.Regexp, error) {
	regexpCacheMu.RLock()
	re, ok := regexpCache[pattern]
	regexpCacheMu.RUnlock()
	if ok {
		return re, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	regexpCacheMu.Lock()
	regexpCache[pattern] = re
	regexpCacheMu.Unlock()
	return re, nil
}

// Prepare validates and normalizes a single input value according to the
// field's Type/Min/Max/RegExp/Options. Unlike the old
// FormatField.Prepare, the float upper bound is actually checked, boolean
// and phone values are no longer silently dropped, and "url" is a real
// type. For a repeatable/multi-select field (IsMultiple), call Prepare
// once per element and use RangeCount to validate the element count
// separately (Prepare itself only validates one value at a time).
func (f Field) Prepare(value any) (result any, err error) {
	if f.IsRequired() && gocast.IsEmpty(value) {
		return nil, ErrFieldIsRequired
	}

	switch f.GetType() {
	case FieldStringType:
		v := gocast.Str(value)
		if err = f.validateRegExp(v); err == nil {
			if f.Min > 0 && float64(len(v)) < f.Min {
				err = fmt.Errorf("adformat: min length is %d", int(f.Min))
			} else if f.Max > 0 && float64(len(v)) > f.Max {
				err = fmt.Errorf("adformat: max length is %d", int(f.Max))
			}
		}
		result = v
	case FieldIntType:
		v := gocast.Int64(value)
		if f.Min > 0 && v < int64(f.Min) {
			err = fmt.Errorf("adformat: min value is %d", int(f.Min))
		} else if f.Max > 0 && v > int64(f.Max) {
			err = fmt.Errorf("adformat: max value is %d", int(f.Max))
		}
		result = v
	case FieldFloatType:
		v := gocast.Float64(value)
		if f.Min > 0 && v < f.Min {
			err = fmt.Errorf("adformat: min value is %.3f", f.Min)
		} else if f.Max > 0 && v > f.Max {
			err = fmt.Errorf("adformat: max value is %.3f", f.Max)
		}
		result = v
	case FieldBoolType:
		result = gocast.Bool(value)
	case FieldPhoneType:
		v := gocast.Str(value)
		if err = f.validatePhone(v); err == nil {
			result = v
		}
	case FieldEmailType:
		v := gocast.Str(value)
		if err = f.validateEmail(v); err == nil {
			result = v
		}
	case FieldURLType:
		v := gocast.Str(value)
		if err = f.validateURL(v); err == nil {
			result = v
		}
	case FieldGeoType:
		geo, gerr := toGeoValue(value)
		if gerr != nil {
			err = gerr
		} else if verr := geo.Validate(); verr != nil {
			err = verr
		} else {
			result = geo
		}
	case FieldScreenRefType:
		result = gocast.Str(value)
	default:
		result = value
	}

	if err == nil && f.HasOptions() && !gocast.IsEmpty(result) {
		if !f.IsValidOption(result) {
			err = fmt.Errorf("adformat: value %v is not a valid option for field %q", result, f.GetName())
		}
	}

	return result, err
}

func (f Field) validateRegExp(v string) error {
	if f.RegExp == "" {
		return nil
	}
	re, err := compileRegexp(f.RegExp)
	if err != nil {
		return fmt.Errorf("adformat: invalid regexp for field %q: %w", f.GetName(), err)
	}
	if !re.MatchString(v) {
		return fmt.Errorf("adformat: value does not match pattern for field %q", f.GetName())
	}
	return nil
}

var basicPhonePattern = regexp.MustCompile(`^[+]?[0-9()\-\s.]{5,20}$`)

func (f Field) validatePhone(v string) error {
	if f.RegExp != "" {
		return f.validateRegExp(v)
	}
	if !basicPhonePattern.MatchString(v) {
		return fmt.Errorf("adformat: %q is not a valid phone number", v)
	}
	return nil
}

func (f Field) validateEmail(v string) error {
	if f.RegExp != "" {
		return f.validateRegExp(v)
	}
	if _, err := mail.ParseAddress(v); err != nil {
		return fmt.Errorf("adformat: %q is not a valid email: %w", v, err)
	}
	return nil
}

func (f Field) validateURL(v string) error {
	if f.RegExp != "" {
		return f.validateRegExp(v)
	}
	u, err := url.Parse(v)
	if err != nil {
		return fmt.Errorf("adformat: %q is not a valid url: %w", v, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("adformat: %q is not an absolute url", v)
	}
	return nil
}

func toGeoValue(value any) (GeoValue, error) {
	switch v := value.(type) {
	case GeoValue:
		return v, nil
	case *GeoValue:
		if v == nil {
			return GeoValue{}, ErrFieldInvalidType
		}
		return *v, nil
	case map[string]any:
		return GeoValue{Lat: gocast.Float64(v["lat"]), Lng: gocast.Float64(v["lng"])}, nil
	}
	return GeoValue{}, fmt.Errorf("%w: expected geo value, got %T", ErrFieldInvalidType, value)
}
