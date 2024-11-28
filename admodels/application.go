//
// @project GeniusRabbit corelib 2016, 2022
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016, 2022
//

package admodels

import (
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/models"
)

// Application model
type Application struct {
	ID        uint64   // Authoincrement key
	Account   *Account // Who have this company
	AccountID uint64   //

	URI      string                // Unical application identificator
	Type     types.ApplicationType `gorm:"type:ApplicationType" json:"type"`
	Platform types.PlatformType    `gorm:"type:PlatformType" json:"platform"`
	Premium  bool                  `json:"premium"`

	Categories   []uint
	RevenueShare float64 // % 1.0 -> 100%, 0.655 -> 65.5%
}

// ApplicationFromModel convert database model to specified model
func ApplicationFromModel(app *models.Application) Application {
	return Application{
		ID:           app.ID,
		AccountID:    app.AccountID,
		URI:          app.URI,
		Type:         app.Type,
		Platform:     app.Platform,
		Premium:      app.Premium,
		Categories:   app.Categories,
		RevenueShare: app.RevenueShare,
	}
}

// ObjectKey of the application
func (a *Application) ObjectKey() string {
	return a.URI
}

// RevenueShareFactor amount %
func (a *Application) RevenueShareFactor() float64 {
	if a.RevenueShare > 0 {
		return a.RevenueShare
	}
	return a.Account.RevenueShareFactor()
}

// CommissionShareFactor which system get from publisher
func (a *Application) CommissionShareFactor() float64 {
	if a.RevenueShare > 0 {
		return (1.0 - a.RevenueShare)
	}
	return a.Account.CommissionShareFactor()
}
