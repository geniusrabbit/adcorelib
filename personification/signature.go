package personification

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/gogeo"
	"github.com/google/uuid"
	"github.com/sspserver/udetect"
	"github.com/valyala/fasthttp"

	"geniusrabbit.dev/adcorelib/gtracing"
	fasthttpext "geniusrabbit.dev/adcorelib/net/fasthttp"
)

// Signeture provides the builder of cookie assigned to the user by HTTP
type Signeture struct {
	UUIDName       string
	SessidName     string
	SessidLifetime time.Duration
	Detector       Client
}

// Whois user information
func (sign *Signeture) Whois(ctx context.Context, req *fasthttp.RequestCtx) (Person, error) {
	var (
		uuidCookie   fasthttp.Cookie
		sessidCookie fasthttp.Cookie
	)

	if span, _ := gtracing.StartSpanFromFastContext(req, "personification.whois"); span != nil {
		defer span.Finish()
	}

	_ = uuidCookie.ParseBytes(
		req.Request.Header.Cookie(sign.UUIDName),
	)

	_ = sessidCookie.ParseBytes(
		req.Request.Header.Cookie(sign.SessidName),
	)

	primaryLanguage, langs := pareseAcceptLanguage(
		string(req.Request.Header.Peek("Accept-Language")),
	)

	uuidObj, _ := uuid.Parse(string(uuidCookie.Value()))
	sessidObj, _ := uuid.Parse(string(sessidCookie.Value()))
	request := &udetect.Request{
		UID:             uuidObj,
		SessID:          sessidObj,
		IP:              fasthttpext.IPAdressByRequestCF(req),
		UA:              string(req.UserAgent()),
		URL:             string(req.Referer()),
		Ref:             string(req.Request.Header.Referer()), // TODO: add additional information
		DNT:             int8(gocast.ToInt(req.Request.Header.Peek("Dnt"))),
		LMT:             int8(gocast.ToInt(req.QueryArgs().Peek("lmt"))),
		Adblock:         int8(gocast.ToInt(req.QueryArgs().Peek("adb"))),
		PrivateBrowsing: int8(gocast.ToInt(req.QueryArgs().Peek("private"))),
		JS:              1,
		Languages:       langs,
		PrimaryLanguage: primaryLanguage,
		FlashVer:        "",
		Width:           0,
		Height:          0,
		Extensions:      nil,
	}

	response, err := sign.Detector.Detect(ctx, request)
	// Init additional information
	if response.Geo == nil || len(response.Geo.IP) == 0 || response.Geo.Country == "" {
		if response.Geo == nil {
			response.Geo = &udetect.Geo{}
		}
		if len(response.Geo.IP) == 0 {
			response.Geo.IP = net.ParseIP(request.IP)
		}
		if response.Geo.Country == "" {
			cc := string(req.Request.Header.Peek("Cf-Ipcountry"))
			country := gogeo.CountryByCode2(cc)
			response.Geo.ID = uint(country.ID)
			response.Geo.Country = country.Code2
		}
	}
	return &person{
		request: request,
		userInfo: UserInfo{
			Device: response.Device,
			Geo:    response.Geo,
		},
	}, err
}

// SignCookie do sign request by traking response
func (sign *Signeture) SignCookie(resp Person, req *fasthttp.RequestCtx) {
	if span, _ := gtracing.StartSpanFromFastContext(req, "personification.sign"); span != nil {
		defer span.Finish()
	}

	if resp == nil {
		return
	}

	if _uuid := resp.UserInfo().UUID(); len(_uuid) > 0 {
		c := &fasthttp.Cookie{}
		c.SetKey(sign.UUIDName)
		c.SetValue(_uuid)
		c.SetHTTPOnly(true)
		c.SetExpire(time.Now().Add(365 * 24 * time.Hour))
		req.Response.Header.SetCookie(c)
	}

	if sessid := resp.UserInfo().SessionID(); len(sessid) > 0 {
		c := &fasthttp.Cookie{}
		c.SetKey(sign.SessidName)
		c.SetValue(sessid)
		c.SetHTTPOnly(true)
		c.SetExpire(time.Now().Add(sign.SessidLifetime))
		req.Response.Header.SetCookie(c)
	}
}

func pareseAcceptLanguage(langs string) (primaryLanguage string, langArr []string) {
	arr := strings.Split(langs, ",")
	for _, lang := range arr {
		lang = strings.TrimSpace(lang)
		if len(lang) < 2 {
			continue
		}
		if primaryLanguage == "" {
			primaryLanguage = lang[:2]
		} else {
			langArr = append(langArr, lang[:2])
		}
	}
	return primaryLanguage, langArr
}
