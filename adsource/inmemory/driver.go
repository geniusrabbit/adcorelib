package inmemory

import (
	"sync/atomic"

	"github.com/demdxx/rpool/v2"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/auction"
	"github.com/geniusrabbit/adcorelib/fasttime"
	"github.com/geniusrabbit/adcorelib/rand"
	"github.com/geniusrabbit/adcorelib/searchtypes"
	"github.com/geniusrabbit/adstorage/accessors/campaignaccessor"
)

type driver struct {
	execPool *rpool.Pool

	adCampaigns *campaignaccessor.CampaignAccessor
	adStore     atomic.Value
}

func (d *driver) init() {
	d.execPool = rpool.NewSinglePool()
	d.reloadAds()
}

// ID of the source driver
func (d *driver) ID() uint64 { return 0 }

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
	if d.adCampaigns.NeedUpdate() {
		d.execPool.Go(d.reloadAds)
	}
	var referee auction.Referee
	for i := 0; i < len(request.Imps); i++ {
		imp := &request.Imps[i]
		if ires := d.bidImp(request, imp); ires != nil {
			referee.Push(ires...)
		}
	}
	itemResp := referee.MatchRequest(request)
	if len(itemResp) == 0 {
		return adtype.NewEmptyResponse(request, d, nil)
	}
	return adtype.NewResponse(request, d, itemResp, nil)
}

// ProcessResponseItem result or error
func (d *driver) ProcessResponseItem(resp adtype.Responser, item adtype.ResponserItem) {
	// Reduce inmemory counters
	adItem := item.(*adtype.ResponseAdItem)
	adItem.Ad.Campaign.Budget -= adItem.BidECPM / 1000
	adItem.Ad.Campaign.DailyBudget -= adItem.BidECPM / 1000
}

func (d *driver) bidImp(request *adtype.BidRequest, imp *adtype.Impression) []adtype.ResponserItemCommon {
	res := make([]adtype.ResponserItemCommon, 0, len(imp.Formats()))
	for _, format := range imp.Formats() {
		if fmtres := d.bidImpFormat(request, imp, format); fmtres != nil {
			res = append(res, fmtres...)
		}
	}
	return res
}

func (d *driver) bidImpFormat(request *adtype.BidRequest, imp *adtype.Impression, format *types.Format) []adtype.ResponserItemCommon {
	ads := d.getAds()[format.ID]
	if len(ads) == 0 {
		return nil
	}
	var (
		filter searchtypes.Bitset64
		res    = make([]adtype.ResponserItemCommon, 0, imp.Count)
	)
	for _, ad := range ads {
		if filter.Has(uint(ad.ID)) {
			continue
		}
		if imp.BidFloor > ad.TargetBid(request).ECPM ||
			!ad.Campaign.Hours.TestTime(fasttime.Now()) ||
			(ad.Campaign.Keywords.Len() > 0 && ad.Campaign.Keywords.OneOf(request.Keywords())) ||
			(ad.Campaign.Zones.Len() > 0 && ad.Campaign.Zones.IndexOf(request.TargetID()) < 0) ||
			(ad.Campaign.Domains.Len() > 0 && ad.Campaign.Domains.OneOf(request.Domain())) ||
			(ad.Campaign.Sex.Len() > 0 && ad.Campaign.Sex.IndexOf(request.Sex()) < 0) ||
			(ad.Campaign.Age.Len() > 0 && ad.Campaign.Age.IndexOf(request.Age()) < 0) || // TODO range processing 0-10 years, 10-20, 20-25 & etc.
			(ad.Campaign.Categories.Len() > 0 && ad.Campaign.Categories.IndexOf(request.GeoID()) < 0) ||
			(ad.Campaign.Cities.Len() > 0 && ad.Campaign.Cities.IndexOf(request.City()) < 0) ||
			(ad.Campaign.Countries.Len() > 0 && ad.Campaign.Countries.IndexOf(request.GeoID()) < 0) ||
			(ad.Campaign.Languages.Len() > 0 && ad.Campaign.Languages.IndexOf(request.LanguageID()) < 0) ||
			(ad.Campaign.Browsers.Len() > 0 && ad.Campaign.Browsers.IndexOf(request.BrowserID()) < 0) ||
			(ad.Campaign.Os.Len() > 0 && ad.Campaign.Os.IndexOf(request.OSID()) < 0) ||
			(ad.Campaign.DeviceTypes.Len() > 0 && ad.Campaign.DeviceTypes.IndexOf(request.DeviceType()) < 0) ||
			(ad.Campaign.Devices.Len() > 0 && ad.Campaign.Devices.IndexOf(request.DeviceID()) < 0) {
			continue // Skip if targeting not matched
		}
		filter.Set(uint(ad.ID))
		targetBid := ad.TargetBid(request)
		res = append(res, &adtype.ResponseAdItem{
			Ctx:         request.Ctx,
			ItemID:      rand.UUID(),
			Src:         d,
			Req:         request,
			Imp:         imp,
			Campaign:    ad.Campaign,
			Ad:          ad,
			AdBid:       targetBid.Bid,
			BidECPM:     targetBid.ECPM,
			BidPrice:    targetBid.BidPrice,
			AdPrice:     targetBid.Price,
			AdLeadPrice: targetBid.LeadPrice,
			SecondAd:    adtype.SecondAd{}, // Competitor advertisement
		})
	}
	return res
}

func (d *driver) reloadAds() {
	list, _ := d.adCampaigns.CampaignList()
	ads := map[uint64][]*admodels.Ad{}

	for _, campaign := range list {
		for _, ad := range campaign.Ads {
			ads[ad.Format.ID] = append(ads[ad.Format.ID], ad)
		}
	}

	d.storeAds(ads)
}

func (d *driver) storeAds(ads map[uint64][]*admodels.Ad) {
	d.adStore.Store(ads)
}

func (d *driver) getAds() map[uint64][]*admodels.Ad {
	return d.adStore.Load().(map[uint64][]*admodels.Ad)
}
