package price

import (
	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/billing"
)

type priceCalculatorItem interface {
	// CommissionShareFactor which is the system get from publisher 0..1 for reducing of the value
	CommissionShareFactor() float64

	// PriceCorrectionReduceFactor which is a potential 0..1 for reducing of the value
	SourceCorrectionFactor() float64

	// TargetCorrectionFactor which is a potential 0..1 for reducing of the value
	TargetCorrectionFactor() float64

	// ECPM returns ECPM value for the item
	ECPM() billing.Money

	// FixedPurchasePrice returns fixed price for the target action
	FixedPurchasePrice(action admodels.Action) billing.Money

	// Price returns price for the target action (view, click, lead, etc)
	Price(action admodels.Action) billing.Money
}

// CalculatePurchasePrice returns purchase price for the target whithout system comission and with any corrections
//
// Formula:
//
//	PurchasePrice = Price - SourceCorrectionFactor[%] - TargetCorrectionFactor[%] - CommissionShareFactor[%]
func CalculatePurchasePrice(item priceCalculatorItem, action admodels.Action) billing.Money {
	if fixedPrice := item.FixedPurchasePrice(action); fixedPrice > 0 {
		return fixedPrice
	}
	price := item.Price(action).Float64()
	return billing.MoneyFloat(
		price /
			(1. + item.SourceCorrectionFactor()) /
			(1. + item.TargetCorrectionFactor()) /
			(1. + item.CommissionShareFactor()),
	)
}

// CalculatePotentialPrice returns the base price without any corrections or commissions
//
// Formula:
//
//	PotentialPrice = Price
//
//go:inline
func CalculatePotentialPrice(item priceCalculatorItem, action admodels.Action) billing.Money {
	return item.Price(action)
}

// CalculateFinalPrice returns final price for the item which is including all possible commissions with all corrections
//
// Formula:
//
//	FinalPrice = Price - SourceCorrectionFactor[%] - TargetCorrectionFactor[%]
func CalculateFinalPrice(item priceCalculatorItem, action admodels.Action) billing.Money {
	sourceCorrection := item.SourceCorrectionFactor()
	targetCorrection := item.TargetCorrectionFactor()
	price := item.Price(action)

	if sourceCorrection == 0 && targetCorrection == 0 {
		return price
	}

	return billing.MoneyFloat(
		price.Float64() /
			(1. + sourceCorrection) /
			(1. + targetCorrection),
	)
}

// CalculateInternalAuctionBid returns price for the internal auction. Normaly it's the ECPM value, with correction (in case of external source) and without any system comission
//
//go:inline
func CalculateInternalAuctionBid(item priceCalculatorItem) billing.Money {
	return item.ECPM()
}
