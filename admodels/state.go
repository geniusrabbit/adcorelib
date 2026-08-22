package admodels

import (
	"time"

	"github.com/geniusrabbit/adcorelib/billing"
)

type State interface {
	Profit() billing.Money
	TotalProfit() billing.Money
	Spend() billing.Money
	TotalSpend() billing.Money
	Imps() uint64
	Directs() uint64
	Views() uint64
	Clicks() uint64
	Leads() uint64
	LastSyncTime() time.Time
	ServerCount() int
}
