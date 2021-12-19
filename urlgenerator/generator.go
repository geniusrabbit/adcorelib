//
// @project GeniusRabbit rotator 2017 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018
//

package urlgenerator

import (
	"net/url"
	"strings"
	"time"

	"geniusrabbit.dev/corelib/admodels"
	"geniusrabbit.dev/corelib/adtype"
	"geniusrabbit.dev/corelib/eventtraking/eventgenerator"
	"geniusrabbit.dev/corelib/eventtraking/events"
	"geniusrabbit.dev/corelib/eventtraking/pixelgenerator"
)

// Generator of URLs
type Generator struct {
	EventGenerator eventgenerator.Generator
	PixelGenerator pixelgenerator.PixelGenerator
	CDNDomain      string
	ClickPattern   string
	DirectPattern  string
	WinPattern     string
}

// CDNURL returns full URL to path
func (g *Generator) CDNURL(path string) string {
	if path == "" {
		return ""
	}
	if path[0] == '/' {
		return "//" + g.CDNDomain + path
	}
	return "//" + g.CDNDomain + "/" + path
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
	var (
		rctx      = response.Request().HTTPRequest()
		code, err = g.EventCode(event, status, item, response)
	)
	if err != nil {
		return "", err
	}
	code = url.QueryEscape(code)
	if !strings.Contains(pattern, "{hostname}") {
		if strings.HasPrefix(pattern, "/") {
			return "//" + string(rctx.Host()) + strings.Replace(pattern, "{code}", code, -1), nil
		}
		return "//" + string(rctx.Host()) + "/" + strings.Replace(pattern, "{code}", code, -1), nil
	}

	return strings.NewReplacer(
		"{code}", code,
		"{hostname}", string(rctx.Host()),
	).Replace(pattern), nil
}
