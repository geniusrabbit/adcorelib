//
// @project GeniusRabbit corelib 2017, 2019, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019, 2024
//

package admodels

import (
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/billing"
)

type AccountBalanceState interface {
	Balance() billing.Money
	Spend() billing.Money
}

// Account model
type Account struct {
	IDval uint64 // Authoincrement key

	MaxDaily     billing.Money
	CurrentState AccountBalanceState

	// RevenueShare it's amount of percent of the raw incode which will be shared with the publisher company
	// For example:
	//   Displayed ads for 100$
	//   Account revenue share 60%
	//   In such case the ad network have 40$
	//   The publisher have 60$
	RevenueShare float64 // % 1.0 -> 100%, 0.655 -> 65.5%
}

// ID of object
func (c *Account) ID() uint64 {
	return c.IDval
}

// ObjectKey of the target
func (c *Account) ObjectKey() uint64 {
	return c.IDval
}

// DailyBudget of the account
func (c *Account) DailyBudget() billing.Money {
	return c.MaxDaily
}

// Balance of the account
func (c *Account) Balance() billing.Money {
	if c.CurrentState == nil {
		return 0
	}
	return c.CurrentState.Balance()
}

// Spend of the account
func (c *Account) Spend() billing.Money {
	if c.CurrentState == nil {
		return 0
	}
	return c.CurrentState.Spend()
}

// CommissionShareFactor which system get from publisher 0..1
func (c *Account) CommissionShareFactor() float64 {
	if c == nil {
		zap.L().Error("account is not inited", zap.Stack("trace"))
		return 0.
	}
	return 1. - c.RevenueShare
}
