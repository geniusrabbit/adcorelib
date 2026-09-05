package trafaret

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/adquery/bidresponse"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFiller(t *testing.T) {
	filler := &Filler{}

	ad1 := blankAd("ad1", "imp1", 1.0)
	ad2 := blankAd("ad2", "imp1", 0.5)
	ad3 := blankAd("ad3", "imp1", 2.0)
	adx1 := blockAd(ad1, ad3)

	filler.Push(0.3, ad1, ad2)
	filler.Push(0.5, ad1, ad2, ad3)
	filler.Push(0.7, adx1, ad3)

	for _, size := range []int{1, 2, 3} {
		ls := filler.Copy().Fill("imp1", size)
		resSize := adsListRealSize(ls)
		assert.Equal(t, size, resSize, "Expected %d items, got %d", size, resSize)
	}
}

func TestFillerPushNilAndEmpty(t *testing.T) {
	filler := &Filler{}

	filler.Push(0.5)
	assert.Equal(t, 0, filler.Len())

	filler.Push(0.5, nil)
	assert.Equal(t, 0, filler.Len())

	var typedNil *bidresponse.ResponseItemBlank
	filler.Push(0.5, typedNil)
	assert.Equal(t, 0, filler.Len())

	ad := blankAd("ad1", "imp1", 1)
	filler.Push(0.5, nil, typedNil, ad)
	assert.Equal(t, 1, filler.Len())
	require.NotNil(t, filler.Block("imp1"))
	assert.Equal(t, 1, filler.Block("imp1").Len())
}

func TestFillerPushSameImpression(t *testing.T) {
	filler := &Filler{}
	ad1 := blankAd("ad1", "imp1", 1)
	ad2 := blankAd("ad2", "imp1", 2)

	filler.Push(0.4, ad1, ad2)

	assert.Equal(t, 1, filler.Len())
	block := filler.Block("imp1")
	require.NotNil(t, block)
	assert.Equal(t, "imp1", block.impid)
	assert.Equal(t, float32(0.4), block.summ)
	assert.Equal(t, 1, block.Len())
	assert.ElementsMatch(t, []string{"ad1", "ad2"}, adIDs(block.ads[0].ads))
}

func TestFillerPushMultipleImpressions(t *testing.T) {
	filler := &Filler{}
	adA := blankAd("a", "impA", 1)
	adB := blankAd("b", "impB", 2)
	adA2 := blankAd("a2", "impA", 3)

	filler.Push(0.5, adA, adB, adA2)

	assert.Equal(t, 2, filler.Len())
	blockA := filler.Block("impA")
	blockB := filler.Block("impB")
	require.NotNil(t, blockA)
	require.NotNil(t, blockB)
	assert.ElementsMatch(t, []string{"a", "a2"}, adIDs(blockA.ads[0].ads))
	assert.Equal(t, []string{"b"}, adIDs(blockB.ads[0].ads))
	assert.Nil(t, filler.Block("missing"))
}

func TestFillerPushAppendsToExistingBlock(t *testing.T) {
	filler := &Filler{}
	ad1 := blankAd("ad1", "imp1", 1)
	ad2 := blankAd("ad2", "imp1", 2)

	filler.Push(0.3, ad1)
	filler.Push(0.5, ad2)

	assert.Equal(t, 1, filler.Len())
	block := filler.Block("imp1")
	require.NotNil(t, block)
	assert.Equal(t, float32(0.8), block.summ)
	assert.Equal(t, 2, block.Len())
}

func TestFillerFillZeroNegativeAndUnknown(t *testing.T) {
	filler := &Filler{}
	filler.Push(1, blankAd("ad1", "imp1", 1))

	assert.Nil(t, filler.Fill("imp1", 0))
	assert.Nil(t, filler.Fill("imp1", -1))
	assert.Nil(t, filler.Fill("missing", 1))
}

func TestFillerFillSinglesUpToSize(t *testing.T) {
	filler := &Filler{}
	filler.Push(1,
		blankAd("ad1", "imp1", 1),
		blankAd("ad2", "imp1", 2),
		blankAd("ad3", "imp1", 3),
	)

	ls := filler.Copy().Fill("imp1", 2)
	require.Len(t, ls, 2)
	for _, ad := range ls {
		assert.Equal(t, 1, adSize(ad))
	}
	assert.Subset(t, []string{"ad1", "ad2", "ad3"}, adIDs(ls))
}

func TestFillerFillDoesNotExceedInventory(t *testing.T) {
	filler := &Filler{}
	filler.Push(1,
		blankAd("ad1", "imp1", 1),
		blankAd("ad2", "imp1", 2),
		blankAd("ad3", "imp1", 3),
	)

	ls := filler.Copy().Fill("imp1", 10)
	assert.Len(t, ls, 3)
	assert.ElementsMatch(t, []string{"ad1", "ad2", "ad3"}, adIDs(ls))
}

func TestFillerCopyIsIndependent(t *testing.T) {
	orig := &Filler{}
	orig.Push(1,
		blankAd("ad1", "imp1", 1),
		blankAd("ad2", "imp1", 2),
	)

	first := orig.Copy().Fill("imp1", 2)
	require.Len(t, first, 2)

	second := orig.Copy().Fill("imp1", 2)
	require.Len(t, second, 2)
	assert.ElementsMatch(t, []string{"ad1", "ad2"}, adIDs(second))

	orig.Fill("imp1", 2)
	assert.Empty(t, orig.Copy().Fill("imp1", 2))
}

func TestFillerPushDoesNotAliasCallerSlice(t *testing.T) {
	orig := []adtype.ResponseItemCommon{
		blankAd("ad1", "imp1", 1),
		blankAd("ad2", "imp1", 2),
	}

	filler := &Filler{}
	filler.Push(1, orig...)
	_ = filler.Fill("imp1", 2)

	require.Len(t, orig, 2)
	require.NotNil(t, orig[0])
	require.NotNil(t, orig[1])
	assert.ElementsMatch(t, []string{"ad1", "ad2"}, adIDs(orig))
}

func TestFillerFillPacksMultipleItems(t *testing.T) {
	t.Run("only block", func(t *testing.T) {
		filler := &Filler{}
		block := blockAd(
			blankAd("b1", "imp1", 5),
			blankAd("b2", "imp1", 5),
		)
		filler.Push(1, block)

		ls := filler.Copy().Fill("imp1", 2)
		require.Len(t, ls, 1)
		assert.Equal(t, 2, adsListRealSize(ls))
	})

	t.Run("block and singles", func(t *testing.T) {
		filler := &Filler{}
		ad1 := blankAd("s1", "imp1", 1)
		ad2 := blankAd("s2", "imp1", 1)
		block := blockAd(
			blankAd("b1", "imp1", 5),
			blankAd("b2", "imp1", 5),
		)
		filler.Push(1, block, ad1, ad2)

		const size = 2
		ls := filler.Copy().Fill("imp1", size)
		require.NotEmpty(t, ls)
		assert.LessOrEqual(t, adsListRealSize(ls), size)
	})
}

func TestFillerFillExhaustsAllBuckets(t *testing.T) {
	high := blankAd("high", "imp1", 10)
	low := make([]adtype.ResponseItemCommon, 0, 8)
	for i := range 8 {
		low = append(low, blankAd(string(rune('a'+i)), "imp1", 1))
	}

	filler := &Filler{}
	filler.Push(0.9, high)
	filler.Push(0.1, low...)

	const want = 9
	ls := filler.Copy().Fill("imp1", want)
	assert.Len(t, ls, want, "Fill should drain remaining ads from other priority buckets")
}

func TestAdSize(t *testing.T) {
	assert.Equal(t, 0, adSize(nil))
	assert.Equal(t, 1, adSize(blankAd("ad1", "imp1", 1)))
	assert.Equal(t, 2, adSize(blockAd(
		blankAd("a", "imp1", 1),
		blankAd("b", "imp1", 2),
	)))
	assert.Equal(t, 0, adSize(blockAd()))
}

func TestFilterNilAds(t *testing.T) {
	ad := blankAd("ad1", "imp1", 1)
	var typedNil *bidresponse.ResponseItemBlank
	got := filterNilAds([]adtype.ResponseItemCommon{nil, ad, typedNil})
	require.Len(t, got, 1)
	assert.Equal(t, "ad1", got[0].ID())

	assert.Empty(t, filterNilAds(nil))
	assert.Empty(t, filterNilAds([]adtype.ResponseItemCommon{nil, typedNil}))
}

func TestPackAdObjectsPrefersHigherCPM(t *testing.T) {
	block := blockAd(
		blankAd("b1", "imp1", 5),
		blankAd("b2", "imp1", 5),
	)
	s1 := blankAd("s1", "imp1", 3)
	s2 := blankAd("s2", "imp1", 2)

	got := packAdObjects([]adtype.ResponseItemCommon{block, s1, s2}, 2)
	require.Len(t, got, 1)
	multi, ok := got[0].(adtype.ResponseMultipleItem)
	require.True(t, ok)
	assert.Equal(t, 2, multi.Count())
	assert.Equal(t, block.InternalAuctionCPMBid(), got[0].InternalAuctionCPMBid())
}

func TestPackAdObjectsExactFit(t *testing.T) {
	block := blockAd(
		blankAd("b1", "imp1", 5),
		blankAd("b2", "imp1", 5),
	)

	got := packAdObjects([]adtype.ResponseItemCommon{block}, 2)
	require.Len(t, got, 1)
	assert.Equal(t, 2, adSize(got[0]))

	s1 := blankAd("s1", "imp1", 1)
	s2 := blankAd("s2", "imp1", 1)
	singles := packAdObjects([]adtype.ResponseItemCommon{s1, s2}, 2)
	assert.ElementsMatch(t, []string{"s1", "s2"}, adIDs(singles))
}

func TestPackAdObjectsEmptyAndZeroSize(t *testing.T) {
	assert.Empty(t, packAdObjects(nil, 2))
	assert.Empty(t, packAdObjects([]adtype.ResponseItemCommon{blankAd("ad1", "imp1", 1)}, 0))
}
