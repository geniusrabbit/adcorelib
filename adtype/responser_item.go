//
// @project GeniusRabbit corelib 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com>
//

package adtype

import (
	"context"
	"errors"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

// Content item names represent various components of an advertisement.
// These constants are used as keys to retrieve specific content items from an ad response.
const (
	ContentItemLink             = "link"               // URL to which the ad redirects
	ContentItemContent          = "content"            // Main content of the ad
	ContentItemIFrameURL        = "iframe_url"         // URL for embedding the ad in an iframe
	ContentItemNotifyWinURL     = "notify_win_url"     // URL to notify when the ad wins an auction
	ContentItemNotifyDisplayURL = "notify_display_url" // URL to notify when the ad is displayed
)

// Predefined errors used within the adtype package.
var (
	// ErrNewAuctionBidIsHigherThenMaxBid is returned when a new auction bid exceeds the maximum allowed bid.
	ErrNewAuctionBidIsHigherThenMaxBid = errors.New("new auction bid is higher than max bid")
)

// ResponserItemCommon defines the common interface for response items in an ad auction.
// It includes methods to access identification, pricing, validation, and contextual data.
type ResponserItemCommon interface {
	// ID returns the unique identifier of the current response item.
	ID() string

	// Impression returns the Impression object associated with this response item.
	Impression() *Impression

	// ImpressionID returns the unique identifier string of the Impression.
	ImpressionID() string

	// ExtImpressionID returns the external unique identifier of the auction bid impression.
	ExtImpressionID() string

	// ExtTargetID returns the external target identifier from the external network.
	ExtTargetID() string

	// InternalAuctionCPMBid provides the maximum possible CPM (Cost Per Mille) bid without any commission.
	// This value is used by the system to select the best item for the auction.
	InternalAuctionCPMBid() billing.Money

	// PriorityFormatType returns the format type that has priority in the current ad.
	PriorityFormatType() types.FormatType

	// Validate checks the validity of the response item.
	// It returns an error if the item is invalid.
	Validate() error

	// Context returns the contextual information associated with the response item.
	// It can accept additional contexts to merge with the existing one.
	Context(ctx ...context.Context) context.Context

	// Get retrieves a value from the response item's extension map by the provided key.
	// It returns nil if the key does not exist.
	Get(key string) any
}

// ResponserItem defines the interface for a single advertisement response item.
// It extends ResponserItemCommon by adding methods specific to individual ads.
type ResponserItem interface {
	ResponserItemCommon

	// NetworkName returns the name of the ad network that provided the response.
	NetworkName() string

	// AdID returns the unique identifier of the advertisement.
	AdID() uint64

	// AccountID returns the unique identifier of the advertiser's account.
	AccountID() uint64

	// CampaignID returns the unique identifier of the advertising campaign.
	CampaignID() uint64

	// Format returns the format information of the advertisement.
	Format() *types.Format

	// PricingModel returns the pricing model used for the advertisement (e.g., CPC, CPM).
	PricingModel() types.PricingModel

	// ContentItem retrieves the ad response data for a given content item name.
	// It returns the content in its original type.
	ContentItem(name string) any

	// ContentItemString retrieves the ad response data for a given content item name as a string.
	ContentItemString(name string) string

	// ContentFields returns a map of all content fields associated with the advertisement.
	ContentFields() map[string]any

	// MainAsset returns the primary asset of the advertisement (e.g., image, video).
	MainAsset() *admodels.AdAsset

	// Assets returns a list of all assets associated with the advertisement.
	Assets() admodels.AdAssets

	// Source returns the source information of the advertisement response.
	Source() Source

	// ImpressionTrackerLinks returns traking links for impression action
	ImpressionTrackerLinks() []string

	// ViewTrackerLinks returns a slice of tracking URLs for view actions.
	ViewTrackerLinks() []string

	// ClickTrackerLinks returns a slice of tracking URLs for click actions.
	ClickTrackerLinks() []string

	// IsDirect indicates whether the advertisement is a direct ad type.
	IsDirect() bool

	// IsBackup indicates whether the advertisement is a backup ad type.
	IsBackup() bool

	// ActionURL returns the target URL for direct and banner click actions.
	ActionURL() string

	// Width returns the width of the advertisement in pixels.
	Width() int

	// Height returns the height of the advertisement in pixels.
	Height() int

	// ECPM returns the Effective Cost Per Mille (thousand impressions) value of the advertisement.
	ECPM() billing.Money

	// PriceTestMode indicates whether the price is in test mode.
	PriceTestMode() bool

	// Price returns the total price for a specific action (e.g., click, lead, view).
	Price(action admodels.Action) billing.Money

	// BidPrice returns the bid price for the external auction source.
	// The bid price is adjusted according to the source correction factor and the commission share factor.
	BidPrice() billing.Money

	// SetBidPrice sets the bid price value for external sources in an auction.
	// The system will pay this bid price.
	SetBidPrice(price billing.Money) error

	// PurchasePrice returns the price of a specific action from an external resource (e.g., site, app, RTB).
	// This represents the cost of the request for the network.
	PurchasePrice(action admodels.Action) billing.Money

	// PotentialPrice returns the potential price that could be received from the source but was marked as a discrepancy.
	PotentialPrice(action admodels.Action) billing.Money

	// FinalPrice returns the final price for the advertisement item, including all possible commissions and corrections.
	// This is the price that will be charged to the advertiser.
	FinalPrice(action admodels.Action) billing.Money

	// Second returns a pointer to the SecondAd, representing secondary campaigns.
	Second() *SecondAd

	// CommissionShareFactor returns the multiplier used for commission calculations.
	// It represents the portion of user revenue the system receives, ranging from 0 to 1.
	CommissionShareFactor() float64

	// SourceCorrectionFactor returns the correction factor applied to the source.
	SourceCorrectionFactor() float64

	// TargetCorrectionFactor returns the correction factor applied to the target.
	TargetCorrectionFactor() float64
}

// ResponserMultipleItem defines the interface for handling multiple advertisement response items.
// It is typically used for complex banners that contain multiple ads.
type ResponserMultipleItem interface {
	ResponserItemCommon

	// Ads returns a slice of ResponserItem interfaces representing individual ads within the response.
	Ads() []ResponserItem

	// Count returns the total number of response items (ads) in the response.
	Count() int
}
