package adformat

import (
	"sort"

	"github.com/geniusrabbit/adcorelib/searchtypes"
)

// SimpleAccessor implements a basic in-memory Accessor (was
// admodels/types.SimpleFormatAccessor). It accepts either a static Go
// slice/var-table of Format values (§3.2.10) or the result of loading
// codec.go-produced YAML/JSON — both paths produce the same Format type,
// so the accessor (and any downstream validation engine) never needs to
// know how a Format was built.
type SimpleAccessor struct {
	formatList []*Format
	directs    *searchtypes.NumberBitset[uint]
}

// NewSimpleAccessor from a fixed list of formats. Accepts value or
// pointer formats for convenience with both Go var-table literals and
// codec.go decode results.
func NewSimpleAccessor(formats ...*Format) *SimpleAccessor {
	acc := &SimpleAccessor{formatList: formats}
	acc.Prepare()
	return acc
}

// Prepare sorts the format list by ID — required for FormatByID's binary
// search to work; called automatically by NewSimpleAccessor, but exposed
// for callers that mutate Formats() directly afterwards.
func (fa *SimpleAccessor) Prepare() {
	sort.Slice(fa.formatList, func(i, j int) bool { return fa.formatList[i].ID < fa.formatList[j].ID })
}

// Formats list collection.
func (fa *SimpleAccessor) Formats() []*Format {
	return fa.formatList
}

// FormatsBySize returns the list of formats whose Sizes accept
// width/height; direct formats are always included, since they carry no
// display area at all.
func (fa *SimpleAccessor) FormatsBySize(width, height int) []*Format {
	list := make([]*Format, 0, len(fa.formatList))
	for _, frmt := range fa.formatList {
		if frmt.IsDirect() || frmt.Suits(width, height) {
			list = append(list, frmt)
		}
	}
	return list
}

// FormatByID of the model — requires Prepare (called by
// NewSimpleAccessor) to have run first, otherwise the binary search may
// silently return the wrong result.
func (fa *SimpleAccessor) FormatByID(id uint64) *Format {
	i := sort.Search(len(fa.formatList), func(i int) bool { return fa.formatList[i].ID >= id })
	if i >= 0 && i < len(fa.formatList) && fa.formatList[i].ID == id {
		return fa.formatList[i]
	}
	return nil
}

// FormatByCode of the model.
func (fa *SimpleAccessor) FormatByCode(code string) *Format {
	for _, f := range fa.formatList {
		if f.Codename == code {
			return f
		}
	}
	return nil
}

// DirectFormatSet to search.
func (fa *SimpleAccessor) DirectFormatSet() *searchtypes.NumberBitset[uint] {
	if fa.directs == nil {
		fa.directs = new(searchtypes.NumberBitset[uint])
		for _, f := range fa.formatList {
			if f.IsDirect() {
				fa.directs.Set(uint(f.ID))
			}
		}
	}
	return fa.directs
}

var _ Accessor = (*SimpleAccessor)(nil)
