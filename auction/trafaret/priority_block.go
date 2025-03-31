package trafaret

import (
	"math/rand/v2"

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
func (b *blockPriority) Pop() (float32, adtype.ResponserItemCommon) {
	if len(b.ads) == 0 {
		return 0, nil
	}
	rv := rand.Float32() * b.summ
	vl := float32(0)
	for i := 0; i < len(b.ads); i++ {
		vl += b.ads[i].priority
		if len(b.ads[i].ads) == 0 {
			continue
		}
		if rv <= vl {
			ad := b.ads[i].Pop()
			return b.ads[i].priority, ad
		}
	}
	return 0, nil
}

func (b *blockPriority) copyFrom(other *blockPriority) {
	b.impid = other.impid
	b.summ = other.summ
	b.ads = make([]adPreority, len(other.ads))
	for i := range other.ads {
		b.ads[i].copyFrom(&other.ads[i])
	}
}
