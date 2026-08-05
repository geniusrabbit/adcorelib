package prices

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// PriceProvider is the minimal contract required by the price calculators.
// It is implemented by [PriceScope] and can be implemented by any other type which
// stores the price of the action.
type PriceProvider interface {
	// PricePerAction returns the current price of one action
	PricePerAction(action adtype.Action) billing.Money

	// MaxPricePerAction returns the maximal possible price of one action
	MaxPricePerAction(action adtype.Action) billing.Money
}

// PotentialPrice returns the maximal price which the advertiser could have paid for
// the action. The difference with the [AdvertiserPrice] is the discrepancy of the
// current placement.
//
// Formula:
//
//	PotentialPrice = MaxPrice
//
//go:inline
func PotentialPrice(scope PriceProvider, action adtype.Action) billing.Money {
	return scope.MaxPricePerAction(action)
}

// AdvertiserPrice returns the price which will be charged from the advertiser for
// the action. The advertiser covers the full cost of running the campaign across
// every connected network, discrepancies included, so none of the source/target
// discrepancy correction factors (nor the commission) are subtracted here. The
// factors argument is accepted only for a uniform signature with the other
// calculators and is intentionally ignored.
//
// Formula:
//
//	AdvertiserPrice = Price
//
//go:inline
func AdvertiserPrice(scope PriceProvider, _ Factors, action adtype.Action) billing.Money {
	return scope.PricePerAction(action)
}

// PublisherPrice returns the price which the system has to pay to the publisher for
// the action: the price reduced by the source and the target discrepancy
// correction factors and by the system commission share. Whatever these factors
// remove stays with the network, see [NetworkProfit]. If the factors value
// implements the [FixedPurchasePricer] interface and defines a positive price of
// the action, that price is used as is.
//
// Formula:
//
//	PublisherPrice = Price - SourceCorrectionFactor[%] - TargetCorrectionFactor[%] - CommissionShareFactor[%]
func PublisherPrice(scope PriceProvider, factors Factors, action adtype.Action) billing.Money {
	if pricer, _ := factors.(FixedPurchasePricer); pricer != nil {
		if fixedPrice := pricer.FixedPurchasePrice(action); fixedPrice > 0 {
			return fixedPrice
		}
	}
	return reduce(scope.PricePerAction(action),
		sourceCorrectionFactor(factors), targetCorrectionFactor(factors), commissionShareFactor(factors))
}

// NetworkProfit returns the profit of the network for the action: the difference
// between what the advertiser was charged and what the publisher was paid. It is
// the system commission share plus whatever the source and the target discrepancy
// corrections have deducted from the publisher payout, since none of that is ever
// subtracted from the advertiser charge.
//
// Formula:
//
//	NetworkProfit = AdvertiserPrice - PublisherPrice
//
//go:inline
func NetworkProfit(scope PriceProvider, factors Factors, action adtype.Action) billing.Money {
	return AdvertiserPrice(scope, factors, action) - PublisherPrice(scope, factors, action)
}

// BidUpPrice returns the advertiser-side bid that must be stored so that
// [PublisherPrice] recovers the given net price after source/target discrepancy
// corrections and the system commission are deducted. It is the true inverse of
// the reduce path used by [PublisherPrice] (not a (1+factor) markup).
//
// Formula:
//
//	BidUpPrice = Price / (1-SourceCorrectionFactor) / (1-TargetCorrectionFactor) / (1-CommissionShareFactor)
//
//go:inline
func BidUpPrice(price billing.Money, factors Factors) billing.Money {
	return grossUp(price, commissionShareFactor(factors),
		targetCorrectionFactor(factors), sourceCorrectionFactor(factors))
}
