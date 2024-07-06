//
// @project GeniusRabbit corelib 2017, 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019
//

package admodels

import (
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/models"
	"go.uber.org/zap"
)

// Company model
type Company struct {
	ID       uint64        // Authoincrement key
	Balance  billing.Money //
	MaxDaily billing.Money //
	Spent    billing.Money // Daily spent

	// RevenueShare it's amount of percent of the raw incode which will be shared with the publisher company
	// For example:
	//   Displayed ads for 100$
	//   Company revenue share 60%
	//   In such case the ad network have 40$
	//   The publisher have 60$
	RevenueShare float64 // % 100_00, 10000 -> 100%, 6550 -> 65.5%
}

// CompanyFromModel convert database model to specified model
func CompanyFromModel(c *models.Company) *Company {
	return &Company{
		ID:           c.ID,
		Balance:      0,
		MaxDaily:     billing.MoneyFloat(c.MaxDaily),
		Spent:        0,
		RevenueShare: c.RevenueShare,
	}
}

// RevenueShareFactor multipler 0..1
func (c *Company) RevenueShareFactor() float64 {
	if c == nil {
		zap.L().Error("company is not inited")
		return 0
	}
	return c.RevenueShare / 100.
}

// ComissionShareFactor which system get from publisher 0..1
func (c *Company) ComissionShareFactor() float64 {
	if c == nil {
		zap.L().Error("company is not inited")
		return 1
	}
	return (100 - c.RevenueShare) / 100.
}

///////////////////////////////////////////////////////////////////////////////
/// Target wrapper
///////////////////////////////////////////////////////////////////////////////

// CompanyTarget wrapper for replac of epsent target object
type CompanyTarget struct {
	Comp *Company
}

// ID of object (Zone OR SmartLink only)
func (c CompanyTarget) ID() uint64 {
	return 0
}

// Size default of target item
func (c CompanyTarget) Size() (w, h int) {
	return w, h
}

// RevenueShareFactor of current target
func (c CompanyTarget) RevenueShareFactor() float64 {
	return c.Comp.RevenueShareFactor()
}

// ComissionShareFactor of current target
func (c CompanyTarget) ComissionShareFactor() float64 {
	return c.Comp.ComissionShareFactor()
}

// Company object
func (c CompanyTarget) Company() *Company {
	return c.Comp
}

// ProjectID number
func (c CompanyTarget) ProjectID() uint64 {
	return 0
}

// CompanyID of current target
func (c CompanyTarget) CompanyID() uint64 {
	return c.Comp.ID
}
