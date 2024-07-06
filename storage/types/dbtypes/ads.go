package dbtypes

import (
	"sort"

	"github.com/demdxx/gocast/v2"
	"gorm.io/gorm"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/models"
)

// AdList represented in DB
type AdList []*models.Ad

// Result of data as a list
func (d AdList) Result() []any {
	return gocast.Slice[any](d)
}

// Reset stored data
func (d *AdList) Reset() {
	if d != nil {
		*d = (*d)[:0]
	}
}

func (d *AdList) PrepareQuery(db *gorm.DB) *gorm.DB {
	return db.Preload("Assets", "processing_status=?", types.ProcessingProcessed.Code())
}

// Target real of the list
func (d *AdList) Target() any {
	return (*[]*models.Ad)(d)
}

// Merge loaded data
func (d *AdList) Merge(l any) {
	newData := make([]*models.Ad, 0, len(*d))
	for _, it := range *d {
		if it.Status.IsApproved() && it.Active.IsActive() {
			newData = append(newData, it)
		}
	}
	for _, it := range l.([]*models.Ad) {
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
