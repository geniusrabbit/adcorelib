package trafaret

import (
	"math/rand/v2"

	"github.com/demdxx/gocast/v2"

	"github.com/geniusrabbit/adcorelib/adtype"
)

// blockPriority represents a collection of adPreority grouped by impression ID.
type blockPriority struct {
	impid string
	summ  float32
	ads   []adPreority
}

// Len returns the number of adPreority blocks.
func (b *blockPriority) Len() int {
	return len(b.ads)
}

// Pop selects an ad from the blocks based on a weighted random selection.
func (b *blockPriority) Pop() (float32, adtype.ResponseItemCommon) {
	if len(b.ads) == 0 {
		return 0, nil
	}

	var liveSum float32
	for i := range b.ads {
		if hasLiveAd(&b.ads[i]) {
			liveSum += b.ads[i].priority
		}
	}
	if liveSum == 0 {
		return 0, nil
	}

	rv := rand.Float32() * liveSum
	vl := float32(0)
	for i := range b.ads {
		if !hasLiveAd(&b.ads[i]) {
			continue
		}
		vl += b.ads[i].priority
		if rv <= vl {
			if ad := b.ads[i].Pop(); !gocast.IsNil(ad) {
				return b.ads[i].priority, ad
			}
			break
		}
	}

	for i := range b.ads {
		if ad := b.ads[i].Pop(); !gocast.IsNil(ad) {
			return b.ads[i].priority, ad
		}
	}
	return 0, nil
}

func hasLiveAd(a *adPreority) bool {
	for _, ad := range a.ads {
		if !gocast.IsNil(ad) {
			return true
		}
	}
	return false
}

func (b *blockPriority) copyFrom(other *blockPriority) {
	b.impid = other.impid
	b.summ = other.summ
	b.ads = make([]adPreority, len(other.ads))
	for i := range other.ads {
		b.ads[i].copyFrom(&other.ads[i])
	}
}
