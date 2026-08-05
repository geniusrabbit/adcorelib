package prices

import (
	"errors"

	"github.com/geniusrabbit/adcorelib/billing"
)

// ErrNegativeBidPrice returns in case of the attempt to set a negative bid price.
var ErrNegativeBidPrice = errors.New("bid price must not be negative")

// ErrUnsupportedAction returns in case of the attempt to set a bid price for the
// action which is not supported by the price scope.
var ErrUnsupportedAction = errors.New("action is not supported by the price scope")

// cpmDelimiter is the amount of actions covered by one CPM value.
const cpmDelimiter = 1000

// PriceFromCPM converts the CPM value into the price of one single action.
//
//go:inline
func PriceFromCPM(cpm billing.Money) billing.Money { return cpm / cpmDelimiter }

// CPMFromPrice converts the price of one single action into the CPM value.
//
//go:inline
func CPMFromPrice(price billing.Money) billing.Money { return price * cpmDelimiter }

// prepareBid limits the value by the maximal bid.
// The zero maximal bid means that there is no upper limit.
func prepareBid(value, maxValue billing.Money) billing.Money {
	if value <= 0 {
		return 0
	}
	if maxValue > 0 && value > maxValue {
		return maxValue
	}
	return value
}

// setBid validates the new bid value and stores it clamped into the [0, maxValue] range.
func setBid(target *billing.Money, value, maxValue billing.Money) error {
	if value < 0 {
		return ErrNegativeBidPrice
	}
	*target = prepareBid(value, maxValue)
	return nil
}
