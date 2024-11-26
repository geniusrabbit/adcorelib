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

	// NetworkName by source
	NetworkName() string

	// Price for specific action if supported `click`, `lead`, `view`
	// returns total price of the action
	Price(action admodels.Action) billing.Money

	// PurchasePrice gives the price of action from external resource.
	// The cost of this request for the network.
	PurchasePrice(action admodels.Action) billing.Money

	// PotentialPrice wich can be received from source but was marked as descrepancy
	PotentialPrice(action admodels.Action) billing.Money

	// FinalPrice returns final price for the item which is including all possible commissions with all corrections
	FinalPrice(action admodels.Action) billing.Money

	// SetAuctionCPMBid value for external sources auction the system will pay
	SetAuctionCPMBid(price billing.Money, includeFactors ...PriceFactor) error

	// AuctionCPMBid value provides price for external sources
	// The prive what we can pay for the action to the external source
	AuctionCPMBid(removeFactors ...PriceFactor) billing.Money

	// Second campaigns
	Second() *SecondAd

	// ECPM item value
	ECPM() billing.Money

	// RevenueShareFactor returns the multipler for company
	// revenue calculation per action from 0 to 1
	// RevenueShareFactor() float64

	// CommissionShareFactor returns the multipler for commission
	// calculation which system get from user revenue from 0 to 1
	CommissionShareFactor() float64

	// SourceCorrectionFactor value for the source
	SourceCorrectionFactor() float64

	// TargetCorrectionFactor value for the target
	TargetCorrectionFactor() float64

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
