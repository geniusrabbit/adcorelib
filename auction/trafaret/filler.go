package trafaret

import (
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// Filler manages a collection of blockPriority objects.
type Filler struct {
	blocks []blockPriority
}

// Push adds ads to the filler collection, grouping them by impression ID.
func (f *Filler) Push(priority float32, ads ...adtype.ResponseItemCommon) {
	if len(ads) == 0 {
		return
	}
	// Check if all ads share the same impression ID.
	impID := ads[0].ImpressionID()
	oneImpression := true
	for i := 1; i < len(ads); i++ {
		if ads[i].ImpressionID() != impID {
			oneImpression = false
			break
		}
	}

	if oneImpression {
		// Handle ads with a single impression ID.
		priorityAds := adPreority{
			priority: priority,
			ads:      ads,
		}
		priorityAds.Sort()

		// Add to an existing block or create a new one.
		for i := range f.blocks {
			if f.blocks[i].impid == impID {
				f.blocks[i].summ += priority
				f.blocks[i].ads = append(f.blocks[i].ads, priorityAds)
				return
			}
		}
		f.blocks = append(f.blocks, blockPriority{
			impid: impID,
			summ:  priority,
			ads:   []adPreority{priorityAds},
		})
		return
	}

	// Handle ads with multiple impression IDs.
	blocks := make(map[string]adPreority, len(ads))
	for _, ad := range ads {
		impID := ad.ImpressionID()
		if blockAds, ok := blocks[impID]; ok {
			blockAds.ads = append(blockAds.ads, ad)
			blocks[impID] = blockAds
		} else {
			blocks[impID] = adPreority{
				priority: priority,
				ads:      []adtype.ResponseItemCommon{ad},
			}
		}
	}

	// Add grouped ads to existing blocks or create new ones.
	for impID, block := range blocks {
		foundBlock := false
		for i := 0; i < len(f.blocks); i++ {
			if f.blocks[i].impid == impID {
				f.blocks[i].summ += block.priority
				f.blocks[i].ads = append(f.blocks[i].ads, block)
				foundBlock = true
				break
			}
		}
		if !foundBlock {
			f.blocks = append(f.blocks, blockPriority{
				impid: impID,
				summ:  block.priority,
				ads:   []adPreority{block},
			})
		}
	}

	// Sort the blocks by impression ID.
	for i := range f.blocks {
		for j := range f.blocks[i].ads {
			f.blocks[i].ads[j].Sort()
		}
	}
}

// Len returns the number of blocks in the filler.
func (f *Filler) Len() int {
	return len(f.blocks)
}

// Fill retrieves ads for a specific impression ID up to the specified size.
func (f *Filler) Fill(impid string, size int) []adtype.ResponseItemCommon {
	if size <= 0 {
		return nil
	}

	block := f.Block(impid)
	if block == nil {
		return nil
	}

	muliadsCount := 0
	result := make([]adtype.ResponseItemCommon, 0, size)

	// Retrieve ads from the block.
	for i := 0; i < size; i++ {
		_, ad := block.Pop()
		if ad == nil {
			break
		}
		result = append(result, ad)
		switch ad.(type) {
		case adtype.ResponseMultipleItem:
			muliadsCount++
			size++
		}
	}

	if muliadsCount == 0 {
		return result
	}

	// Postprocess ads if multiple ads are present.
	return packAdObjects(result, size)
}

// Block retrieves a blockPriority by impression ID.
func (f *Filler) Block(impid string) *blockPriority {
	for i := range f.blocks {
		if f.blocks[i].impid == impid {
			return &f.blocks[i]
		}
	}
	return nil
}

// Copy creates a deep copy of the Filler object.
func (f *Filler) Copy() *Filler {
	copyBlocks := make([]blockPriority, len(f.blocks))
	for i := range f.blocks {
		copyBlocks[i].copyFrom(&f.blocks[i])
	}
	return &Filler{blocks: copyBlocks}
}

// packAdObjects optimizes the selection of ads to fit within the maximum size.
func packAdObjects(objects []adtype.ResponseItemCommon, maxSize int) []adtype.ResponseItemCommon {
	n := len(objects)
	blockSize := maxSize + 1
	dp := make([]billing.Money, (n+1)*blockSize)

	// Fill the dp table for the knapsack problem.
	for i := 1; i <= n; i++ {
		for w := 0; w <= maxSize; w++ {
			dpIdx := (i - 1) * blockSize
			if objSize := adSize(objects[i-1]); objSize < w {
				withItem := dp[dpIdx+w-objSize] + objects[i-1].InternalAuctionCPMBid()
				withoutItem := dp[dpIdx+w]
				dp[i*blockSize+w] = max(withItem, withoutItem)
			} else {
				dp[i*blockSize+w] = dp[dpIdx+w]
			}
		}
	}

	// Restore selected objects.
	res := []adtype.ResponseItemCommon{}
	w := maxSize
	for i := n; i > 0; i-- {
		if dp[i*blockSize+w] != dp[(i-1)*blockSize+w] {
			res = append([]adtype.ResponseItemCommon{objects[i-1]}, res...)
			if w -= adSize(objects[i-1]); w <= 0 {
				break
			}
		}
	}

	return res
}

// adSize calculates the size of an ad.
func adSize(ad adtype.ResponseItemCommon) int {
	switch adv := ad.(type) {
	case nil:
		return 0
	case adtype.ResponseMultipleItem:
		return adv.Count()
	default:
		return 1
	}
}
