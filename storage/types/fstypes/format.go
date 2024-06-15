package fstypes

import (
	"geniusrabbit.dev/adcorelib/models"
)

// FormatData represented in files
type FormatData struct {
	Formats []*models.Format `json:"formats"`
	Format  *models.Format   `json:"format"`
}

// Merge two data items
func (d *FormatData) Merge(it interface{}) {
	v := it.(*FormatData)
	if v.Format != nil {
		d.Formats = append(d.Formats, v.Format)
	}
	d.Formats = append(d.Formats, v.Formats...)
}

// Result of data as a list
func (d *FormatData) Result() []interface{} {
	data := make([]interface{}, 0, len(d.Formats))
	for _, it := range d.Formats {
		data = append(data, it)
	}
	return data
}

// Reset stored data
func (d *FormatData) Reset() {
	d.Format = nil
	d.Formats = d.Formats[:0]
}
