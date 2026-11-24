package personification

import (
	"context"
	"net"
	"net/netip"
	"strings"
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/demdxx/langlib"
	"github.com/geniusrabbit/udetect"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"

	adgeo "github.com/geniusrabbit/adcorelib/geo"
	"github.com/geniusrabbit/adcorelib/gtracing"
	fasthttpext "github.com/geniusrabbit/adcorelib/net/fasthttp"
)

// Signature provides the builder of cookie assigned to the user by HTTP
type Signature struct {
	UUIDName       string
	SessidName     string
	SessidLifetime time.Duration
	Detector       Client
	// CountryByIP is an optional IP to ISO-3166-1 alpha-2 lookup.
	// Nil leaves country resolution to the detector and Cloudflare.
	CountryByIP CountryByIP
}

// CountryByIP resolves an IP to an ISO-3166-1 alpha-2 code.
// LookupCountry must not allocate. ok is false when the address is unknown.
type CountryByIP interface {
	LookupCountry(addr netip.Addr) (iso2 [2]byte, ok bool)
}

// Whois user information
func (sign *Signature) Whois(ctx context.Context, req *fasthttp.RequestCtx) (Person, error) {
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
		DNT:             gocast.Int8(req.Request.Header.Peek("Dnt")),
		LMT:             gocast.Int8(req.QueryArgs().Peek("lmt")),
		AdBlock:         gocast.Int8(req.QueryArgs().Peek("adb")),
		PrivateBrowsing: gocast.Int8(req.QueryArgs().Peek("private")),
		JS:              1,
		Languages:       langs,
		PrimaryLanguage: primaryLanguage,
		FlashVer:        "",
		Width:           0,
		Height:          0,
		Extensions:      nil,
	}

	response, err := sign.Detector.Detect(ctx, request)
	if response == nil {
		response = &udetect.Response{}
	}
	response.Geo = cloneResponseGeo(response.Geo)
	if len(response.Geo.IP) == 0 {
		response.Geo.IP = net.ParseIP(request.IP)
	}
	ApplyCloudflare(response.Geo, &req.Request.Header)
	return &PersonType{
		Request: request,
		UserInfoValue: UserInfo{
			Device: response.Device,
			Geo:    response.Geo,
		},
	}, err
}

// SignCookie do sign request by traking response
func (sign *Signature) SignCookie(resp Person, req *fasthttp.RequestCtx) {
	if sign == nil {
		return
	}

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

func pareseAcceptLanguage(langs string) (primaryLanguage langlib.Code, langArr []langlib.Code) {
	arr := strings.SplitSeq(langs, ",")
	for lang := range arr {
		code := langlib.CodeFromString(strings.TrimSpace(lang))
		if code.IsUndefined() {
			continue
		}
		if primaryLanguage.IsUndefined() {
			primaryLanguage = code
		} else {
			langArr = append(langArr, code)
		}
	}
	return primaryLanguage, langArr
}

// ApplyCountry sets geo.Country from geo.IP.
// A hit replaces the country. A miss clears it, so a country learned for a
// different address is not kept. A nil resolver does nothing.
func (sign *Signature) ApplyCountry(geo *udetect.Geo) {
	if sign == nil || sign.CountryByIP == nil || geo == nil || len(geo.IP) == 0 {
		return
	}
	addr, ok := netip.AddrFromSlice(geo.IP)
	if !ok {
		geo.Country = adgeo.UndefinedCountryCode2
		return
	}
	code, ok := sign.CountryByIP.LookupCountry(addr)
	if !ok {
		geo.Country = adgeo.UndefinedCountryCode2
		return
	}
	geo.Country = adgeo.CountryCode2{code[0], code[1]}
}

func cloneResponseGeo(g *udetect.Geo) *udetect.Geo {
	if g == nil {
		return &udetect.Geo{}
	}
	geo := *g
	if g.Carrier != nil {
		carrier := *g.Carrier
		geo.Carrier = &carrier
	}
	if len(g.IP) > 0 {
		geo.IP = append(net.IP(nil), g.IP...)
	}
	return &geo
}
