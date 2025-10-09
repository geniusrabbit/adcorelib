package adtype

import (
	"strings"

	"github.com/demdxx/gocast/v2"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

// PrepareURL by event
func PrepareURL(url string, response Response, it ResponseItem) string {
	var (
		req        = response.Request()
		imp        = it.Impression()
		targetID   uint64
		targetCode string
	)
	if imp != nil && imp.Target != nil {
		targetID = imp.Target.ID()
		targetCode = imp.Target.Codename()
	}
	replacer := strings.NewReplacer(
		"{country}", req.GeoInfo().Country,
		"{city}", req.GeoInfo().City,
		"{lang}", req.BrowserInfo().PrimaryLanguage,
		"{domain}", req.DomainName(),
		"{impid}", it.ImpressionID(),
		"{aucid}", response.AuctionID(),
		"{auctype}", response.AuctionType().Name(),
		"{platform}", "",
		"{unit_id}", gocast.Str(targetID),
		"{unit_code}", targetCode,
		"{jumper_id}", "",
		"{pm}", it.PricingModel().Name(),
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
	)
	return replacer.Replace(url)
}
