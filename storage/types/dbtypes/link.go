package dbtypes

import (
	"sort"

	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/adcorelib/models"
)

// LinkList represented in DB
type LinkList []*models.AdLink

// Result of data as a list
func (d LinkList) Result() []any {
	return gocast.Slice[any](d)
}

// Reset stored data
func (d *LinkList) Reset() {
	if d != nil {
		*d = (*d)[:0]
	}
}

// Target real of the list
func (d *LinkList) Target() any {
	return (*[]*models.AdLink)(d)
}

// Merge loaded data
func (d *LinkList) Merge(l any) {
	newData := make([]*models.AdLink, 0, len(*d))
	for _, it := range *d {
		if it.Status.IsApproved() && it.Active.IsActive() {
			newData = append(newData, it)
		}
	}
	for _, it := range l.([]*models.AdLink) {
		if !it.Status.IsApproved() || !it.Active.IsActive() {
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
