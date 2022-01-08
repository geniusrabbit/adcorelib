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
type PriceFactor int

const (
	AllPriceFactors              PriceFactor = 0xffffffff
	SourcePriceFactor            PriceFactor = 0x0001
	SystemComissionPriceFactor   PriceFactor = 0x0002
	TargetSgareReducePriceFactor PriceFactor = 0x0004
)

// Calck new price
func (f PriceFactor) Calc(price billing.Money, it ResponserItem) billing.Money {
	var newPrice billing.Money
	if f&SourcePriceFactor == SourcePriceFactor {
		newPrice = PriceSourceFactors(price, it.Source())
	}
	if f&SystemComissionPriceFactor == SystemComissionPriceFactor {
		newPrice += PriceSystemComission(price, it)
	}
	if f&TargetSgareReducePriceFactor == TargetSgareReducePriceFactor {
		newPrice += PriceRevenueShareReduceFactors(price, it.Impression().Target)
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
