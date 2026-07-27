package adformat

// Condition describes when a Field, AssetRequirement or nested Config
// (group/param) is active (shown + eligible to be required), based on the
// current value of another field/asset in the same Config. Absent
// Condition = always active. Zero value is a valid empty leaf (matches
// nothing in particular — always true) so Condition can be composed by
// value, not just by pointer.
//
// It replaces the old, never-implemented FormatField.Exclude: instead of
// declaring "if I am filled — hide them" from the controlling field's
// side, the dependent field declares "I am active when it equals X" from
// its own side — easier to read top-to-bottom when generating a form (the
// field's own description already tells you when it appears).
type Condition struct {
	// Leaf predicate: name of the controlling field/asset + one of
	// In/NotIn/Present. Ignored when All/Any/Not below are non-empty.
	Field   string `json:"field,omitempty" yaml:"field,omitempty"`
	In      []any  `json:"in,omitempty" yaml:"in,omitempty"`
	NotIn   []any  `json:"not_in,omitempty" yaml:"not_in,omitempty"`
	Present *bool  `json:"present,omitempty" yaml:"present,omitempty"`

	// Composition: combine leaf/nested conditions with boolean logic. The
	// plain leaf form above still works unchanged for the common
	// single-field case.
	All []Condition `json:"all,omitempty" yaml:"all,omitempty"`
	Any []Condition `json:"any,omitempty" yaml:"any,omitempty"`
	Not *Condition  `json:"not,omitempty" yaml:"not,omitempty"`
}

// IsZero reports whether c carries no predicate at all (always true).
func (c Condition) IsZero() bool {
	return c.Field == "" && len(c.In) == 0 && len(c.NotIn) == 0 && c.Present == nil &&
		len(c.All) == 0 && len(c.Any) == 0 && c.Not == nil
}

// Eval evaluates the condition tree against the current set of
// field/asset values (keyed by name, same map shape as used for
// Config.ActiveFields/ActiveAssets). It is a pure function without side
// effects — the schema (Condition) and the engine (Eval) stay separate,
// same principle as Field.Prepare.
//
// If the referenced value is a slice (i.e. the field it points to is a
// repeatable/multi-select Field or AssetRequirement, see MinCount/MaxCount
// and IsMultiple), In/NotIn match element-wise ("at least one of the
// selected elements intersects with the given set", like an element-wise
// SQL IN rather than a whole-slice comparison) and Present means "the
// slice is non-empty".
func (c Condition) Eval(values map[string]any) bool {
	if len(c.All) > 0 {
		for _, sub := range c.All {
			if !sub.Eval(values) {
				return false
			}
		}
		return true
	}
	if len(c.Any) > 0 {
		for _, sub := range c.Any {
			if sub.Eval(values) {
				return true
			}
		}
		return false
	}
	if c.Not != nil {
		return !c.Not.Eval(values)
	}
	if c.Field == "" {
		// Empty leaf — matches nothing in particular, always true.
		return true
	}

	v, present := values[c.Field]
	items, isSlice := asValueSlice(v)

	if c.Present != nil {
		if isSlice {
			return (len(items) > 0) == *c.Present
		}
		return (present && !isEmptyValue(v)) == *c.Present
	}
	if len(c.In) > 0 {
		if isSlice {
			return sliceIntersects(items, c.In)
		}
		return present && containsValue(c.In, v)
	}
	if len(c.NotIn) > 0 {
		if isSlice {
			return !sliceIntersects(items, c.NotIn)
		}
		return !present || !containsValue(c.NotIn, v)
	}
	return true
}

// FieldIn builds a Condition that is true when the named field's value is
// (or, for a repeatable field, includes) one of values.
func FieldIn(field string, values ...any) Condition {
	return Condition{Field: field, In: values}
}

// FieldNotIn builds a Condition that is true when the named field's value
// is not (or, for a repeatable field, does not include) any of values.
func FieldNotIn(field string, values ...any) Condition {
	return Condition{Field: field, NotIn: values}
}

// FieldPresent builds a Condition that is true when the named field has a
// non-empty value.
func FieldPresent(field string) Condition {
	t := true
	return Condition{Field: field, Present: &t}
}

// FieldAbsent builds a Condition that is true when the named field has no
// value (or an empty one).
func FieldAbsent(field string) Condition {
	f := false
	return Condition{Field: field, Present: &f}
}

// And combines conditions with boolean AND.
func And(conds ...Condition) Condition {
	return Condition{All: conds}
}

// Or combines conditions with boolean OR.
func Or(conds ...Condition) Condition {
	return Condition{Any: conds}
}

// Not negates c.
func Not(c Condition) Condition {
	return Condition{Not: &c}
}

func asValueSlice(v any) ([]any, bool) {
	switch s := v.(type) {
	case []any:
		return s, true
	case nil:
		return nil, false
	}
	return nil, false
}

func isEmptyValue(v any) bool {
	switch vv := v.(type) {
	case nil:
		return true
	case string:
		return vv == ""
	case []any:
		return len(vv) == 0
	}
	return false
}

func containsValue(list []any, v any) bool {
	for _, item := range list {
		if item == v {
			return true
		}
	}
	return false
}

func sliceIntersects(items, list []any) bool {
	for _, item := range items {
		if containsValue(list, item) {
			return true
		}
	}
	return false
}
