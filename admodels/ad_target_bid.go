package admodels

import (
	"math/rand"

	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/adcorelib/billing"
)

// TargetBid object
type TargetBid struct {
	TestMode bool

	Ad *Ad
	// Bid *AdBid

	ECPM          billing.Money
	BidPrice      billing.Money // Max price per one View (used in DSP auction)
	TestViewPrice billing.Money // Price for test period per one view (as CPM mode)
	Price         billing.Money // Price per one view or click
	LeadPrice     billing.Money // Price per one lead
}

// Less then other target
func (t *TargetBid) Less(tb TargetBid) bool {
	return t.ECPM < tb.ECPM
}

// FixedPurchasePrice returns fixed price for the target action
func (t *TargetBid) PricePerAction(action Action) billing.Money {
	switch action {
	case ActionView:
		if t.TestMode && t.TestViewPrice > 0 {
			return t.TestViewPrice
		}
		if t.Ad.PricingModel.IsCPM() {
			return t.Price
		}
	case ActionClick:
		if t.Ad.PricingModel.IsCPC() {
			return t.Price
		}
	case ActionLead:
		if t.Ad.PricingModel.IsCPA() {
			return t.LeadPrice
		}
	}
	return 0
}

// MaxBidPrice returns max bid price for the target
func (t *TargetBid) MaxBidPrice() billing.Money {
	if t.BidPrice > 0 {
		return t.BidPrice
	}
	if t.Ad.Campaign.MaxBid > 0 {
		return t.Ad.Campaign.MaxBid
	}
	return gocast.IfThen(t.ECPM > 0, t.ECPM/1000, t.PricePerAction(ActionView))
}

// CalcBidPrice returns calculated bid price for the target
func (t *TargetBid) CalcBidPrice() (bidPrice billing.Money) {
	if t.BidPrice > 0 {
		bidPrice = t.BidPrice
	} else {
		bidPrice = t.PricePerAction(ActionView)
	}
	return min(bidPrice, t.MaxBidPrice())
}

// TargetBidList object list representation
type TargetBidList []TargetBid

// Random target from list
func (list TargetBidList) Random() TargetBid {
	if len(list) == 0 {
		return TargetBid{}
	}
	return list[rand.Intn(len(list))]
}

// Weighted returns target by Ad.Weight
func (list TargetBidList) Weighted() TargetBid {
	if len(list) == 0 {
		return TargetBid{}
	}
	if len(list) == 1 {
		return list[0]
	}
	// Calculate min weight
	minWeight := uint8(255)
	for _, tb := range list {
		minWeight = min(minWeight, tb.Ad.Weight)
	}
	// Get random target no less then minWeight
	j := rand.Intn(len(list))
	for i := j; i < len(list); i++ {
		if list[i].Ad.Weight >= minWeight {
			return list[i]
		}
	}
	for i := 0; i < j; i++ {
		if list[i].Ad.Weight >= minWeight {
			return list[i]
		}
	}
	return list[0]
}
