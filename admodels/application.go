//
// @project GeniusRabbit corelib 2016, 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016, 2022
//

package admodels

import "github.com/geniusrabbit/adcorelib/models"

// Application model
type Application struct {
	ID           uint64   // Authoincrement key
	Account      *Account // Who have this company
	AccountID    uint64   //
	Opt          [8]uint8 // Platform, Premium, Type
	Categories   []uint   //
	RevenueShare float64  // From 0 to 100 percents
}

// ApplicationFromModel convert database model to specified model
func ApplicationFromModel(app *models.Application) Application {
	return Application{
		ID:           app.ID,
		AccountID:    app.AccountID,
		Categories:   app.Categories,
		RevenueShare: app.RevenueShare,
	}
}

// RevenueShareFactor amount %
func (a *Application) RevenueShareFactor() float64 {
	if a.RevenueShare > 0 {
		return a.RevenueShare / 100.0
	}
	return a.Account.RevenueShareFactor()
}

// ComissionShareFactor which system get from publisher
func (a *Application) ComissionShareFactor() float64 {
	if a.RevenueShare > 0 {
		return (100.0 - a.RevenueShare) / 100.0
	}
	return a.Account.ComissionShareFactor()
}
