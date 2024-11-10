package adtype

import (
	"github.com/geniusrabbit/adcorelib/billing"
)

type systemComissionFactorer interface {
	// ComissionShareFactor which system get from publisher 0..1
	ComissionShareFactor() float64
}

type revenueShareReducerFactorer interface {
	// RevenueShareReduceFactor the percentage of revenue that the owner of the Zone/SmartLink/AdAccessPoint will be excluded as potential descrepancy.
	RevenueShareReduceFactor() float64
}

// PriceFactor defines action to calculate the factor
type PriceFactor uint32

const (
	NonePriceFactor            PriceFactor = 0x0000
	AllPriceFactors            PriceFactor = 0xffffffff
	SourcePriceFactor          PriceFactor = 0x0001
	SystemComissionPriceFactor PriceFactor = 0x0002
	TargetReducePriceFactor    PriceFactor = 0x0004
)

// PriceFactorFromList create factor from list
func PriceFactorFromList(factors ...PriceFactor) (f PriceFactor) {
	for _, factor := range factors {
		f |= factor
	}
	return f
}

// AddComission to price and rеturns comissions with positive sign `+`
func (f PriceFactor) AddComission(price billing.Money, it ResponserItem) (comissions billing.Money) {
	if f == NonePriceFactor || price <= 0 {
		return 0
	}
	if f&SourcePriceFactor == SourcePriceFactor {
		pValue := PriceSourceFactors(price, it.Source(), false)
		comissions += pValue
		price += pValue
	}
	if f&SystemComissionPriceFactor == SystemComissionPriceFactor {
		pValue := PriceSystemComission(price, it, false)
		comissions += pValue
	}
	if f&TargetReducePriceFactor == TargetReducePriceFactor {
		pValue := PriceRevenueShareReduceFactors(price, it.Impression().Target, false)
		comissions += pValue
	}
	return comissions
}

// RemoveComission from price and rеturns comissions with negative sign `-`
func (f PriceFactor) RemoveComission(price billing.Money, it ResponserItem) (comissions billing.Money) {
	if f == NonePriceFactor || price <= 0 {
		return 0
	}
	if f&TargetReducePriceFactor == TargetReducePriceFactor {
		pValue := PriceRevenueShareReduceFactors(price, it.Impression().Target, true)
		comissions += pValue
		price += pValue
	}
	if f&SystemComissionPriceFactor == SystemComissionPriceFactor {
		pValue := PriceSystemComission(price, it, true)
		comissions += pValue
		price += pValue
	}
	if f&SourcePriceFactor == SourcePriceFactor {
		pValue := PriceSourceFactors(price, it.Source(), true)
		comissions += pValue
	}
	return comissions
}

// PriceSourceFactors currection to reduce descreancy
func PriceSourceFactors(price billing.Money, src Source, remove bool) billing.Money {
	if src == nil || price <= 0 {
		return 0
	}
	if reduce := src.PriceCorrectionReduceFactor(); reduce > 0 {
		if remove {
			return price/100*billing.Money((1-reduce)*100) - price
		}
		return price / 100 * billing.Money(reduce*100)
	}
	return 0
}

// PriceSystemComission = 1. - `TrafficSourceComission`
func PriceSystemComission(price billing.Money, item systemComissionFactorer, remove bool) billing.Money {
	if item == nil || price <= 0 {
		return 0
	}
	if reduce := item.ComissionShareFactor(); reduce > 0 {
		if remove {
			return price/100*billing.Money((1-reduce)*100) - price
		}
		return price / 100 * billing.Money(reduce*100)
	}
	return 0
}

// PriceRevenueShareReduceFactors cirrection to reduce descreancy
func PriceRevenueShareReduceFactors(price billing.Money, rs revenueShareReducerFactorer, remove bool) billing.Money {
	if rs == nil || price <= 0 {
		return 0
	}
	if reduce := rs.RevenueShareReduceFactor(); reduce > 0 {
		if remove {
			return price/100*billing.Money((1-reduce)*100) - price
		}
		return price / 100 * billing.Money(reduce*100)
	}
	return 0
}
