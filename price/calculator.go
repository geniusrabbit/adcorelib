package price

import (
	"github.com/geniusrabbit/adcorelib/adtype"
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
	FixedPurchasePrice(action adtype.Action) billing.Money

	// BidPrice returns bid price for the external auction source.
	// The current bid price will be adjusted according to the source correction factor and the commission share factor
	BidPrice() billing.Money

	// Price returns price for the target action (view, click, lead, etc)
	Price(action adtype.Action) billing.Money
}

// CalculateNewBidPrice returns new bid price for the target with system comission and with source corrections
func CalculateNewBidPrice(price billing.Money, item priceCalculatorItem) billing.Money {
	return billing.MoneyFloat(
		price.Float64() /
			(1 + item.CommissionShareFactor()) /
			(1 + item.SourceCorrectionFactor()),
	)
}

// CalculatePurchasePrice returns purchase price for the target whithout system comission and with any corrections
//
// Formula:
//
//	PurchasePrice = Price - SourceCorrectionFactor[%] - TargetCorrectionFactor[%] - CommissionShareFactor[%]
func CalculatePurchasePrice(item priceCalculatorItem, action adtype.Action) billing.Money {
	if fixedPrice := item.FixedPurchasePrice(action); fixedPrice > 0 {
		return fixedPrice
	}
	if action == adtype.ActionView {
		if bidPrice := item.BidPrice(); bidPrice > 0 {
			return bidPrice
		}
	}
	price := item.Price(action)
	return billing.MoneyFloat(
		price.Float64() /
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
func CalculatePotentialPrice(item priceCalculatorItem, action adtype.Action) billing.Money {
	val := item.Price(action)
	if val == 0 && action.IsView() {
		val = item.ECPM()
	}
	return val
}

// CalculateFinalPrice returns final price for the item which is including all possible commissions with all corrections
//
// Formula:
//
//	FinalPrice = Price - SourceCorrectionFactor[%] - TargetCorrectionFactor[%]
func CalculateFinalPrice(item priceCalculatorItem, action adtype.Action) billing.Money {
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
