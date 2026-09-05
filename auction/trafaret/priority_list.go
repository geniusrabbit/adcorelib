package trafaret

import (
	"math/rand/v2"
	"slices"

	"github.com/demdxx/gocast/v2"

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
	size := len(a.ads)
	offset := rand.IntN(size)
	for i := range size {
		idx := (i + offset) % size
		if ad := a.ads[idx]; !gocast.IsNil(ad) {
			a.ads[idx] = nil
			return ad
		}
	}
	return nil
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

// Shuffle randomizes the order of the ads in the collection.
func (a *adPreority) Shuffle() {
	rand.Shuffle(len(a.ads), func(i, j int) {
		a.ads[i], a.ads[j] = a.ads[j], a.ads[i]
	})
}

func (a *adPreority) copyFrom(other *adPreority) {
	a.priority = other.priority
	a.ads = append(a.ads[:0], other.ads...)
}
