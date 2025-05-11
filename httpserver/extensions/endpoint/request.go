package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/geniusrabbit/udetect"
	"github.com/valyala/fasthttp"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	fasthttpext "github.com/geniusrabbit/adcorelib/net/fasthttp"
	"github.com/geniusrabbit/adcorelib/personification"
	"github.com/geniusrabbit/adcorelib/rand"
)

// NewRequestFor specific person
func NewRequestFor(ctx context.Context,
	app *admodels.Application,
	target adtype.Target,
	person personification.Person,
	opt *RequestOptions, formatAccessor types.FormatsAccessor) *adtype.BidRequest {
	var (
		userInfo         = person.UserInfo()
		ageStart, ageEnd = userInfo.Ages()
		referer          = string(opt.Request.Referer())
		requestID        = rand.UUID()
		stateFlags       adtype.BidRequestFlags
	)
	if fasthttpext.IsSecureCF(opt.Request) {
		stateFlags |= adtype.BidRequestFlagSecure
	}
	if brwsr := userInfo.DeviceInfo().Browser; brwsr != nil {
		if brwsr.IsRobot == 1 {
			stateFlags |= adtype.BidRequestFlagBot
		}
		if brwsr.PrivateBrowsing == 1 {
			stateFlags |= adtype.BidRequestFlagPrivateBrowsing
		}
		if brwsr.Adblock == 1 {
			stateFlags |= adtype.BidRequestFlagAdblock
		}
	}

	req := &adtype.BidRequest{
		ID:         requestID,
		Debug:      opt.Debug,
		RequestCtx: opt.Request,
		StateFlags: stateFlags,
		Device:     userInfo.DeviceInfo(),
		AppTarget:  app,
		Imps: []adtype.Impression{
			{
				ID:          rand.UUID(), // Impression ID
				Target:      target,
				FormatTypes: opt.GetFormatTypes(),
				FormatCodes: opt.FormatCodes,
				Count:       max(opt.Count, 1),
				X:           opt.X,
				Y:           opt.Y,
				Width:       opt.Width,
				Height:      opt.Height,
				WidthMax:    opt.WidthMax,
				HeightMax:   opt.HeightMax,
				SubID1:      opt.SubID1,
				SubID2:      opt.SubID2,
				SubID3:      opt.SubID3,
				SubID4:      opt.SubID4,
				SubID5:      opt.SubID5,
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
			ExtID:         "",              // External ID
			Domain:        domain(referer), //
			Cat:           nil,             // Array of categories
			PrivacyPolicy: 0,               // Default: 1 ("1": has a privacy policy)
			Keywords:      opt.Keywords,    // Comma separated list of keywords about the site.
			Page:          referer,         // URL of the page
			Referrer:      referer,         // Referrer URL
			Search:        "",              // Search string that caused navigation
			Mobile:        0,               // Mobile ("1": site is mobile optimised)
		},
		Person:   person,
		Ctx:      ctx,
		Timemark: time.Now(),
	}
	req.Init(formatAccessor)
	return req
}

// NewRequestByContext from request
func NewRequestByContext(ctx context.Context, ectx *fasthttp.RequestCtx) (*adtype.BidRequest, error) {
	request := &adtype.BidRequest{RequestCtx: ectx, Timemark: time.Now(), Ctx: ctx}
	if err := json.NewDecoder(bytes.NewBuffer(ectx.Request.Body())).Decode(request); err != nil {
		return nil, err
	}
	return request, nil
}
