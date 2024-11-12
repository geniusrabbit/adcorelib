//
// @project GeniusRabbit corelib 2016 – 2019, 2021
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019, 2021
//

package admodels

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/geniusrabbit/gosql/v2"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/fasttime"
	"github.com/geniusrabbit/adcorelib/geo"
	"github.com/geniusrabbit/adcorelib/i18n/languages"
	"github.com/geniusrabbit/adcorelib/models"
	"github.com/geniusrabbit/adcorelib/searchtypes"
)

// Errors set
var (
	ErrInvalidCampaignAds = errors.New("invalid campaigns ads")
)

// CampaignCamparator interface for index
type CampaignCamparator interface {
	CompareCampaign(c *Campaign) int
}

// Campaign model represents a campaign with ads and targeting settings for the ad network.
type Campaign struct {
	id    uint64
	Acc   *Account
	AccID uint64

	Weight uint8
	Flags  CampaignFlagType

	DailyTestBudget billing.Money // Test money amount a day (it stops automaticaly if not profit for this amount)
	TestBudget      billing.Money // Test money amount for the whole period
	DailyBudget     billing.Money
	Budget          billing.Money
	MaxBid          billing.Money // Max bid for the campaign (in RTB auction)
	CurrentState    State         `json:"-"`

	// List of ads
	Ads   []*Ad
	Links []AdLink

	// Targeting
	FormatSet   searchtypes.NumberBitset[uint]           //
	Context     gosql.NullableJSON[map[string]any]       //
	Tags        gosql.NullableStringArray                //
	Zones       gosql.NullableOrderedNumberArray[uint64] //
	Domains     gosql.NullableStringArray                // site domains or application bundels
	Sex         gosql.NullableOrderedNumberArray[uint]   //
	Age         gosql.NullableOrderedNumberArray[uint]   //
	Categories  gosql.NullableOrderedNumberArray[uint64] //
	Countries   gosql.NullableOrderedNumberArray[uint64] //
	Cities      gosql.NullableStringArray                //
	Languages   gosql.NullableOrderedNumberArray[uint64] //
	Browsers    gosql.NullableOrderedNumberArray[uint64] //
	OS          gosql.NullableOrderedNumberArray[uint64] //
	DeviceTypes gosql.NullableOrderedNumberArray[uint64] //
	Devices     gosql.NullableOrderedNumberArray[uint64] //
	Hours       types.Hours                              // len(24) * bitmask in week days

	// DEBUG
	Trace        gosql.NullableStringArray
	TracePercent int
}

// CampaignFromModel convert database model to specified model
func CampaignFromModel(camp *models.Campaign, formats types.FormatsAccessor) *Campaign {
	var (
		// bids, _      = gocast.ToSiMap(camp.Bids.GetValue(), "", false)
		// geoBids      = parseGeoBids(billing.Money(camp.MaxBid), gocast.ToInterfaceSlice(mapDef(bids, "geo", nil)))
		hours, err = types.HoursByString(camp.Hours.String)
		flags      CampaignFlagType
	)

	if err != nil {
		return nil
	}

	if camp.DeletedAt.Valid {
		flags = CampaignFlagDeleted
	} else if camp.Active.IsActive() && camp.Status.IsApproved() {
		flags = CampaignFlagActive
	}

	if camp.Private.IsPrivate() {
		flags |= CampaignFlagPrivate
	}

	// Order ext bids
	// sort.Sort(geoBids)

	campaign := &Campaign{
		id:     camp.ID,
		AccID:  camp.AccountID,
		Weight: 0, // camp.Weight,
		Flags:  flags,

		DailyBudget:     billing.MoneyFloat(camp.DailyBudget),
		Budget:          billing.MoneyFloat(camp.Budget),
		DailyTestBudget: billing.MoneyFloat(camp.DailyTestBudget),
		TestBudget:      billing.MoneyFloat(camp.TestBudget),
		// MaxBid:          billing.Money(camp.MaxBid),

		Ads:   nil,
		Links: nil,

		Context:      camp.Context,
		Tags:         nil,
		Zones:        camp.Zones.Ordered(),
		Domains:      camp.Domains,
		Categories:   camp.Categories.Ordered(),
		Countries:    geo.CountryCodes2IDs(camp.Geos),
		Languages:    languages.LangCodes2IDs(camp.Languages),
		Browsers:     camp.Browsers.Ordered(),
		OS:           camp.OS.Ordered(),
		DeviceTypes:  camp.DeviceTypes.Ordered(),
		Devices:      camp.Devices.Ordered(),
		Hours:        hours,
		Sex:          camp.Sex.Ordered(),
		Age:          camp.Age.Ordered(),
		Trace:        camp.Trace,
		TracePercent: camp.TracePercent,
	}

	campaign.Ads = parseAds(campaign, camp, formats)
	if len(camp.Links) > 0 {
		campaign.Links = make([]AdLink, 0, len(camp.Links))

		// Assign links
		for _, link := range camp.Links {
			campaign.Links = append(campaign.Links, AdLink{
				ID:   link.ID,
				Link: link.Link,
			})
		}
	}

	// supported types
	for _, ad := range campaign.Ads {
		campaign.FormatSet.Set(uint(ad.Format.ID))
	}

	return campaign
}

// ID of the object
func (c *Campaign) ID() uint64 {
	return c.id
}

// ObjectKey of the object
func (c *Campaign) ObjectKey() uint64 {
	return c.id
}

// Account object
func (c *Campaign) Account() *Account {
	return c.Acc
}

// AccountID of current target
func (c *Campaign) AccountID() uint64 {
	return c.AccID
}

// SetAccount for target
func (c *Campaign) SetAccount(acc *Account) {
	c.Acc = acc
}

// ProjectID number
func (c *Campaign) ProjectID() uint64 {
	return 0
}

// Test campaign by pointer target
func (c *Campaign) Test(pointer types.TargetPointer) bool {
	// Skip if invalid campaign
	if c == nil || c.Acc == nil {
		return false
	}

	// Check tags
	if c.Tags.Len() > 0 && !c.Tags.OneOf(pointer.Tags()) {
		return false
	}

	// Check zones
	if c.Zones.Len() > 0 && c.Zones.IndexOf(pointer.TargetID()) == -1 {
		return false
	}

	// Check domains
	if c.Domains.Len() > 0 && !c.Domains.OneOf(pointer.Domain()) {
		return false
	}

	// Check gender
	if c.Sex.Len() > 0 && c.Sex.IndexOf(pointer.Sex()) == -1 {
		return false
	}

	// Check age (TODO: implement range processing, e.g., 0-10 years, 10-20, 20-25, etc.)
	if c.Age.Len() > 0 && c.Age.IndexOf(pointer.Age()) == -1 {
		return false
	}

	// Check categories
	if c.Categories.Len() > 0 && c.Categories.IndexOf(pointer.GeoID()) == -1 {
		return false
	}

	// Check cities
	if c.Cities.Len() > 0 && c.Cities.IndexOf(pointer.City()) == -1 {
		return false
	}

	// Check countries
	if c.Countries.Len() > 0 && c.Countries.IndexOf(pointer.GeoID()) == -1 {
		return false
	}

	// Check languages
	if c.Languages.Len() > 0 && c.Languages.IndexOf(pointer.LanguageID()) == -1 {
		return false
	}

	// Check browsers
	if c.Browsers.Len() > 0 && c.Browsers.IndexOf(pointer.BrowserID()) == -1 {
		return false
	}

	// Check operating systems
	if c.OS.Len() > 0 && c.OS.IndexOf(pointer.OSID()) == -1 {
		return false
	}

	// Check device types
	if c.DeviceTypes.Len() > 0 && c.DeviceTypes.IndexOf(pointer.DeviceType()) == -1 {
		return false
	}

	// Check devices
	if c.Devices.Len() > 0 && c.Devices.IndexOf(pointer.DeviceID()) == -1 {
		return false
	}

	// Check if the current time is within the campaign's active hours
	if !c.Hours.IsAllActive() && !c.Hours.TestTime(fasttime.Now()) {
		return false
	}

	// All checks passed
	return true
}

// State of the campaign
func (c *Campaign) State() State {
	if c == nil {
		return nil
	}
	return c.CurrentState
}

// TargetBid by targeting pointer
func (c *Campaign) TargetBid(pointer types.TargetPointer) TargetBid {
	var (
		list    []TargetBid
		tb, tb2 TargetBid
	)

	for _, ad := range c.Ads {
		hasIt := false

		// Check if format equaly replaceble or equal to the target format
		if pointer.FormatBitset().Has(uint(ad.Format.ID)) {
			hasIt = true
		} else {
			for _, f := range pointer.Formats() {
				if ad.Format.SuitsCompare(f) == 0 {
					hasIt = true
					break
				}
			}
		}

		if !hasIt {
			continue
		}
		if tb2 = ad.TargetBid(pointer); tb2.Ad == nil {
			continue
		}

		if tb.Ad == nil || tb.ECPM < tb2.ECPM-tb2.ECPM/20 {
			tb = tb2
			list = list[:0]
		} else if tb.ECPM-tb.ECPM/20 <= tb2.ECPM {
			if len(list) == 0 {
				list = append(list, tb, tb2)
			} else {
				list = append(list, tb2)
			}

			if tb.ECPM < tb2.ECPM {
				tb = tb2
			}
		}
	}

	// Choise random weighted banner
	// TODO: add weight ad check list flag
	// if c.Flags.Has(CampaignFlagHasWeighedAds) {
	// 	return TargetBidList(list).Random()
	// }
	return TargetBidList(list).Weighted()
}

// VirtualAd by targeting pointer
func (c *Campaign) VirtualAd(pointer types.TargetPointer) *VirtualAd {
	if bid := c.TargetBid(pointer); bid.Ad != nil {
		return &VirtualAd{
			Ad:       bid.Ad,
			Campaign: c,
			Bid:      bid,
		}
	}
	return nil
}

// VirtualAds for target
func (c *Campaign) VirtualAds(pointer types.TargetPointer) (ads *VirtualAds) {
	for ad := range c.VirtualAdsList(pointer) {
		if ads == nil {
			ads = &VirtualAds{Campaign: c}
		}
		ads.Bids = append(ads.Bids, ad.Bid)
	}
	return ads
}

// VirtualAdsList stream
func (c *Campaign) VirtualAdsList(pointer types.TargetPointer) <-chan *VirtualAd {
	ch := make(chan *VirtualAd, 2)

	go func() {
		for _, ad := range c.Ads {
			if !pointer.FormatBitset().Has(uint(ad.Format.ID)) {
				continue
			}

			suitable := false
			if suitable = !ad.Format.IsCloned(); !suitable {
				w, h := pointer.Size()
				suitable = ad.Format.SuitsCompareSize(w, h, 0, 0) == 0
			}

			if !suitable {
				continue
			}

			if bid := ad.TargetBid(pointer); bid.Ad != nil {
				ch <- &VirtualAd{Ad: bid.Ad, Campaign: c, Bid: bid}
			}
		}
		close(ch)
	}()

	return (<-chan *VirtualAd)(ch)
}

///////////////////////////////////////////////////////////////////////////////
/// Base Actions
///////////////////////////////////////////////////////////////////////////////

// RandomAd from list
//
//go:inline
func (c *Campaign) RandomAd() *Ad {
	return c.Ads[rand.Intn(len(c.Ads))]
}

///////////////////////////////////////////////////////////////////////////////
/// Checks
///////////////////////////////////////////////////////////////////////////////

// Validate campaign
func (c *Campaign) Validate() error {
	if len(c.Ads) < 1 {
		return ErrInvalidCampaignAds
	}

	for _, ad := range c.Ads {
		if err := ad.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Active campaign
//
//go:inline
func (c *Campaign) Active() bool { return c.Flags.IsActive() }

// Deleted campaign
//
//go:inline
func (c *Campaign) Deleted() bool { return c.Flags.IsDeleted() }

// Private campaign
//
//go:inline
func (c *Campaign) Private() bool { return c.Flags.IsPrivate() }

// Premium campaign
//
//go:inline
func (c *Campaign) Premium() bool { return c.Flags.IsPremium() }

// TestHour active
//
//go:inline
func (c *Campaign) TestHour(t time.Time) bool { return c.Hours.TestTime(t) }

// TestFormatSet of the campaign
//
//go:inline
func (c *Campaign) TestFormatSet(formatIDSet *searchtypes.NumberBitset[uint]) bool {
	return c.FormatSet.Mask()&formatIDSet.Mask() != 0
}

// TraceExperiment state
//
//go:inline
func (c *Campaign) TraceExperiment(experiment string) bool {
	return c.Trace.IndexOf(experiment) >= 0 &&
		(c.TracePercent <= 0 || rand.Intn(100) <= c.TracePercent)
}

///////////////////////////////////////////////////////////////////////////////
/// Balance
///////////////////////////////////////////////////////////////////////////////

// TestBalanceState returns true if the budget is valid for the specified format
func (c *Campaign) TestBalanceState(formatIDSet *searchtypes.NumberBitset[uint], secure bool) bool {
	if c.FormatSet.Mask()&formatIDSet.Mask() == 0 || !c.TestBudgetValue() {
		return false
	}
	for _, ad := range c.Ads {
		if formatIDSet.Has(uint(ad.Format.ID)) && (!secure || ad.Secure()) && ad._TestBudgetValue() {
			return true
		}
	}
	return false
}

// TestTestBudgetValue returns true if test budget is valid
func (c *Campaign) TestTestBudgetValue() bool {
	if c.CurrentState == nil {
		return true // Not balance state provided
	}
	return true &&
		// Total test with profit
		(c.TestBudget <= 0 || c.TestBudget >= c.CurrentState.TotalSpend()) &&
		// test daily with profit
		(c.DailyTestBudget <= 0 || c.DailyTestBudget >= c.CurrentState.Spend())
}

// TestBudgetValue return true if budget is valid
func (c *Campaign) TestBudgetValue() bool {
	if c.CurrentState == nil {
		return true // Not balance state provided
	}
	return true &&
		// Total budget test
		(c.Budget <= 0 || c.Budget >= c.CurrentState.TotalSpend()) &&
		// Daily budget test
		(c.DailyBudget <= 0 || c.DailyBudget >= c.CurrentState.Spend())
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func parseAds(newCampaign *Campaign, camp *models.Campaign, formats types.FormatsAccessor) (ads []*Ad) {
	ads = make([]*Ad, 0, len(camp.Ads))
	for _, adBase := range camp.Ads {
		if ad, err := parseAd(newCampaign, adBase, formats); err == nil {
			ad.Campaign = newCampaign
			ads = append(ads, ad)
		} else {
			log.Print("[parseAds]", err)
		}
	}
	return ads
}
