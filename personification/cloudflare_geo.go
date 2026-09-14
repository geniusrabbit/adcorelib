package personification

import (
	"strconv"
	"time"

	"github.com/geniusrabbit/gogeo"
	"github.com/geniusrabbit/udetect"
	"github.com/valyala/fasthttp"
)

// applyCloudflareGeo fills empty/undefined udetect.Geo fields from CloudFlare
// request headers. Detector values win: each field is written only when it is
// still empty. There is no outer “country missing” gate — city/region can be
// filled while country is already set.
func applyCloudflareGeo(geo *udetect.Geo, h *fasthttp.RequestHeader) {
	if geo == nil || h == nil {
		return
	}

	if geo.Country.ISO2() == gogeo.UndefinedCountryCodeISO2 {
		if country, ok := cloudflareCountry(peekHeader(h, "CF-IPCountry")); ok {
			geo.ID = uint(country.ID)
			geo.Country = country.Code2
		}
	}

	if geo.Region == gogeo.UndefinedRegionCode {
		if rc := cloudflareRegionCode(peekHeader(h, "CF-Region-Code"), geo.Country); rc != gogeo.UndefinedRegionCode {
			geo.Region = rc
		}
	}

	if geo.City == "" {
		geo.City = peekHeader(h, "CF-IPCity")
	}
	if geo.ZIP == "" {
		geo.ZIP = peekHeader(h, "CF-Postal-Code")
	}
	if geo.Metro == "" {
		geo.Metro = peekHeader(h, "CF-Metro-Code")
	}

	if geo.Lat == 0 && geo.Lon == 0 {
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
	}

	if geo.UTCOffset == 0 {
		if off, ok := timezoneOffsetHours(peekHeader(h, "CF-Timezone")); ok {
			geo.UTCOffset = off
		}
	}
}

func peekHeader(h *fasthttp.RequestHeader, key string) string {
	// fasthttp Peek normalizes the key (Cf-Ipcity); HTTP header names are case-insensitive.
	return string(h.Peek(key))
}

func cloudflareCountry(code string) (*gogeo.Country, bool) {
	switch code {
	case "", "XX", "T1":
		return nil, false
	}
	country := gogeo.CountryByCode2(code)
	if country.ISO2() == gogeo.UndefinedCountryCodeISO2 {
		return nil, false
	}
	return country, true
}

func cloudflareRegionCode(code string, country gogeo.Code2) gogeo.RegionCode {
	if code == "" {
		return gogeo.UndefinedRegionCode
	}
	if rc := gogeo.RegionCodeByString(code); rc != gogeo.UndefinedRegionCode {
		return rc
	}
	// OpenRTB / CF often send the subdivision only ("BA"), not "SK-BA".
	if country.ISO2() != gogeo.UndefinedCountryCodeISO2 {
		return gogeo.RegionCodeByString(country.ISO2() + "-" + code)
	}
	return gogeo.UndefinedRegionCode
}

func timezoneOffsetHours(name string) (int, bool) {
	if name == "" {
		return 0, false
	}
	loc, err := time.LoadLocation(name)
	if err != nil {
		return 0, false
	}
	_, seconds := time.Now().In(loc).Zone()
	return seconds / 3600, true
}
