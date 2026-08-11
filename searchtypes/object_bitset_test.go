package searchtypes

import (
	"reflect"
	"testing"
)

type objItem struct {
	id   uint
	name string
}

func objLess(a, b any) bool {
	return a.(objItem).id < b.(objItem).id
}

func objValue(a any) uint {
	return a.(objItem).id
}

func TestObjectBitset_SetHasValues(t *testing.T) {
	b := NewObjectBitset(objLess, objValue,
		objItem{3, "c"},
		objItem{1, "a"},
		objItem{2, "b"},
	)

	if b.Len() != 3 {
		t.Fatalf("Len=%d want 3", b.Len())
	}
	if !b.Has(objItem{2, "b"}) {
		t.Fatal("Has(2) failed")
	}
	if b.Has(objItem{9, "x"}) {
		t.Fatal("Has missing should be false")
	}

	want := []any{objItem{1, "a"}, objItem{2, "b"}, objItem{3, "c"}}
	if !reflect.DeepEqual(want, b.Values()) {
		t.Fatalf("Values=%v want sorted by id", b.Values())
	}
}

func TestObjectBitset_DuplicateSet(t *testing.T) {
	item := objItem{1, "a"}
	b := NewObjectBitset(objLess, objValue, item, item)
	if b.Len() != 1 {
		t.Fatalf("Len=%d want 1", b.Len())
	}
}

func TestObjectBitset_Unset(t *testing.T) {
	b := NewObjectBitset(objLess, objValue,
		objItem{1, "a"},
		objItem{2, "b"},
		objItem{3, "c"},
	)
	b2 := b.Unset(objItem{2, "b"})
	if b2.Has(objItem{2, "b"}) {
		t.Fatal("Unset failed")
	}
	if !reflect.DeepEqual([]any{objItem{1, "a"}, objItem{3, "c"}}, b2.Values()) {
		t.Fatalf("Values=%v", b2.Values())
	}
	// Original unchanged (Unset returns new set when changed).
	if !b.Has(objItem{2, "b"}) {
		t.Fatal("original should still have item 2")
	}
}

func TestObjectBitset_UnsetMissingNoop(t *testing.T) {
	b := NewObjectBitset(objLess, objValue, objItem{1, "a"})
	same := b.Unset(objItem{9, "x"})
	if same != b {
		t.Fatal("Unset missing should return same pointer")
	}
}

func TestObjectBitset_NilReceiver(t *testing.T) {
	var b *ObjectBitset
	if b.Len() != 0 {
		t.Fatalf("nil Len=%d", b.Len())
	}
	if b.Values() != nil {
		t.Fatalf("nil Values=%v", b.Values())
	}
	if b.Has(objItem{1, "a"}) {
		t.Fatal("nil Has should be false")
	}
}

func TestObjectBitset_Mod64Collision(t *testing.T) {
	b := NewObjectBitset(objLess, objValue,
		objItem{1, "one"},
		objItem{65, "sixtyfive"},
	)
	if !b.Has(objItem{1, "one"}) || !b.Has(objItem{65, "sixtyfive"}) {
		t.Fatal("colliding ids must both be present")
	}
	b2 := b.Unset(objItem{1, "one"})
	if b2.Has(objItem{1, "one"}) {
		t.Fatal("1 should be removed")
	}
	if !b2.Has(objItem{65, "sixtyfive"}) {
		t.Fatal("65 must remain")
	}
}
