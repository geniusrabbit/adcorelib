//
// @project GeniusRabbit corelib 2017 - 2019, 2021
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019, 2021
//

package admodels

import (
	"strconv"

	"github.com/geniusrabbit/gosql/v2"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/models"
)

// Smartlink model
type Smartlink struct {
	id       uint64
	StringID string

	Acc   *Account
	AccID uint64

	FixedPurchasePrice billing.Money // The cost of single view

	AllowedTypes      gosql.NullableOrderedNumberArray[int]
	AllowedSources    gosql.NullableOrderedNumberArray[int]
	DisallowedSources gosql.NullableOrderedNumberArray[int]
	Campaigns         gosql.NullableOrderedNumberArray[int]
	DefaultCode       map[string]string
}

// SmartlinkFromModel convert database model to specified model
func SmartlinkFromModel(zone *models.Zone, account *Account) *Smartlink {
	return &Smartlink{
		id:                 zone.ID,
		StringID:           strconv.FormatUint(zone.ID, 10),
		FixedPurchasePrice: zone.FixedPurchasePrice,
		Acc:                account,
		AccID:              zone.AccountID,
		AllowedTypes:       zone.AllowedTypes,
		AllowedSources:     zone.AllowedSources,
		DisallowedSources:  zone.DisallowedSources,
		Campaigns:          zone.Campaigns,
		DefaultCode:        *zone.DefaultCode.Data,
	}
}

// ID of object
func (l *Smartlink) ID() uint64 {
	return l.id
}

// Codename of the target (equal to tagid)
func (l *Smartlink) Codename() string {
	return l.StringID
}

// PricingModel of the target
func (l *Smartlink) PricingModel() types.PricingModel {
	return types.PricingModelUndefined
}

// AlternativeAdCode returns URL or any code (HTML, XML, etc)
func (l *Smartlink) AlternativeAdCode(key string) string {
	if l.DefaultCode == nil {
		return ""
	}
	return l.DefaultCode[key]
}

// PurchasePrice gives the price of view from external resource
func (l *Smartlink) PurchasePrice(action Action) billing.Money {
	if action.IsImpression() {
		return l.FixedPurchasePrice
	}
	return 0
}

// Account object
func (l *Smartlink) Account() *Account {
	return l.Acc
}

// AccountID of current target
func (l *Smartlink) AccountID() uint64 {
	return l.AccID
}

// SetAccount for target
func (l *Smartlink) SetAccount(acc *Account) {
	l.Acc = acc
}

// RevenueShareFactor amount of % which company get from publisher
func (l *Smartlink) RevenueShareFactor() float64 {
	return l.Acc.RevenueShareFactor()
}

// ComissionShareFactor which system get from publisher
func (l *Smartlink) ComissionShareFactor() float64 {
	return l.Acc.ComissionShareFactor()
}

// RevenueShareReduceFactor correction factor to reduce target proce of the access point to avoid descrepancy
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (l *Smartlink) RevenueShareReduceFactor() float64 { return 0 }

// IsAllowedSource for targeting
func (l *Smartlink) IsAllowedSource(id uint64, types []int) bool {
	if len(l.AllowedSources) > 0 {
		return l.AllowedSources.IndexOf(int(id)) >= 0
	}
	if len(l.DisallowedSources) > 0 {
		return l.DisallowedSources.IndexOf(int(id)) < 0
	}
	return true
}
