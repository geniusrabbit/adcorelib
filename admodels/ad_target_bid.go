package admodels

import (
	"math/rand"

	"github.com/geniusrabbit/adcorelib/billing"
)

// TargetBid object
type TargetBid struct {
	Ad        *Ad
	Bid       *AdBid
	ECPM      billing.Money
	BidPrice  billing.Money // Max price per one View (used in DSP auction)
	Price     billing.Money // Price per one view or click
	LeadPrice billing.Money // Price per one lead
}

// Less then other target
func (t TargetBid) Less(tb TargetBid) bool {
	return t.ECPM < tb.ECPM
}

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
