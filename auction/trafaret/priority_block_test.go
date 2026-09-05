package trafaret

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockPriorityPopEmpty(t *testing.T) {
	var empty blockPriority
	prio, ad := empty.Pop()
	assert.Equal(t, float32(0), prio)
	assert.Nil(t, ad)
}

func TestBlockPriorityPopSingleBucket(t *testing.T) {
	b := blockPriority{
		impid: "imp1",
		summ:  1,
		ads: []adPreority{{
			priority: 1,
			ads: []adtype.ResponseItemCommon{
				blankAd("a", "imp1", 1),
				blankAd("b", "imp1", 2),
			},
		}},
	}

	got := make([]string, 0, 2)
	for range 2 {
		prio, ad := b.Pop()
		require.NotNil(t, ad)
		assert.Equal(t, float32(1), prio)
		got = append(got, ad.ID())
	}
	prio, ad := b.Pop()
	assert.Equal(t, float32(0), prio)
	assert.Nil(t, ad)
	assert.ElementsMatch(t, []string{"a", "b"}, got)
}

func TestBlockPriorityPopRespectsWeights(t *testing.T) {
	const n = 4000
	src := blockPriority{
		impid: "imp1",
		summ:  1,
		ads: []adPreority{
			{priority: 0.9, ads: []adtype.ResponseItemCommon{blankAd("high", "imp1", 1)}},
			{priority: 0.1, ads: []adtype.ResponseItemCommon{blankAd("low", "imp1", 1)}},
		},
	}

	high := 0
	for range n {
		var copy blockPriority
		copy.copyFrom(&src)
		_, ad := copy.Pop()
		require.NotNil(t, ad)
		if ad.ID() == "high" {
			high++
		}
	}

	ratio := float64(high) / n
	assert.Greater(t, ratio, 0.75, "high-priority bucket should win most of the time, got %v", ratio)
	assert.Less(t, ratio, 0.98, "low-priority bucket should still be selected sometimes, got %v", ratio)
}

func TestBlockPriorityCopyFromIndependent(t *testing.T) {
	src := blockPriority{
		impid: "imp1",
		summ:  0.5,
		ads: []adPreority{{
			priority: 0.5,
			ads:      []adtype.ResponseItemCommon{blankAd("a", "imp1", 1)},
		}},
	}
	var dst blockPriority
	dst.copyFrom(&src)

	assert.Equal(t, src.impid, dst.impid)
	assert.Equal(t, src.summ, dst.summ)
	require.Len(t, dst.ads, 1)
	assert.Equal(t, []string{"a"}, adIDs(dst.ads[0].ads))

	dst.ads[0].ads[0] = nil
	assert.Equal(t, "a", src.ads[0].ads[0].ID())
}
