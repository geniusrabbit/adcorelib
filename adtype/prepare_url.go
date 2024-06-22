package adtype

import (
	"strings"

	"geniusrabbit.dev/adcorelib/admodels/types"
	"github.com/demdxx/gocast/v2"
)

// PrepareURL by event
func PrepareURL(url string, response Responser, it ResponserItem) string {
	var (
		r      = response.Request()
		imp    = it.Impression()
		zoneID uint64
	)
	if imp != nil && imp.Target != nil {
		zoneID = imp.Target.ID()
	}
	replacer := strings.NewReplacer(
		"{country}", r.GeoInfo().Country,
		"{city}", r.GeoInfo().City,
		"{lang}", r.BrowserInfo().PrimaryLanguage,
		"{domain}", r.DomainName(),
		"{impid}", it.ImpressionID(),
		"{aucid}", response.AuctionID(),
		"{auctype}", response.AuctionType().Name(),
		"{platform}", "",
		"{zone_id}", gocast.Str(zoneID),
		"{jumper_id}", "",
		"{pm}", it.PricingModel().Name(),
		"{udid}", r.DeviceInfo().IFA,
		"{uuid}", r.UserInfo().ID,
		"{sessid}", r.UserInfo().SessionID,
		"{fingerprint}", r.UserInfo().FingerPrintID,
		"{etag}", r.UserInfo().ETag,
		"{ip}", r.GeoInfo().IP.String(),
		"{carrier_id}", "",
		"{latitude}", "",
		"{longitude}", "",
		"{device_type}", types.PlatformType(r.DeviceType()).Name(),
		"{device_id}", gocast.Str(r.DeviceID()),
		"{os_id}", gocast.Str(r.OSID()),
		"{browser_id}", gocast.Str(r.BrowserID()),
	)
	return replacer.Replace(url)
}
