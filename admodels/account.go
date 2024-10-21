//
// @project GeniusRabbit corelib 2017, 2019, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017, 2019, 2024
//

package admodels

import (
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/billing"
)

// Account model
type Account struct {
	IDval    uint64        // Authoincrement key
	Balance  billing.Money //
	MaxDaily billing.Money //
	Spent    billing.Money // Daily spent

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

// RevenueShareFactor multipler 0..1 that publisher get from system
func (c *Account) RevenueShareFactor() float64 {
	if c == nil {
		zap.L().Error("account is not inited", zap.Stack("trace"))
		return 0.
	}
	return c.RevenueShare
}

// ComissionShareFactor which system get from publisher 0..1
func (c *Account) ComissionShareFactor() float64 {
	return 1. - c.RevenueShareFactor()
}
