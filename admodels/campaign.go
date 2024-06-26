//
// @project GeniusRabbit rotator 2016 – 2019, 2021
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019, 2021
//

package admodels

import (
	"errors"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/geniusrabbit/gogeo"
	"github.com/geniusrabbit/gosql/v2"

	"geniusrabbit.dev/adcorelib/admodels/types"
	"geniusrabbit.dev/adcorelib/billing"
	"geniusrabbit.dev/adcorelib/i18n/languages"
	"geniusrabbit.dev/adcorelib/models"
	"geniusrabbit.dev/adcorelib/searchtypes"
)

// Errors set
var (
	ErrInvalidCampaignAds = errors.New("invalid campaigns ads")
)

// Flags set
const (
	CampaignFlagActive  = 1 << iota // 0x01
	CampaignFlagDeleted             //
	CampaignFlagPrivate             // Private campaigns not avalable for publick usage
	CampaignFlagPremium
)

// CampaignCamparator interface for index
type CampaignCamparator interface {
	CompareCampaign(c *Campaign) int
}

// Campaign model
type Campaign struct {
	ID        uint64
	Company   *Company
	CompanyID uint64

	Weight uint8
	Flags  uint8

	DailyTestBudget billing.Money // Test money amount a day (it stops automaticaly if not profit for this amount)
	TestBudget      billing.Money // Test money amount for the whole period
	DailyBudget     billing.Money
	Budget          billing.Money
	// State           balance.State

	// List of ads
	Ads   []*Ad
	Links []AdLink

	// Targeting
	FormatSet   searchtypes.UIntBitset                 //
	Context     gosql.NullableJSON[map[string]any]     //
	Keywords    gosql.NullableStringArray              //
	Zones       gosql.NullableOrderedNumberArray[uint] //
	Domains     gosql.NullableStringArray              // site domains or application bundels
	Sex         gosql.NullableOrderedNumberArray[uint] //
	Age         gosql.NullableOrderedNumberArray[uint] //
	Categories  gosql.NullableOrderedNumberArray[uint] //
	Countries   gosql.NullableOrderedNumberArray[uint] //
	Cities      gosql.NullableStringArray              //
	Languages   gosql.NullableOrderedNumberArray[uint] //
	Browsers    gosql.NullableOrderedNumberArray[uint] //
	Os          gosql.NullableOrderedNumberArray[uint] //
	DeviceTypes gosql.NullableOrderedNumberArray[uint] //
	Devices     gosql.NullableOrderedNumberArray[uint] //
	Hours       types.Hours                            // len(24) * bitmask in week days

	// DEBUG
	Trace        gosql.NullableStringArray
	TracePercent int
}

// CampaignFromModel convert database model to specified model
func CampaignFromModel(camp *models.Campaign, formats types.FormatsAccessor) *Campaign {
	var (
		countriesArr gosql.NullableOrderedNumberArray[uint]
		languagesArr gosql.NullableOrderedNumberArray[uint]
		// bids, _      = gocast.ToSiMap(camp.Bids.GetValue(), "", false)
		// geoBids      = parseGeoBids(billing.Money(camp.MaxBid), gocast.ToInterfaceSlice(mapDef(bids, "geo", nil)))
		hours, err = types.HoursByString(camp.Hours.String)
		flags      uint8
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

	// Countries filter
	if camp.Geos.Len() > 0 {
		seted := map[string]bool{}
		for _, cc := range camp.Geos {
			cc = strings.ToUpper(cc)
			if !seted[cc] {
				seted[cc] = true
				switch cc {
				case "EU", "AS", "AF", "OC", "SA", "NA", "AN":
					for _, country := range gogeo.Countries {
						if country.Continent == cc {
							seted[country.Code2] = true
							countriesArr = append(countriesArr, uint(country.ID))
						}
					}
				default: // ** - as undefined
					countriesArr = append(countriesArr, uint(gogeo.CountryByCode2(cc).ID))
				}
			}
		}
		countriesArr.Sort()
	}

	// Languages filter
	if len(camp.Languages) > 0 {
		for _, lg := range camp.Languages {
			languagesArr = append(languagesArr, languages.GetLanguageIdByCodeString(lg))
		}
		languagesArr.Sort()
	}

	// Order ext bids
	// sort.Sort(geoBids)

	campaign := &Campaign{
		// MaxBid:      billing.Money(camp.MaxBid),
		ID:        camp.ID,
		CompanyID: camp.CompanyID,
		Weight:    0, // camp.Weight,
		Flags:     flags,

		DailyBudget:     billing.MoneyFloat(camp.DailyBudget),
		Budget:          billing.MoneyFloat(camp.Budget),
		DailyTestBudget: billing.MoneyFloat(camp.DailyTestBudget),
		TestBudget:      billing.MoneyFloat(camp.TestBudget),

		Ads:   nil,
		Links: nil,

		Context:      camp.Context,
		Keywords:     nil,
		Zones:        camp.Zones.Ordered(),
		Domains:      camp.Domains,
		Categories:   camp.Categories.Ordered(),
		Countries:    countriesArr,
		Languages:    languagesArr,
		Browsers:     camp.Browsers.Ordered(),
		Os:           camp.Os.Ordered(),
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

// GetID of the object
func (c *Campaign) GetID() uint64 {
	return c.ID
}

// ProjectID number
func (c *Campaign) ProjectID() uint64 {
	return 0
}

// // State of the campaign
// func (c *Campaign) State() State {
// 	if c == nil {
// 		return nil
// 	}
// 	return c.State
// }

// TargetBid by targeting pointer
func (c *Campaign) TargetBid(pointer types.TargetPointer) TargetBid {
	var (
		list []TargetBid
		tb   TargetBid
	)

	for _, ad := range c.Ads {
		hasIt := false
		for _, f := range pointer.Formats() {
			if ad.Format.SuitsCompare(f) == 0 {
				hasIt = true
				break
			}
		}

		if !hasIt || pointer.FormatBitset().Has(uint(ad.Format.ID)) {
			continue
		}

		tb2 := ad.TargetBid(pointer)
		if tb2.Ad == nil {
			continue
		}

		if tb.Ad == nil || tb.ECPM < tb2.ECPM-tb2.ECPM/20 {
			tb = tb2
			list = nil
		} else if tb.ECPM-tb.ECPM/20 <= tb2.ECPM {
			if list == nil {
				list = []TargetBid{tb, tb2}
			} else {
				list = append(list, tb2)
			}

			if tb.ECPM < tb2.ECPM {
				tb = tb2
			}
		}
	}

	// Choise random banner
	if len(list) > 0 {
		if len(list) < 2 {
			tb = list[0]
		} else {
			tb = list[rand.Intn(len(list))]
		}
	}

	return tb
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
func (c *Campaign) VirtualAds(pointer types.TargetPointer) *VirtualAds {
	var ads *VirtualAds
	for ad := range c.VirtualAdsList(pointer) {
		if ads == nil {
			ads = &VirtualAds{Campaign: ad.Campaign}
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

			var suitable bool
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
func (c *Campaign) Active() bool {
	return c.Flags&CampaignFlagActive != 0
}

// Deleted campaign
func (c *Campaign) Deleted() bool {
	return c.Flags&CampaignFlagDeleted != 0
}

// Private campaign
func (c *Campaign) Private() bool {
	return c.Flags&CampaignFlagPrivate != 0
}

// Premium campaign
func (c *Campaign) Premium() bool {
	return c.Flags&CampaignFlagPremium != 0
}

// TestHour active
func (c *Campaign) TestHour(t time.Time) bool {
	return c.Hours.TestTime(t)
}

// // TestMoneyState of the campaign
// func (c *Campaign) TestMoneyState(formatIDSet *searchtypes.UIntBitset, secure bool) bool {
// 	if c.FormatSet.Mask()&formatIDSet.Mask() == 0 || !c.TestBudgetValue() || !c.TestProfit() {
// 		return false
// 	}
// 	for _, ad := range c.Ads {
// 		if formatIDSet.Has(uint(ad.Format.ID)) && ad.TestBudgetValues() && ad.TestProfit() && (!secure || ad.Secure()) {
// 			return true
// 		}
// 	}
// 	return false
// }

// TestFormatSet of the campaign
func (c *Campaign) TestFormatSet(formatIDSet *searchtypes.UIntBitset) bool {
	return c.FormatSet.Mask()&formatIDSet.Mask() != 0
}

// // TestBudgetValue of campaign
// func (c *Campaign) TestBudgetValue() bool {
// 	return (c.GetDailyBudget() <= 0 || c.GetSpent() < c.GetDailyBudget()) &&
// 		(c.Budget <= 0 || c.GetTotalSpent() < c.Budget)
// }

// // TestProfit of the campaign
// func (c *Campaign) TestProfit() bool {
// 	return true &&
// 		// Total test with profit
// 		(c.TestBudget <= 0 || c.TestBudget >= c.GetTotalSpent()-c.GetTotalProfit()) &&
// 		// test daily with profit
// 		(c.DailyTestBudget <= 0 || c.DailyTestBudget >= c.GetSpent()-c.GetProfit())
// }

// // UpdateBalance from ads
// func (c *Campaign) UpdateBalance() {
// 	var spent billing.Money
// 	for _, ad := range c.Ads {
// 		spent += ad.GetSpent()
// 	}
// 	c.State.SetSpent(spent)
// }

// TraceExperiment state
func (c *Campaign) TraceExperiment(experiment string) bool {
	return c.Trace.IndexOf(experiment) >= 0 &&
		(c.TracePercent <= 0 || rand.Intn(100) <= c.TracePercent)
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
