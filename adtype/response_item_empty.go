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
type ResponseItemEmpty struct {
	Imp     *Impression
	Req     *BidRequest
	Src     Source
	context context.Context
}

// ID of current response item (unique code of current response)
func (*ResponseItemEmpty) ID() string { return "" }

// AuctionID returns ID of the current auction
func (r *ResponseItemEmpty) AuctionID() string {
	if r.Req == nil {
		return ""
	}
	return r.Req.ID
}

// Impression place object
func (it *ResponseItemEmpty) Impression() *Impression {
	if it != nil && it.Imp != nil {
		return it.Imp
	}
	return &emptyImpression
}

// ImpressionID unique code string
func (it *ResponseItemEmpty) ImpressionID() string { return it.Impression().ID }

// ExtImpressionID it's unique code of the auction bid impression
func (it *ResponseItemEmpty) ExtImpressionID() string { return it.Impression().ExternalID }

// ExtTargetID of the external network
func (it *ResponseItemEmpty) ExtTargetID() string { return it.Impression().ExternalTargetID }

// InternalAuctionCPMBid value provides maximal possible price without any comission
// According to this value the system can choice the best item for the auction
func (*ResponseItemEmpty) InternalAuctionCPMBid() billing.Money { return 0 }

// PriorityFormatType from current Ad
func (it *ResponseItemEmpty) PriorityFormatType() types.FormatType {
	format := it.Impression().FormatTypes
	if formatType := format.HasOneType(); formatType > types.FormatUndefinedType {
		return formatType
	}
	return format.FirstType()
}

// Validate item
func (*ResponseItemEmpty) Validate() error { return ErrResponseItemEmpty }

// Context value
func (it *ResponseItemEmpty) Context(ctx ...context.Context) context.Context {
	if len(ctx) > 0 && ctx[0] != nil {
		it.context = ctx[0]
	}
	if it.context == nil {
		it.context = context.Background()
	}
	return it.context
}

// Get ext field
func (*ResponseItemEmpty) Get(key string) any { return nil }

// AdID number
func (*ResponseItemEmpty) AdID() uint64 { return 0 }

// AccountID number
func (*ResponseItemEmpty) AccountID() uint64 { return 0 }

// CampaignID number
func (*ResponseItemEmpty) CampaignID() uint64 { return 0 }

// Format object
func (it *ResponseItemEmpty) Format() *types.Format {
	if it != nil {
		if fmtLst := it.Imp.Formats(); len(fmtLst) > 0 {
			return fmtLst[0]
		}
	}
	return &emptyFormat
}

// PricingModel of advertisement
func (*ResponseItemEmpty) PricingModel() types.PricingModel { return types.PricingModelUndefined }

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
func (it *ResponseItemEmpty) Source() Source {
	if it == nil || it.Src == nil {
		return nil
	}
	return it.Src
}

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

// PrepareBidViewPrice prepares the price for the action
// The price is adjusted according to the source correction factor and the commission share factor
func (*ResponseItemEmpty) PrepareBidViewPrice(price billing.Money) billing.Money { return price }

// PurchasePrice gives the price of view from external resource.
// The cost of this request.
func (*ResponseItemEmpty) PurchasePrice(action Action) billing.Money { return 0 }

// PotentialPrice wich can be received from source but was marked as descrepancy
func (*ResponseItemEmpty) PotentialPrice(action Action) billing.Money { return 0 }

// FinalPrice returns final price for the item which is including all possible commissions with all corrections
func (*ResponseItemEmpty) FinalPrice(action Action) billing.Money { return 0 }

// Second competitor campaign
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
func (it *ResponseItemEmpty) CommissionShareFactor() float64 {
	return it.Impression().CommissionShareFactor()
}

// IsDirect AD type
func (it *ResponseItemEmpty) IsDirect() bool { return it.PriorityFormatType().IsDirect() }

// IsBackup indicates whether the advertisement is a backup ad type.
func (it *ResponseItemEmpty) IsBackup() bool { return false }

// ActionURL returns target resource link for direct and banner click as well
func (*ResponseItemEmpty) ActionURL() string { return "" }

// Width of item
func (*ResponseItemEmpty) Width() int { return 0 }

// Height of item
func (*ResponseItemEmpty) Height() int { return 0 }

var _ ResponserItem = (*ResponseItemEmpty)(nil)
