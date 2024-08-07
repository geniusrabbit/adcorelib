//
// @project GeniusRabbit corelib 2016 – 2018, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2018, 2024
//

package admodels

import (
	"errors"
	"fmt"
	"math/rand"
	"path/filepath"
	"strings"

	"github.com/demdxx/gocast/v2"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/models"
)

// Errors
var (
	ErrUndefinedAdContext   = errors.New("[advertisement model]: undefined AD context")
	ErrTooManyErrorsInTheAd = errors.New("[advertisement model]: too many errors")
	ErrInvalidAdFormat      = errors.New("[advertisement model]: invalid ad format")
)

// AdFlag option state type
type AdFlag uint8

// Ad flag option types
const (
	AdFlagIsPremium AdFlag = 0x01
	AdFlagActive    AdFlag = 0x02
	AdFlagInsecure  AdFlag = 0x04
	AdFlagAsPopover AdFlag = 0x08
)

const (
	proxyIFrameURL = "iframe_url"
)

// Ad model
type Ad struct {
	ID uint64 // Ad ID

	// Data
	Content map[string]any // Extend data
	Assets  AdAssets

	PricingModel     types.PricingModel
	Weight           uint8
	FrequencyCapping uint8
	Flags            AdFlag
	Campaign         *Campaign `json:"-" xml:"-"`
	Bids             AdBids

	// Some advertisement formats could be streacheble but min/max width of heights
	// very significant for some types of advertisement where that ad can be look
	// wierd in some cases
	Format *types.Format

	// State           balance.State // Balance and counters state
	BidPrice        billing.Money // Max price per one View (used in DSP auction)
	Price           billing.Money // Price per one view or click
	LeadPrice       billing.Money // Price per one lead
	DailyBudget     billing.Money //
	Budget          billing.Money //
	DailyTestBudget billing.Money // Test money amount a day (it stops automaticaly if not profit for this amount)
	TestBudget      billing.Money // Test money amount for the whole period
	DefaultLink     string

	// Targeting
	Hours types.Hours // len(24) * bitmask in week days

	errors []error
}

// GetID of the object
func (a *Ad) GetID() uint64 {
	return a.ID
}

// ContentItem by name
func (a *Ad) ContentItem(name string) any {
	if a.Content != nil {
		return a.Content[name]
	}
	return nil
}

// ContentItemString by name
func (a *Ad) ContentItemString(name string) string {
	return gocast.Str(a.ContentItem(name))
}

// MainAsset field
func (a *Ad) MainAsset() *AdAsset {
	return a.Asset(types.FormatAssetMain)
}

// Asset by name
func (a *Ad) Asset(name string) *AdAsset {
	return a.Assets.Asset(name)
}

// RandomAdLink from ad model
func (a *Ad) RandomAdLink() AdLink {
	if a.DefaultLink != "" {
		return AdLink{Link: a.DefaultLink}
	}
	if a.Campaign != nil {
		if count := len(a.Campaign.Links); count > 0 {
			return a.Campaign.Links[rand.Int()%count]
		}
	}
	return AdLink{}
}

// Validate ad
func (a *Ad) Validate() error {
	if a.Format == nil {
		return ErrInvalidAdFormat
	}
	if len(a.errors) > 0 {
		if len(a.errors) == 1 {
			return a.errors[0]
		}
		return ErrTooManyErrorsInTheAd
	}
	if a.Format.Config == nil {
		return nil
	}
	for _, asset := range a.Format.Config.Assets {
		if asset.IsRequired() {
			if a.Asset(asset.GetName()) == nil {
				return fmt.Errorf(`[Advertisement model] asset "%s" is not present in Ad%d format "%s"`,
					asset.GetName(), a.ID, a.Format.Codename)
			}
		}
	}
	return nil
}

// SetPricingModel strategy
func (a *Ad) SetPricingModel(pm types.PricingModel) {
	a.PricingModel = pm
}

// SetWeight of the AD
func (a *Ad) SetWeight(w uint8) {
	a.Weight = w
}

// SetFrequencyCapping of the AD
func (a *Ad) SetFrequencyCapping(frequencyCapping uint8) {
	a.FrequencyCapping = frequencyCapping
}

func (a *Ad) ProxyURL() string {
	return a.ContentItemString(proxyIFrameURL)
}

// // GetTotalSpent of campaign
// func (a *Ad) GetTotalSpent() billing.Money {
// 	if a.State == nil {
// 		return 0
// 	}
// 	return a.State.TotalSpent()
// }

// // GetSpent of campaign
// func (a *Ad) GetSpent() billing.Money {
// 	if a.State == nil {
// 		return 0
// 	}
// 	return a.State.Spent()
// }

// // GetTotalProfit of campaign
// func (a *Ad) GetTotalProfit() billing.Money {
// 	if a.State == nil {
// 		return 0
// 	}
// 	return a.State.TotalProfit()
// }

// // GetProfit of campaign
// func (a *Ad) GetProfit() billing.Money {
// 	if a.State == nil {
// 		return 0
// 	}
// 	return a.State.Profit()
// }

// // GetDailyBudget of campaign
// func (a *Ad) GetDailyBudget() billing.Money {
// 	if a.State == nil {
// 		return a.DailyBudget
// 	}
// 	return a.State.DailyBudget()
// }

// // Imps returns count of impressions
// func (a *Ad) Imps() uint64 {
// 	if a.State == nil {
// 		return 0
// 	}
// 	return a.State.Imps()
// }

// // Clicks returns count of clicks
// func (a *Ad) Clicks() uint64 {
// 	if a.State == nil {
// 		return 0
// 	}
// 	return a.State.Clicks()
// }

// // Leads returns count of leads
// func (a *Ad) Leads() uint64 {
// 	if a.State == nil {
// 		return 0
// 	}
// 	return a.State.Leads()
// }

///////////////////////////////////////////////////////////////////////////////
/// Check methods
///////////////////////////////////////////////////////////////////////////////

// TargetBid by targeting pointer
func (a *Ad) TargetBid(pointer types.TargetPointer) TargetBid {
	if bid := a.Bids.Bid(pointer); bid != nil {
		return TargetBid{
			Ad:        a,
			Bid:       bid,
			BidPrice:  bid.BidPrice,
			Price:     bid.Price,
			LeadPrice: bid.LeadPrice,
			ECPM:      a.ecpm(pointer, bid.Price),
		}
	}

	return TargetBid{
		Ad:        a,
		Bid:       nil,
		BidPrice:  a.BidPrice,
		Price:     a.Price,
		LeadPrice: a.LeadPrice,
		ECPM:      a.ecpm(pointer, a.Price),
	}
}

// // TestBudgetValues of campaign
// func (a *Ad) TestBudgetValues() bool {
// 	return a.Campaign.TestBudgetValue() &&
// 		(a.GetDailyBudget() <= 0 || a.GetSpent() < a.GetDailyBudget()) &&
// 		(a.Budget <= 0 || a.GetTotalSpent() < a.Budget)
// }

// // TestProfit of the campaign
// func (a *Ad) TestProfit() bool {
// 	return a.Campaign.TestProfit() &&
// 		// test daily with profit
// 		(a.DailyTestBudget <= 0 || a.DailyTestBudget >= a.GetSpent()-a.GetProfit()) &&
// 		// Total test with profit
// 		(a.TestBudget <= 0 || a.TestBudget >= a.GetTotalSpent()-a.GetTotalProfit())
// }

///////////////////////////////////////////////////////////////////////////////
/// Status methods
///////////////////////////////////////////////////////////////////////////////

// Active ad
func (a *Ad) Active() bool {
	return a.Flags&AdFlagActive != 0
}

// Secure ad
func (a *Ad) Secure() bool {
	return a.Flags&AdFlagInsecure == 0
}

// AsPopover ad
func (a *Ad) AsPopover() bool {
	return a.Flags&AdFlagAsPopover != 0
}

// IsPremium ad
func (a *Ad) IsPremium() bool {
	return a.Flags&AdFlagIsPremium != 0
}

// AddFlag to model
func (a *Ad) AddFlag(flag AdFlag) {
	a.Flags |= flag
}

///////////////////////////////////////////////////////////////////////////////
/// Extra errors state
///////////////////////////////////////////////////////////////////////////////

// SetError named error
func (a *Ad) SetError(name string, err error) {
	nilIndex := -1

	if err != nil {
		err = NewNamedErrorWrapper(name, err)
	}

	for i, e := range a.errors {
		if e == nil {
			nilIndex = i
			continue
		}
		switch er := e.(type) {
		case NamedErrorWrapper:
			if er.GetName() == name {
				a.errors[i] = err
				return
			}
		}
	}
	if nilIndex >= 0 {
		a.errors[nilIndex] = err
	} else if err != nil {
		a.errors = append(a.errors, err)
	}
}

// ErrorByName returns one error with the name of object
func (a *Ad) ErrorByName(name string) error {
	for _, e := range a.errors {
		switch er := e.(type) {
		case NamedErrorWrapper:
			if er.GetName() == name {
				return e
			}
		}
	}
	return nil
}

// Errors object array
func (a *Ad) Errors() []error {
	return a.errors
}

///////////////////////////////////////////////////////////////////////////////
/// Internal methods
///////////////////////////////////////////////////////////////////////////////

func (a *Ad) ecpm( /*pointer*/ types.TargetPointer /*price*/, billing.Money) billing.Money {
	// if a.State != nil {
	// 	if ecpm := a.State.ECPM().Value(pointer); ecpm > 0 {
	// 		return ecpm
	// 	}
	// }

	switch a.PricingModel {
	case types.PricingModelCPM:
		return a.Price * 1000
	default:
		return 0
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func parseAd(camp *Campaign, adBase *models.Ad, formats types.FormatsAccessor) (ad *Ad, err error) {
	var (
		bids   []AdBid
		hours  types.Hours
		flags  AdFlag
		format *types.Format
	)

	// Preprocess info
	{
		if hours, err = types.HoursByString(adBase.Hours); err != nil {
			return ad, err
		}

		if _bids := adBase.Bids.Data; _bids != nil && len(*_bids) > 0 {
			for _, bid := range *_bids {
				bids = append(bids, AdBid{
					BidPrice:    bid.BidPrice,
					Price:       bid.Price,
					LeadPrice:   bid.LeadPrice,
					Tags:        bid.Tags,
					Zones:       bid.Zones,
					Domains:     bid.Domains,
					Sex:         bid.Sex,
					Age:         0, //bid.Age,
					Categories:  bid.Categories,
					Countries:   nil, // bid.Countries,
					Cities:      bid.Cities,
					Languages:   nil, // bid.Languages,
					DeviceTypes: bid.DeviceTypes,
					Devices:     bid.Devices,
					Os:          bid.Os,
					Browsers:    bid.Browsers,
					Hours:       bid.Hours,
				})
			}
		}
	}

	if adBase.Active.IsActive() && adBase.Status.IsApproved() {
		flags |= AdFlagActive
	}

	if adBase.FormatID != 0 {
		format = formats.FormatByID(adBase.FormatID)
	} else if adBase.Format != nil {
		if adBase.Format.ID > 0 {
			format = formats.FormatByID(adBase.Format.ID)
		} else {
			format = formats.FormatByCode(adBase.Format.Codename)
		}
	}

	ad = &Ad{
		ID:               adBase.ID,
		Format:           format,
		Assets:           nil,
		PricingModel:     adBase.PricingModel,
		FrequencyCapping: uint8(adBase.FrequencyCapping),
		Weight:           uint8(adBase.Weight),
		Flags:            flags,
		Bids:             bids,
		Price:            billing.MoneyFloat(adBase.Price),
		BidPrice:         billing.MoneyFloat(adBase.BidPrice),
		LeadPrice:        billing.MoneyFloat(adBase.LeadPrice),
		DailyBudget:      billing.MoneyFloat(adBase.DailyBudget),
		Budget:           billing.MoneyFloat(adBase.Budget),
		DailyTestBudget:  billing.MoneyFloat(adBase.DailyTestBudget),
		TestBudget:       billing.MoneyFloat(adBase.TestBudget),
		Hours:            hours,
		Campaign:         camp,
	}

	if ad.Format == nil {
		return nil, fmt.Errorf("ad[%d] undefined format ID: %d", adBase.ID, adBase.FormatID)
	}

	for _, as := range adBase.Assets {
		adFile := &AdAsset{
			ID:          as.ID,
			Name:        as.Name.String,
			Path:        filepath.Join(as.ObjectID, as.Meta.Data.Main.Name),
			Type:        types.AdAssetType(as.Type),
			ContentType: as.ContentType,
			Width:       as.Meta.Data.Main.Width,
			Height:      as.Meta.Data.Main.Height,
			Thumbs:      make([]AdAssetThumb, 0, len(as.Meta.Data.Items)),
		}
		for _, thmb := range as.Meta.Data.Items {
			if thmb.Type == types.AdAssetUndefinedType {
				continue
			}
			adFile.Thumbs = append(adFile.Thumbs, AdAssetThumb{
				Path:        filepath.Join(as.ObjectID, thmb.Name),
				Type:        types.AdAssetType(thmb.Type),
				Width:       thmb.Width,
				Height:      thmb.Height,
				ContentType: thmb.ContentType,
				Ext:         thmb.Ext,
			})
		}
		ad.Assets = append(ad.Assets, adFile)
	}

	// Add restriction of minimal-maximal dementions
	if ad.Format.IsStretch() {
		if adBase.MinWidth > 0 || adBase.MinHeight > 0 || adBase.MaxWidth > 0 || adBase.MaxHeight > 0 {
			ad.Format = ad.Format.CloneWithSize(adBase.MaxWidth, adBase.MaxHeight, adBase.MinWidth, adBase.MinHeight)
		}
	}

	if adBase.Context.Data != nil {
		ad.Content = *adBase.Context.Data
		ad.DefaultLink = ad.ContentItemString("default_link")
	}

	// Up secure flag by iframe URL or content
	urlFields := []string{proxyIFrameURL, "url"}
	for _, key := range urlFields {
		url := ad.ContentItemString(key)
		if url == "" {
			continue
		}
		if strings.HasPrefix(url, "http://") {
			ad.Flags |= AdFlagInsecure
			break
		}
	}

	return ad, err
}
