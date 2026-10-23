package searchtypes

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestNumberBitset(t *testing.T) {
	bits := NewNumberBitset[uint](1, 3, 5, 7)

	bits.Set(2)
	if !bits.Has(2) {
		t.Error("bits Set/Has not working")
	}

	bits = bits.Unset(3, 1).Unset(3, 7)
	if bits.Has(3) || bits.Has(1) || bits.Has(7) {
		t.Error("bits Unset/Has not working")
	}

	if !reflect.DeepEqual([]uint{2, 5}, bits.Values()) {
		t.Errorf("bits Values problems: %v", bits.Values())
	}
}

func TestNumberBitset_NilReceiver(t *testing.T) {
	var bits *NumberBitset[uint]
	if bits.Len() != 0 {
		t.Fatalf("nil Len=%d want 0", bits.Len())
	}
	if bits.Values() != nil {
		t.Fatalf("nil Values=%v want nil", bits.Values())
	}
	if bits.Has(1) {
		t.Fatal("nil Has should be false")
	}
}

func TestNumberBitset_Empty(t *testing.T) {
	bits := NewNumberBitset[uint]()
	if bits.Len() != 0 {
		t.Fatalf("empty Len=%d", bits.Len())
	}
	if bits.Has(0) {
		t.Fatal("empty should not Has(0)")
	}
	if bits.Mask() != 0 {
		t.Fatalf("empty Mask=%d want 0", bits.Mask())
	}
}

func TestNumberBitset_DuplicatesAndSort(t *testing.T) {
	bits := NewNumberBitset(5, 1, 5, 3, 1)
	if bits.Len() != 3 {
		t.Fatalf("Len=%d want 3", bits.Len())
	}
	if !reflect.DeepEqual([]int{1, 3, 5}, bits.Values()) {
		t.Fatalf("Values=%v want sorted unique", bits.Values())
	}
}

func TestNumberBitset_Mod64Collision(t *testing.T) {
	// 1 and 65 share the same mask bit; Has must still distinguish them.
	bits := NewNumberBitset[uint](1, 65)
	if !bits.Has(1) || !bits.Has(65) {
		t.Fatal("both colliding values must be present")
	}
	if bits.Has(129) {
		t.Fatal("129 not inserted")
	}
	bits = bits.Unset(1)
	if bits.Has(1) {
		t.Fatal("1 should be removed")
	}
	if !bits.Has(65) {
		t.Fatal("65 must remain after Unset(1)")
	}
}

func TestNumberBitset_UnsetMissingNoop(t *testing.T) {
	orig := NewNumberBitset[uint](1, 2)
	same := orig.Unset(99)
	if same != orig {
		t.Fatal("Unset missing should return same pointer")
	}
	if !reflect.DeepEqual([]uint{1, 2}, same.Values()) {
		t.Fatalf("Values changed: %v", same.Values())
	}
}

func TestNumberBitset_UnsetEdges(t *testing.T) {
	t.Run("first", func(t *testing.T) {
		b := NewNumberBitset[uint](1, 2, 3).Unset(1)
		if !reflect.DeepEqual([]uint{2, 3}, b.Values()) {
			t.Fatalf("got %v", b.Values())
		}
	})
	t.Run("last", func(t *testing.T) {
		b := NewNumberBitset[uint](1, 2, 3).Unset(3)
		if !reflect.DeepEqual([]uint{1, 2}, b.Values()) {
			t.Fatalf("got %v", b.Values())
		}
	})
	t.Run("middle", func(t *testing.T) {
		b := NewNumberBitset[uint](1, 2, 3).Unset(2)
		if !reflect.DeepEqual([]uint{1, 3}, b.Values()) {
			t.Fatalf("got %v", b.Values())
		}
	})
	t.Run("all", func(t *testing.T) {
		b := NewNumberBitset[uint](1, 2).Unset(1, 2)
		if b.Len() != 0 {
			t.Fatalf("Len=%d want 0", b.Len())
		}
	})
	t.Run("only", func(t *testing.T) {
		b := NewNumberBitset[uint](7).Unset(7)
		if b.Len() != 0 || b.Has(7) {
			t.Fatalf("single Unset failed: %v", b.Values())
		}
	})
}

func TestNumberBitset_Reset(t *testing.T) {
	bits := NewNumberBitset[uint](1, 2, 3)
	bits.Reset()
	if bits.Len() != 0 || bits.Mask() != 0 || bits.Has(1) {
		t.Fatalf("Reset incomplete: len=%d mask=%d has1=%v", bits.Len(), bits.Mask(), bits.Has(1))
	}
	bits.Set(9)
	if !bits.Has(9) || bits.Len() != 1 {
		t.Fatal("Set after Reset failed")
	}
}

func TestNumberBitset_ContainsAllFrom(t *testing.T) {
	a := NewNumberBitset[uint](1, 2, 3)
	b := NewNumberBitset[uint](1, 2, 3, 4)
	c := NewNumberBitset[uint](1, 9)

	if !a.ContainsAllFrom(b) {
		t.Fatal("a should be subset of b")
	}
	if b.ContainsAllFrom(a) {
		t.Fatal("b should not be subset of a")
	}
	if a.ContainsAllFrom(c) {
		t.Fatal("a should not be subset of c")
	}
	//nolint:staticcheck
	if !NewNumberBitset[uint]().ContainsAllFrom(a) {
		// empty set: mask 0, 0&set.mask==0, then no values to check → true
		// Document actual behavior.
	}
	empty := NewNumberBitset[uint]()
	if !empty.ContainsAllFrom(a) {
		t.Fatal("empty ContainsAllFrom non-empty should be true (vacuous)")
	}
	if a.ContainsAllFrom(nil) {
		t.Fatal("ContainsAllFrom(nil) should be false")
	}
	if a.ContainsAllFrom(empty) {
		t.Fatal("non-empty should not be subset of empty")
	}
}

func TestNumberBitset_HasOneOf(t *testing.T) {
	var nilBits *NumberBitset[uint]
	if nilBits.HasOneOf(1) {
		t.Fatal("nil HasOneOf should be false")
	}

	empty := NewNumberBitset[uint]()
	if empty.HasOneOf(1) || empty.HasOneOf() {
		t.Fatal("empty HasOneOf should be false")
	}

	bits := NewNumberBitset[uint](1, 3, 5)
	if bits.HasOneOf() {
		t.Fatal("no args should be false")
	}
	if !bits.HasOneOf(9, 4, 3, 8) {
		t.Fatal("expected a hit on 3")
	}
	if bits.HasOneOf(2, 4, 6) {
		t.Fatal("all misses should be false")
	}

	onlyOne := NewNumberBitset[uint](1)
	if onlyOne.HasOneOf(65) {
		t.Fatal("65 shares a mask bit with 1 but is not in the set")
	}
	with65 := NewNumberBitset[uint](65)
	if !with65.HasOneOf(1, 65) {
		t.Fatal("65 should match despite the shared mask bit with 1")
	}

	// A long query uses the same binary search as a short one.
	if !bits.HasOneOf(10, 11, 12, 13, 14, 15, 16, 17, 5) {
		t.Fatal("long query should find 5")
	}
	if bits.HasOneOf(10, 11, 12, 13, 14, 15, 16, 17, 65) {
		t.Fatal("long query must not treat a mask collision as a hit")
	}
}

func TestNumberBitset_Int64(t *testing.T) {
	bits := NewNumberBitset[int64](1, 0, 2)
	if bits.Len() != 3 {
		t.Fatalf("Len=%d", bits.Len())
	}
	if !bits.Has(0) || !bits.Has(1) || !bits.Has(2) {
		t.Fatalf("Has failed: %v", bits.Values())
	}
	if !reflect.DeepEqual([]int64{0, 1, 2}, bits.Values()) {
		t.Fatalf("Values=%v", bits.Values())
	}
}

func BenchmarkNumberBitset(b *testing.B) {
	var (
		cursor   int
		variants []uint
		bits     = NewNumberBitset[uint](1, 3)
	)

	for i := 0; i < 10000; i++ {
		variants = append(variants, uint(rand.Uint64()%1000))
	}

	for i := 0; b.Loop(); i++ {
		switch {
		case i%2 == 0 || i%37 == 0:
			bits.Set(variants[cursor])
		case i%3 == 0 || i%5 == 0 || i%11 == 0:
			bits = bits.Unset(variants[cursor])
		default:
			_ = bits.Has(variants[cursor])
		}
		if cursor++; cursor >= len(variants) {
			cursor = 0
		}
	}
}

var sink bool

func BenchmarkNumberBitset_Has(b *testing.B) {
	bits := NewNumberBitset[uint](1, 2, 3)
	b.Run("Hit", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sink = bits.Has(2)
		}
	})
	b.Run("MaskMiss", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sink = bits.Has(64)
		}
	})
}

func BenchmarkNumberBitset_HasOneOf(b *testing.B) {
	bits := NewNumberBitset[uint](1, 2, 3)
	shortHit := []uint{8, 9, 10, 3}
	maskMiss := []uint{64, 128, 192}
	collisionSet := NewNumberBitset[uint](1)
	collision := []uint{65}

	longMiss := make([]uint, 128)
	for i := range longMiss {
		longMiss[i] = uint(10 + i)
	}
	longHit := append([]uint(nil), longMiss...)
	longHit[len(longHit)-1] = 2

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
