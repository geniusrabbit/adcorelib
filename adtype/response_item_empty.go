package adtype

import (
	"context"

	"geniusrabbit.dev/corelib/admodels"
	"geniusrabbit.dev/corelib/admodels/types"
	"geniusrabbit.dev/corelib/billing"
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

// AuctionCPMBid value price without any comission
func (*ResponseItemEmpty) AuctionCPMBid() billing.Money { return 0 }

// PriorityFormatType from current Ad
func (*ResponseItemEmpty) PriorityFormatType() types.FormatType { return 0 }

// Validate item
func (*ResponseItemEmpty) Validate() error { return nil }

// Context value
func (*ResponseItemEmpty) Context(ctx ...context.Context) context.Context { return nil }

// Get ext field
func (*ResponseItemEmpty) Get(key string) interface{} { return nil }

// AdID number
func (*ResponseItemEmpty) AdID() uint64 { return 0 }

// CompanyID number
func (*ResponseItemEmpty) CompanyID() uint64 { return 0 }

// CampaignID number
func (*ResponseItemEmpty) CampaignID() uint64 { return 0 }

// Format object
func (*ResponseItemEmpty) Format() *types.Format { return &emptyFormat }

// PricingModel of advertisement
func (*ResponseItemEmpty) PricingModel() types.PricingModel { return 0 }

// ContentItem returns the ad response data
func (*ResponseItemEmpty) ContentItem(name string) interface{} { return nil }

// ContentItemString from the ad
func (*ResponseItemEmpty) ContentItemString(name string) string { return "" }

// ContentFields from advertisement object
func (*ResponseItemEmpty) ContentFields() map[string]interface{} { return nil }

// MainAsset from response
func (*ResponseItemEmpty) MainAsset() *admodels.AdFile { return nil }

// Asset by name
func (*ResponseItemEmpty) Asset(name string) *admodels.AdFile { return nil }

// Assets list
func (*ResponseItemEmpty) Assets() []admodels.AdFile { return nil }

// Source of response
func (*ResponseItemEmpty) Source() Source { return nil }

// ViewTrackerLinks returns traking links for view action
func (*ResponseItemEmpty) ViewTrackerLinks() []string { return nil }

// ClickTrackerLinks returns traking links for click action
func (*ResponseItemEmpty) ClickTrackerLinks() []string { return nil }

// NetworkName by source
func (*ResponseItemEmpty) NetworkName() string { return "" }

// Price summ
func (*ResponseItemEmpty) Price(action admodels.Action) billing.Money { return 0 }

// SetCPMPrice update of DSP auction value
func (*ResponseItemEmpty) SetCPMPrice(price billing.Money, includeFactors ...bool) {}

// CPMPrice value price value for DSP auction
func (*ResponseItemEmpty) CPMPrice(removeFactors ...bool) billing.Money { return 0 }

// Second campaigns
func (*ResponseItemEmpty) Second() *SecondAd { return nil }

// Revenue money (in percents)
func (*ResponseItemEmpty) Revenue() float64 { return 0 }

// Potential money (in percents)
func (*ResponseItemEmpty) Potential() float64 { return 0 }

// ECPM item value
func (*ResponseItemEmpty) ECPM() billing.Money { return 0 }

// RevenueShareFactor returns the multipler for company
// revenue calculation per action from 0 to 1
func (*ResponseItemEmpty) RevenueShareFactor() float64 { return 0 }

// ComissionShareFactor returns the multipler for commission
// calculation which system get from user revenue from 0 to 1
func (*ResponseItemEmpty) ComissionShareFactor() float64 { return 0 }

// IsDirect AD type
func (*ResponseItemEmpty) IsDirect() bool { return false }

// ActionURL returns target resource link for direct and banner click as well
func (*ResponseItemEmpty) ActionURL() string { return "" }

// Width of item
func (*ResponseItemEmpty) Width() int { return 0 }

// Height of item
func (*ResponseItemEmpty) Height() int { return 0 }

var _ ResponserItem = (*ResponseItemEmpty)(nil)