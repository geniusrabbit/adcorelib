package adtype

import (
	"geniusrabbit.dev/adcorelib/billing"
)

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

// Calc remove or additiona comissions
func (f PriceFactor) Calc(price billing.Money, it ResponserItem, remove bool) (comissions billing.Money) {
	if remove {
		if f&TargetReducePriceFactor == TargetReducePriceFactor {
			pValue := PriceRevenueShareReduceFactors(price, it.Impression().Target, true)
			comissions += pValue
			price -= pValue
		}
		if f&SystemComissionPriceFactor == SystemComissionPriceFactor {
			pValue := PriceSystemComission(price, it, true)
			comissions += pValue
			price -= pValue
		}
		if f&SourcePriceFactor == SourcePriceFactor {
			pValue := PriceSourceFactors(price, it.Source(), true)
			comissions += pValue
			price -= pValue
		}
	} else {
		if f&SourcePriceFactor == SourcePriceFactor {
			pValue := PriceSourceFactors(price, it.Source(), false)
			comissions += pValue
			price += pValue
		}
		if f&SystemComissionPriceFactor == SystemComissionPriceFactor {
			pValue := PriceSystemComission(price, it, false)
			comissions += pValue
			price += pValue
		}
		if f&TargetReducePriceFactor == TargetReducePriceFactor {
			pValue := PriceRevenueShareReduceFactors(price, it.Impression().Target, false)
			comissions += pValue
			price += pValue
		}
	}
	return comissions
}

// PriceFactorList defines list of factors before calck
type PriceFactorList []PriceFactor

// Calc list of factors
func (l PriceFactorList) Calc(price billing.Money, it ResponserItem, remove bool) (comissions billing.Money) {
	if len(l) == 0 {
		return 0
	}
	var factors PriceFactor
	for _, f := range l {
		factors |= f
	}
	return factors.Calc(price, it, remove)
}

// PriceSourceFactors currection to reduce descreancy
func PriceSourceFactors(price billing.Money, src Source, remove bool) billing.Money {
	if src != nil {
		if reduce := src.PriceCorrectionReduceFactor(); reduce > 0 {
			if remove {
				return price/100*billing.Money((1-reduce)*100) - price
			}
			return price / 100 * billing.Money(reduce*100)
		}
	}
	return 0
}

// PriceSystemComission = 1. - `TrafficSourceComission`
func PriceSystemComission(price billing.Money, item systemComissionFactorer, remove bool) billing.Money {
	if item != nil {
		if reduce := item.ComissionShareFactor(); reduce > 0 {
			if remove {
				return price/100*billing.Money((1-reduce)*100) - price
			}
			return price / 100 * billing.Money(reduce*100)
		}
	}
	return 0
}

// PriceRevenueShareReduceFactors cirrection to reduce descreancy
func PriceRevenueShareReduceFactors(price billing.Money, rs revenueShareReducerFactorer, remove bool) billing.Money {
	if rs != nil {
		if reduce := rs.RevenueShareReduceFactor(); reduce > 0 {
			if remove {
				return price/100*billing.Money((1-reduce)*100) - price
			}
			return price / 100 * billing.Money(reduce*100)
		}
	}
	return 0
}
