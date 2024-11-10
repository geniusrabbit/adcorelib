package inmemory

import (
	"sort"
	"sync/atomic"

	"github.com/demdxx/gocast/v2"
	"github.com/demdxx/rpool/v2"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adsource/inmemory/adresponse"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/auction"
	"github.com/geniusrabbit/adcorelib/rand"
	"github.com/geniusrabbit/adstorage/accessors/campaignaccessor"
)

type driver struct {
	execPool *rpool.Pool

	adCampaigns    *campaignaccessor.CampaignAccessor
	balanceManager balanceManager
	adStore        atomic.Value
}

func newDriver(adCampaigns *campaignaccessor.CampaignAccessor, balanceManager balanceManager) *driver {
	return &driver{
		adCampaigns:    adCampaigns,
		balanceManager: balanceManager,
		execPool:       rpool.NewSinglePool(),
	}
}

// ID of the source driver
func (d *driver) ID() uint64 { return 0 }

// ObjectKey of the source driver
func (d *driver) ObjectKey() uint64 { return d.ID() }

// Protocol of the source driver
func (d *driver) Protocol() string { return protocol }

// Test request before processing
func (d *driver) Test(request *adtype.BidRequest) bool { return true }

// PriceCorrectionReduceFactor which is a potential
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (d *driver) PriceCorrectionReduceFactor() float64 {
	return 0
}

// RequestStrategy description
func (d *driver) RequestStrategy() adtype.RequestStrategy {
	return adtype.AsynchronousRequestStrategy
}

// Bid request for standart system filter
func (d *driver) Bid(request *adtype.BidRequest) adtype.Responser {
	// Reload ads if needed
	if d.adCampaigns.NeedUpdate() {
		d.execPool.Go(d.reloadAds)
	}

	// Create a referee for the auction
	// TODO: reduce allocations with sync.Pool
	var referee auction.Referee

	// Request ads for each impression block in the request
	for i := 0; i < len(request.Imps); i++ {
		imp := &request.Imps[i]

		// Bid request for each impression
		if ires := d.bidImp(request, imp); len(ires) > 0 {
			// Add the response items to the referee team
			referee.Push(ires...)
		}
	}

	// Match the request with the ads according to the auction competition rules
	itemResp := referee.MatchRequest(request)

	// If there is no response items, return an empty response
	if len(itemResp) == 0 {
		return adtype.NewEmptyResponse(request, d, nil)
	}

	// Return the response with the items
	return adtype.NewResponse(request, d, itemResp, nil)
}

// ProcessResponseItem result or error
func (d *driver) ProcessResponseItem(resp adtype.Responser, item adtype.ResponserItem) {
	if d.balanceManager != nil {
		// Reduce inmemory balance counters
		adItem := item.(*adresponse.ResponseAdItem)
		d.balanceManager.MakeVirtualView(adItem.PriceScope.TestViewBudget,
			adItem.Ad, adItem.PriceScope.ViewPrice) //adItem.BidECPM)
	}
}

// process bid request for the impression
func (d *driver) bidImp(request *adtype.BidRequest, imp *adtype.Impression) []adtype.ResponserItemCommon {
	res := make([]adtype.ResponserItemCommon, 0, len(imp.Formats())*max(imp.Count, 1))

	// Bid request for each format in the impression block
	for _, format := range imp.Formats() {
		if fmtres := d.bidImpFormat(request, imp, format); fmtres != nil {
			res = append(res, fmtres...)
		}
	}

	return res
}

func (d *driver) bidImpFormat(request *adtype.BidRequest, imp *adtype.Impression, format *types.Format) []adtype.ResponserItemCommon {
	ads := d.getAdsByFormat(format.ID)
	if len(ads) == 0 {
		return nil
	}

	// Allocate the response items for the ads
	res := make([]adtype.ResponserItemCommon, 0, max(imp.Count, 1))

	// Process each advertisement in the list one by one
	for _, ad := range ads {
		// Skip if Campaign is not matched with the request
		if !ad.Campaign.Test(request) {
			continue
		}

		// Skip if Ad is not matched with the request
		targetBid := ad.TargetBid(request)
		if imp.BidFloor > targetBid.ECPM {
			continue
		}

		// Calculate the price for the view action, and the maximum bid price
		testPrice := false
		viewPrice := gocast.IfThen(ad.PricingModel.IsCPM(), targetBid.Price, 0)
		if ad.TestTestBudgetValue() { // If test budget is active
			if targetBid.Bid != nil && targetBid.Bid.TestPrice > 0 && targetBid.ECPM < targetBid.Bid.TestPrice*1000 {
				viewPrice = targetBid.Bid.TestPrice
				testPrice = true
			} else if targetBid.ECPM < ad.TestViewPrice*1000 {
				viewPrice = ad.TestViewPrice
				testPrice = true
			}
		}
		maxBidPrice := gocast.IfThen(ad.Campaign.MaxBid > 0, ad.Campaign.MaxBid, max(targetBid.BidPrice, targetBid.ECPM/1000))
		maxBidPrice = gocast.IfThen(maxBidPrice > 0, maxBidPrice, viewPrice)
		bidPrice := gocast.IfThen(targetBid.BidPrice > 0, targetBid.BidPrice, viewPrice)
		bidPrice = gocast.IfThen(bidPrice > 0, bidPrice, targetBid.ECPM/1000)
		bidPrice = min(maxBidPrice, bidPrice)

		// Fill the response item with the advertisement data and the price scope
		// For the specific Ad object
		res = append(res, &adresponse.ResponseAdItem{
			Ctx:      request.Ctx,
			ItemID:   rand.UUID(), // Unique ID of the item
			Src:      d,           // Source of the advertisement
			Req:      request,     // Request object
			Imp:      imp,         // Impression object
			Campaign: ad.Campaign, // Campaign object of the advertisement
			Ad:       ad,          // Ad object
			AdBid:    targetBid.Bid,
			PriceScope: adtype.PriceScope{
				TestViewBudget: testPrice,
				MaxBidPrice:    maxBidPrice,
				BidPrice:       bidPrice,
				ViewPrice:      viewPrice,
				ClickPrice:     gocast.IfThen(ad.PricingModel.IsCPC(), targetBid.Price, 0),
				LeadPrice:      gocast.IfThen(ad.PricingModel.IsCPA(), targetBid.Price, 0),
				ECPM:           targetBid.ECPM,
			},
			// BidECPM:     targetBid.ECPM,
			// BidPrice:    targetBid.BidPrice,
			// AdPrice:     targetBid.Price,
			// AdLeadPrice: targetBid.LeadPrice,
		})

		if len(res) >= imp.Count {
			break
		}
	}
	return res
}

func (d *driver) reloadAds() {
	list, _ := d.adCampaigns.CampaignList()
	ads := map[uint64][]*admodels.Ad{}

	// Collect all active ads into ad inmemory storage
	for _, campaign := range list {
		if !campaign.Active() || !campaign.Deleted() {
			continue
		}
		for _, ad := range campaign.Ads {
			if ad.Active() {
				ads[ad.Format.ID] = append(ads[ad.Format.ID], ad)
			}
		}
	}

	// Descending sort of the ads by the weight (ECPM)
	for _, adList := range ads {
		sort.Slice(adList, func(i, j int) bool {
			return adList[i].ECPM() > adList[j].ECPM()
		})
	}

	d.storeAds(ads)
}

//go:inline
func (d *driver) storeAds(ads map[uint64][]*admodels.Ad) {
	d.adStore.Store(ads)
}

//go:inline
func (d *driver) getAds() map[uint64][]*admodels.Ad {
	return d.adStore.Load().(map[uint64][]*admodels.Ad)
}

//go:inline
func (d *driver) getAdsByFormat(formatID uint64) []*admodels.Ad {
	return d.getAds()[formatID]
}
