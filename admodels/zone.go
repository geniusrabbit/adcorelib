//
// @project GeniusRabbit rotator 2016 – 2019, 2021
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019, 2021
//

package admodels

import (
	"strconv"

	"github.com/geniusrabbit/gosql"

	"geniusrabbit.dev/corelib/admodels/types"
	"geniusrabbit.dev/corelib/billing"
	"geniusrabbit.dev/corelib/models"
)

// Zone model
type Zone struct {
	id                uint64
	StringID          string
	MinECPM           float64
	MinECPMByGeo      GeoBidSlice
	Price             billing.Money // CPM of source
	Comp              *Company
	CompID            uint64
	AllowedTypes      gosql.NullableOrderedIntArray
	AllowedSources    gosql.NullableOrderedIntArray
	DisallowedSources gosql.NullableOrderedIntArray
	DefaultCode       map[string]string
}

// ZoneFromModel convert database model to specified model
func ZoneFromModel(zone models.Zone) *Zone {
	var code map[string]string
	_ = zone.DefaultCode.UnmarshalTo(&code)

	return &Zone{
		id:                zone.ID,
		StringID:          strconv.FormatUint(zone.ID, 10),
		Price:             billing.MoneyFloat(zone.Price),
		Comp:              nil,
		CompID:            zone.CompanyID,
		MinECPM:           zone.MinECPM,
		MinECPMByGeo:      nil,
		AllowedTypes:      zone.AllowedTypes,
		AllowedSources:    zone.AllowedSources,
		DisallowedSources: zone.DisallowedSources,
		DefaultCode:       code,
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

// Company object
func (z *Zone) Company() *Company {
	return z.Comp
}

// CompanyID of current target
func (z *Zone) CompanyID() uint64 {
	return z.CompID
}

// RevenueShareFactor amount %
func (z *Zone) RevenueShareFactor() float64 {
	return z.Comp.RevenueShareFactor()
}

// ComissionShareFactor which system get from publisher
func (z *Zone) ComissionShareFactor() float64 {
	return z.Comp.ComissionShareFactor()
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
