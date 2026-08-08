//
// @project GeniusRabbit corelib 2026
// @author Dmitry Ponomarev <demdxx@gmail.com>
//

package preprocessors

import (
	"cmp"
	"slices"

	"github.com/geniusrabbit/adcorelib/adsource"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// SecondPrice adjusts already-selected response bids to second-price (GSP) clearing.
//
// Wire with adsource.WithResponsePreprocessor(preprocessors.SecondPrice{}).
//
// Algorithm:
//  1. Collect flat ads via IterAds and sort by InternalAuctionCPMBid descending.
//  2. For each ad, clearing price is (in order): Second().GetPrice() if > 0,
//     else the next ad's Price(ActionImpression), else the single/tail rule:
//     - BidFloorCPM > 0 → BidFloor/1000 + half the gap to current (tail capped
//     by min of other cleared prices)
//     - no BidFloor, MinimumPrice > 0 → MinimumPrice
//     - BidFloor and MinimumPrice both ≤ 0 → 10% of current, not above the
//     previous (higher-ranked) ad's cleared price
//  3. Apply via SetBidPrice(..., withCommission=true). Skip zero/negative.
type SecondPrice struct{}

var _ adsource.ResponsePreprocessor = SecondPrice{}

// PreprocessResponse implements adsource.ResponsePreprocessor.
func (SecondPrice) PreprocessResponse(response adtype.Response) (adtype.Response, error) {
	if response == nil || response.Count() < 1 {
		return response, nil
	}

	ads := collectAds(response)
	if len(ads) == 0 {
		return response, nil
	}

	slices.SortFunc(ads, func(a, b adtype.ResponseItem) int {
		return cmp.Compare(b.InternalAuctionCPMBid(), a.InternalAuctionCPMBid())
	})

	// Snapshot original impression prices before any SetBidPrice mutation.
	origPrices := make([]billing.Money, len(ads))
	for i, ad := range ads {
		origPrices[i] = ad.Price(adtype.ActionImpression)
	}

	cleared := make([]billing.Money, len(ads))
	for i := range ads {
		cleared[i] = resolveClearedPrice(ads, origPrices, cleared, i)
	}

	for i, ad := range ads {
		if cleared[i] <= 0 {
			continue
		}
		_ = ad.SetBidPrice(adtype.ActionImpression, cleared[i], true)
	}

	return response, nil
}

func collectAds(response adtype.Response) []adtype.ResponseItem {
	var ads []adtype.ResponseItem
	for ad := range response.IterAds() {
		if ad != nil {
			ads = append(ads, ad)
		}
	}
	return ads
}

func resolveClearedPrice(
	ads []adtype.ResponseItem,
	origPrices []billing.Money,
	cleared []billing.Money,
	i int,
) billing.Money {
	ad := ads[i]

	if second := ad.Second(); second != nil {
		if p := second.GetPrice(); p > 0 {
			return p
		}
	}

	if i+1 < len(ads) {
		return origPrices[i+1]
	}

	return singleOrTailPrice(ad, i, cleared)
}

// singleOrTailPrice clears the sole or last ad when there is no Second/next competitor.
func singleOrTailPrice(ad adtype.ResponseItem, i int, cleared []billing.Money) billing.Money {
	current := ad.Price(adtype.ActionImpression)

	var (
		bidFloorCPM billing.Money
		minPrice    billing.Money
	)
	if imp := ad.Impression(); imp != nil {
		bidFloorCPM = imp.BidFloorCPM
		minPrice = imp.MinimumPrice(adtype.ActionImpression)
	}

	if bidFloorCPM <= 0 {
		if minPrice > 0 {
			return minPrice
		}
		// No BidFloor and no MinimumPrice → 10% of current, ≤ previous ad.
		price := current / 10
		if i > 0 {
			if prev := cleared[i-1]; prev > 0 && price > prev {
				price = prev
			}
		}
		return price
	}

	floor := bidFloorCPM / 1000
	price := floor + (current-floor)/2
	if i > 0 {
		if cap := minPositive(cleared[:i]); cap > 0 && price > cap {
			price = cap
		}
	}
	return price
}

func minPositive(prices []billing.Money) billing.Money {
	var (
		min    billing.Money
		hasMin bool
	)
	for _, p := range prices {
		if p <= 0 {
			continue
		}
		if !hasMin || p < min {
			min = p
			hasMin = true
		}
	}
	return min
}
