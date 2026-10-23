package searchtypes

import (
	"reflect"
	"testing"
)

type objItem struct {
	id   uint
	name string
}

func objCmp(a, b any) int {
	ida, idb := a.(objItem).id, b.(objItem).id
	if ida < idb {
		return -1
	}
	if ida > idb {
		return 1
	}
	return 0
}

func objValue(a any) uint {
	return a.(objItem).id
}

func TestObjectBitset_SetHasValues(t *testing.T) {
	b := NewObjectBitset(objCmp, objValue,
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
	b := NewObjectBitset(objCmp, objValue, item, item)
	if b.Len() != 1 {
		t.Fatalf("Len=%d want 1", b.Len())
	}
}

func TestObjectBitset_Unset(t *testing.T) {
	b := NewObjectBitset(objCmp, objValue,
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
	b := NewObjectBitset(objCmp, objValue, objItem{1, "a"})
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

func TestObjectBitset_HasOneOf(t *testing.T) {
	var nilBits *ObjectBitset
	if nilBits.HasOneOf(objItem{1, "a"}) {
		t.Fatal("nil HasOneOf should be false")
	}

	empty := NewObjectBitset(objCmp, objValue)
	if empty.HasOneOf(objItem{1, "a"}) || empty.HasOneOf() {
		t.Fatal("empty HasOneOf should be false")
	}

	b := NewObjectBitset(objCmp, objValue,
		objItem{1, "a"},
		objItem{3, "c"},
	)
	if b.HasOneOf() {
		t.Fatal("no args should be false")
	}
	if !b.HasOneOf(objItem{9, "x"}, objItem{3, "c"}, objItem{4, "d"}) {
		t.Fatal("expected a hit on 3")
	}
	if b.HasOneOf(objItem{2, "b"}, objItem{4, "d"}) {
		t.Fatal("all misses should be false")
	}
	if b.HasOneOf(objItem{1, "other"}) {
		t.Fatal("same id with a different name is not a member")
	}

	onlyOne := NewObjectBitset(objCmp, objValue, objItem{1, "one"})
	if onlyOne.HasOneOf(objItem{65, "sixtyfive"}) {
		t.Fatal("65 shares a mask bit with 1 but is not in the set")
	}
	with65 := NewObjectBitset(objCmp, objValue, objItem{65, "sixtyfive"})
	if !with65.HasOneOf(objItem{1, "one"}, objItem{65, "sixtyfive"}) {
		t.Fatal("65 should match despite the shared mask bit with 1")
	}
}

func TestObjectBitset_Mod64Collision(t *testing.T) {
	b := NewObjectBitset(objCmp, objValue,
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

func BenchmarkObjectBitset_Has(b *testing.B) {
	bits := NewObjectBitset(objCmp, objValue,
		objItem{1, "a"},
		objItem{2, "b"},
		objItem{3, "c"},
	)
	hit := any(objItem{2, "b"})
	maskMiss := any(objItem{64, "x"})
	b.Run("Hit", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sink = bits.Has(hit)
		}
	})
	b.Run("MaskMiss", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sink = bits.Has(maskMiss)
		}
	})
}

func BenchmarkObjectBitset_HasOneOf(b *testing.B) {
	bits := NewObjectBitset(objCmp, objValue,
		objItem{1, "a"},
		objItem{2, "b"},
		objItem{3, "c"},
	)
	shortHit := []any{objItem{8, "x"}, objItem{9, "y"}, objItem{10, "z"}, objItem{3, "c"}}
	maskMiss := []any{objItem{64, "a"}, objItem{128, "b"}, objItem{192, "c"}}
	collisionSet := NewObjectBitset(objCmp, objValue, objItem{1, "one"})
	collision := []any{objItem{65, "sixtyfive"}}

	longMiss := make([]any, 128)
	for i := range longMiss {
		longMiss[i] = objItem{uint(10 + i), "miss"}
	}
	longHit := append([]any(nil), longMiss...)
	longHit[len(longHit)-1] = objItem{2, "b"}

	b.Run("ShortHit", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sink = bits.HasOneOf(shortHit...)
		}
	})
	b.Run("MaskMiss", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sink = bits.HasOneOf(maskMiss...)
		}
	})
	b.Run("Collision", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sink = collisionSet.HasOneOf(collision...)
		}
	})
	b.Run("LongMiss", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sink = bits.HasOneOf(longMiss...)
		}
	})
	b.Run("LongHit", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sink = bits.HasOneOf(longHit...)
		}
	})
}
