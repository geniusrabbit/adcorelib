package adtype

import "geniusrabbit.dev/corelib/billing"

type systemComissionFactorer interface {
	// ComissionShareFactor which system get from publisher 0..1
	ComissionShareFactor() float64
}

type revenueShareReducerFactorer interface {
	RevenueShareReduceFactor() float64
}

// PriceFactor defines action to calculate the factor
type PriceFactor uint32

const (
	AllPriceFactors            PriceFactor = 0xffffffff
	SourcePriceFactor          PriceFactor = 0x0001
	SystemComissionPriceFactor PriceFactor = 0x0002
	TargetReducePriceFactor    PriceFactor = 0x0004
)

// Calck new price
func (f PriceFactor) Calc(price billing.Money, it ResponserItem, remove bool) billing.Money {
	var newPrice billing.Money
	if remove {
		if f&TargetReducePriceFactor == TargetReducePriceFactor {
			pValue := PriceRevenueShareReduceFactors(price, it.Impression().Target)
			newPrice += pValue
			price -= pValue
		}
		if f&SystemComissionPriceFactor == SystemComissionPriceFactor {
			pValue := PriceSystemComission(price, it)
			newPrice += pValue
			price -= pValue
		}
		if f&SourcePriceFactor == SourcePriceFactor {
			pValue := PriceSourceFactors(price, it.Source())
			newPrice += pValue
			price -= pValue
		}
	} else {
		if f&SourcePriceFactor == SourcePriceFactor {
			pValue := PriceSourceFactors(price, it.Source())
			newPrice += pValue
			price += pValue
		}
		if f&SystemComissionPriceFactor == SystemComissionPriceFactor {
			pValue := PriceSystemComission(price, it)
			newPrice += pValue
			price += pValue
		}
		if f&TargetReducePriceFactor == TargetReducePriceFactor {
			pValue := PriceRevenueShareReduceFactors(price, it.Impression().Target)
			newPrice += pValue
			price += pValue
		}
	}
	return newPrice
}

// PriceSourceFactors currection to reduce descreancy
func PriceSourceFactors(price billing.Money, src Source) billing.Money {
	if src != nil {
		if reduce := src.PriceCorrectionReduceFactor(); reduce > 0 {
			return price / 100 * billing.Money(reduce*100)
		}
	}
	return 0
}

// PriceSystemComission = 1. - `TrafficSourceComission`
func PriceSystemComission(price billing.Money, item systemComissionFactorer) billing.Money {
	if item != nil {
		if reduce := item.ComissionShareFactor(); reduce > 0 {
			return price / 100 * billing.Money(reduce*100)
		}
	}
	return 0
}

// PriceRevenueShareReduceFactors cirrection to reduce descreancy
func PriceRevenueShareReduceFactors(price billing.Money, rs revenueShareReducerFactorer) billing.Money {
	if rs != nil {
		if reduce := rs.RevenueShareReduceFactor(); reduce > 0 {
			return price / 100 * billing.Money(reduce*100)
		}
	}
	return 0
}
