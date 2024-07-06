package dbtypes

import (
	"sort"

	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/adcorelib/models"
)

// AdAssetList represented in DB
type AdAssetList []*models.AdAsset

// Result of data as a list
func (d AdAssetList) Result() []any {
	return gocast.Slice[any](d)
}

// Reset stored data
func (d *AdAssetList) Reset() {
	if d != nil {
		*d = (*d)[:0]
	}
}

// Target real of the list
func (d *AdAssetList) Target() any {
	return (*[]*models.AdAsset)(d)
}

// Merge loaded data
func (d *AdAssetList) Merge(l any) {
	newData := make([]*models.AdAsset, 0, len(*d)+len(l.([]*models.AdAsset)))
	newData = append(newData, *d...)
	for _, it := range l.([]*models.AdAsset) {
		if i := d.IndexOf(it.ID); i >= 0 {
			newData[i] = it
		} else {
			newData = append(newData, it)
		}
	}
	sort.Slice(newData, func(i, j int) bool { return newData[i].ID < newData[j].ID })
	*d = newData
}

func (d AdAssetList) IndexOf(id uint64) int {
	i := sort.Search(len(d), func(i int) bool { return d[i].ID >= id })
	if i >= 0 && i < len(d) && d[i].ID == id {
		return i
	}
	return -1
}
