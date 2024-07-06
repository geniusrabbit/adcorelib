package dbtypes

import (
	"sort"

	"github.com/demdxx/gocast/v2"

	"github.com/geniusrabbit/adcorelib/models"
)

// CampaignList represented in DB
type CampaignList []*models.Campaign

// Result of data as a list
func (d CampaignList) Result() []any {
	return gocast.Slice[any](d)
}

// Reset stored data
func (d *CampaignList) Reset() {
	if d != nil {
		*d = (*d)[:0]
	}
}

// Target real of the list
func (d *CampaignList) Target() any {
	return (*[]*models.Campaign)(d)
}

// Merge loaded data
func (d *CampaignList) Merge(l any) {
	newData := make([]*models.Campaign, 0, len(*d))
	for _, it := range *d {
		if it.Status.IsApproved() && it.Active.IsActive() {
			newData = append(newData, it)
		}
	}
	for _, it := range l.([]*models.Campaign) {
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
