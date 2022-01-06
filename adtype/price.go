package adtype

import "geniusrabbit.dev/corelib/billing"

type systemComissionFactorer interface {
	// ComissionShareFactor which system get from publisher 0..1
	ComissionShareFactor() float64
}

type revenueSgareReducerFactorer interface {
	RevenueShareReduceFactor() float64
}

// PriceSourceFactors currection to reduce descreancy
func PriceSourceFactors(price billing.Money, src Source) billing.Money {
	if src != nil {
		if reduce := src.PriceCorrectionReduceFactor(); reduce > 0 {
			return price / 100 * billing.MoneyFloat(reduce)
		}
	}
	return 0
}

// PriceSystemComission = 1. - `TrafficSourceComission`
func PriceSystemComission(price billing.Money, item systemComissionFactorer) billing.Money {
	if item != nil {
		if reduce := item.ComissionShareFactor(); reduce > 0 {
			return price / 100 * billing.MoneyFloat(reduce)
		}
	}
	return 0
}

// PriceRevenueShareReduceFactors cirrection to reduce descreancy
func PriceRevenueShareReduceFactors(price billing.Money, rs revenueSgareReducerFactorer) billing.Money {
	if rs != nil {
		if reduce := rs.RevenueShareReduceFactor(); reduce > 0 {
			return price / 100 * billing.MoneyFloat(reduce)
		}
	}
	return 0
}
