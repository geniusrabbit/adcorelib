//
// @project GeniusRabbit corelib 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package adtype

import (
	"context"
	"errors"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

// Content item names
const (
	ContentItemLink             = "link"
	ContentItemContent          = "content"
	ContentItemIFrameURL        = "iframe_url"
	ContentItemNotifyWinURL     = "notify_win_url"
	ContentItemNotifyDisplayURL = "notify_display_url"
)

var (
	ErrNewAuctionBidIsHigherThenMaxBid = errors.New("new auction bid is higher then max bid")
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

	// InternalAuctionCPMBid value provides maximal possible price without any comission
	// According to this value the system can choice the best item for the auction
	InternalAuctionCPMBid() billing.Money

	// PriorityFormatType from current Ad
	PriorityFormatType() types.FormatType

	// Validate item
	Validate() error

	// Context value
	Context(ctx ...context.Context) context.Context

	// Get ext field
	Get(key string) any
}

// ResponserItem for single AD
type ResponserItem interface {
	ResponserItemCommon

	// NetworkName by source
	NetworkName() string

	// AdID number
	AdID() uint64

	// AccountID number
	AccountID() uint64

	// CampaignID number
	CampaignID() uint64

	// Format object
	Format() *types.Format

	// PricingModel of the response advertisement
	PricingModel() types.PricingModel

	// ContentItem returns the ad response data
	ContentItem(name string) any

	// ContentItemString from the ad
	ContentItemString(name string) string

	// ContentFields from advertisement object
	ContentFields() map[string]any

	// MainAsset from response
	MainAsset() *admodels.AdAsset

	// Assets list
	Assets() admodels.AdAssets

	// Source of response
	Source() Source

	// ViewTrackerLinks returns traking links for view action
	ViewTrackerLinks() []string

	// ClickTrackerLinks returns traking links for click action
	ClickTrackerLinks() []string

	// IsDirect AD type
	IsDirect() bool

	// ActionURL returns target resource link for direct and banner click as well
	ActionURL() string

	// Width of item
	Width() int

	// Height of item
	Height() int

	// ECPM item value
	ECPM() billing.Money

	// PriceTestMode returns true if the price is in test mode
	PriceTestMode() bool

	// Price for specific action if supported `click`, `lead`, `view`
	// returns total price of the action
	Price(action admodels.Action) billing.Money

	// BidPrice returns bid price for the external auction source.
	// The current bid price will be adjusted according to the source correction factor and the commission share factor
	BidPrice() billing.Money

	// SetBidPrice value for external sources auction the system will pay
	SetBidPrice(price billing.Money) error

	// PurchasePrice gives the price of action from external resource (site, app, rtb, etc)
	// The cost of this request for the network.
	PurchasePrice(action admodels.Action) billing.Money

	// PotentialPrice wich can be received from source but was marked as descrepancy
	PotentialPrice(action admodels.Action) billing.Money

	// FinalPrice returns final price for the item which is including all possible commissions with all corrections
	// This price will be charged from advertiser
	FinalPrice(action admodels.Action) billing.Money

	// Second campaigns
	Second() *SecondAd

	// CommissionShareFactor returns the multipler for commission
	// calculation which system get from user revenue from 0 to 1
	CommissionShareFactor() float64

	// SourceCorrectionFactor value for the source
	SourceCorrectionFactor() float64

	// TargetCorrectionFactor value for the target
	TargetCorrectionFactor() float64
}

// ResponserMultipleItem interface for complex banners
type ResponserMultipleItem interface {
	ResponserItemCommon

	// Ads list response
	Ads() []ResponserItem

	// Count of response items
	Count() int
}
