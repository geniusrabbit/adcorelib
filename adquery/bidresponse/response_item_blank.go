package bidresponse

import (
	"context"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/price"
)

// ResponseItemBlank value
type ResponseItemBlank struct {
	ItemID string
	Imp    *adtype.Impression
	Src    adtype.Source

	FormatVal *types.Format

	PricingModelVal types.PricingModel
	PriceScope      price.PriceScope
}

// ID of current response item (unique code of current response)
func (it *ResponseItemBlank) ID() string { return it.ItemID }

// Impression place object
func (it *ResponseItemBlank) Impression() *adtype.Impression { return it.Imp }

// ImpressionID unique code string
func (it *ResponseItemBlank) ImpressionID() string { return it.Imp.ID }

// ExtImpressionID it's unique code of the auction bid impression
func (it *ResponseItemBlank) ExtImpressionID() string { return it.Imp.ExternalID }

// ExtTargetID of the external network
func (it *ResponseItemBlank) ExtTargetID() string { return it.Imp.ExternalTargetID }

// NetworkName by source
func (*ResponseItemBlank) NetworkName() string { return "" }

// Validate item
func (*ResponseItemBlank) Validate() error { return nil }

// AdID number
func (*ResponseItemBlank) AdID() uint64 { return 0 }

// AccountID number
func (*ResponseItemBlank) AccountID() uint64 { return 0 }

// CampaignID number
func (*ResponseItemBlank) CampaignID() uint64 { return 0 }

// Source of response
func (it *ResponseItemBlank) Source() adtype.Source { return it.Src }

// PriorityFormatType from current Ad
func (it *ResponseItemBlank) PriorityFormatType() types.FormatType {
	// If it's the only format type
	if formatType := it.FormatVal.Types.HasOneType(); formatType > types.FormatUndefinedType {
		return formatType
	}

	intersection := it.Imp.FormatTypes.Intersec(it.FormatVal.Types)
	if formatType := intersection.HasOneType(); formatType > types.FormatUndefinedType {
		return formatType
	}

	return intersection.FirstType()
}

// Format object
func (it *ResponseItemBlank) Format() *types.Format { return it.FormatVal }

///////////////////////////////////////////////////////////////////////////////
// Content data accessor method
///////////////////////////////////////////////////////////////////////////////

// ContentItem returns the ad response data
func (*ResponseItemBlank) ContentItem(name string) any { return nil }

// ContentItemString from the ad
func (*ResponseItemBlank) ContentItemString(name string) string { return "" }

// ContentFields from advertisement object
func (*ResponseItemBlank) ContentFields() map[string]any { return nil }

// MainAsset from response
func (*ResponseItemBlank) MainAsset() *admodels.AdAsset { return nil }

// Asset by name
func (*ResponseItemBlank) Asset(name string) *admodels.AdAsset { return nil }

// Assets list
func (*ResponseItemBlank) Assets() admodels.AdAssets { return nil }

// ImpressionTrackerLinks returns traking links for impression action
func (*ResponseItemBlank) ImpressionTrackerLinks() []string { return nil }

// ViewTrackerLinks returns traking links for view action
func (*ResponseItemBlank) ViewTrackerLinks() []string { return nil }

// ClickTrackerLinks returns traking links for click action
func (*ResponseItemBlank) ClickTrackerLinks() []string { return nil }

///////////////////////////////////////////////////////////////////////////////
// Price calculation methods
///////////////////////////////////////////////////////////////////////////////

// PricingModel of advertisement
func (it *ResponseItemBlank) PricingModel() types.PricingModel { return it.PricingModelVal }

// FixedPurchasePrice returns the fixed price of the action
func (it *ResponseItemBlank) FixedPurchasePrice(action adtype.Action) billing.Money {
	return it.Imp.PurchasePrice(action)
}

// ECPM returns the effective cost per mille of the item.
func (it *ResponseItemBlank) ECPM() billing.Money { return it.PriceScope.ECPM }

// PriceTestMode returns true if the price is in test mode
func (it *ResponseItemBlank) PriceTestMode() bool { return it.PriceScope.TestMode }

// Price per specific action type (view, click, lead, etc)
func (it *ResponseItemBlank) Price(action adtype.Action) billing.Money {
	return it.PriceScope.PricePerAction(action)
}

// BidViewPrice returns bid price for the external auction source.
// The current bid price will be adjusted according to the source correction factor and the commission share factor
func (it *ResponseItemBlank) BidViewPrice() billing.Money {
	return it.PriceScope.BidViewPrice
}

// SetBidViewPrice value for external sources auction the system will pay
func (it *ResponseItemBlank) SetBidViewPrice(bid billing.Money) error {
	if !it.PriceScope.SetBidViewPrice(bid, false) {
		return adtype.ErrNewAuctionBidIsHigherThenMaxBid
	}
	return nil
}

// PrepareBidViewPrice prepares the price for the action
// The price is adjusted according to the source correction factor and the commission share factor
func (it *ResponseItemBlank) PrepareBidViewPrice(price billing.Money) billing.Money {
	return it.PriceScope.PrepareBidViewPrice(price)
}

// PurchasePrice gives the price of view from external resource.
// The cost of this request.
func (it *ResponseItemBlank) PurchasePrice(action adtype.Action) billing.Money {
	return price.CalculatePurchasePrice(it, action)
}

// PotentialPrice wich can be received from source but was marked as descrepancy
func (it *ResponseItemBlank) PotentialPrice(action adtype.Action) billing.Money {
	return price.CalculatePotentialPrice(it, action)
}

// FinalPrice returns final price for the item which is including all possible commissions with all corrections
func (it *ResponseItemBlank) FinalPrice(action adtype.Action) billing.Money {
	return price.CalculateFinalPrice(it, action)
}

// InternalAuctionCPMBid value provides maximal possible price without any comission
// According to this value the system can choice the best item for the auction
func (it *ResponseItemBlank) InternalAuctionCPMBid() billing.Money {
	return price.CalculateInternalAuctionBid(it)
}

// Second campaigns
func (*ResponseItemBlank) Second() *adtype.SecondAd { return nil }

///////////////////////////////////////////////////////////////////////////////
// Revenue share/comission methods
///////////////////////////////////////////////////////////////////////////////

// CommissionShareFactor returns the commission share percentage which system gets from publisher.
func (it *ResponseItemBlank) CommissionShareFactor() float64 {
	return it.Imp.CommissionShareFactor()
}

// SourceCorrectionFactor value for the source
func (it *ResponseItemBlank) SourceCorrectionFactor() float64 {
	return it.Src.PriceCorrectionReduceFactor()
}

// TargetCorrectionFactor value for the target
func (it *ResponseItemBlank) TargetCorrectionFactor() float64 {
	return it.Imp.Target.RevenueShareReduceFactor()
}

///////////////////////////////////////////////////////////////////////////////
// Context methods
///////////////////////////////////////////////////////////////////////////////

// Context value
func (*ResponseItemBlank) Context(ctx ...context.Context) context.Context { return nil }

// Get ext field
func (*ResponseItemBlank) Get(key string) any { return nil }

///////////////////////////////////////////////////////////////////////////////
// Other methods
///////////////////////////////////////////////////////////////////////////////

// IsDirect AD type
func (*ResponseItemBlank) IsDirect() bool { return false }

// IsBackup indicates whether the advertisement is a backup ad type.
func (*ResponseItemBlank) IsBackup() bool { return false }

// ActionURL returns target resource link for direct and banner click as well
func (*ResponseItemBlank) ActionURL() string { return "" }

// Width of item
func (*ResponseItemBlank) Width() int { return 0 }

// Height of item
func (*ResponseItemBlank) Height() int { return 0 }

var _ adtype.ResponserItem = (*ResponseItemBlank)(nil)
