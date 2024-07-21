//
// @project GeniusRabbit corelib 2016 – 2019, 2021, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019, 2021, 2024
//

package admodels

import (
	"strconv"

	"github.com/geniusrabbit/gosql/v2"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/models"
)

// Zone model
type Zone struct {
	id       uint64
	StringID string

	Acc   *Account
	AccID uint64

	Price        billing.Money // CPM of source
	MinECPM      float64
	MinECPMByGeo GeoBidSlice

	AllowedTypes      gosql.NullableOrderedNumberArray[int]
	AllowedSources    gosql.NullableOrderedNumberArray[int]
	DisallowedSources gosql.NullableOrderedNumberArray[int]
	DefaultCode       map[string]string
}

// ZoneFromModel convert database model to specified model
func ZoneFromModel(zone *models.Zone, account *Account) *Zone {
	return &Zone{
		id:                zone.ID,
		StringID:          strconv.FormatUint(zone.ID, 10),
		Price:             billing.MoneyFloat(zone.Price),
		Acc:               account,
		AccID:             zone.AccountID,
		MinECPM:           zone.MinECPM,
		MinECPMByGeo:      nil,
		AllowedTypes:      zone.AllowedTypes,
		AllowedSources:    zone.AllowedSources,
		DisallowedSources: zone.DisallowedSources,
		DefaultCode:       zone.DefaultCode.DataOr(nil),
	}
}

// ID of object
func (z *Zone) ID() uint64 {
	return z.id
}

// Codename of the target (equal to tagid)
func (z *Zone) Codename() string {
	return z.StringID
}

// PricingModel of the target
func (z *Zone) PricingModel() types.PricingModel {
	return types.PricingModelUndefined
}

// AlternativeAdCode returns URL or any code (HTML, XML, etc)
func (z *Zone) AlternativeAdCode(key string) string {
	if z.DefaultCode == nil {
		return ""
	}
	return z.DefaultCode[key]
}

// PurchasePrice gives the price of view from external resource
func (z *Zone) PurchasePrice(action Action) billing.Money {
	if action.IsImpression() {
		return z.Price
	}
	return 0
}

// Account object
func (z *Zone) Account() *Account {
	return z.Acc
}

// AccountID of current target
func (z *Zone) AccountID() uint64 {
	return z.AccID
}

// SetAccount for target
func (z *Zone) SetAccount(acc *Account) {
	z.Acc = acc
}

// RevenueShareFactor amount %
func (z *Zone) RevenueShareFactor() float64 {
	return z.Acc.RevenueShareFactor()
}

// ComissionShareFactor which system get from publisher
func (z *Zone) ComissionShareFactor() float64 {
	return z.Acc.ComissionShareFactor()
}

// RevenueShareReduceFactor correction factor to reduce target proce of the access point to avoid descrepancy
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (z *Zone) RevenueShareReduceFactor() float64 { return 0 }

// IsAllowedSource for targeting
func (z *Zone) IsAllowedSource(id uint64, types []int) bool {
	if len(z.AllowedSources) > 0 {
		return z.AllowedSources.IndexOf(int(id)) >= 0
	}
	if len(z.DisallowedSources) > 0 {
		return z.DisallowedSources.IndexOf(int(id)) < 0
	}
	return true
}
