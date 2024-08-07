package urlgenerator

import (
	"strings"

	"github.com/demdxx/gocast/v2"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
)

// PrepareURL by event
func PrepareURL(url string, event *events.Event) string {
	replacer := strings.NewReplacer(
		"{country}", event.Country,
		"{city}", event.City,
		"{lang}", event.Language,
		"{domain}", event.Domain,
		"{impid}", event.ImpID,
		"{aucid}", event.AuctionID,
		"{auctype}", types.AuctionType(event.AuctionType).Name(),
		"{platform}", types.PlatformType(event.Platform).Name(),
		"{zone_id}", gocast.Str(event.Zone),
		"{jumper_id}", gocast.Str(event.Jumper),
		"{pm}", types.PricingModel(event.PricingModel).Name(),
		"{udid}", event.UDID,
		"{uuid}", event.UUID,
		"{sessid}", event.SessionID,
		"{fingerprint}", event.Fingerprint,
		"{etag}", event.ETag,
		"{ip}", event.IPString,
		"{carrier_id}", gocast.Str(event.Carrier),
		"{latitude}", event.Latitude,
		"{longitude}", event.Longitude,
		"{device_type}", types.PlatformType(event.DeviceType).Name(),
		"{device_id}", gocast.Str(event.Device),
		"{os_id}", gocast.Str(event.OS),
		"{browser_id}", gocast.Str(event.Browser),
	)
	return replacer.Replace(url)
}
