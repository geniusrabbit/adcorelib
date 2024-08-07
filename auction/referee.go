//
// @project GeniusRabbit corelib 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018
//

package auction

import (
	"math/rand"
	"sort"

	"github.com/geniusrabbit/adcorelib/adtype"
)

// Referee which could combine and select most profitable advertisement bids
type Referee struct {
	// normalized data
	normalized bool

	// equipment which used in auction competition
	equipment []adtype.ResponserItemCommon
}

// Push items into equipment
func (r *Referee) Push(its ...adtype.ResponserItemCommon) {
	if len(its) > 0 {
		r.normalized = false
		r.equipment = append(r.equipment, its...)
	}
}

// TotalCapacity of the equipment
func (r *Referee) TotalCapacity() (capacity int) {
	for _, it := range r.equipment {
		switch a := it.(type) {
		case adtype.ResponserMultipleItem:
			capacity += a.Count()
		default:
			capacity++
		}
	}
	return capacity
}

// Match point O(N * K * 2)
func (r *Referee) Match(rings ...Ring) (resp []adtype.ResponserItemCommon) {
	if len(rings) < 1 {
		return resp
	}

	r.normalize()

	// Borrow counters array from pool
	var (
		capacity         = ringsCapacity(rings)
		counters         = borrowCounters()
		tail             = make([]adtype.ResponserItemCommon, 0, capacity)
		tailCount        int
		multipleResponse bool
	)
	defer returnCounter(counters)

	// First fill loop, complexity O(N * k)
	for i, it := range r.equipment {
		if v, _ := it.(adtype.ResponserMultipleItem); v != nil {
			if capacity < v.Count() {
				continue
			}

			// I bealive that this is could be rare case and for improving
			// performance we need mark this case
			multipleResponse = true

			added := true
			for idx, rn := range rings {
				count := adsCountByImpID(rn.ID, v.Ads())
				if count <= 0 {
					continue
				}

				if rn.Count-counters.count(idx) >= count {
					counters.inc(idx, count)
					capacity -= count
				} else { // Revert counters
					for i := 0; i < idx; i++ {
						count := adsCountByImpID(rings[i].ID, v.Ads())
						counters.inc(i, -count)
						capacity += count
					}
					added = false
					break
				}
			}

			if added {
				resp = append(resp, v)
			}
		} else {
			idx, rn := ringByID(it.ImpressionID(), rings)
			if idx >= 0 && rn.Count-counters.count(idx) > 0 {
				resp = append(resp, it)
				counters.inc(idx, 1)
				capacity--
			}
		}

		if capacity < 1 {
			if multipleResponse {
				tail = r.equipment[i:]
			}
			break
		}
	}

	// Here we are cheking do we have extra equipment
	// for optimisation of response
	if len(tail) < 1 || capacity < r.TotalCapacity() {
		return resp
	}

	// Normalize ordering. Most cheaper multiresponses in front
	normalizeResponseForOptimization(resp)
	tailCount = len(resp)

	// Optimizing fill loop
	for i := 0; i < len(resp); i++ {
		it := resp[i]
		if v, _ := it.(adtype.ResponserMultipleItem); v != nil {
			// Not enought advertisement for filling this spcace
			if v.Count() > tailCount {
				continue
			}

			if replacement := collectReplacement(v, tail); len(replacement) > 0 {
				resp[i] = replacement[0]
				if len(replacement) > 1 {
					resp = append(resp, replacement[1:]...)
				}
				tailCount -= len(replacement)
				returnResponseList(replacement)
			}
		} else {
			tailCount = 0
		}

		if tailCount <= 0 {
			break
		}
	}

	return resp
}

// MatchRequest response by request
func (r *Referee) MatchRequest(req *adtype.BidRequest) []adtype.ResponserItemCommon {
	var rings []Ring
	for _, imp := range req.Imps {
		rings = append(rings, Ring{ID: imp.ID, Count: imp.Count})
	}
	return r.Match(rings...)
}

// Equipment list
func (r *Referee) Equipment() []adtype.ResponserItemCommon {
	return r.equipment
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

// normalize data for competition
func (r *Referee) normalize() {
	if !r.normalized {
		// Shuffle data before processing
		rand.Shuffle(len(r.equipment), func(i, j int) { r.equipment[i], r.equipment[j] = r.equipment[j], r.equipment[i] })
		sort.Sort(equipmentSlice(r.equipment))
		r.normalized = true
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func normalizeResponseForOptimization(resp []adtype.ResponserItemCommon) {
	sort.Slice(resp, func(i, j int) bool {
		e1, _ := resp[i].(adtype.ResponserMultipleItem)
		e2, _ := resp[j].(adtype.ResponserMultipleItem)
		if e1 != nil && e2 == nil {
			return true
		}
		if e1 == nil && e2 != nil {
			return false
		}
		return avgBid(e1, resp[i]) < avgBid(e2, resp[j])
	})
}

// collectReplacement for multiple response O(N * K)
func collectReplacement(target adtype.ResponserMultipleItem, items []adtype.ResponserItemCommon) (resp []adtype.ResponserItemCommon) {
	for _, it := range target.Ads() {
		for i, ad := range items {
			if ad == nil {
				continue
			}
			if _, ok := ad.(adtype.ResponserMultipleItem); !ok {
				if ad.ImpressionID() == it.ImpressionID() {
					if resp == nil {
						resp = borrowResponseList()
					}
					resp = append(resp, ad)
					items[i] = nil
					break
				}
			} // end if
		}
	}

	if len(resp) != target.Count() {
		for i, j := 0, 0; i < len(resp); i++ {
			for ; j < len(items); j++ {
				if items[j] == nil {
					items[j] = resp[i]
					break
				}
			}
		}
		returnResponseList(resp)
		resp = nil
	}
	return resp
}

func ringsCapacity(rings []Ring) (total int) {
	for _, ring := range rings {
		total += ring.Count
	}
	return total
}

func ringByID(id string, rings []Ring) (int, Ring) {
	for i, ring := range rings {
		if ring.ID == id {
			return i, ring
		}
	}
	return -1, Ring{}
}

func adsCountByImpID(impID string, ads any) (count int) {
	switch nads := ads.(type) {
	case []adtype.ResponserItemCommon:
		for _, ad := range nads {
			if ad.ImpressionID() == impID {
				count++
			}
		}
	case []adtype.ResponserItem:
		for _, ad := range nads {
			if ad.ImpressionID() == impID {
				count++
			}
		}
	}
	return count
}
