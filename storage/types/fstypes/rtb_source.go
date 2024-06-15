package fstypes

import "geniusrabbit.dev/adcorelib/models"

// RTBSourceData represented in files
type RTBSourceData struct {
	Sources []*models.RTBSource `json:"rtb_sources" yaml:"rtb_sources"`
	Source  *models.RTBSource   `json:"rtb_source" yaml:"rtb_source"`
}

// Merge two data items
func (d *RTBSourceData) Merge(it interface{}) {
	v := it.(*RTBSourceData)
	if v.Source != nil {
		d.Sources = append(d.Sources, v.Source)
	}
	d.Sources = append(d.Sources, v.Sources...)
}

// Result of data as a list
func (d *RTBSourceData) Result() []interface{} {
	data := make([]interface{}, 0, len(d.Sources))
	for _, it := range d.Sources {
		data = append(data, it)
	}
	return data
}

// Reset stored data
func (d *RTBSourceData) Reset() {
	d.Source = nil
	d.Sources = d.Sources[:0]
}
