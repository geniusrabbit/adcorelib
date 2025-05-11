package bidrequest

import (
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/searchtypes"
)

type AdFormats struct {
	formats        []*types.Format                // Cached formats
	formatBitset   searchtypes.NumberBitset[uint] // Bitset for format IDs
	formatTypeMask types.FormatTypeBitset         // Bitmask for format types
}

func (a *AdFormats) Reset() {
	a.formats = a.formats[:0]
	a.formatBitset.Reset()
	a.formatTypeMask.Reset()
}

func (a *AdFormats) Add(format *types.Format) bool {
	if format == nil || a.formatBitset.Has(uint(format.ID)) {
		return false
	}
	a.formats = append(a.formats, format)
	a.formatBitset.Set(uint(format.ID))
	a.formatTypeMask.SetOneBitset(format.Types)
	return true
}

func (a *AdFormats) List() []*types.Format {
	return a.formats
}

func (a *AdFormats) Bitset() *searchtypes.NumberBitset[uint] {
	return &a.formatBitset
}

func (a *AdFormats) TypeMask() types.FormatTypeBitset {
	return a.formatTypeMask
}
