// Package trafaret fills impression slots from a priority-weighted pool of ads.
//
// A source response is pushed into a [Filler] together with that source's
// priority. Ads are grouped by impression ID. Fill then draws from that pool
// until the requested number of slots is occupied.
//
// The typical consumer is a multi-source bid wrapper: each upstream source
// reports ads at its own priority, then each impression is filled independently.
//
// # Data model
//
//	Filler
//	└─ blockPriority          one per impression ID
//	   ├─ summ               sum of bucket priorities (informational)
//	   └─ ads[] adPreority   one bucket per Push for this impression
//	      ├─ priority        weight of this source / push (0..1)
//	      └─ ads[]           shuffled ResponseItemCommon
//
// A single item occupies one slot. A [adtype.ResponseMultipleItem] occupies
// Count() slots (a bundled banner).
//
// # Push
//
//	filler.Push(priority, ads...)
//
// Nil and typed-nil ads are dropped. Ads that share an impression ID land in
// the same blockPriority; mixed IDs are split. A later Push for an already
// seen impression appends a new priority bucket and adds to summ. Each
// bucket is shuffled so equal-priority ads are not taken in input order.
//
// # Fill
//
//	items := filler.Copy().Fill(impID, size)
//
// Fill is destructive: Pop nils taken slots. Copy before Fill when the same
// pool must be reused (for example filling several impressions, or retrying).
//
// For one impression the loop is:
//
//  1. Weighted random Pop from live priority buckets (exhausted buckets are
//     skipped so remaining inventory is still used).
//  2. A multiple item increments the fetch budget so extra candidates are
//     collected; packing still uses the original size.
//  3. If any multiple item was drawn, a 0-1 knapsack (packAdObjects) keeps the
//     subset with the highest InternalAuctionCPMBid that fits in size slots.
//     Otherwise the drawn list is returned as-is.
//
// size <= 0 or an unknown impression ID yields nil. If inventory runs out
// earlier, Fill returns whatever was drawn.
//
// # Weighted pop
//
// Among buckets that still have a non-nil ad, the chance of drawing from
// bucket i is:
//
//	P(i) = priority_i / sum(priority of live buckets)
//
// Inside a bucket Pop starts at a random offset and walks the slice, skipping
// already taken (nil) slots.
//
// # Knapsack
//
// packAdObjects is a classic 0-1 knapsack: item weight is adSize (1 or Count()),
// value is InternalAuctionCPMBid, capacity is the original Fill size. An item
// that exactly fills the remaining capacity is eligible (weight <= capacity).
package trafaret
