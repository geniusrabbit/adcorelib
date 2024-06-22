package fstypes

import "geniusrabbit.dev/adcorelib/models"

// CompanyData represented in files
type CompanyData struct {
	Companies []*models.Company `json:"companies" yaml:"companies"`
	Company   *models.Company   `json:"company" yaml:"company"`
}

// Merge two data items
func (d *CompanyData) Merge(it any) {
	v := it.(*CompanyData)
	if v.Company != nil {
		d.Companies = append(d.Companies, v.Company)
	}
	d.Companies = append(d.Companies, v.Companies...)
}

// Result of data as a list
func (d *CompanyData) Result() []any {
	data := make([]any, 0, len(d.Companies))
	for _, it := range d.Companies {
		data = append(data, it)
	}
	return data
}

// Reset stored data
func (d *CompanyData) Reset() {
	d.Company = nil
	d.Companies = d.Companies[:0]
}
