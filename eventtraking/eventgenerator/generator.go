//
// @project GeniusRabbit corelib 2018 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2018 - 2019
//

package eventgenerator

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/demdxx/gocast/v2"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
)

// Errors set
var (
	ErrInvalidMultipleItemAsSingle = errors.New("can`t convert multipleitem to single action")
)

// Generator object
type Generator interface {
	// Event object by response
	Event(event events.Type, status uint8, response adtype.Responser, it adtype.ResponserItem) (*events.Event, error)

	// Events object list
	Events(event events.Type, status uint8, response adtype.Responser, it adtype.ResponserItemCommon) []*events.Event

	// UserInfo event object by response
	UserInfo(response adtype.Responser, it adtype.ResponserItem) (*events.UserInfo, error)
}

type generator struct {
	service string
}

// New generator object
func New(service string) Generator {
	return generator{service: service}
}

// Event object by response
func (g generator) Event(event events.Type, status uint8, response adtype.Responser, it adtype.ResponserItem) (*events.Event, error) {
	var (
		r             = response.Request()
		imp           = it.Impression()
		sourceID      uint64
		zoneID        uint64
		accessPointID uint64
	)

	if src := it.Source(); src != nil {
		sourceID = src.ID()
	}

	if sourceID == 0 && response.Source() != nil {
		sourceID = response.Source().ID()
	}

	if imp != nil && imp.Target != nil {
		zoneID = imp.Target.ID()
	}

	if _, ok := it.(adtype.ResponserMultipleItem); ok {
		return nil, ErrInvalidMultipleItemAsSingle
	}

	if response.Request().AccessPoint != nil {
		accessPointID = response.Request().AccessPoint.ID()
	}

	fmt.Println(">> EVENT", event, status, "\n\t",
		"PurchaseViewPrice", fmt.Sprintf("%.7f", it.PurchasePrice(admodels.ActionView).Float64()), "\n\t",
		"PurchaseClickPrice", fmt.Sprintf("%.7f", it.PurchasePrice(admodels.ActionClick).Float64()), "\n\t",
		"PurchaseLeadPrice", fmt.Sprintf("%.7f", it.PurchasePrice(admodels.ActionLead).Float64()), "\n\t",
		"PotentialViewPrice", fmt.Sprintf("%.7f", it.PotentialPrice(admodels.ActionView).Float64()), "\n\t",
		"PotentialClickPrice", fmt.Sprintf("%.7f", it.PotentialPrice(admodels.ActionClick).Float64()), "\n\t",
		"PotentialLeadPrice", fmt.Sprintf("%.7f", it.PotentialPrice(admodels.ActionLead).Float64()), "\n\t",
		"ViewPrice", fmt.Sprintf("%.7f", it.Price(admodels.ActionView).Float64()), "\n\t",
		"ClickPrice", fmt.Sprintf("%.7f", it.Price(admodels.ActionClick).Float64()), "\n\t",
		"LeadPrice", fmt.Sprintf("%.7f", it.Price(admodels.ActionLead).Float64()), "\n\t",
	)

	return &events.Event{
		Time:     time.Now().UnixNano(),
		Delay:    0,
		Duration: 0,         //
		Service:  g.service, // Service
		Event:    event,     // Action code (tech param, Do not store)
		Status:   status,    //

		// Accounts link information
		Project:           0,               // Project network ID
		PublisherAccount:  imp.AccountID(), // -- // --
		AdvertiserAccount: it.AccountID(),  // -- // --

		// Source
		AuctionID:     r.ID,                          // ID of last auction
		AuctionType:   uint8(response.AuctionType()), // Aution type 1 - First price, 2 - Second price
		ImpID:         it.ImpressionID(),             // Sub ID of request for paticular impression spot
		ImpAdID:       it.ID(),                       // Specific ID for paticular ad impression
		ExtAuctionID:  r.ExtID,                       // External auction ID
		ExtImpID:      it.ExtImpressionID(),          // External auction Imp ID
		SourceID:      sourceID,                      // Advertisement Source ID
		Network:       it.NetworkName(),              // Source Network Name or Domain (Cross sails)
		AccessPointID: accessPointID,                 // Access Point ID to own Advertisement

		// State Location
		Platform:      0,                 // Where displaid? 0 – undefined, 1 – web site, 2 – native app, 3 – game
		Domain:        r.DomainName(),    //
		ApplicationID: uint64(r.AppID()), // Place target
		ZoneID:        zoneID,            // -- // --
		CampaignID:    it.CampaignID(),   // Campaign info
		FormatID:      it.Format().ID,    // Format object ID
		AdID:          it.AdID(),         // -- // --
		AdWidth:       it.Width(),        // -- // --
		AdHeight:      it.Height(),       // -- // --
		SourceURL:     "",                // Advertisement source URL (iframe, image, video, direct)
		WinURL:        "",                // Win URL used for RTB confirmation
		URL:           it.ActionURL(),    // Non modified target URL
		JumperID:      0,                 // Jumper Page ID

		// Money
		PricingModel:         it.PricingModel().UInt(),                        // Display As CPM/CPC/CPA/CPI
		TestPriceMode:        it.PriceTestMode(),                              // Test price mode
		PurchaseViewPrice:    it.PurchasePrice(admodels.ActionView).Int64(),   // Price of of the view of source traffic cost
		PurchaseClickPrice:   it.PurchasePrice(admodels.ActionClick).Int64(),  // Price of of the click of source traffic cost
		PurchaseLeadPrice:    it.PurchasePrice(admodels.ActionLead).Int64(),   // Price of of the lead of source traffic cost
		PotentialViewPrice:   it.PotentialPrice(admodels.ActionView).Int64(),  // Price of of the view of source traffic cost including descrepancy correction
		PotentialClickPrice:  it.PotentialPrice(admodels.ActionClick).Int64(), // Price of of the click of source traffic cost including descrepancy correction
		PotentialLeadPrice:   it.PotentialPrice(admodels.ActionLead).Int64(),  // Price of of the lead of source traffic cost including descrepancy correction
		ViewPrice:            it.FinalPrice(admodels.ActionView).Int64(),      // Price per view with total comissions and with descrepancy correction
		ClickPrice:           it.FinalPrice(admodels.ActionClick).Int64(),     // Price per click with total comissions and with descrepancy correction
		LeadPrice:            it.FinalPrice(admodels.ActionLead).Int64(),      // Price per lead with total comissions and with descrepancy correction
		CompetitorCampaignID: it.Second().GetCampaignID(),                     // Competitor compaign ID
		CompetitorSourceID:   it.Second().GetSourceID(),                       // Competitor source ID
		CompetitorECPM:       it.Second().GetECPM().Float64(),                 // Competitor ECPM or auction

		// User IDENTITY
		UDID:        r.DeviceInfo().IFA,         // Unique Device ID (IDFA)
		UUID:        r.UserInfo().ID,            // User
		SessionID:   r.UserInfo().SessionID,     // -- // --
		Fingerprint: r.UserInfo().FingerPrintID, //
		ETag:        r.UserInfo().ETag,          //

		// Targeting
		CarrierID:       r.CarrierInfo().ID,
		Country:         r.GeoInfo().Country,
		City:            r.GeoInfo().City,
		Language:        r.BrowserInfo().PrimaryLanguage,
		Referer:         r.BrowserInfo().Ref,
		IPString:        r.GeoInfo().IP.String(),
		UserAgent:       r.BrowserInfo().UA,
		DeviceID:        r.DeviceInfo().ID,
		OSID:            r.DeviceInfo().OS.ID,
		BrowserID:       uint(r.BrowserInfo().ID),
		Categories:      "",
		Adblock:         b2u(r.IsAdblock()),
		PrivateBrowsing: b2u(r.IsPrivateBrowsing()),
		Robot:           0,
		Proxy:           0,
		Backup:          0,
		X:               imp.X,
		Y:               imp.Y,
		Width:           r.Width(),
		Height:          r.Height(),

		SubID1: imp.SubID1,
		SubID2: imp.SubID2,
		SubID3: imp.SubID3,
		SubID4: imp.SubID4,
		SubID5: imp.SubID5,
	}, nil
}

// Events object list
func (g generator) Events(event events.Type, status uint8, response adtype.Responser, it adtype.ResponserItemCommon) (events []*events.Event) {
	if mit, _ := it.(adtype.ResponserMultipleItem); mit != nil {
		for _, it := range mit.Ads() {
			if event, err := g.Event(event, status, response, it); err == nil {
				events = append(events, event)
			}
		}
	} else if event, err := g.Event(event, status, response, it.(adtype.ResponserItem)); err == nil {
		events = append(events, event)
	}
	return events
}

// UserInfo event object by response
func (g generator) UserInfo(response adtype.Responser, it adtype.ResponserItem) (*events.UserInfo, error) {
	var (
		r       = response.Request()
		imp     = it.Impression()
		user    = r.UserInfo()
		geo     = r.GeoInfo()
		browser = r.BrowserInfo()
	)
	if user.Email == "" {
		return nil, nil
	}
	return &events.UserInfo{
		Time:      time.Now().UnixNano(),
		AuctionID: r.ID, // ID of last auction
		// User IDENTITY
		UDID:      r.DeviceInfo().IFA, // Unique Device ID (IDFA)
		UUID:      user.ID,            // User
		SessionID: user.SessionID,     // -- // --
		// Personal information
		Age:           user.AvgAge(),
		Gender:        byte(user.Sex()),
		SearchGender:  sex(gocast.Str(imp.Get("search_gender"))),
		Email:         user.Email,
		Phone:         user.GetDataItemOrDefault("phone", ""),
		MessangerType: user.GetDataItemOrDefault("messanger_type", ""),
		Messanger:     user.GetDataItemOrDefault("messanger", ""),
		Postcode:      geo.ZIP,
		Facebook:      user.GetDataItemOrDefault("sn.facebook", ""),
		Twitter:       user.GetDataItemOrDefault("sn.twitter", ""),
		Linkedin:      user.GetDataItemOrDefault("sn.linkedin", ""),
		// Location info
		Country:   geo.Country,                               // Country Code ISO-2
		City:      geo.City,                                  // City Code
		Latitude:  strconv.FormatFloat(geo.Lat, 'G', -1, 64), // -- // --
		Longitude: strconv.FormatFloat(geo.Lon, 'G', -1, 64), // -- // --
		Language:  browser.PrimaryLanguage,                   // en-US
	}, nil
}

func b2u(v bool) uint {
	if v {
		return 1
	}
	return 0
}

func sex(s string) (sx byte) {
	switch s {
	case "male", "m", "M":
		sx = byte(adtype.UserSexMale)
	case "female", "f", "F":
		sx = byte(adtype.UserSexFemale)
	}
	return sx
}
