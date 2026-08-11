//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package searchtypes

import (
	"sort"
)

// ObjectBitset any numbers
type ObjectBitset struct {
	less   func(o1, o2 any) bool
	value  func(o1 any) uint
	values []any
	mask   uint64
}

// NewObjectBitset from numbers
func NewObjectBitset(less func(o1, o2 any) bool, value func(o1 any) uint, vals ...any) (b *ObjectBitset) {
	return (&ObjectBitset{less: less, value: value}).Set(vals...)
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
			sort.Slice(b.values, func(i, j int) bool {
				return b.less(b.values[i], b.values[j])
			})
		}
	}
	return b
}

// Unset type values
func (b *ObjectBitset) Unset(vals ...any) *ObjectBitset {
	newVals := b.values
	changed := false
	for _, v := range vals {
		idx := sort.Search(len(newVals), func(i int) bool {
			return !b.less(newVals[i], v) // first >= v
		})
		if idx < len(newVals) && newVals[idx] == v {
			newVals = append(newVals[:idx:idx], newVals[idx+1:]...)
			changed = true
		}
	}
	if !changed {
		return b
	}
	return NewObjectBitset(b.less, b.value, newVals...)
}

// Has type in bitset
func (b *ObjectBitset) Has(v any) bool {
	if b != nil && b.mask&(1<<uint64(b.value(v)%64)) != 0 {
		idx := sort.Search(b.Len(), func(i int) bool {
			return !b.less(b.values[i], v) // first >= v
		})
		return idx < b.Len() && b.values[idx] == v
	}
	return false
}
