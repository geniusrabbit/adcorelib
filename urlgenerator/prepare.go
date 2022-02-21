package urlgenerator

import (
	"net/url"
	"strings"

	"geniusrabbit.dev/corelib/admodels/types"
	"geniusrabbit.dev/corelib/eventtraking/events"
	"github.com/demdxx/gocast"
)

// PrepareURL by event
func PrepareURL(url string, event *events.Event, params url.Values) string {
	replacer := strings.NewReplacer(
		"{country}", event.Country,
		"{city}", event.City,
		"{lang}", event.Language,
		"{domain}", event.Domain,
		"{impid}", event.ImpID,
		"{aucid}", event.AuctionID,
		"{platform}", types.PlatformType(event.Platform).Name(),
		"{zone_id}", gocast.ToString(event.Zone),
		"{jumper_id}", gocast.ToString(event.Jumper),
		"{pm}", types.PricingModel(event.PricingModel).Name(),
		"{udid}", event.UDID,
		"{uuid}", event.UUID,
		"{sessid}", event.SessionID,
		"{fingerprint}", event.Fingerprint,
		"{etag}", event.ETag,
		"{ip}", event.IPString,
		"{carrier_id}", gocast.ToString(event.Carrier),
		"{latitude}", event.Latitude,
		"{longitude}", event.Longitude,
		"{longitude}", event.Longitude,
		"{device_type}", types.PlatformType(event.DeviceType).Name(),
		"{device}", gocast.ToString(event.Device),
		"{os_id}", gocast.ToString(event.OS),
		"{browser_id}", gocast.ToString(event.Browser),
	)
	return replacer.Replace(url)
}
