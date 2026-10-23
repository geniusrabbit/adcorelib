//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package searchtypes

import (
	"slices"
)

// ObjectBitset any numbers
type ObjectBitset struct {
	cmp    func(a, b any) int // <0, 0, >0
	value  func(o1 any) uint
	values []any
	mask   uint64
}

// NewObjectBitset from numbers
func NewObjectBitset(cmp func(a, b any) int, value func(o1 any) uint, vals ...any) (b *ObjectBitset) {
	return (&ObjectBitset{cmp: cmp, value: value}).Set(vals...)
}

// Len of the elements
func (b *ObjectBitset) Len() int {
	if b == nil {
		return 0
	}
	return len(b.values)
}

// Values list
func (b *ObjectBitset) Values() []any {
	if b == nil {
		return nil
	}
	return b.values
}

// Set type values
func (b *ObjectBitset) Set(vals ...any) *ObjectBitset {
	for _, v := range vals {
		if !b.Has(v) {
			b.mask |= 1 << uint64(b.value(v)%64)
			b.values = append(b.values, v)
			slices.SortFunc(b.values, b.cmp)
		}
	}
	return b
}

// Unset type values
func (b *ObjectBitset) Unset(vals ...any) *ObjectBitset {
	newVals := b.values
	changed := false
	for _, v := range vals {
		idx, found := slices.BinarySearchFunc(newVals, v, b.cmp)
		if found {
			newVals = append(newVals[:idx:idx], newVals[idx+1:]...)
			changed = true
		}
	}
	if !changed {
		return b
	}
	return NewObjectBitset(b.cmp, b.value, newVals...)
}

// Has type in bitset
func (b *ObjectBitset) Has(v any) bool {
	if b == nil || b.mask&(1<<uint64(b.value(v)%64)) == 0 {
		return false
	}
	return b.hasSorted(v)
}

// HasOneOf reports whether the set contains any of vals.
// Membership is the same == check as Has, not cmp == 0.
func (b *ObjectBitset) HasOneOf(vals ...any) bool {
	if b == nil || len(b.values) == 0 || len(vals) == 0 {
		return false
	}
	mask := b.mask
	for _, v := range vals {
		if mask&(1<<uint64(b.value(v)%64)) == 0 {
			continue
		}
		if b.hasSorted(v) {
			return true
		}
	}
	return false
}

func (b *ObjectBitset) hasSorted(v any) bool {
	values := b.values
	lo, hi := 0, len(values)
	for lo < hi {
		mid := int(uint(lo+hi) >> 1)
		if b.cmp(values[mid], v) < 0 {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo < len(values) && values[lo] == v
}
