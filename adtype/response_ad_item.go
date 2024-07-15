//
// @project GeniusRabbit corelib 2016 – 2019, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019, 2024
//

package adtype

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

// ResponseAdItem for select from storage
type ResponseAdItem struct {
	Ctx    context.Context `json:"-"`
	ItemID string          `json:"id"`

	Src Source      `json:"source,omitempty"`
	Req *BidRequest `json:"request,omitempty"`
	Imp *Impression `json:"impression,omitempty"` // Impression Unique

	Campaign *admodels.Campaign `json:"campaign,omitempty"`
	Ad       *admodels.Ad       `json:"ad,omitempty"`
	AdBid    *admodels.AdBid    `json:"ad_bid,omitempty"`
	AdLink   admodels.AdLink    `json:"ad_link,omitempty"`

	BidECPM     billing.Money `json:"bid_ecpm,omitempty"`   //
	BidPrice    billing.Money `json:"bid_price,omitempty"`  // Max RTB bid price (CPM only)
	AdPrice     billing.Money `json:"price,omitempty"`      // New price of advertisement target action (click, lead, impression)
	AdLeadPrice billing.Money `json:"lead_price,omitempty"` //
	CPMBidPrice billing.Money `json:"cpm_bid,omitempty"`    // This param can update only price predictor
	SecondAd    SecondAd      `json:"second_ad,omitempty"`  //
}

// ID of current response item (unique code of current response)
func (it *ResponseAdItem) ID() string {
	return it.ItemID
}

// AuctionID response
func (it *ResponseAdItem) AuctionID() string {
	return it.Req.ID
}

// Impression place object
func (it *ResponseAdItem) Impression() *Impression {
	return it.Imp
}

// ImpressionID unique code string
func (it *ResponseAdItem) ImpressionID() string {
	if it.Imp == nil {
		return ""
	}
	return it.Imp.ID
}

// ExtImpressionID unique code of RTB response
func (it *ResponseAdItem) ExtImpressionID() string {
	if it.Imp == nil {
		return ""
	}
	return it.Imp.ExtID
}

// ExtTargetID of the external network
func (it *ResponseAdItem) ExtTargetID() string {
	if it.Imp == nil {
		return ""
	}
	return it.Imp.ExtTargetID
}

// Source of response
func (it *ResponseAdItem) Source() Source {
	return it.Src
}

// NetworkName by source
func (it *ResponseAdItem) NetworkName() string {
	return ""
}

// PriorityFormatType from current Ad
func (it *ResponseAdItem) PriorityFormatType() types.FormatType {
	// If it's the only one format type
	if formatType := it.Ad.Format.Types.HasOneType(); formatType > types.FormatUndefinedType {
		return formatType
	}

	intersection := it.Imp.FormatTypes.Intersec(it.Ad.Format.Types)
	if formatType := intersection.HasOneType(); formatType > types.FormatUndefinedType {
		return formatType
	}

	return intersection.FirstType()
}

// Second campaigns
func (it *ResponseAdItem) Second() *SecondAd {
	return &it.SecondAd
}

// Request information
func (it *ResponseAdItem) Request() *BidRequest {
	return it.Req
}

// AdDirectLink of the ad
func (it *ResponseAdItem) AdDirectLink() string {
	if it.AdLink.ID == 0 {
		it.AdLink = it.Ad.RandomAdLink()
	}
	return it.AdLink.Link
}

// ContentItemString from the ad
func (it *ResponseAdItem) ContentItemString(name string) string {
	switch name {
	case ContentItemLink:
		if !it.Ad.Format.IsProxy() {
			return it.processParameters(it.AdDirectLink())
		}
	case ContentItemIFrameURL:
		if it.Ad.Format.IsProxy() {
			return it.processParameters(it.Ad.ProxyURL())
		}
	}
	return it.processParameters(it.Ad.ContentItemString(name))
}

// ContentItem returns the ad response data
func (it *ResponseAdItem) ContentItem(name string) any {
	switch name {
	case ContentItemLink:
		if !it.Ad.Format.IsProxy() {
			return it.processParameters(it.AdDirectLink())
		}
	case ContentItemIFrameURL:
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

// ContentFields from advertisement object
func (it *ResponseAdItem) ContentFields() map[string]any {
	return it.Ad.Content
}

// ViewTrackerLinks returns traking links for view action
func (it *ResponseAdItem) ViewTrackerLinks() []string {
	return nil
}

// ClickTrackerLinks returns third-party tracker URLs to be fired on click of the URL
func (it *ResponseAdItem) ClickTrackerLinks() []string {
	return nil
}

// MainAsset from response
func (it *ResponseAdItem) MainAsset() *admodels.AdAsset {
	return it.Ad.MainAsset()
}

// Asset by name
func (it *ResponseAdItem) Asset(name string) *admodels.AdAsset {
	return it.Ad.Asset(name)
}

// Assets list
func (it *ResponseAdItem) Assets() admodels.AdAssets {
	return it.Ad.Assets
}

// Width of AD
func (it *ResponseAdItem) Width() int {
	if it.Imp.W <= 0 {
		return it.Imp.WMax
	}
	return it.Imp.W
}

// Height of AD
func (it *ResponseAdItem) Height() int {
	if it.Imp.H <= 0 {
		return it.Imp.HMax
	}
	return it.Imp.H
}

// TargetID number
func (it *ResponseAdItem) TargetID() uint64 {
	if it == nil || it.Imp == nil || it.Imp.Target == nil {
		return 0
	}
	return it.Imp.Target.ID()
}

// TargetIDString number
func (it *ResponseAdItem) TargetIDString() string {
	return strconv.FormatInt(int64(it.TargetID()), 10)
}

// AdID number
func (it *ResponseAdItem) AdID() uint64 {
	if it == nil || it.Ad == nil {
		return 0
	}
	return it.Ad.ID
}

// AdIDString References the ad to be served if the bid wins.
func (it *ResponseAdItem) AdIDString() string {
	return strconv.FormatUint(it.AdID(), 10)
}

// ProjectID number
func (it *ResponseAdItem) ProjectID() uint64 {
	if it == nil || it.Campaign == nil {
		return 0
	}
	return it.Campaign.ProjectID()
}

// CompanyID number
func (it *ResponseAdItem) CompanyID() uint64 {
	if it == nil || it.Campaign == nil {
		return 0
	}
	return it.Campaign.CompanyID
}

// CampaignIDString that appears with the Ad markup.
func (it *ResponseAdItem) CampaignIDString() string {
	return fmt.Sprintf("%d", it.Campaign.ID)
}

// CampaignID number
func (it *ResponseAdItem) CampaignID() uint64 {
	if it == nil || it.Campaign == nil {
		return 0
	}
	return it.Campaign.ID
}

// CreativeIDString for reporting content issues or defects.
// This could also be used as a reference to a creative ID
// that is posted with an exchange.
func (it *ResponseAdItem) CreativeIDString() string {
	return it.AdIDString()
}

// Format object model
func (it *ResponseAdItem) Format() *types.Format {
	if it == nil || it.Ad == nil {
		return nil
	}
	return it.Ad.Format
}

// // LeadCode value
// func (it *ResponseAdItem) LeadCode() events.Code {
// 	return (&events.LeadCode{
// 		AuctionID:  it.AuctionID(),
// 		ImpAdID:    it.ID(),
// 		SourceID:   0,
// 		ProjectID:  0,
// 		CampaignID: it.CampaignID(),
// 		AdID:       it.AdID(),
// 		Price:      it.Price(models.ActionLead).Int64(),
// 		Timestamp:  time.Now().Unix(),
// 	}).Pack().Compress()
// }

// PricingModel of advertisement
func (it *ResponseAdItem) PricingModel() types.PricingModel {
	if it == nil || it.Ad == nil {
		return types.PricingModelUndefined
	}
	return it.Ad.PricingModel
}

// IsDirect response item
func (it *ResponseAdItem) IsDirect() bool {
	return it != nil && it.Ad.Format.IsDirect()
}

// ActionURL for direct baners
func (it *ResponseAdItem) ActionURL() string {
	return it.processParameters(it.AdDirectLink())
}

// ECPM item value
func (it *ResponseAdItem) ECPM() billing.Money {
	if it == nil {
		return 0
	}
	if it.PricingModel().IsCPM() {
		if it.CPMBidPrice > 0 {
			return it.CPMBidPrice
		}
	}
	return it.BidECPM
}

// Price for specific action if supported `click`, `lead`, `view`
// returns total price of the action
func (it *ResponseAdItem) Price(action admodels.Action, removeFactors ...PriceFactor) (price billing.Money) {
	if it == nil || it.Ad == nil {
		return 0
	}
	if price = it.AdPrice; price <= 0 {
		price = it.Ad.Price
	}
	switch action {
	case admodels.ActionImpression:
		if !it.PricingModel().IsCPM() {
			price = 0
		}
		if it.BidPrice > 0 {
			price = it.BidPrice
		}
	case admodels.ActionClick:
		if !it.PricingModel().IsCPC() {
			price = 0
		}
	case admodels.ActionLead:
		if it.AdLeadPrice > 0 {
			price = it.AdLeadPrice
		} else {
			price = it.Ad.LeadPrice
		}
	}
	price += PriceFactorFromList(removeFactors...).Remove(price, it)
	return price
}

// SetCPMPrice update of DSP auction value
func (it *ResponseAdItem) SetCPMPrice(price billing.Money, includeFactors ...PriceFactor) {
	price += PriceFactorFromList(includeFactors...).Add(price, it)
	if price < it.ECPM() || price < it.Ad.BidPrice {
		it.CPMBidPrice = price
	}
}

// CPMPrice value price value for DSP auction
func (it *ResponseAdItem) CPMPrice(removeFactors ...PriceFactor) (price billing.Money) {
	if it.CPMBidPrice > 0 {
		price = it.CPMBidPrice
	} else if it.PricingModel().IsCPM() {
		price = it.Price(admodels.ActionImpression) * 1000
	} else {
		price = it.ECPM()
	}

	price = it.prepareMaxBidPrice(price, true)

	// Remove system commision from the price
	price += PriceFactorFromList(removeFactors...).Remove(price, it)
	return price
}

// AuctionCPMBid value price without any comission
// Can be replaced on comission only
func (it *ResponseAdItem) AuctionCPMBid() billing.Money {
	return it.CPMPrice(AllPriceFactors)
}

// PurchasePrice gives the price of view from external resource.
// The cost of this request.
func (it *ResponseAdItem) PurchasePrice(action admodels.Action, removeFactors ...PriceFactor) billing.Money {
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
		removeFactors = []PriceFactor{^TargetReducePriceFactor}
	}
	switch action {
	case admodels.ActionImpression:
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
func (it *ResponseAdItem) PotentialPrice(action admodels.Action) billing.Money {
	return -SourcePriceFactor.Remove(it.Price(action), it)
}

// Validate item
func (it *ResponseAdItem) Validate() error {
	return nil
}

// RevenueShareFactor amount %
func (it *ResponseAdItem) RevenueShareFactor() float64 {
	return it.Imp.RevenueShareFactor()
}

// ComissionShareFactor which system get from publisher
func (it *ResponseAdItem) ComissionShareFactor() float64 {
	return it.Imp.ComissionShareFactor()
}

// Context value
func (it *ResponseAdItem) Context(ctx ...context.Context) (c context.Context) {
	c = it.Ctx
	if len(ctx) > 0 {
		it.Ctx = ctx[0]
	}
	return c
}

// Get ext field
func (it *ResponseAdItem) Get(key string) any {
	if it.Ctx == nil {
		return nil
	}
	return it.Ctx.Value(key)
}

func (it *ResponseAdItem) prepareMaxBidPrice(price billing.Money, maxIfZero bool) billing.Money {
	switch {
	case it.BidPrice > 0:
		if price > it.BidPrice || (maxIfZero && price <= 0) {
			return it.BidPrice
		}
	case it.Ad.BidPrice > 0 && (price > it.Ad.BidPrice || (maxIfZero && price <= 0)):
		return it.Ad.BidPrice
	}
	return price
}

func (it *ResponseAdItem) reset() {
	*it = ResponseAdItem{}
}

// // ResetAdSpent values
// func (it *ResponseAdItem) ResetAdSpent(amount billing.Money) {
// 	it.Ad.State.SetSpent(amount)
// 	it.Ad.Campaign.UpdateBalance()
// 	it.Ad.Campaign.Company.Spent = it.Ad.Campaign.GetSpent()
// }

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
	_ ResponserItem = &ResponseAdItem{}
)
