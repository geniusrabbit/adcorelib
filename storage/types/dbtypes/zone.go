package dbtypes

import (
	"sort"

	"geniusrabbit.dev/adcorelib/models"
	"github.com/demdxx/gocast/v2"
)

// ZoneList represented in DB
type ZoneList []*models.Zone

// Result of data as a list
func (d ZoneList) Result() []any {
	return gocast.Slice[any](d)
}

// Reset stored data
func (d *ZoneList) Reset() {
	if d != nil {
		*d = (*d)[:0]
	}
}

// Target real of the list
func (d *ZoneList) Target() any {
	return (*[]*models.Zone)(d)
}

// Merge loaded data
func (d *ZoneList) Merge(l any) {
	newData := make([]*models.Zone, 0, len(*d))
	for _, it := range *d {
		if it.Status.IsApproved() && it.Active.IsActive() {
			newData = append(newData, it)
		}
	}
	for _, it := range l.([]*models.Zone) {
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
