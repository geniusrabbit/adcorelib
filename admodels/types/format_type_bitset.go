//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package types

// FormatTypeBitset of types
type FormatTypeBitset uint

const (
	// FormatTypeBitsetEmpty by default
	FormatTypeBitsetEmpty FormatTypeBitset = 0
	// FormatTypeBitsetDirect for direct type
	FormatTypeBitsetDirect FormatTypeBitset = 1 << FormatTypeBitset(FormatDirectType)
)

// NewFormatTypeBitset from types
func NewFormatTypeBitset(types ...FormatType) *FormatTypeBitset {
	return new(FormatTypeBitset).Set(types...)
}

// Set format type values
func (b *FormatTypeBitset) Set(types ...FormatType) *FormatTypeBitset {
	for _, t := range types {
		*b |= 1 << FormatTypeBitset(t)
	}
	return b
}

// SetOne format type value
func (b *FormatTypeBitset) SetOne(t FormatType) *FormatTypeBitset {
	if t != FormatInvalidType && t != FormatUndefinedType {
		*b |= 1 << FormatTypeBitset(t)
	}
	return b
}

// SetBitset intersection
func (b *FormatTypeBitset) SetBitset(bitset ...FormatTypeBitset) *FormatTypeBitset {
	for _, v := range bitset {
		*b |= v
	}
	return b
}

// SetOneBitset set one bitset
func (b *FormatTypeBitset) SetOneBitset(bitset FormatTypeBitset) *FormatTypeBitset {
	*b |= bitset
	return b
}

// SetFromFormats type values
func (b *FormatTypeBitset) SetFromFormats(formats ...*Format) *FormatTypeBitset {
	for _, f := range formats {
		b.SetBitset(f.Types)
	}
	return b
}

// Reset type value
func (b *FormatTypeBitset) Reset() *FormatTypeBitset {
	*b = 0
	return b
}

// Unset type values
func (b *FormatTypeBitset) Unset(types ...FormatType) *FormatTypeBitset {
	for _, t := range types {
		*b &= ^(1 << FormatTypeBitset(t))
	}
	return b
}

// Has type in bitset
func (b *FormatTypeBitset) Has(t FormatType) bool {
	return b != nil && (*b)&(1<<FormatTypeBitset(t)) != 0
}

// Is type only
func (b *FormatTypeBitset) Is(t FormatType) bool {
	return b != nil && *b == 1<<FormatTypeBitset(t)
}

// Intersec between two bitsets
func (b FormatTypeBitset) Intersec(set FormatTypeBitset) FormatTypeBitset {
	return b & set
}

// IsIntersec between two bitsets
func (b *FormatTypeBitset) IsIntersec(set FormatTypeBitset) bool {
	return b != nil && *b&set != 0
}

// Types list from bites
func (b FormatTypeBitset) Types() (list []FormatType) {
	for _, t := range FormatTypeList {
		if b&(1<<FormatTypeBitset(t)) != 0 {
			list = append(list, t)
		}
	}
	return list
}

// HasOneType in the mask
func (b FormatTypeBitset) HasOneType() FormatType {
	for i := 0; i < 32; i++ {
		if item := b & (1 << FormatTypeBitset(i)); item == b {
			return FormatType(i)
		} else if item != 0 {
			break
		}
	}
	return FormatInvalidType
}

// FirstType from the bitset
func (b FormatTypeBitset) FirstType() FormatType {
	for i := 1; i < 32; i++ {
		if item := b & (1 << FormatTypeBitset(i)); item != 0 {
			return FormatType(i)
		}
	}
	return FormatUndefinedType
}

// IsEmpty format type
func (b FormatTypeBitset) IsEmpty() bool {
	return b == FormatTypeBitsetEmpty
}
