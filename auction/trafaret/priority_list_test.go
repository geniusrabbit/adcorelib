package trafaret

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdPriorityPopEmpty(t *testing.T) {
	var empty adPreority
	assert.Equal(t, 0, empty.Len())
	assert.Nil(t, empty.Pop())

	allNil := adPreority{ads: []adtype.ResponseItemCommon{nil, nil}}
	assert.Nil(t, allNil.Pop())
}

func TestAdPriorityPopEachAdOnce(t *testing.T) {
	p := adPreority{
		priority: 1,
		ads: []adtype.ResponseItemCommon{
			blankAd("a", "imp1", 1),
			blankAd("b", "imp1", 2),
			blankAd("c", "imp1", 3),
		},
	}

	got := make([]string, 0, 3)
	for range 3 {
		ad := p.Pop()
		require.NotNil(t, ad)
		got = append(got, ad.ID())
	}
	assert.Nil(t, p.Pop())
	assert.ElementsMatch(t, []string{"a", "b", "c"}, got)
}

func TestAdPriorityPopSkipsNils(t *testing.T) {
	keep := blankAd("keep", "imp1", 1)
	p := adPreority{
		ads: []adtype.ResponseItemCommon{nil, keep, nil},
	}

	assert.Equal(t, "keep", p.Pop().ID())
	assert.Nil(t, p.Pop())
}

func TestAdPrioritySortByCPMAscending(t *testing.T) {
	p := adPreority{
		ads: []adtype.ResponseItemCommon{
			blankAd("c", "imp1", 3),
			blankAd("a", "imp1", 1),
			blankAd("b", "imp1", 2),
			blankAd("a2", "imp1", 1),
		},
	}
	p.Sort()

	for i := 1; i < len(p.ads); i++ {
		assert.LessOrEqual(t, p.ads[i-1].InternalAuctionCPMBid(), p.ads[i].InternalAuctionCPMBid())
	}
}

func TestAdPriorityShuffleKeepsAds(t *testing.T) {
	p := adPreority{
		ads: []adtype.ResponseItemCommon{
			blankAd("a", "imp1", 1),
			blankAd("b", "imp1", 2),
			blankAd("c", "imp1", 3),
		},
	}
	before := adIDs(p.ads)
	p.Shuffle()
	assert.ElementsMatch(t, before, adIDs(p.ads))
}

func TestAdPriorityCopyFromIndependent(t *testing.T) {
	src := adPreority{
		priority: 0.7,
		ads: []adtype.ResponseItemCommon{
			blankAd("a", "imp1", 1),
			blankAd("b", "imp1", 2),
		},
	}
	var dst adPreority
	dst.copyFrom(&src)

	assert.Equal(t, src.priority, dst.priority)
	assert.Equal(t, adIDs(src.ads), adIDs(dst.ads))

	dst.ads[0] = nil
	assert.Equal(t, "a", src.ads[0].ID())
}
