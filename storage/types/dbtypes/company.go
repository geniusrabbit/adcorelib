package dbtypes

import (
	"sort"

	"github.com/demdxx/gocast/v2"

	"geniusrabbit.dev/adcorelib/models"
)

// CompanyList represented in DB
type CompanyList []*models.Company

// Result of data as a list
func (d CompanyList) Result() []any {
	return gocast.Slice[any](d)
}

// Reset stored data
func (d *CompanyList) Reset() {
	if d != nil {
		*d = (*d)[:0]
	}
}

// Target real of the list
func (d *CompanyList) Target() any {
	return (*[]*models.Company)(d)
}

// Merge loaded data
func (d *CompanyList) Merge(l any) {
	newData := make([]*models.Company, 0, len(*d))
	for _, it := range *d {
		if it.Status.IsApproved() {
			newData = append(newData, it)
		}
	}
	for _, it := range l.([]*models.Company) {
		if !it.Status.IsApproved() {
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
