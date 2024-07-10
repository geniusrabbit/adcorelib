package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/udetect"
	"github.com/valyala/fasthttp"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	fasthttpext "github.com/geniusrabbit/adcorelib/net/fasthttp"
	"github.com/geniusrabbit/adcorelib/personification"
	"github.com/geniusrabbit/adcorelib/rand"
)

// RequestOptions prepare
type RequestOptions struct {
	Debug   bool
	Request *fasthttp.RequestCtx
	Count   int
	X, Y    int
	W, WMax int
	H, HMax int
	Page    string
	SubID1  string
	SubID2  string
	SubID3  string
	SubID4  string
	SubID5  string
}

// NewRequestOptions prepare
func NewRequestOptions(ctx *fasthttp.RequestCtx) *RequestOptions {
	var (
		queryArgs        = ctx.QueryArgs()
		w, h, minW, minH = getSizeByCtx(ctx)
		debug, _         = strconv.ParseBool(string(queryArgs.Peek("debug")))
	)
	return &RequestOptions{
		Debug:   debug,
		Request: ctx,
		X:       gocast.Int(string(queryArgs.Peek("x"))),
		Y:       gocast.Int(string(queryArgs.Peek("y"))),
		W:       minW,
		WMax:    ifPositiveNumber(w, -1),
		H:       minH,
		HMax:    ifPositiveNumber(h, -1),
		SubID1:  peekOneFromQuery(queryArgs, "subid1", "subid", "s1"),
		SubID2:  peekOneFromQuery(queryArgs, "subid2", "s2"),
		SubID3:  peekOneFromQuery(queryArgs, "subid3", "s3"),
		SubID4:  peekOneFromQuery(queryArgs, "subid4", "s4"),
		SubID5:  peekOneFromQuery(queryArgs, "subid5", "s5"),
		Count:   gocast.Int(peekOneFromQuery(queryArgs, "count")),
	}
}

// NewDirectRequestOptions prepare
func NewDirectRequestOptions(ctx *fasthttp.RequestCtx) *RequestOptions {
	var (
		queryArgs = ctx.QueryArgs()
		debug, _  = strconv.ParseBool(string(queryArgs.Peek("debug")))
	)
	return &RequestOptions{
		Debug:   debug,
		Request: ctx,
		X:       gocast.Int(string(queryArgs.Peek("x"))),
		Y:       gocast.Int(string(queryArgs.Peek("y"))),
		W:       -1,
		H:       -1,
		SubID1:  peekOneFromQuery(queryArgs, "subid1", "subid", "s1"),
		SubID2:  peekOneFromQuery(queryArgs, "subid2", "s2"),
		SubID3:  peekOneFromQuery(queryArgs, "subid3", "s3"),
		SubID4:  peekOneFromQuery(queryArgs, "subid4", "s4"),
		SubID5:  peekOneFromQuery(queryArgs, "subid5", "s5"),
	}
}

// NewRequestFor person
func NewRequestFor(ctx context.Context, target admodels.Target, person personification.Person,
	opt *RequestOptions, formatAccessor types.FormatsAccessor) *adtype.BidRequest {
	var (
		userInfo         = person.UserInfo()
		ageStart, ageEnd = userInfo.Ages()
		referer          = string(opt.Request.Referer())
		requestID        = rand.UUID()
	)
	req := &adtype.BidRequest{
		ID:         requestID,
		Debug:      opt.Debug,
		RequestCtx: opt.Request,
		Secure:     b2i(fasthttpext.IsSecureCF(opt.Request)),
		Device:     userInfo.DeviceInfo(),
		Imps: []adtype.Impression{
			{
				ID:          rand.UUID(), // Impression ID
				ExtTargetID: "",
				Target:      target,
				FormatTypes: directTypeMask(opt.W == -1 && opt.H == -1),
				Count:       minInt(opt.Count, 1),
				X:           opt.X,
				Y:           opt.Y,
				W:           opt.W,
				H:           opt.H,
				WMax:        opt.WMax,
				HMax:        opt.HMax,
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
			PrivacyPolicy: 1,               // Default: 1 ("1": has a privacy policy)
			Keywords:      "",              // Comma separated list of keywords about the site.
			Page:          referer,         // URL of the page
			Ref:           referer,         // Referrer URL
			Search:        "",              // Search string that caused naviation
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
