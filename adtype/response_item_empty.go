package adtype

import (
	"context"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
)

var (
	emptyFormat     = types.Format{}
	emptyImpression = Impression{}
)

// ResponseItemEmpty value
type ResponseItemEmpty struct{}

// ID of current response item (unique code of current response)
func (*ResponseItemEmpty) ID() string { return "" }

// Impression place object
func (*ResponseItemEmpty) Impression() *Impression { return &emptyImpression }

// ImpressionID unique code string
func (*ResponseItemEmpty) ImpressionID() string { return "" }

// ExtImpressionID it's unique code of the auction bid impression
func (*ResponseItemEmpty) ExtImpressionID() string { return "" }

// ExtTargetID of the external network
func (*ResponseItemEmpty) ExtTargetID() string { return "" }

// InternalAuctionCPMBid value provides maximal possible price without any comission
// According to this value the system can choice the best item for the auction
func (*ResponseItemEmpty) InternalAuctionCPMBid() billing.Money { return 0 }

// PriorityFormatType from current Ad
func (*ResponseItemEmpty) PriorityFormatType() types.FormatType { return 0 }

// Validate item
func (*ResponseItemEmpty) Validate() error { return nil }

// Context value
func (*ResponseItemEmpty) Context(ctx ...context.Context) context.Context { return nil }

// Get ext field
func (*ResponseItemEmpty) Get(key string) any { return nil }

// AdID number
func (*ResponseItemEmpty) AdID() uint64 { return 0 }

// AccountID number
func (*ResponseItemEmpty) AccountID() uint64 { return 0 }

// CampaignID number
func (*ResponseItemEmpty) CampaignID() uint64 { return 0 }

// Format object
func (*ResponseItemEmpty) Format() *types.Format { return &emptyFormat }

// PricingModel of advertisement
func (*ResponseItemEmpty) PricingModel() types.PricingModel { return 0 }

// ContentItem returns the ad response data
func (*ResponseItemEmpty) ContentItem(name string) any { return nil }

// ContentItemString from the ad
func (*ResponseItemEmpty) ContentItemString(name string) string { return "" }

// ContentFields from advertisement object
func (*ResponseItemEmpty) ContentFields() map[string]any { return nil }

// MainAsset from response
func (*ResponseItemEmpty) MainAsset() *admodels.AdAsset { return nil }

// Asset by name
func (*ResponseItemEmpty) Asset(name string) *admodels.AdAsset { return nil }

// Assets list
func (*ResponseItemEmpty) Assets() admodels.AdAssets { return nil }

// Source of response
func (*ResponseItemEmpty) Source() Source { return nil }

// ImpressionTrackerLinks returns traking links for impression action
func (*ResponseItemEmpty) ImpressionTrackerLinks() []string { return nil }

// ViewTrackerLinks returns traking links for view action
func (*ResponseItemEmpty) ViewTrackerLinks() []string { return nil }

// ClickTrackerLinks returns traking links for click action
func (*ResponseItemEmpty) ClickTrackerLinks() []string { return nil }

// NetworkName by source
func (*ResponseItemEmpty) NetworkName() string { return "" }

// ECPM item value
func (*ResponseItemEmpty) ECPM() billing.Money { return 0 }

// PriceTestMode returns true if the price is in test mode
func (*ResponseItemEmpty) PriceTestMode() bool { return false }

// Price per specific action type (view, click, lead, etc)
func (*ResponseItemEmpty) Price(action Action) billing.Money { return 0 }

// BidViewPrice returns bid price for the external auction source.
// The current bid price will be adjusted according to the source correction factor and the commission share factor
func (*ResponseItemEmpty) BidViewPrice() billing.Money { return 0 }

// SetBidViewPrice value for external sources auction the system will pay
func (*ResponseItemEmpty) SetBidViewPrice(price billing.Money) error { return nil }

// PurchasePrice gives the price of view from external resource.
// The cost of this request.
func (*ResponseItemEmpty) PurchasePrice(action Action) billing.Money { return 0 }

// PotentialPrice wich can be received from source but was marked as descrepancy
func (*ResponseItemEmpty) PotentialPrice(action Action) billing.Money { return 0 }

// FinalPrice returns final price for the item which is including all possible commissions with all corrections
func (*ResponseItemEmpty) FinalPrice(action Action) billing.Money { return 0 }

// Second campaigns
func (*ResponseItemEmpty) Second() *SecondAd { return nil }

// RevenuePercent money
// func (*ResponseItemEmpty) RevenuePercent() float64 { return 0 }

// PotentialPercent money
func (*ResponseItemEmpty) PotentialPercent() float64 { return 0 }

// SourceCorrectionFactor value for the source
func (it *ResponseItemEmpty) SourceCorrectionFactor() float64 { return 0 }

// TargetCorrectionFactor value for the target
func (it *ResponseItemEmpty) TargetCorrectionFactor() float64 { return 0 }

// CommissionShareFactor returns the multipler for commission
// calculation which system get from user revenue from 0 to 1
func (*ResponseItemEmpty) CommissionShareFactor() float64 { return 0 }

// IsDirect AD type
func (*ResponseItemEmpty) IsDirect() bool { return false }

// IsBackup indicates whether the advertisement is a backup ad type.
func (*ResponseItemEmpty) IsBackup() bool { return false }

// ActionURL returns target resource link for direct and banner click as well
func (*ResponseItemEmpty) ActionURL() string { return "" }

// Width of item
func (*ResponseItemEmpty) Width() int { return 0 }

// Height of item
func (*ResponseItemEmpty) Height() int { return 0 }

var _ ResponserItem = (*ResponseItemEmpty)(nil)
