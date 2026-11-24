package personification

import (
	"strconv"
	"time"

	"github.com/geniusrabbit/udetect"
	"github.com/valyala/fasthttp"

	adgeo "github.com/geniusrabbit/adcorelib/geo"
)

// ApplyCloudflare writes CloudFlare geo headers onto geo.
// A header with usable data replaces the current value. An empty or unusable
// header leaves the existing value. country is true when CF-IPCountry set the country.
// There is no outer “country missing” gate: city and region are applied on their own.
func ApplyCloudflare(geo *udetect.Geo, h *fasthttp.RequestHeader) (country bool) {
	if geo == nil || h == nil {
		return false
	}

	if found, ok := cloudflareCountry(peekHeader(h, "CF-IPCountry")); ok {
		geo.Country = found.Code2
		country = true
	}

	if rc := adgeo.RegionCodeByPartial(peekHeader(h, "CF-Region-Code"), geo.Country); rc != adgeo.UndefinedRegionCode {
		geo.Region = rc
	}

	if city := peekHeader(h, "CF-IPCity"); city != "" {
		geo.City = city
	}
	if zip := peekHeader(h, "CF-Postal-Code"); zip != "" {
		geo.ZIP = zip
	}
	if metro := peekHeader(h, "CF-Metro-Code"); metro != "" {
		geo.Metro = metro
	}

	latS := peekHeader(h, "CF-IPLatitude")
	lonS := peekHeader(h, "CF-IPLongitude")
	if latS != "" && lonS != "" {
		lat, latErr := strconv.ParseFloat(latS, 64)
		lon, lonErr := strconv.ParseFloat(lonS, 64)
		if latErr == nil && lonErr == nil {
			geo.Lat = lat
			geo.Lon = lon
		}
	}

	if off, ok := timezoneOffsetMinutes(peekHeader(h, "CF-Timezone")); ok {
		geo.SetUTCOffsetWithAutocorrection(off)
	}
	return country
}

func peekHeader(h *fasthttp.RequestHeader, key string) string {
	// fasthttp Peek normalizes the key (Cf-Ipcity); HTTP header names are case-insensitive.
	return string(h.Peek(key))
}

func cloudflareCountry(code string) (*adgeo.Country, bool) {
	switch code {
	case "", "XX", "T1":
		return nil, false
	}
	country := adgeo.CountryByCode2(code)
	if country.ISO2() == adgeo.UndefinedCountryCodeISO2 {
		return nil, false
	}
	return country, true
}

func timezoneOffsetMinutes(name string) (int, bool) {
	if name == "" {
		return 0, false
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return 0, false
	}
	_, seconds := time.Now().In(loc).Zone()
	return seconds / 60, true
}
