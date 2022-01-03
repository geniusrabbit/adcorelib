//
// @project GeniusRabbit rotator 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package adtype

import (
	"context"

	"geniusrabbit.dev/corelib/admodels"
	"geniusrabbit.dev/corelib/admodels/types"
	"geniusrabbit.dev/corelib/billing"
)

// Content item names
const (
	ContentItemLink             = "link"
	ContentItemContent          = "content"
	ContentItemIFrameURL        = "iframe_url"
	ContentItemNotifyWinURL     = "notify_win_url"
	ContentItemNotifyDisplayURL = "notify_display_url"
)

// ResponserItemCommon interface
type ResponserItemCommon interface {
	// ID of current response item (unique code of current response)
	ID() string

	// Impression place object
	Impression() *Impression

	// ImpressionID unique code string
	ImpressionID() string

	// ExtImpressionID it's unique code of the auction bid impression
	ExtImpressionID() string

	// ExtTargetID of the external network
	ExtTargetID() string

	// AuctionCPMBid value price without any comission
	AuctionCPMBid() billing.Money

	// PriorityFormatType from current Ad
	PriorityFormatType() types.FormatType

	// Validate item
	Validate() error

	// Context value
	Context(ctx ...context.Context) context.Context

	// Get ext field
	Get(key string) interface{}
}

// ResponserItem for single AD
type ResponserItem interface {
	ResponserItemCommon

	// AdID number
	AdID() uint64

	// CompanyID number
	CompanyID() uint64

	// CampaignID number
	CampaignID() uint64

	// Format object
	Format() *types.Format

	// PricingModel of advertisement
	PricingModel() types.PricingModel

	// ContentItem returns the ad response data
	ContentItem(name string) interface{}

	// ContentItemString from the ad
	ContentItemString(name string) string

	// ContentFields from advertisement object
	ContentFields() map[string]interface{}

	// MainAsset from response
	MainAsset() *admodels.AdFile

	// Asset by name
	Asset(name string) *admodels.AdFile

	// Assets list
	Assets() []admodels.AdFile

	// Source of response
	Source() Source

	// ViewTrackerLinks returns traking links for view action
	ViewTrackerLinks() []string

	// ClickTrackerLinks returns traking links for click action
	ClickTrackerLinks() []string

	// NetworkName by source
	NetworkName() string

	// Price for specific action if supported `click`, `lead`, `view`
	Price(action admodels.Action) billing.Money

	// SetAuctionCPMPrice update of DSP auction value
	SetAuctionCPMPrice(price billing.Money, includeFactors ...bool)

	// CPMPrice value price value for DSP auction
	AuctionCPMPrice(removeFactors ...bool) billing.Money

	// Second campaigns
	Second() *SecondAd

	// Revenue money (in percents)
	Revenue() float64

	// Potential money (in percents)
	Potential() float64

	// ECPM item value
	ECPM() billing.Money

	// RevenueShareFactor returns the multipler for company
	// revenue calculation per action from 0 to 1
	RevenueShareFactor() float64

	// ComissionShareFactor returns the multipler for commission
	// calculation which system get from user revenue from 0 to 1
	ComissionShareFactor() float64

	// IsDirect AD type
	IsDirect() bool

	// ActionURL returns target resource link for direct and banner click as well
	ActionURL() string

	// Width of item
	Width() int

	// Height of item
	Height() int
}

// ResponserMultipleItem interface for complex banners
type ResponserMultipleItem interface {
	ResponserItemCommon

	// Ads list response
	Ads() []ResponserItem

	// Count of response items
	Count() int
}
