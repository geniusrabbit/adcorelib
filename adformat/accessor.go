package adformat

import "github.com/geniusrabbit/adcorelib/searchtypes"

// Accessor object interface (was admodels/types.FormatsAccessor).
type Accessor interface {
	// Formats list collection.
	Formats() []*Format

	// FormatsBySize returns the list of formats whose Sizes accept
	// width/height (§3.2.13); direct formats are always included, same
	// convention as the old model.
	FormatsBySize(width, height int) []*Format

	// FormatByID of the model.
	FormatByID(id uint64) *Format

	// FormatByCode of the model.
	FormatByCode(code string) *Format

	// DirectFormatSet to search.
	DirectFormatSet() *searchtypes.NumberBitset[uint]
}
