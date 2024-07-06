package dbtypes

import (
	"sort"

	"github.com/demdxx/gocast/v2"

	"github.com/geniusrabbit/adcorelib/models"
)

// FormatList represented in DB
type FormatList []*models.Format

// Result of data as a list
func (d FormatList) Result() []any {
	return gocast.Slice[any](d)
}

// Reset stored data
func (d *FormatList) Reset() {
	if d != nil {
		*d = (*d)[:0]
	}
}

// Target real of the list
func (d *FormatList) Target() any {
	return (*[]*models.Format)(d)
}

// Merge loaded data
func (d *FormatList) Merge(l any) {
	newData := make([]*models.Format, 0, len(*d))
	for _, it := range *d {
		if it.Active.IsActive() {
			newData = append(newData, it)
		}
	}
	for _, it := range l.([]*models.Format) {
		if !it.Active.IsActive() {
			continue
		}
		i := sort.Search(len(newData), func(i int) bool { return newData[i].ID >= it.ID })
		if i >= 0 && i < len(newData) && newData[i].ID == it.ID {
			newData[i] = it
		} else {
			newData = append(newData, it)
		}
	}
	sort.Slice(newData, func(i, j int) bool { return newData[i].ID < newData[j].ID })
	*d = newData
}
