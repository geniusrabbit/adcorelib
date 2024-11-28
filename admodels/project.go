//
// @project GeniusRabbit corelib 2016 – 2017, 2021, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2017, 2021, 2024
//

package admodels

import "github.com/geniusrabbit/adcorelib/billing"

// Project model
type Project struct {
	ID     uint64
	UserID uint64

	Balance  billing.Money
	MaxDaily billing.Money
	Spent    billing.Money

	RevenueShare float64 // From 0 to 1 -> 100%
}

// CommissionShareFactor which system get from publisher 0..1
func (p *Project) CommissionShareFactor() float64 {
	return 1.0 - p.RevenueShare
}
