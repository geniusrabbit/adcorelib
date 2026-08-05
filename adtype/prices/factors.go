package prices

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// Factors provides the commission and the discrepancy correction values resolved
// from the request context: the ad source, the target (zone/smartlink/access point)
// and the account. All values are percents in the range from 0 to 1.
type Factors interface {
	// CommissionShareFactor which the system gets from the publisher revenue
	CommissionShareFactor() float64

	// SourceCorrectionFactor of the advertisement source which is excluded as a
	// potential discrepancy
	SourceCorrectionFactor() float64

	// TargetCorrectionFactor of the target endpoint which is excluded as a
	// potential discrepancy
	TargetCorrectionFactor() float64
}

// FixedPurchasePricer is an optional extension of the [Factors] interface for the
// targets which define a fixed price of the traffic acceptance. If the value passed
// to [PublisherPrice] implements it and returns a positive price, that price
// overrides any calculation.
type FixedPurchasePricer interface {
	// FixedPurchasePrice returns the fixed price of the action or 0 if not defined
	FixedPurchasePrice(action adtype.Action) billing.Money
}

// StaticFactors is a plain value implementation of the [Factors] interface.
type StaticFactors struct {
	// Commission share which the system gets from the publisher revenue
	Commission float64 `json:"commission,omitempty"`

	// Source correction factor of the advertisement source
	Source float64 `json:"source,omitempty"`

	// Target correction factor of the target endpoint
	Target float64 `json:"target,omitempty"`
}

// CommissionShareFactor implementation of the [Factors] interface.
//
//go:inline
func (f StaticFactors) CommissionShareFactor() float64 { return f.Commission }

// SourceCorrectionFactor implementation of the [Factors] interface.
//
//go:inline
func (f StaticFactors) SourceCorrectionFactor() float64 { return f.Source }

// TargetCorrectionFactor implementation of the [Factors] interface.
//
//go:inline
func (f StaticFactors) TargetCorrectionFactor() float64 { return f.Target }

var _ Factors = StaticFactors{}

// commissionShareFactor of the factors or 0 if the factors are not defined.
func commissionShareFactor(f Factors) float64 {
	if f == nil {
		return 0
	}
	return f.CommissionShareFactor()
}

// sourceCorrectionFactor of the factors or 0 if the factors are not defined.
func sourceCorrectionFactor(f Factors) float64 {
	if f == nil {
		return 0
	}
	return f.SourceCorrectionFactor()
}

// targetCorrectionFactor of the factors or 0 if the factors are not defined.
func targetCorrectionFactor(f Factors) float64 {
	if f == nil {
		return 0
	}
	return f.TargetCorrectionFactor()
}

// reduce the value by the list of factors.
//
// Formula:
//
//	value * (1-factor[0]) * (1-factor[1]) * ...
func reduce(value billing.Money, factors ...float64) billing.Money {
	multiplier := 1.
	for _, factor := range factors {
		if factor != 0 {
			multiplier *= max(1.-factor, 0.)
		}
	}
	if multiplier == 1. {
		return value
	}
	return billing.MoneyFloat(value.Float64() * multiplier)
}

// grossUp is the inverse of [reduce]: it restores the pre-reduction value so that
// reduce(grossUp(v, f...), f...) == v (within money float precision).
//
// Formula:
//
//	value / (1-factor[0]) / (1-factor[1]) / ...
//
// If any factor is >= 1 (so that 1-factor <= 0), the result is 0 — the same
// overflow posture as [reduce].
func grossUp(value billing.Money, factors ...float64) billing.Money {
	divisor := 1.
	for _, factor := range factors {
		if factor == 0 {
			continue
		}
		part := 1. - factor
		if part <= 0 {
			return 0
		}
		divisor *= part
	}
	if divisor == 1. {
		return value
	}
	return billing.MoneyFloat(value.Float64() / divisor)
}
