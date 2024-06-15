package fstypes

import "geniusrabbit.dev/adcorelib/models"

// RTBAccessPointData represented in files
type RTBAccessPointData struct {
	AccessPoints []*models.RTBAccessPoint `json:"rtb_access_points" yaml:"rtb_access_points"`
	AccessPoint  *models.RTBAccessPoint   `json:"rtb_access_point" yaml:"rtb_access_point"`
}

// Merge two data items
func (d *RTBAccessPointData) Merge(it interface{}) {
	v := it.(*RTBAccessPointData)
	if v.AccessPoint != nil {
		d.AccessPoints = append(d.AccessPoints, v.AccessPoint)
	}
	d.AccessPoints = append(d.AccessPoints, v.AccessPoints...)
}

// Result of data as a list
func (d *RTBAccessPointData) Result() []interface{} {
	data := make([]interface{}, 0, len(d.AccessPoints))
	for _, it := range d.AccessPoints {
		data = append(data, it)
	}
	return data
}

// Reset stored data
func (d *RTBAccessPointData) Reset() {
	d.AccessPoint = nil
	d.AccessPoints = d.AccessPoints[:0]
}
