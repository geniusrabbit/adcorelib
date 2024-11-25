//
// @project GeniusRabbit corelib 2016 – 2019, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019, 2024
//

package adresponse

import (
	"context"
	"strconv"
	"strings"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// ResponseAdItem represents an advertisement item selected from storage.
type ResponseAdItem struct {
	Ctx    context.Context `json:"-"`
	ItemID string          `json:"id"`

	Src adtype.Source      `json:"source,omitempty"`
	Req *adtype.BidRequest `json:"request,omitempty"`
	Imp *adtype.Impression `json:"impression,omitempty"` // Unique impression

	Campaign *admodels.Campaign `json:"campaign,omitempty"`
	Ad       *admodels.Ad       `json:"ad,omitempty"`
	AdBid    *admodels.AdBid    `json:"ad_bid,omitempty"`
	AdLink   admodels.AdLink    `json:"ad_link,omitempty"`

	PriceScope adtype.PriceScope `json:"price_scope,omitempty"`

	// BidECPM billing.Money `json:"bid_ecpm,omitempty"` // Bid's effective CPM
	// BidPrice billing.Money `json:"bid_price,omitempty"` // Max RTB bid price (CPM only)
	// AdPrice     billing.Money   `json:"price,omitempty"`      // New price of advertisement target action (click, lead, impression)
	// AdLeadPrice billing.Money   `json:"lead_price,omitempty"` // Lead price for the ad
	// CPMBidPrice billing.Money   `json:"cpm_bid,omitempty"`   // Updated only by price predictor
	SecondAd adtype.SecondAd `json:"second_ad,omitempty"` // Secondary ad information
}

// ID returns the unique identifier of the current response item.
func (it *ResponseAdItem) ID() string {
	return it.ItemID
}

// AuctionID returns the auction identifier from the request.
func (it *ResponseAdItem) AuctionID() string {
	return it.Req.ID
}

// Validate checks the validity of the item.
func (it *ResponseAdItem) Validate() error {
	return nil
}

// Impression returns the impression associated with the response item.
func (it *ResponseAdItem) Impression() *adtype.Impression {
	return it.Imp
}

// ImpressionID returns the unique identifier of the impression.
func (it *ResponseAdItem) ImpressionID() string {
	if it.Imp == nil {
		return ""
	}
	return it.Imp.ID
}

// ExtImpressionID returns the external impression identifier.
func (it *ResponseAdItem) ExtImpressionID() string {
	if it.Imp == nil {
		return ""
	}
	return it.Imp.ExtID
}

// ExtTargetID returns the external target identifier.
func (it *ResponseAdItem) ExtTargetID() string {
	if it.Imp == nil {
		return ""
	}
	return it.Imp.ExtTargetID
}

// Source returns the source of the response.
func (it *ResponseAdItem) Source() adtype.Source {
	return it.Src
}

// NetworkName returns the network name associated with the source.
func (it *ResponseAdItem) NetworkName() string {
	return ""
}

// PriorityFormatType returns the primary format type from the current ad.
func (it *ResponseAdItem) PriorityFormatType() types.FormatType {
	// If it's the only format type
	if formatType := it.Ad.Format.Types.HasOneType(); formatType > types.FormatUndefinedType {
		return formatType
	}

	intersection := it.Imp.FormatTypes.Intersec(it.Ad.Format.Types)
	if formatType := intersection.HasOneType(); formatType > types.FormatUndefinedType {
		return formatType
	}

	return intersection.FirstType()
}

// Second returns the secondary advertisement information.
func (it *ResponseAdItem) Second() *adtype.SecondAd {
	return &it.SecondAd
}

// Request returns the bid request information.
func (it *ResponseAdItem) Request() *adtype.BidRequest {
	return it.Req
}

// AdDirectLink returns the direct link of the advertisement.
func (it *ResponseAdItem) AdDirectLink() string {
	if it.AdLink.ID == 0 {
		it.AdLink = it.Ad.RandomAdLink()
	}
	return it.AdLink.Link
}

// ContentItemString returns the content item as a string from the advertisement.
func (it *ResponseAdItem) ContentItemString(name string) string {
	switch name {
	case adtype.ContentItemLink:
		if !it.Ad.Format.IsProxy() {
			return it.processParameters(it.AdDirectLink())
		}
	case adtype.ContentItemIFrameURL:
		if it.Ad.Format.IsProxy() {
			return it.processParameters(it.Ad.ProxyURL())
		}
	}
	return it.processParameters(it.Ad.ContentItemString(name))
}

// ContentItem returns the ad response data for the specified content item.
func (it *ResponseAdItem) ContentItem(name string) any {
	switch name {
	case adtype.ContentItemLink:
		if !it.Ad.Format.IsProxy() {
			return it.processParameters(it.AdDirectLink())
		}
	case adtype.ContentItemIFrameURL:
		if it.Ad.Format.IsProxy() {
			return it.processParameters(it.Ad.ProxyURL())
		}
	}
	item := it.Ad.ContentItem(name)
	switch s := item.(type) {
	case string:
		return it.processParameters(s)
	}
	return item
}

// ContentFields returns the content fields from the advertisement object.
func (it *ResponseAdItem) ContentFields() map[string]any {
	return it.Ad.Content
}

// ViewTrackerLinks returns tracking links for view action.
func (it *ResponseAdItem) ViewTrackerLinks() []string {
	return nil
}

// ClickTrackerLinks returns third-party tracker URLs to be fired on click of the URL.
func (it *ResponseAdItem) ClickTrackerLinks() []string {
	return nil
}

// MainAsset returns the main asset from the advertisement.
func (it *ResponseAdItem) MainAsset() *admodels.AdAsset {
	return it.Ad.MainAsset()
}

// Assets returns the list of assets from the advertisement.
func (it *ResponseAdItem) Assets() admodels.AdAssets {
	return it.Ad.Assets
}

// Width returns the width of the ad.
func (it *ResponseAdItem) Width() int {
	if it.Imp.Width <= 0 {
		return it.Imp.WidthMax
	}
	return it.Imp.Width
}

// Height returns the height of the ad.
func (it *ResponseAdItem) Height() int {
	if it.Imp.Height <= 0 {
		return it.Imp.HeightMax
	}
	return it.Imp.Height
}

// TargetID returns the target identifier.
func (it *ResponseAdItem) TargetID() uint64 {
	if it == nil || it.Imp == nil || it.Imp.Target == nil {
		return 0
	}
	return it.Imp.Target.ID()
}

// TargetIDString returns the target identifier as a string.
func (it *ResponseAdItem) TargetIDString() string {
	return strconv.FormatUint(it.TargetID(), 10)
}

// AdID returns the advertisement identifier.
func (it *ResponseAdItem) AdID() uint64 {
	if it == nil || it.Ad == nil {
		return 0
	}
	return it.Ad.ID
}

// AdIDString returns the advertisement identifier as a string.
func (it *ResponseAdItem) AdIDString() string {
	return strconv.FormatUint(it.AdID(), 10)
}

// ProjectID returns the project identifier.
func (it *ResponseAdItem) ProjectID() uint64 {
	if it == nil || it.Campaign == nil {
		return 0
	}
	return it.Campaign.ProjectID()
}

// AccountID returns the account identifier.
func (it *ResponseAdItem) AccountID() uint64 {
	if it == nil || it.Campaign == nil {
		return 0
	}
	return it.Campaign.AccountID()
}

// CampaignIDString returns the campaign identifier as a string.
func (it *ResponseAdItem) CampaignIDString() string {
	return strconv.FormatUint(it.Campaign.ID(), 10)
}

// CampaignID returns the campaign identifier.
func (it *ResponseAdItem) CampaignID() uint64 {
	if it == nil || it.Campaign == nil {
		return 0
	}
	return it.Campaign.ID()
}

// CreativeIDString returns the creative identifier as a string for reporting content issues or defects.
func (it *ResponseAdItem) CreativeIDString() string {
	return it.AdIDString()
}

// Format returns the format object model.
func (it *ResponseAdItem) Format() *types.Format {
	if it == nil || it.Ad == nil {
		return nil
	}
	return it.Ad.Format
}

// PricingModel returns the pricing model of the advertisement.
func (it *ResponseAdItem) PricingModel() types.PricingModel {
	if it == nil || it.Ad == nil {
		return types.PricingModelUndefined
	}
	return it.Ad.PricingModel
}

// IsDirect returns true if the response item is direct.
func (it *ResponseAdItem) IsDirect() bool {
	return it != nil && it.Ad.Format.IsDirect()
}

// ActionURL returns the action URL for direct banners.
func (it *ResponseAdItem) ActionURL() string {
	return it.processParameters(it.AdDirectLink())
}

///////////////////////////////////////////////////////////////////////////////
// Price calculation methods
///////////////////////////////////////////////////////////////////////////////

// ECPM returns the effective cost per mille of the item.
func (it *ResponseAdItem) ECPM() billing.Money {
	if it == nil {
		return 0
	}
	// if it.PricingModel().IsCPM() {
	// 	if it.CPMBidPrice > 0 {
	// 		return it.CPMBidPrice
	// 	}
	// }
	// return it.BidECPM
	return it.PriceScope.ECPM
}

// Price returns the total price for a specific action, if supported (click, lead, view).
func (it *ResponseAdItem) Price(action admodels.Action, removeFactors ...adtype.PriceFactor) billing.Money {
	if it == nil || it.Ad == nil {
		return 0
	}
	// if price = it.AdPrice; price <= 0 {
	// 	price = it.Ad.Price
	// }
	// switch action {
	// case admodels.ActionView:
	// 	if !it.PricingModel().IsCPM() {
	// 		price = 0
	// 	}
	// 	if it.BidPrice > 0 {
	// 		price = it.BidPrice
	// 	}
	// case admodels.ActionClick:
	// 	if !it.PricingModel().IsCPC() {
	// 		price = 0
	// 	}
	// case admodels.ActionLead:
	// 	if it.AdLeadPrice > 0 {
	// 		price = it.AdLeadPrice
	// 	} else {
	// 		price = it.Ad.LeadPrice
	// 	}
	// }
	price := it.PriceScope.PricePerAction(action)
	price += adtype.PriceFactorFromList(removeFactors...).RemoveComission(price, it)
	return price
}

// InternalAuctionCPMBid value provides maximal possible price without any comission
// According to this value the system can choice the best item for the auction
func (it *ResponseAdItem) InternalAuctionCPMBid() billing.Money {
	return it.AuctionCPMBid(adtype.AllPriceFactors)
}

// PurchasePrice gives the price of view from external resource.
func (it *ResponseAdItem) PurchasePrice(action admodels.Action, removeFactors ...adtype.PriceFactor) billing.Money {
	if it == nil {
		return 0
	}
	// Some sources can have the fixed price of buying
	if pPrice := it.Imp.PurchasePrice(action); pPrice > 0 {
		return pPrice
	}
	if len(removeFactors) == 0 {
		removeFactors = []adtype.PriceFactor{^adtype.TargetReducePriceFactor}
	}
	switch action {
	case admodels.ActionImpression:
		if it.Imp.Target.PricingModel().Or(it.PricingModel()).IsCPM() {
			return it.AuctionCPMBid(removeFactors...) / 1000 // Price per one impression
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

// PotentialPrice returns the price which can be received from source but was marked as discrepancy.
func (it *ResponseAdItem) PotentialPrice(action admodels.Action) billing.Money {
	return -adtype.SourcePriceFactor.RemoveComission(it.Price(action), it)
}

// SetAuctionCPMBid value for external sources auction the system will pay
func (it *ResponseAdItem) SetAuctionCPMBid(price billing.Money, includeFactors ...adtype.PriceFactor) error {
	if len(includeFactors) > 0 {
		price += adtype.PriceFactorFromList(includeFactors...).AddComission(price, it)
	}
	if !it.PriceScope.SetBidPrice(price/1000, false) {
		return adtype.ErrNewAuctionBidIsHigherThenMaxBid
	}
	return nil
}

// AuctionCPMBid value provides price for external sources
// The prive what we can pay for the action to the external source
func (it *ResponseAdItem) AuctionCPMBid(removeFactors ...adtype.PriceFactor) billing.Money {
	price := it.PriceScope.BidPrice * 1000

	// Remove commissions from the price if provided
	if len(removeFactors) > 0 {
		price += adtype.PriceFactorFromList(removeFactors...).RemoveComission(price, it)
	}
	return price
}

///////////////////////////////////////////////////////////////////////////////
// Revenue share/comission methods
///////////////////////////////////////////////////////////////////////////////

// RevenueShareFactor returns the revenue share percentage.
func (it *ResponseAdItem) RevenueShareFactor() float64 {
	return it.Imp.RevenueShareFactor()
}

// ComissionShareFactor returns the commission share percentage which system gets from publisher.
func (it *ResponseAdItem) ComissionShareFactor() float64 {
	return it.Imp.ComissionShareFactor()
}

///////////////////////////////////////////////////////////////////////////////
// Context methods
///////////////////////////////////////////////////////////////////////////////

// Context returns or sets the context value.
func (it *ResponseAdItem) Context(ctx ...context.Context) (c context.Context) {
	c = it.Ctx
	if len(ctx) > 0 {
		it.Ctx = ctx[0]
	}
	return c
}

// Get retrieves the value associated with the key from the context.
func (it *ResponseAdItem) Get(key string) any {
	if it.Ctx == nil {
		return nil
	}
	return it.Ctx.Value(key)
}

// Release the response item and all linked objects.
func (it *ResponseAdItem) Release() {
	if it != nil {
		it.reset()
	}
}

func (it *ResponseAdItem) reset() {
	*it = ResponseAdItem{}
}

func (it *ResponseAdItem) processParameters(s string) string {
	if strings.Contains(s, "${click_id}") {
		return strings.NewReplacer(
			// "${click_id}", it.LeadCode().URLEncode().String(),
			"${imp_id}", it.ImpressionID(),
			"${domain}", it.Request().DomainName(),
			"${zone}", it.TargetIDString(),
			"${country_code}", it.Request().GeoCode(),
			"${language}", it.Request().BrowserInfo().PrimaryLanguage,
			"${format}", it.Format().Codename,
			"${width}", strconv.Itoa(it.Width()),
			"${height}", strconv.Itoa(it.Height()),
		).Replace(s)
	}

	return strings.NewReplacer(
		"${imp_id}", it.ImpressionID(),
		"${domain}", it.Request().DomainName(),
		"${zone}", it.TargetIDString(),
		"${country_code}", it.Request().GeoCode(),
		"${language}", it.Request().BrowserInfo().PrimaryLanguage,
		"${format}", it.Format().Codename,
		"${width}", strconv.Itoa(it.Width()),
		"${height}", strconv.Itoa(it.Height()),
	).Replace(s)
}

var (
	_ adtype.ResponserItem = &ResponseAdItem{}
)
