//
// @project GeniusRabbit corelib 2017 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018
//

package urlgenerator

import (
	"net/url"
	"strings"
	"time"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
	"github.com/geniusrabbit/adcorelib/eventtraking/pixelgenerator"
)

// Generator of URLs
type Generator struct {
	EventGenerator       eventgenerator.Generator
	PixelGenerator       pixelgenerator.PixelGenerator
	Schema               string
	ServiceDomain        string
	CDNDomain            string
	LibDomain            string
	ClickPattern         string
	DirectPattern        string
	WinPattern           string
	BillingNoticePattern string
}

func (g *Generator) Init() *Generator {
	if g.Schema != "" && !strings.HasSuffix(g.Schema, "://") {
		g.Schema = strings.TrimRight(g.Schema, ":/") + "://"
	}
	if !(false ||
		strings.HasPrefix(g.CDNDomain, "http://") ||
		strings.HasPrefix(g.CDNDomain, "https://") ||
		strings.HasPrefix(g.CDNDomain, "//")) {
		g.CDNDomain = "//" + strings.TrimRight(g.CDNDomain, "/")
	}
	if !(false ||
		strings.HasPrefix(g.LibDomain, "http://") ||
		strings.HasPrefix(g.LibDomain, "https://") ||
		strings.HasPrefix(g.LibDomain, "//")) {
		g.LibDomain = "//" + strings.TrimRight(g.LibDomain, "/")
	}
	return g
}

// CDNURL returns full URL to path
func (g *Generator) CDNURL(path string) string {
	if path == "" || isFullURL(path) {
		return path
	}
	if path[0] == '/' {
		return g.CDNDomain + path
	}
	return g.CDNDomain + "/" + path
}

// LibURL returns full URL to lib file path
func (g *Generator) LibURL(path string) string {
	if path == "" {
		return path
	}
	if path[0] == '/' {
		return g.LibDomain + path
	}
	return g.LibDomain + "/" + path
}

// PixelURL generator from response of item
func (g *Generator) PixelURL(event events.Type, status uint8, item adtype.ResponserItem, response adtype.Responser, js bool) (string, error) {
	ev, err := g.EventGenerator.Event(event, status, response, item)
	if err != nil {
		return "", err
	}
	return g.PixelGenerator.Event(ev, js)
}

// Lead URL traking for lead type of event
func (g *Generator) PixelLead(item adtype.ResponserItem, response adtype.Responser, js bool) (string, error) {
	var sourceID uint64
	if item.Source() != nil {
		sourceID = item.Source().ID()
	}
	return g.PixelGenerator.Lead(&events.LeadCode{
		AuctionID:  response.Request().ID,
		ImpAdID:    item.ID(),
		SourceID:   sourceID,
		ProjectID:  response.Request().ProjectID(),
		CampaignID: item.CampaignID(),
		AdID:       item.AdID(),
		Price:      item.Price(admodels.ActionLead).Int64(),
		Timestamp:  time.Now().Unix(),
	})
}

// PixelDirectURL generator from response of item
func (g *Generator) PixelDirectURL(event events.Type, status uint8, item adtype.ResponserItem, response adtype.Responser, direct string) (string, error) {
	ev, err := g.EventGenerator.Event(event, status, response, item)
	if err != nil {
		return "", err
	}
	return g.PixelGenerator.EventDirect(ev, direct)
}

// ClickURL generator from respponse of item
func (g *Generator) ClickURL(item adtype.ResponserItem, response adtype.Responser) (string, error) {
	return g.encodeURL(g.ClickPattern, events.Click, events.StatusSuccess, item, response)
}

// MustClickURL generator from respponse of item
func (g *Generator) MustClickURL(item adtype.ResponserItem, response adtype.Responser) string {
	res, _ := g.ClickURL(item, response)
	return res
}

// ClickRouterURL returns router pattern
func (g *Generator) ClickRouterURL() string {
	urls := strings.Split(g.ClickPattern, "?")
	return urls[0]
}

// DirectURL generator from respponse of item
func (g *Generator) DirectURL(event events.Type, item adtype.ResponserItem, response adtype.Responser) (string, error) {
	if event == events.Undefined {
		event = events.Direct
	}
	return g.encodeURL(g.DirectPattern, event, events.StatusSuccess, item, response)
}

// DirectRouterURL returns router pattern
func (g *Generator) DirectRouterURL() string {
	urls := strings.Split(g.DirectPattern, "?")
	return urls[0]
}

// WinURL generator from response of item
func (g *Generator) WinURL(event events.Type, status uint8, item adtype.ResponserItem, response adtype.Responser) (string, error) {
	if event == events.Undefined {
		event = events.AccessPointWin
	}
	return g.encodeURL(g.WinPattern, event, events.StatusSuccess, item, response)
}

// BillingNoticeURL generator from response of item
func (g *Generator) BillingNoticeURL(event events.Type, status uint8, item adtype.ResponserItem, response adtype.Responser) (string, error) {
	if event == events.Undefined {
		event = events.AccessPointBillingNotice
	}
	return g.encodeURL(g.BillingNoticePattern, event, status, item, response)
}

// WinRouterURL returns router pattern
func (g *Generator) WinRouterURL() string {
	urls := strings.Split(g.WinPattern, "?")
	return urls[0]
}

// EventCode generator
func (g *Generator) EventCode(event events.Type, status uint8, item adtype.ResponserItem, response adtype.Responser) (string, error) {
	ev, err := g.EventGenerator.Event(event, status, response, item)
	if err != nil {
		return "", err
	}
	code := ev.Pack().Compress().URLEncode()
	return code.String(), code.ErrorObj()
}

func (g *Generator) encodeURL(pattern string, event events.Type, status uint8, item adtype.ResponserItem, response adtype.Responser) (string, error) {
	if pattern == "" {
		return "", nil
	}
	var (
		code, err = g.EventCode(event, status, item, response)
		urlVal    string
	)
	if err != nil {
		return "", err
	}

	code = url.QueryEscape(code)
	if !strings.Contains(pattern, "{hostname}") {
		if strings.HasPrefix(pattern, "/") {
			urlVal = g.hostSchema() + g.hostDomain(response) + strings.Replace(pattern, "{code}", code, -1)
		} else {
			urlVal = g.hostSchema() + g.hostDomain(response) + "/" + strings.Replace(pattern, "{code}", code, -1)
		}
	} else {
		urlVal = strings.NewReplacer(
			"{schema}", g.hostSchema(),
			"{code}", code,
			"{hostname}", g.hostDomain(response),
		).Replace(pattern)
	}

	if response.Request().AuctionType.IsSecondPrice() {
		if strings.Contains(urlVal, "?") {
			urlVal += "&"
		} else {
			urlVal += "?"
		}
		// ${AUCTION_PRICE} - Clearing price using the same currency and units as
		// the bid. Note that this macro is currently not supported in AMP ads.
		urlVal += "price=${AUCTION_PRICE}"
	}
	return urlVal, nil
}

func (g *Generator) hostSchema() string {
	if g.Schema == "" {
		return "//"
	}
	return g.Schema
}

func (g *Generator) hostDomain(response adtype.Responser) string {
	if g.ServiceDomain != "" {
		return g.ServiceDomain
	}
	return string(response.Request().HTTPRequest().Host())
}

func isFullURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "//")
}
