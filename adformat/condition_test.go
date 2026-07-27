package adformat

import "testing"

func TestConditionEvalLeaf(t *testing.T) {
	cond := FieldIn("adtype", "business")

	if !cond.Eval(map[string]any{"adtype": "business"}) {
		t.Error("expected condition to match adtype=business")
	}
	if cond.Eval(map[string]any{"adtype": "game"}) {
		t.Error("expected condition not to match adtype=game")
	}
	if cond.Eval(map[string]any{}) {
		t.Error("expected condition not to match when field is absent")
	}
}

func TestConditionEvalNotIn(t *testing.T) {
	cond := FieldNotIn("adtype", "business")
	if !cond.Eval(map[string]any{"adtype": "game"}) {
		t.Error("expected NotIn to match a different value")
	}
	if cond.Eval(map[string]any{"adtype": "business"}) {
		t.Error("expected NotIn not to match the excluded value")
	}
	if !cond.Eval(map[string]any{}) {
		t.Error("expected NotIn to match when the field is absent")
	}
}

func TestConditionEvalPresent(t *testing.T) {
	present := FieldPresent("phone")
	if !present.Eval(map[string]any{"phone": "123"}) {
		t.Error("expected Present to match a non-empty value")
	}
	if present.Eval(map[string]any{"phone": ""}) {
		t.Error("expected Present not to match an empty string")
	}
	if present.Eval(map[string]any{}) {
		t.Error("expected Present not to match a missing field")
	}

	absent := FieldAbsent("phone")
	if !absent.Eval(map[string]any{}) {
		t.Error("expected Absent to match a missing field")
	}
}

func TestConditionEvalComposition(t *testing.T) {
	cond := And(
		FieldIn("adtype", "business"),
		FieldIn("contact_method", "phone"),
	)
	if !cond.Eval(map[string]any{"adtype": "business", "contact_method": "phone"}) {
		t.Error("expected And to match when both are true")
	}
	if cond.Eval(map[string]any{"adtype": "business", "contact_method": "email"}) {
		t.Error("expected And not to match when one is false")
	}

	orCond := Or(FieldIn("a", "1"), FieldIn("b", "2"))
	if !orCond.Eval(map[string]any{"b": "2"}) {
		t.Error("expected Or to match when one branch is true")
	}

	notCond := Not(FieldIn("a", "1"))
	if !notCond.Eval(map[string]any{"a": "2"}) {
		t.Error("expected Not to invert its inner condition")
	}
}

func TestConditionEvalZeroValueAlwaysTrue(t *testing.T) {
	var zero Condition
	if !zero.Eval(map[string]any{}) {
		t.Error("expected zero-value Condition to always evaluate true")
	}
	if !zero.IsZero() {
		t.Error("expected zero-value Condition to report IsZero")
	}
}

// TestConditionEvalMultipleValues covers §3.2.15's slice semantics: when
// the referenced field is itself repeatable (IsMultiple), In/NotIn match
// element-wise and Present means "the slice is non-empty".
func TestConditionEvalMultipleValues(t *testing.T) {
	values := map[string]any{"categories": []any{"fashion", "home"}}

	if !FieldIn("categories", "fashion").Eval(values) {
		t.Error("expected In to match when one selected element intersects")
	}
	if FieldIn("categories", "electronics").Eval(values) {
		t.Error("expected In not to match when no selected element intersects")
	}
	if !FieldNotIn("categories", "electronics").Eval(values) {
		t.Error("expected NotIn to match when no selected element intersects the excluded set")
	}
	if !FieldPresent("categories").Eval(values) {
		t.Error("expected Present to match a non-empty slice")
	}

	empty := map[string]any{"categories": []any{}}
	if FieldPresent("categories").Eval(empty) {
		t.Error("expected Present not to match an empty slice")
	}
}
