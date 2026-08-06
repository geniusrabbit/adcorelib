package adtype

import (
	"strings"

	"github.com/demdxx/gocast/v2"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// MacroMapper is a struct that contains the macro mapper for the ad content
type MacroMapper struct {
	AuctionID   string
	AuctionType string
	ImpID       string
	ClickID     string

	CampaignID string
	AdID       string

	Platform   string
	PriceModel string

	Country   string
	City      string
	Latitude  string
	Longitude string
	Language  string

	Domain   string
	AppID    string
	ZoneID   string
	ZoneCode string

	UUID      string
	SessionID string
	UserIP    string
	CarrierID string

	DeviceType  string
	DeviceID    string
	DeviceName  string
	OSID        string
	OSName      string
	BrowserID   string
	BrowserName string

	Price float64
}

func (m *MacroMapper) Prepare(extra map[string]string) *strings.Replacer {
	args := []string{
		// Auction identifiers
		"{{auctionid}}", m.AuctionID,
		"{{auction.id}}", m.AuctionID,
		"{{auc_id}}", m.AuctionID,
		"{{auction_type}}", m.AuctionType,
		"{{impid}}", m.ImpID,
		"{{imp_id}}", m.ImpID,
		"{{click_id}}", m.ClickID,
		"{{click.id}}", m.ClickID,

		// Campaign and Ad
		"{{camp_id}}", m.CampaignID,
		"{{camp.id}}", m.CampaignID,
		"{{campaign_id}}", m.CampaignID,
		"{{campaign.id}}", m.CampaignID,
		"{{ad_id}}", m.AdID,
		"{{ad.id}}", m.AdID,

		// Platform and Price Model
		"{{platform}}", m.Platform,
		"{{pm}}", m.PriceModel,

		// Country, City, Latitude, Longitude, Language, Domain
		"{{cc}}", m.Country,
		"{{country}}", m.Country,
		"{{city}}", m.City,
		"{{lang}}", m.Language,
		"{{lat}}", m.Latitude,
		"{{latitude}}", m.Latitude,
		"{{lon}}", m.Longitude,
		"{{longitude}}", m.Longitude,

		// Placement identifiers
		"{{domain}}", m.Domain,
		"{{app_id}}", m.AppID,
		"{{appid}}", m.AppID,
		"{{zone_id}}", m.ZoneID,
		"{{zoneid}}", m.ZoneID,
		"{{zone_code}}", m.ZoneCode,
		"{{zonecode}}", m.ZoneCode,

		// User identifiers
		"{{uuid}}", m.UUID,
		"{{sessid}}", m.SessionID,
		"{{ip}}", m.UserIP,

		// Device identifiers
		"{{carrier_id}}", m.CarrierID,
		"{{device_type}}", m.DeviceType,
		"{{device_id}}", m.DeviceID,
		"{{device_name}}", m.DeviceName,
		"{{os_id}}", m.OSID,
		"{{os_name}}", m.OSName,
		"{{os}}", m.OSName,
		"{{browser_id}}", m.BrowserID,
		"{{browser_name}}", m.BrowserName,
		"{{browser}}", m.BrowserName,
	}
	for k, v := range extra {
		args = append(args, "{{"+k+"}}", v)
	}
	return strings.NewReplacer(args...)
}

// ContentPreparer returns a strings.Replacer that can be used to replace macros in ad content with actual values from the response and item. It includes standard macros like {impid}, {aucid}, {adid}, as well as any custom mappings provided by the item.
func ContentPreparer(response Response, item ResponseItem) *strings.Replacer {
	var (
		mapping    map[string]string
		req        = response.Request()
		imp        = item.Impression()
		targetID   = ""
		targetCode = ""
		prep, _    = item.(interface{ ContentMapping() map[string]string })
	)
	if imp != nil && imp.Target != nil {
		targetID = gocast.Str(imp.Target.ID())
		targetCode = imp.Target.Codename()
	}
	if prep != nil {
		mapping = prep.ContentMapping()
	}
	mapper := MacroMapper{
		// Auction identifiers
		AuctionID:   response.AuctionID(),
		AuctionType: response.AuctionType().Name(),
		ImpID:       item.ImpressionID(),
		ClickID:     item.ImpressionID(),

		// Campaign and Ad
		CampaignID: gocast.Str(item.CampaignID()),
		AdID:       item.AdID(),

		// Platform and Price Model
		Platform:   types.PlatformType(req.DeviceInfo().DeviceType).Name(),
		PriceModel: item.PricingModel().Name(),

		// Country, City, Latitude, Longitude, Language
		Country:   req.GeoInfo().Country,
		City:      req.GeoInfo().City,
		Language:  req.BrowserInfo().PrimaryLanguage,
		Latitude:  gocast.Str(req.GeoInfo().Lat),
		Longitude: gocast.Str(req.GeoInfo().Lon),

		// Placement identifiers
		Domain:   req.DomainName(),
		AppID:    gocast.Str(req.AppID()),
		ZoneID:   targetID,
		ZoneCode: targetCode,

		// User identifiers
		UUID:      req.UserInfo().ID,
		SessionID: req.UserInfo().SessionID,
		UserIP:    req.GeoInfo().IP.String(),

		// Device identifiers
		CarrierID:   gocast.Str(req.CarrierInfo().ID),
		DeviceType:  types.PlatformType(req.DeviceInfo().DeviceType).Name(),
		DeviceID:    gocast.Str(req.DeviceInfo().ID),
		DeviceName:  "",
		OSID:        gocast.Str(req.OSInfo().ID),
		OSName:      req.OSInfo().Name,
		BrowserID:   gocast.Str(req.BrowserInfo().ID),
		BrowserName: req.BrowserInfo().Name,

		Price: item.Price(ActionImp).Float64(),
	}
	return mapper.Prepare(mapping)
}

// ContentMappingPreparer returns a strings.Replacer that can be used to replace macros in ad content with actual values from the item's content mapping. It looks for a ContentMapping method on the item and uses it to create the replacer. If no ContentMapping is provided, it returns nil.
func ContentMappingPreparer(response Response, item ResponseItem) *strings.Replacer {
	if prep := item.(interface{ ContentMapping() map[string]string }); prep != nil {
		mapping := prep.ContentMapping()
		args := make([]string, 0, len(mapping)*2)
		for k, v := range mapping {
			args = append(args, k, v)
		}
		return strings.NewReplacer(args...)
	}
	return nil
}
