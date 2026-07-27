package adformat

// builder.go is the Go-ergonomic layer over Node/Config/AssetRequirement/
// Field/Condition/Param/SizeOption (§3.2.10): constructors plus a chain
// of value-receiver modifiers (each returns a modified copy, so they
// compose freely inside a slice literal) that let a format be described
// as a plain Go var/const table, with zero YAML/JSON/DB involved.
// Config.WithAssets/WithFields/WithGroups/WithParams wrap their arguments
// into Node internally — calling code never has to touch the NodeKind
// discriminator or take the address of the right field itself.
//
// Deviation from the literal plan text (forced by Go's language rules,
// not a design choice): a method cannot share its name with a struct
// field, and a top-level function cannot share its name with a type.
// This affects exactly two modifier names and one constructor name:
//   - Field.Title (struct field) → the modifier is WithTitle, not Title.
//   - AssetRequirement.NavigateTo / Field.NavigateTo (struct fields) →
//     the modifier is WithNavigateTo, not NavigateTo.
//   - the Param struct type → the constructor is NewParam, not Param.
//   - AssetRequirement.AspectRatio (struct field) → the modifier is
//     WithAspectRatio, not AspectRatio (§3.2.1 gap #1).
// Every other constructor/modifier keeps exactly the name from the plan.

// Asset starts a new AssetRequirement builder chain with the given name.
func Asset(name string) AssetRequirement {
	return AssetRequirement{Name: name}
}

// Require marks the asset as mandatory.
func (a AssetRequirement) Require() AssetRequirement {
	a.Required = true
	return a
}

// Size sets the allowed width/height range of the asset file.
func (a AssetRequirement) Size(minW, maxW, minH, maxH int) AssetRequirement {
	a.MinWidth, a.MaxWidth = minW, maxW
	a.MinHeight, a.MaxHeight = minH, maxH
	return a
}

// ExactSize sets a single preferred exact size (Width/Height only, no
// range) — the common case for a classic fixed-size asset.
func (a AssetRequirement) ExactSize(width, height int) AssetRequirement {
	a.Width, a.Height = width, height
	return a
}

// Types sets the allowed asset content types (short category or exact
// MIME type, §3.2.1).
func (a AssetRequirement) Types(types ...string) AssetRequirement {
	a.AllowedTypes = types
	return a
}

// Count sets the MinCount/MaxCount instance-count bounds (§3.2, п.2).
func (a AssetRequirement) Count(min, max int) AssetRequirement {
	a.MinCount, a.MaxCount = min, max
	return a
}

// Focal enables a focal point on this asset with the given default
// anchor (§3.2.1).
func (a AssetRequirement) Focal(x, y float64) AssetRequirement {
	a.FocalPoint = &FocalPoint{Enabled: true, Default: [2]float64{x, y}}
	return a
}

// WithAspectRatio sets a direct "W:H" proportion constraint, e.g.
// "1.91:1" (§3.2.1 gap #1). Named WithAspectRatio, not AspectRatio,
// because the latter would collide with the struct field of the same
// name (see the package-level deviation note above).
func (a AssetRequirement) WithAspectRatio(ratio string) AssetRequirement {
	a.AspectRatio = ratio
	return a
}

// WithSafeZone sets the edge-inset metadata, symmetric to Focal (§3.2.1
// gap #2). Insets are fractions (0..1) of the asset's own width/height.
func (a AssetRequirement) WithSafeZone(top, right, bottom, left float64) AssetRequirement {
	a.SafeZone = &SafeZone{Top: top, Right: right, Bottom: bottom, Left: left}
	return a
}

// RequireAltText marks this asset as needing an accessibility text
// description of at most maxLength characters (0 = unbounded), §3.2.1
// gap #3.
func (a AssetRequirement) RequireAltText(maxLength int) AssetRequirement {
	a.AltText = &AltTextRequirement{Required: true, MaxLength: maxLength}
	return a
}

// Bitrate sets the allowed video/audio bitrate range in kbps (§3.2.1
// gap #4).
func (a AssetRequirement) Bitrate(min, max int) AssetRequirement {
	a.BitrateMin, a.BitrateMax = min, max
	return a
}

// Codecs restricts the allowed video/audio codecs (§3.2.1 gap #4).
func (a AssetRequirement) Codecs(codecs ...string) AssetRequirement {
	a.AllowedCodecs = codecs
	return a
}

// Framerates restricts the allowed video frame rates in fps (§3.2.1
// gap #4).
func (a AssetRequirement) Framerates(fps ...float64) AssetRequirement {
	a.AllowedFramerates = fps
	return a
}

// MaxSize sets the maximum allowed file weight in bytes (§3.2.2).
func (a AssetRequirement) MaxSize(bytes int64) AssetRequirement {
	a.MaxFileSize = bytes
	return a
}

// Duration sets the allowed video/audio duration range in seconds.
func (a AssetRequirement) Duration(min, max float64) AssetRequirement {
	a.DurationMin, a.DurationMax = min, max
	return a
}

// Animated allows animated content for this asset.
func (a AssetRequirement) Animated() AssetRequirement {
	a.AllowAnimated = true
	return a
}

// Sound allows sound for this asset.
func (a AssetRequirement) Sound() AssetRequirement {
	a.AllowSound = true
	return a
}

// When attaches a Condition controlling this asset's activity (§3.2.5).
func (a AssetRequirement) When(c Condition) AssetRequirement {
	a.Condition = &c
	return a
}

// WithNavigateTo marks activating this asset (e.g. a tappable hotspot
// image) as a fixed transition to the named sibling screen — or to the
// reserved "$back"/"$close" commands (§3.2.11). Named WithNavigateTo, not
// NavigateTo, because the latter would collide with the struct field of
// the same name (see the package-level deviation note above).
func (a AssetRequirement) WithNavigateTo(screen string) AssetRequirement {
	a.NavigateTo = screen
	return a
}

// Bind attaches an external-protocol association to this asset (§3.2.12).
func (a AssetRequirement) Bind(b Binding) AssetRequirement {
	a.Bindings = append(a.Bindings, b)
	return a
}

// newField starts a new Field builder chain with the given name and type.
func newField(name string, t FieldType) Field {
	return Field{Name: name, Type: t}
}

// StringField starts a string Field builder chain.
func StringField(name string) Field { return newField(name, FieldStringType) }

// IntField starts an int Field builder chain.
func IntField(name string) Field { return newField(name, FieldIntType) }

// FloatField starts a float Field builder chain.
func FloatField(name string) Field { return newField(name, FieldFloatType) }

// BoolField starts a bool Field builder chain.
func BoolField(name string) Field { return newField(name, FieldBoolType) }

// PhoneField starts a phone Field builder chain.
func PhoneField(name string) Field { return newField(name, FieldPhoneType) }

// EmailField starts an email Field builder chain.
func EmailField(name string) Field { return newField(name, FieldEmailType) }

// URLField starts a url Field builder chain.
func URLField(name string) Field { return newField(name, FieldURLType) }

// GeoField starts a geo Field builder chain.
func GeoField(name string) Field { return newField(name, FieldGeoType) }

// ScreenRefField starts a screen_ref Field builder chain — the advertiser
// picks the target screen from Options (§3.2.11).
func ScreenRefField(name string) Field { return newField(name, FieldScreenRefType) }

// Require marks the field as mandatory (when active, see When).
func (f Field) Require() Field {
	f.Required = true
	return f
}

// Len sets the Min/Max length bounds for a string-like field.
func (f Field) Len(min, max int) Field {
	f.Min, f.Max = float64(min), float64(max)
	return f
}

// Range sets the Min/Max value bounds for an int/float field.
func (f Field) Range(min, max float64) Field {
	f.Min, f.Max = min, max
	return f
}

// WithOptions sets the list of selectable values (§3.2.3).
func (f Field) WithOptions(opts ...Option) Field {
	f.Options = opts
	return f
}

// When attaches a Condition controlling this field's activity (§3.2.5).
func (f Field) When(c Condition) Field {
	f.Condition = &c
	return f
}

// WithTitle sets the UI title of the field. Named WithTitle, not Title,
// because the latter would collide with the struct field of the same
// name (see the package-level deviation note above).
func (f Field) WithTitle(t string) Field {
	f.Title = t
	return f
}

// Bind attaches an external-protocol association to this field
// (§3.2.12).
func (f Field) Bind(b Binding) Field {
	f.Bindings = append(f.Bindings, b)
	return f
}

// Count sets the MinCount/MaxCount instance-count bounds — repeatable
// text variants or a multi-select out of Options (§3.2.15).
func (f Field) Count(min, max int) Field {
	f.MinCount, f.MaxCount = min, max
	return f
}

// WithRegExp sets an additional/overriding validation pattern (§3.2.4).
func (f Field) WithRegExp(pattern string) Field {
	f.RegExp = pattern
	return f
}

// WithMask sets the UI-only input mask hint.
func (f Field) WithMask(mask string) Field {
	f.Mask = mask
	return f
}

// WithNavigateTo marks activating this field (e.g. a button) as a fixed
// transition to the named sibling screen — or to the reserved
// "$back"/"$close" commands (§3.2.11). Named WithNavigateTo, not
// NavigateTo, because the latter would collide with the struct field of
// the same name (see the package-level deviation note above).
func (f Field) WithNavigateTo(screen string) Field {
	f.NavigateTo = screen
	return f
}

// Opt builds a single Option for Field.WithOptions.
func Opt(value any, title string) Option {
	return Option{Value: value, Title: title}
}

// Group starts a new nested Config ("screen"/group) builder chain with
// the given name.
func Group(name string) Config {
	return Config{Name: name}
}

// Count sets the MinCount/MaxCount instance-count bounds of this group
// (§3.4/§3.5).
func (c Config) Count(min, max int) Config {
	c.MinCount, c.MaxCount = min, max
	return c
}

// When attaches a Condition controlling this group's activity (§3.2.5).
func (c Config) When(cond Condition) Config {
	c.Condition = &cond
	return c
}

// WithTitle sets the UI title of the group.
func (c Config) WithTitle(t string) Config {
	c.Title = t
	return c
}

// AsEntry marks this group as the starting screen of a multi-screen
// creative (§3.2.11).
func (c Config) AsEntry() Config {
	c.Entry = true
	return c
}

// WithAssets appends the given asset requirements to Children, wrapping
// each into a Node — the NodeKind discriminator never surfaces to the
// caller.
func (c Config) WithAssets(assets ...AssetRequirement) Config {
	for i := range assets {
		c.Children = append(c.Children, AssetNode(assets[i]))
	}
	return c
}

// WithFields appends the given fields to Children, wrapping each into a
// Node.
func (c Config) WithFields(fields ...Field) Config {
	for i := range fields {
		c.Children = append(c.Children, FieldNode(fields[i]))
	}
	return c
}

// WithGroups appends the given nested Config groups to Children, wrapping
// each into a Node.
func (c Config) WithGroups(groups ...Config) Config {
	for i := range groups {
		c.Children = append(c.Children, GroupNode(groups[i]))
	}
	return c
}

// WithParams appends the given params to Children, wrapping each into a
// Node (§3.2.14).
func (c Config) WithParams(params ...Param) Config {
	for i := range params {
		c.Children = append(c.Children, ParamNode(params[i]))
	}
	return c
}

// NewParam builds a Param with the given name/value. Named NewParam, not
// Param, because a top-level function cannot share its name with the
// Param type (see the package-level deviation note above).
func NewParam(name string, value any) Param {
	return Param{Name: name, Value: value}
}

// When attaches a Condition controlling this param's activity (§3.2.5).
func (p Param) When(c Condition) Param {
	p.Condition = &c
	return p
}

// FixedSize builds a single exact-size SizeOption (§3.2.13).
func FixedSize(title string, width, height int) SizeOption {
	return SizeOption{Title: title, Width: width, Height: height}
}

// FlexibleSize builds a flexible SizeOption with optional bounds — pass 0
// for an unbounded side (§3.2.13).
func FlexibleSize(title string, minW, maxW, minH, maxH int) SizeOption {
	return SizeOption{
		Title: title, Flexible: true,
		MinWidth: minW, MaxWidth: maxW,
		MinHeight: minH, MaxHeight: maxH,
	}
}

// AsDefault marks this SizeOption as pre-selected when a Format has more
// than one.
func (s SizeOption) AsDefault() SizeOption {
	s.Default = true
	return s
}
