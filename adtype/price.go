package adtype

import (
	"github.com/geniusrabbit/adcorelib/billing"
)

type systemComissionFactorer interface {
	// CommissionShareFactor which system get from publisher 0..1
	CommissionShareFactor() float64
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
	if f&SourcePriceFactor != 0 {
		pValue := PriceSourceFactors(price, it.Source(), false)
		comissions += pValue
		price += pValue
	}
	if f&SystemComissionPriceFactor != 0 {
		pValue := PriceSystemCommission(price, it, false)
		comissions += pValue
	}
	if f&TargetReducePriceFactor != 0 {
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
	if f&TargetReducePriceFactor != 0 {
		pValue := PriceRevenueShareReduceFactors(price, it.Impression().Target, true)
		comissions += pValue
		price += pValue
	}
	if f&SystemComissionPriceFactor != 0 {
		pValue := PriceSystemCommission(price, it, true)
		comissions += pValue
		price += pValue
	}
	if f&SourcePriceFactor != 0 {
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
	return AdjustPrice(price, src.PriceCorrectionReduceFactor(), remove)
}

// PriceSystemCommission = 1. - `TrafficSourceCommission`
func PriceSystemCommission(price billing.Money, item systemComissionFactorer, remove bool) billing.Money {
	if item == nil || price <= 0 {
		return 0
	}
	return AdjustPrice(price, item.CommissionShareFactor(), remove)
}

// PriceRevenueShareReduceFactors correction to reduce descreancy
func PriceRevenueShareReduceFactors(price billing.Money, rs revenueShareReducerFactorer, remove bool) billing.Money {
	if rs == nil || price <= 0 {
		return 0
	}
	return AdjustPrice(price, rs.RevenueShareReduceFactor(), remove)
}

func AdjustPrice(price billing.Money, factor float64, remove bool) billing.Money {
	if price <= 0 || factor <= 0 {
		return 0
	}

	fPrice := price.Float64()

	if remove {
		// Calculate the original price by dividing by (1 + factor)
		originalPrice := fPrice / (1 + factor)
		adjustment := fPrice - originalPrice
		return -billing.MoneyFloat(adjustment) // Round to the nearest integer
	}

	// Calculate adjustment for adding commission
	adjustment := fPrice * factor
	return billing.MoneyFloat(adjustment) // Round to the nearest integer
}
