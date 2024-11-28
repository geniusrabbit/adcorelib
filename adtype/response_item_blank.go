package adtype

import (
	"context"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/price"
)

// ResponseItemBlank value
type ResponseItemBlank struct {
	ItemID string
	Imp    *Impression
	Src    Source

	FormatVal *types.Format

	PricingModelVal types.PricingModel
	PriceScope      price.PriceScope
}

// ID of current response item (unique code of current response)
func (it *ResponseItemBlank) ID() string { return it.ItemID }

// Impression place object
func (it *ResponseItemBlank) Impression() *Impression { return it.Imp }

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
func (it *ResponseItemBlank) Source() Source { return it.Src }

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
func (it *ResponseItemBlank) FixedPurchasePrice(action admodels.Action) billing.Money {
	return it.Imp.PurchasePrice(action)
}

// ECPM returns the effective cost per mille of the item.
func (it *ResponseItemBlank) ECPM() billing.Money { return it.PriceScope.ECPM }

// PriceTestMode returns true if the price is in test mode
func (it *ResponseItemBlank) PriceTestMode() bool { return it.PriceScope.TestMode }

// Price per specific action type (view, click, lead, etc)
func (it *ResponseItemBlank) Price(action admodels.Action) billing.Money {
	return it.PriceScope.PricePerAction(action)
}

// BidPrice returns bid price for the external auction source.
// The current bid price will be adjusted according to the source correction factor and the commission share factor
func (it *ResponseItemBlank) BidPrice() billing.Money {
	return it.PriceScope.BidPrice
}

// SetBidPrice value for external sources auction the system will pay
func (it *ResponseItemBlank) SetBidPrice(bid billing.Money) error {
	if !it.PriceScope.SetBidPrice(bid, false) {
		return ErrNewAuctionBidIsHigherThenMaxBid
	}
	return nil
}

// PurchasePrice gives the price of view from external resource.
// The cost of this request.
func (it *ResponseItemBlank) PurchasePrice(action admodels.Action) billing.Money {
	return price.CalculatePurchasePrice(it, action)
}

// PotentialPrice wich can be received from source but was marked as descrepancy
func (it *ResponseItemBlank) PotentialPrice(action admodels.Action) billing.Money {
	return price.CalculatePotentialPrice(it, action)
}

// FinalPrice returns final price for the item which is including all possible commissions with all corrections
func (it *ResponseItemBlank) FinalPrice(action admodels.Action) billing.Money {
	return price.CalculateFinalPrice(it, action)
}

// InternalAuctionCPMBid value provides maximal possible price without any comission
// According to this value the system can choice the best item for the auction
func (it *ResponseItemBlank) InternalAuctionCPMBid() billing.Money {
	return price.CalculateInternalAuctionBid(it)
}

// SetAuctionCPMBid value for external sources auction the system will pay
func (it *ResponseItemBlank) SetAuctionCPMBid(price billing.Money, includeFactors ...PriceFactor) error {
	if len(includeFactors) > 0 {
		price += PriceFactorFromList(includeFactors...).AddComission(price, it)
	}
	if !it.PriceScope.SetBidPrice(price/1000, false) {
		return ErrNewAuctionBidIsHigherThenMaxBid
	}
	return nil
}

// Second campaigns
func (*ResponseItemBlank) Second() *SecondAd { return nil }

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

// ActionURL returns target resource link for direct and banner click as well
func (*ResponseItemBlank) ActionURL() string { return "" }

// Width of item
func (*ResponseItemBlank) Width() int { return 0 }

// Height of item
func (*ResponseItemBlank) Height() int { return 0 }

var _ ResponserItem = (*ResponseItemBlank)(nil)
