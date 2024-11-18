//
// @project GeniusRabbit corelib 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package adresponse

import (
	"context"
	"strings"

	"github.com/demdxx/gocast/v2"

	"github.com/bsm/openrtb"
	natresp "github.com/bsm/openrtb/native/response"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// ResponseBidItem value
type ResponseBidItem struct {
	ItemID string

	// Request and impression data
	Src adtype.Source
	Req *adtype.BidRequest
	Imp *adtype.Impression

	// Format of response advertisement item
	FormatType types.FormatType
	RespFormat *types.Format

	// External response data from RTB source
	Bid        *openrtb.Bid
	Native     *natresp.Response
	ActionLink string

	// External bidding price
	BidPrice    billing.Money // Bid price per one action (can be updated by price predictor)
	CPMBidPrice billing.Money // This param can update only price predictor

	// Competitive second AD
	SecondAd adtype.SecondAd

	Data    responseDataAccessor
	assets  admodels.AdAssets
	context context.Context
}

// ID of current response item (unique code of current response)
func (it *ResponseBidItem) ID() string {
	return it.ItemID
}

// Source of response
func (it *ResponseBidItem) Source() adtype.Source {
	return it.Src
}

// NetworkName by source
func (it *ResponseBidItem) NetworkName() string {
	return ""
}

// ContentItemString from the ad
func (it *ResponseBidItem) ContentItemString(name string) string {
	val := it.ContentItem(name)
	if val != nil {
		return gocast.Str(val)
	}
	return ""
}

// ContentItem returns the ad response data
func (it *ResponseBidItem) ContentItem(name string) any {
	if it.Data != nil {
		return it.Data.Get(name)
	}

	formatType := it.PriorityFormatType()

	switch name {
	case adtype.ContentItemContent, adtype.ContentItemIFrameURL:
		if formatType.IsBanner() {
			switch name {
			case adtype.ContentItemIFrameURL:
				if strings.HasPrefix(it.Bid.AdMarkup, "http://") ||
					strings.HasPrefix(it.Bid.AdMarkup, "https://") ||
					(strings.HasPrefix(it.Bid.AdMarkup, "//") && !strings.ContainsAny(it.Bid.AdMarkup, "\n\t")) {
					return it.Bid.AdMarkup
				}
			case adtype.ContentItemContent:
				return it.Bid.AdMarkup
			}
		}
	case adtype.ContentItemLink:
		switch {
		case it.Native != nil:
			return it.Native.Link.URL
		case formatType.IsDirect():
			// In this case here have to be the advertisement link
			return it.Bid.AdMarkup
		}
	case adtype.ContentItemNotifyWinURL:
		if it.Bid != nil {
			return it.Bid.NURL
		}
	case adtype.ContentItemNotifyDisplayURL:
		if it.Bid != nil {
			return it.Bid.BURL
		}
	case types.FormatFieldTitle:
		if it.Native != nil {
			for _, asset := range it.Native.Assets {
				if asset.Title != nil {
					return asset.Title.Text
				}
			}
		}
	default:
		if it.Native != nil {
			for _, asset := range it.Native.Assets {
				if asset.Data != nil && asset.Data.Label == name {
					return asset.Data.Value
				}
			}
		}
	}
	return nil
}

// ContentFields from advertisement object
func (it *ResponseBidItem) ContentFields() map[string]any {
	if it.Format().Config == nil {
		return nil
	}
	fields := map[string]any{}
	config := it.Format().Config
	for _, field := range config.Fields {
		for _, asset := range it.Native.Assets {
			if field.ID != asset.ID {
				continue
			}
			switch {
			case asset.Title != nil:
				fields[field.Name] = asset.Title.Text
			case asset.Link != nil:
				fields[field.Name] = asset.Link.URL
			case asset.Data != nil:
				fields[field.Name] = asset.Data.Value
			}
			break
		}
	}
	return fields
}

// ViewTrackerLinks returns traking links for view action
func (it *ResponseBidItem) ViewTrackerLinks() []string {
	if it.Native == nil {
		return nil
	}
	return it.Native.ImpTrackers
}

// ClickTrackerLinks returns third-party tracker URLs to be fired on click of the URL
func (it *ResponseBidItem) ClickTrackerLinks() []string {
	if it.Native == nil {
		return nil
	}
	return it.Native.Link.ClickTrackers
}

// MainAsset from response
func (it *ResponseBidItem) MainAsset() *admodels.AdAsset {
	mainAsset := it.Format().Config.MainAsset()
	if mainAsset == nil {
		return nil
	}
	for _, asset := range it.Assets() {
		if int(asset.ID) == mainAsset.ID {
			return asset
		}
	}
	return nil
}

// Asset by name
func (it *ResponseBidItem) Asset(name string) *admodels.AdAsset {
	return it.Assets().Asset(name)
}

// Assets list
func (it *ResponseBidItem) Assets() (assets admodels.AdAssets) {
	if it.assets != nil || it.Format().Config == nil {
		return it.assets
	}

	config := it.Format().Config
	for _, configAsset := range config.Assets {
		for _, asset := range it.Native.Assets {
			if asset.ID != configAsset.ID {
				continue
			}
			newAsset := &admodels.AdAsset{
				ID:   uint64(asset.ID),
				Name: configAsset.GetName(),
			}
			switch {
			case asset.Image != nil:
				newAsset.Path = asset.Image.URL
				newAsset.Type = types.AdAssetImageType
				newAsset.ContentType = ""
				newAsset.Width = asset.Image.Width
				newAsset.Height = asset.Image.Height
			// case asset.Video != nil:
			// 	newAsset.Path = asset.Video.URL
			// 	newAsset.Type = models.AdAssetVideoType
			default:
				// TODO error generation
			}
			it.assets = append(it.assets, newAsset)
			break
		}
	}
	return it.assets
}

// Format object model
func (it *ResponseBidItem) Format() *types.Format {
	if it == nil {
		return nil
	}
	return it.RespFormat
}

// PriorityFormatType from current Ad
func (it *ResponseBidItem) PriorityFormatType() types.FormatType {
	if it.FormatType != types.FormatUndefinedType {
		return it.FormatType
	}
	format := it.Imp.FormatTypes
	if formatType := format.HasOneType(); formatType > types.FormatUndefinedType {
		return formatType
	}
	return format.FirstType()
}

// Impression place object
func (it *ResponseBidItem) Impression() *adtype.Impression {
	return it.Imp
}

// ImpressionID unique code string
func (it *ResponseBidItem) ImpressionID() string {
	if it.Imp == nil {
		return ""
	}
	return it.Imp.ID
}

// ExtImpressionID unique code of RTB response
func (it *ResponseBidItem) ExtImpressionID() string {
	if it.Imp == nil {
		return ""
	}
	return it.Imp.ExtID
}

// ExtTargetID of the external network
func (it *ResponseBidItem) ExtTargetID() string {
	return it.Imp.ExtTargetID
}

// AdID number
func (it *ResponseBidItem) AdID() uint64 {
	return 0
}

// AccountID number
func (it *ResponseBidItem) AccountID() uint64 {
	if it.Src != nil {
		type accountIDGetter interface {
			AccountID() uint64
		}
		if src, _ := it.Src.(accountIDGetter); src != nil {
			return src.AccountID()
		}
	}
	return 0
}

// CampaignID number
func (it *ResponseBidItem) CampaignID() uint64 {
	return 0
}

// PricingModel of advertisement
func (it *ResponseBidItem) PricingModel() types.PricingModel {
	return types.PricingModelCPM
}

// ECPM item value
func (it *ResponseBidItem) ECPM() billing.Money {
	if it == nil || it.Bid == nil {
		return 0
	}
	if it.PricingModel().IsCPM() && it.CPMBidPrice > 0 {
		return it.CPMBidPrice
	}
	return billing.MoneyFloat(it.Bid.Price)
}

// Price for specific action if supported `click`, `lead`, `view`
// returns total price of the action
func (it *ResponseBidItem) Price(action admodels.Action, removeFactors ...adtype.PriceFactor) (price billing.Money) {
	if it == nil || it.Bid == nil {
		return 0
	}
	if action.IsImpression() {
		if it.BidPrice > 0 {
			price = it.BidPrice
		} else {
			price = billing.MoneyFloat(it.Bid.Price / 1000)
		}
	}
	return price + adtype.PriceFactorFromList(removeFactors...).RemoveComission(price, it)
}

// SetCPMPrice update of DSP auction value
func (it *ResponseBidItem) SetCPMPrice(price billing.Money, includeFactors ...adtype.PriceFactor) {
	price += adtype.PriceFactorFromList(includeFactors...).AddComission(price, it)
	if it != nil && price < it.ECPM() {
		it.CPMBidPrice = price
	}
}

// CPMPrice value price value for DSP auction
func (it *ResponseBidItem) CPMPrice(removeFactors ...adtype.PriceFactor) (price billing.Money) {
	if it.PricingModel().IsCPM() {
		price = it.Price(admodels.ActionImpression) * 1000
	} else {
		price = it.ECPM()
	}
	// Here we have check for the price from the source wasn't exceeded the ad aficiency price
	if it.CPMBidPrice > 0 && it.CPMBidPrice < price {
		price = it.CPMBidPrice
	}
	return price + adtype.PriceFactorFromList(removeFactors...).RemoveComission(price, it)
}

// AuctionCPMBid value price without any comission
func (it *ResponseBidItem) AuctionCPMBid() billing.Money {
	return it.CPMPrice(adtype.AllPriceFactors)
}

// PurchasePrice gives the price of view from external resource.
// The cost of this request for the system.
func (it *ResponseBidItem) PurchasePrice(action admodels.Action, removeFactors ...adtype.PriceFactor) billing.Money {
	if it == nil {
		return 0
	}
	// Some sources can have the fixed price of buying
	if action.IsImpression() && it.Imp.PurchaseViewPrice > 0 {
		return it.Imp.PurchaseViewPrice
	}
	if it.Imp.Target != nil {
		if pPrice := it.Imp.Target.PurchasePrice(action); pPrice > 0 {
			return pPrice
		}
	}
	if len(removeFactors) == 0 {
		removeFactors = []adtype.PriceFactor{^adtype.TargetReducePriceFactor}
	}
	switch action {
	case admodels.ActionImpression: // Equal to admodels.ActionView
		// As we buying from some source we can consider that we will loose approximately
		// target gate reduce factor percent, but anyway price will be higher for X% of that descepancy
		// to protect system from overspands
		if it.Imp.Target.PricingModel().Or(it.PricingModel()).IsCPM() {
			return it.CPMPrice(removeFactors...) / 1000 // Price per One Impression
		}
	case admodels.ActionClick:
		if it.Imp.Target.PricingModel().Or(it.PricingModel()).IsCPC() {
			return it.Price(action, removeFactors...)
		}
	case admodels.ActionLead:
		if it.Imp.Target.PricingModel().Or(it.PricingModel()).IsCPA() {
			return it.Price(action, removeFactors...)
		}
	}
	return 0
}

// PotentialPrice wich can be received from source but was marked as descrepancy
func (it *ResponseBidItem) PotentialPrice(action admodels.Action) billing.Money {
	return -adtype.SourcePriceFactor.RemoveComission(it.Price(action), it)
}

// Second campaigns
func (it *ResponseBidItem) Second() *adtype.SecondAd {
	return &it.SecondAd
}

// RevenueShareFactor value for the publisher company
func (it *ResponseBidItem) RevenueShareFactor() float64 {
	return it.Imp.RevenueShareFactor()
}

// ComissionShareFactor which system get from publisher 0..1
func (it *ResponseBidItem) ComissionShareFactor() float64 {
	return it.Imp.ComissionShareFactor()
}

// IsDirect AD format
func (it *ResponseBidItem) IsDirect() bool {
	return it.Imp.IsDirect()
}

// ActionURL for direct ADS
func (it *ResponseBidItem) ActionURL() string {
	return it.ActionLink
}

// Validate item
func (it *ResponseBidItem) Validate() error {
	if it.Src == nil || it.Req == nil || it.Imp == nil || it.Bid == nil {
		return adtype.ErrInvalidItemInitialisation
	}
	return it.Bid.Validate()
}

// Width of item
func (it *ResponseBidItem) Width() int {
	return 0
}

// Height of item
func (it *ResponseBidItem) Height() int {
	return 0
}

// Markup advertisement
func (it *ResponseBidItem) Markup() (string, error) {
	return "", nil
}

// Context value
func (it *ResponseBidItem) Context(ctx ...context.Context) context.Context {
	if len(ctx) > 0 {
		it.context = ctx[0]
	}
	return it.context
}

// Get ext field
func (it *ResponseBidItem) Get(key string) (res any) {
	if it.context == nil {
		return res
	}
	return it.context.Value(key)
}

var (
	_ adtype.ResponserItem = &ResponseBidItem{}
)
