package endpoint

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/geniusrabbit/udetect"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adquery/bidrequest"
	"github.com/geniusrabbit/adcorelib/adtype"
	fasthttpext "github.com/geniusrabbit/adcorelib/net/fasthttp"
	"github.com/geniusrabbit/adcorelib/personification"
)

// NewRequestFor specific person
func NewRequestFor(
	ctx context.Context,
	app *admodels.Application,
	target adtype.Target,
	person personification.Person,
	opt *RequestOptions,
	formatAccessor types.FormatsAccessor,
) adtype.BidRequester {
	var (
		userInfo         = person.UserInfo()
		ageStart, ageEnd = userInfo.Ages()
		referer          = string(opt.Request.Referer())
		pageURL          = referer
		xPagePath        = string(opt.Request.Request.Header.Peek("X-Page-Path"))
		xPageDomain      = string(opt.Request.Request.Header.Peek("X-Page-Domain"))
		xPageURL         = string(opt.Request.Request.Header.Peek("X-Page-URL"))
		refDomainName    = domain(referer)
		refPath          = urlPath(referer)
		refScheme        = domainScheme(referer)
		stateFlags       bidrequest.BidRequestFlags
	)

	// If the referer is empty or root, try to construct it from X-Page-* headers
	if refPath == "" || refPath == "/" {
		if len(xPagePath) > 0 && (xPageDomain == "" || xPageDomain == refDomainName) {
			pageURL = refScheme + "://" + refDomainName + xPagePath
		} else if xPageURL != "" && strings.HasPrefix(xPageURL, refScheme+"://"+refDomainName) {
			pageURL = xPageURL
		}
	}

	if fasthttpext.IsSecureCF(opt.Request) {
		stateFlags |= bidrequest.BidRequestFlagSecure
	}
	if brwsr := userInfo.DeviceInfo().Browser; brwsr != nil {
		if brwsr.IsRobot == 1 {
			stateFlags |= bidrequest.BidRequestFlagBot
		}
		if brwsr.PrivateBrowsing == 1 {
			stateFlags |= bidrequest.BidRequestFlagPrivateBrowsing
		}
		if brwsr.AdBlock == 1 {
			stateFlags |= bidrequest.BidRequestFlagAdBlock
		}
	}

	req := &bidrequest.BidRequest{
		IDVal:      adtype.NewRequestID(),
		Debug:      opt.Debug,
		RequestCtx: opt.Request,
		StateFlags: stateFlags,
		Device:     userInfo.DeviceInfo(),
		AppTarget:  app,
		Imps: []*adtype.Impression{
			{
				ID:           adtype.NewImpressionID(), // Impression ID
				Target:       target,
				FormatTypes:  opt.GetFormatTypes(),
				FormatCodes:  opt.FormatCodes,
				Count:        max(opt.Count, 1),
				Interstitial: opt.Interstitial,
				Push:         opt.Push,
				X:            opt.X,
				Y:            opt.Y,
				Width:        opt.Width,
				Height:       opt.Height,
				WidthMax:     opt.WidthMax,
				HeightMax:    opt.HeightMax,
				SubID1:       opt.SubID1,
				SubID2:       opt.SubID2,
				SubID3:       opt.SubID3,
				SubID4:       opt.SubID4,
				SubID5:       opt.SubID5,
			},
		},
		User: &adtype.User{
			ID:            userInfo.UUID(),                     // Unique User ID
			SessionID:     userInfo.SessionID(),                // Unique session ID
			FingerPrintID: userInfo.Fingerprint(),              //
			ETag:          userInfo.ETag(),                     //
			AgeStart:      ageStart,                            // Year of birth from
			AgeEnd:        ageEnd,                              // Year of birth from
			Gender:        sexFrom(userInfo.MostPossibleSex()), // Gender ("M": male, "F" female, "O" Other)
			Keywords:      userInfo.Keywords(),                 // Comma separated list of keywords, interests, or intent
			Geo:           userInfo.GeoInfo(),
		},
		Site: &udetect.Site{
			ExtID:         "",            // External ID
			Domain:        refDomainName, //
			Cat:           nil,           // Array of categories
			PrivacyPolicy: 0,             // Default: 1 ("1": has a privacy policy)
			Keywords:      opt.Keywords,  // Comma separated list of keywords about the site.
			Page:          pageURL,       // URL of the page
			Referrer:      referer,       // Referrer URL
			Search:        "",            // Search string that caused navigation
			Mobile:        0,             // Mobile ("1": site is mobile optimised)
		},
		Person:   person,
		Ctx:      ctx,
		Timemark: time.Now(),
	}

	// Debug overrides for testing purposes
	if req.IsDebug() {
		ipStr := opt.Request.QueryArgs().Peek("ip")
		if len(ipStr) > 0 {
			req.User.Geo.IP = net.ParseIP(string(ipStr))
		}
		ccStr := opt.Request.QueryArgs().Peek("cc")
		if len(ccStr) > 0 {
			req.User.Geo.Country = string(ccStr)
		}
		secureStr := opt.Request.QueryArgs().Peek("secure")
		if len(secureStr) > 0 {
			if secure, _ := strconv.ParseBool(string(secureStr)); secure {
				req.StateFlags |= bidrequest.BidRequestFlagSecure
			} else {
				req.StateFlags &^= bidrequest.BidRequestFlagSecure
			}
		}
		domainStr := opt.Request.QueryArgs().Peek("domain")
		if len(domainStr) > 0 {
			req.Site.Domain = string(domainStr)
		}
		// Replace hostname for url and referrer if domain override is provided
		if len(domainStr) > 0 {
			req.Site.Page = replaceDomain(req.Site.Page, string(domainStr), req.IsSecure())
			req.Site.Referrer = replaceDomain(req.Site.Referrer, string(domainStr), req.IsSecure())
		}
	}

	// Prepare bid request with categories and tags
	_ = req.PrepareRequest(0, nil)

	return req.WithFormats(formatAccessor)
}

func replaceDomain(urlStr, newDomain string, secure bool) string {
	if urlStr == "" {
		return urlStr
	}
	parts := strings.SplitN(urlStr, "://", 2)
	if len(parts) != 2 {
		return urlStr
	}
	scheme, rest := parts[0], parts[1]
	pathIndex := strings.Index(rest, "/")
	var path string
	if pathIndex != -1 {
		path = rest[pathIndex:]
	} else {
		path = ""
	}
	if secure {
		scheme = "https"
	} else {
		scheme = "http"
	}
	return scheme + "://" + newDomain + path
}
