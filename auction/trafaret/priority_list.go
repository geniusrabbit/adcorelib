package trafaret

import (
	"slices"

	"github.com/geniusrabbit/adcorelib/adtype"
)

// adPreority represents a collection of ads with a specific priority.
type adPreority struct {
	priority float32 // As percentage from 0 to 1 (0.0 - 100%)
	ads      []adtype.ResponseItemCommon
}

// Len returns the number of ads in the collection.
func (a *adPreority) Len() int {
	return len(a.ads)
}

// Pop removes and returns the last ad in the collection.
func (a *adPreority) Pop() adtype.ResponseItemCommon {
	if len(a.ads) == 0 {
		return nil
	}
	idx := len(a.ads) - 1
	ad := a.ads[idx]
	a.ads = a.ads[:idx]
	return ad
}

// Sort orders the ads in ascending order based on their CPM bid.
func (a *adPreority) Sort() {
	slices.SortFunc(a.ads, func(a, b adtype.ResponseItemCommon) int {
		bid1 := a.InternalAuctionCPMBid()
		bid2 := b.InternalAuctionCPMBid()
		if bid1 == bid2 {
			return 0
		}
		if bid1 < bid2 {
			return -1
		}
		return 1
	})
}

func (a *adPreority) copyFrom(other *adPreority) {
	a.priority = other.priority
	a.ads = append(a.ads[:0], other.ads...)
}
