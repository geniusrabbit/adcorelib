package bidresponse

import (
	"context"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/adtype/prices"
	"github.com/geniusrabbit/adcorelib/billing"
)

// ResponseItemBlank value
type ResponseItemBlank struct {
	ItemID string
	Imp    *adtype.Impression
	Src    adtype.Source

	FormatVal *types.Format

	PricingModelVal types.PricingModel
	PriceScope      prices.PriceScope
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

// TargetCodename of the target placement codename
func (it *ResponseItemBlank) TargetCodename() string { return it.Imp.TargetCodename() }

// NetworkName by source
func (*ResponseItemBlank) NetworkName() string { return "" }

// Validate item
func (*ResponseItemBlank) Validate() error { return nil }

// AccountID returns the unique identifier of the advertiser's account.
func (*ResponseItemBlank) AccountID() uint64 { return 0 }

// CampaignID returns the unique identifier of the advertising campaign.
func (*ResponseItemBlank) CampaignID() uint64 { return 0 }

// AdID returns the unique identifier of the advertisement.
func (*ResponseItemBlank) AdID() string { return "" }

// CreativeID returns the unique identifier of the creative.
func (*ResponseItemBlank) CreativeID() string { return "" }

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
func (*ResponseItemBlank) MainAsset() *admodels.AdFileAsset { return nil }

// Asset by name
func (*ResponseItemBlank) Asset(name string) *admodels.AdFileAsset { return nil }

// Assets list
func (*ResponseItemBlank) Assets() admodels.AdFileAssets { return nil }

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
func (it *ResponseItemBlank) ECPM() billing.Money { return it.PriceScope.EffectiveCPM() }

// Price per specific action type (view, click, lead, etc)
func (it *ResponseItemBlank) Price(action adtype.Action) billing.Money {
	return it.PriceScope.PricePerAction(action)
}

// SetBidPrice sets the current bid price for the given action. If withCommission
// is true, the price is grossed up by the discrepancy corrections and the
// commission share before being stored.
func (it *ResponseItemBlank) SetBidPrice(action adtype.Action, bid billing.Money, withCommission bool) error {
	return it.PriceScope.SetBidPrice(action, bid, it, withCommission)
}

// PrepareBidPrice prepares the bid price for the given action by clamping it
// to the maximal allowed bid of that action (if defined).
func (it *ResponseItemBlank) PrepareBidPrice(action adtype.Action, p billing.Money) billing.Money {
	return it.PriceScope.PrepareBidPerAction(action, p)
}

// PurchasePrice gives the price of view from external resource.
// The cost of this request.
func (it *ResponseItemBlank) PurchasePrice(action adtype.Action) billing.Money {
	return it.PriceScope.PublisherPricePerAction(action, it)
}

// PotentialPrice wich can be received from source but was marked as descrepancy
func (it *ResponseItemBlank) PotentialPrice(action adtype.Action) billing.Money {
	return it.PriceScope.PotentialPricePerAction(action)
}

// FinalPrice returns final price for the item which is including all possible commissions with all corrections
func (it *ResponseItemBlank) FinalPrice(action adtype.Action) billing.Money {
	return it.PriceScope.AdvertiserPricePerAction(action, it)
}

// InternalAuctionCPMBid value provides maximal possible price without any comission
// According to this value the system can choice the best item for the auction
func (it *ResponseItemBlank) InternalAuctionCPMBid() billing.Money {
	return it.ECPM()
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
	if it.Src == nil {
		return 0
	}
	return it.Src.PriceCorrectionReduceFactor()
}

// TargetCorrectionFactor value for the target
func (it *ResponseItemBlank) TargetCorrectionFactor() float64 {
	if it.Imp == nil || it.Imp.Target == nil {
		return 0
	}
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

var (
	_ adtype.ResponseItem        = (*ResponseItemBlank)(nil)
	_ prices.Factors             = (*ResponseItemBlank)(nil)
	_ prices.FixedPurchasePricer = (*ResponseItemBlank)(nil)
)
