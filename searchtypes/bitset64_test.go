//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package searchtypes

import (
	"reflect"
	"testing"
)

func TestBitset64(t *testing.T) {
	bits := NewBitset64(1, 3)

	bits.Set(2)
	if !bits.Has(2) {
		t.Error("bits Set/Has not working")
	}

	bits.Unset(3)
	if bits.Has(3) {
		t.Error("bits Unet/Has not working")
	}

	if !reflect.DeepEqual([]uint{1, 2}, bits.Numbers()) {
		t.Errorf("bits Numbers problems: %v", bits.Numbers())
	}
}

func TestBitset64_Empty(t *testing.T) {
	var bits Bitset64
	if bits.Has(0) {
		t.Fatal("empty Has(0) should be false")
	}
	if got := bits.Numbers(); len(got) != 0 {
		t.Fatalf("empty Numbers=%v", got)
	}
}

func TestBitset64_Mod64(t *testing.T) {
	bits := NewBitset64(1, 65) // 65%64 == 1 — same bit
	if !bits.Has(1) || !bits.Has(65) {
		t.Fatal("mod64 aliases should share Has")
	}
	bits.Unset(65)
	if bits.Has(1) {
		t.Fatal("Unset(65) clears bit 1 as well")
	}
}

func TestBitset64_AllBits(t *testing.T) {
	var bits Bitset64
	for i := uint(0); i < 64; i++ {
		bits.Set(i)
	}
	nums := bits.Numbers()
	if len(nums) != 64 {
		t.Fatalf("Numbers len=%d want 64", len(nums))
	}
	for i := uint(0); i < 64; i++ {
		if nums[i] != i {
			t.Fatalf("Numbers[%d]=%d", i, nums[i])
		}
		if !bits.Has(i) {
			t.Fatalf("missing bit %d", i)
		}
	}
	for i := uint(0); i < 64; i++ {
		bits.Unset(i)
	}
	if bits != 0 {
		t.Fatalf("after Unset all, bits=%d", bits)
	}
}

func TestBitset64_NilPointerHas(t *testing.T) {
	var bits *Bitset64
	if bits.Has(1) {
		t.Fatal("nil pointer Has should be false")
	}
}
