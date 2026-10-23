//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package searchtypes

import (
	"slices"
)

// Integer type constraint for generic NumberBitset
type Integer interface {
	~uint | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// NumberBitset any numbers
type NumberBitset[T Integer] struct {
	values []T
	mask   uint64
}

// NewNumberBitset from numbers
func NewNumberBitset[T Integer](vals ...T) (b *NumberBitset[T]) {
	return (&NumberBitset[T]{}).Set(vals...)
}

// Len of the elements
func (b *NumberBitset[T]) Len() int {
	if b == nil {
		return 0
	}
	return len(b.values)
}

// Mask of the set
func (b *NumberBitset[T]) Mask() uint64 {
	return b.mask
}

// Values list
func (b *NumberBitset[T]) Values() []T {
	if b == nil {
		return nil
	}
	return b.values
}

// Set type values
func (b *NumberBitset[T]) Set(vals ...T) *NumberBitset[T] {
	for _, v := range vals {
		if !b.Has(v) {
			b.mask |= 1 << uint64(v%64)
			b.values = append(b.values, v)
			// Keep sorted so subsequent Has() in this call sees a valid order.
			slices.Sort(b.values)
		}
	}
	return b
}

// Unset type values
func (b *NumberBitset[T]) Unset(vals ...T) *NumberBitset[T] {
	newVals := b.values
	changed := false
	for _, v := range vals {
		idx, found := slices.BinarySearch(newVals, v)
		if found {
			newVals = append(newVals[:idx:idx], newVals[idx+1:]...)
			changed = true
		}
	}
	if !changed {
		return b
	}
	return NewNumberBitset(newVals...)
}

// Has type in bitset
func (b *NumberBitset[T]) Has(v T) bool {
	if b == nil || b.mask&(1<<uint64(v%64)) == 0 {
		return false
	}
	return b.hasSorted(v)
}

// HasOneOf reports whether the set contains any of vals.
func (b *NumberBitset[T]) HasOneOf(vals ...T) bool {
	if b == nil || len(b.values) == 0 || len(vals) == 0 {
		return false
	}
	mask := b.mask
	for _, v := range vals {
		if mask&(1<<uint64(v%64)) == 0 {
			continue
		}
		if b.hasSorted(v) {
			return true
		}
	}
	return false
}

func (b *NumberBitset[T]) hasSorted(v T) bool {
	values := b.values
	lo, hi := 0, len(values)
	for lo < hi {
		mid := int(uint(lo+hi) >> 1)
		if values[mid] < v {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo < len(values) && values[lo] == v
}

// Reset bitset value
func (b *NumberBitset[T]) Reset() *NumberBitset[T] {
	b.mask = 0
	if b.values != nil {
		b.values = b.values[:0]
	}
	return b
}

// ContainsAllFrom items from the set
func (b *NumberBitset[T]) ContainsAllFrom(set *NumberBitset[T]) (res bool) {
	if set != nil && b.mask&set.mask == b.mask {
		res = true
		for _, v := range b.values {
			if res = set.Has(v); !res {
				break
			}
		}
	}
	return res
}
