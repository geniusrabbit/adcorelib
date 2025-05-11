package adtype

import "github.com/geniusrabbit/adcorelib/billing"

// Account type represents the account model
type Account interface {
	ID() uint64

	// CommissionShareFactor which system get from publisher 0..1
	CommissionShareFactor() float64
}

// AccountWithBudget type represents the account model
// It has additional methods for testing the budget
type AccountWithBudget interface {
	Account

	// TestBudget of the account retuns true if the account has sufficient budget
	// and greater than the specified value
	TestBudget(vol billing.Money) bool
}
