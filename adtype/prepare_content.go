package adtype

import (
	"strings"

	"github.com/demdxx/gocast/v2"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// ContentPreparer returns a strings.Replacer that can be used to replace macros in ad content with actual values from the response and item. It includes standard macros like {impid}, {aucid}, {adid}, as well as any custom mappings provided by the item.
func ContentPreparer(response Response, item ResponseItem) *strings.Replacer {
	var (
		req        = response.Request()
		imp        = item.Impression()
		aucID      = response.AuctionID()
		prep       = item.(interface{ ContentMapping() map[string]string })
		targetID   uint64
		targetCode string
		campaignID = gocast.Str(item.CampaignID())
		args       = []string{
			"{impid}", item.ImpressionID(),
			"{aucid}", aucID,
			"{auctionid}", aucID,
			"{auction.id}", aucID,
			"{auc.id}", aucID,
			"{auctype}", response.AuctionType().Name(),
			"{adid}", item.AdID(),
			"{ad.id}", item.AdID(),
			"{campid}", campaignID,
			"{camp.id}", campaignID,
			"{campaignid}", campaignID,
			"{campaign.id}", campaignID,

			"{platform}", "",
			"{unit_id}", gocast.Str(targetID),
			"{unit_code}", targetCode,
			"{jumper_id}", "",
			"{pm}", item.PricingModel().Name(),

			"{country}", req.GeoInfo().Country,
			"{city}", req.GeoInfo().City,
			"{lang}", req.BrowserInfo().PrimaryLanguage,
			"{domain}", req.DomainName(),

			"{udid}", req.DeviceInfo().IFA,
			"{uuid}", req.UserInfo().ID,
			"{sessid}", req.UserInfo().SessionID,
			"{fingerprint}", req.UserInfo().FingerPrintID,
			"{etag}", req.UserInfo().ETag,
			"{ip}", req.GeoInfo().IP.String(),
			"{carrier_id}", "",
			"{latitude}", "",
			"{longitude}", "",
			"{device_type}", types.PlatformType(req.DeviceInfo().DeviceType).Name(),
			"{device_id}", gocast.Str(req.DeviceInfo().ID),
			"{os_id}", gocast.Str(req.OSInfo().ID),
			"{browser_id}", gocast.Str(req.BrowserInfo().ID),
		}
	)
	if imp != nil && imp.Target != nil {
		targetID = imp.Target.ID()
		targetCode = imp.Target.Codename()
	}
	if prep != nil {
		mapping := prep.ContentMapping()
		for k, v := range mapping {
			args = append(args, k, v)
		}
	}
	return strings.NewReplacer(args...)
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
