//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package searchtypes

import (
	"sort"
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
			sort.Slice(b.values, func(i, j int) bool { return b.values[i] < b.values[j] })
		}
	}
	return b
}

// Unset type values
func (b *NumberBitset[T]) Unset(vals ...T) *NumberBitset[T] {
	newVals := b.values
	changed := false
	for _, v := range vals {
		idx := sort.Search(len(newVals), func(i int) bool {
			return newVals[i] >= v
		})
		if idx < len(newVals) && newVals[idx] == v {
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
	if b != nil && b.mask&(1<<uint64(v%64)) != 0 {
		idx := sort.Search(b.Len(), func(i int) bool {
			return b.values[i] >= v
		})
		return idx >= 0 && idx < b.Len() && b.values[idx] == v
	}
	return false
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
