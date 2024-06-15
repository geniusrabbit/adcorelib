//
// @project GeniusRabbit rotator 2016, 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016, 2022
//

package admodels

import "geniusrabbit.dev/adcorelib/models"

// Application model
type Application struct {
	ID           uint64   // Authoincrement key
	Company      *Company // Who have this company
	CompanyID    uint64   //
	Opt          [8]uint8 // Platform, Premium, Type
	Categories   []uint   //
	RevenueShare float64  // From 0 to 100 percents
}

// ApplicationFromModel convert database model to specified model
func ApplicationFromModel(app *models.Application) Application {
	return Application{
		ID:           app.ID,
		CompanyID:    app.CompanyID,
		Categories:   app.Categories,
		RevenueShare: app.RevenueShare,
	}
}

// RevenueShareFactor amount %
func (a *Application) RevenueShareFactor() float64 {
	if a.RevenueShare > 0 {
		return a.RevenueShare / 100.0
	}
	return a.Company.RevenueShareFactor()
}

// ComissionShareFactor which system get from publisher
func (a *Application) ComissionShareFactor() float64 {
	if a.RevenueShare > 0 {
		return (100.0 - a.RevenueShare) / 100.0
	}
	return a.Company.ComissionShareFactor()
}
