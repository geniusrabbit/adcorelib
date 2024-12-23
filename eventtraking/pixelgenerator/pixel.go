//
// @project geniusrabbit::archivarious 2017 - 2018, 2021, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2018, 2021, 2024
//

package pixelgenerator

import (
	"fmt"
	"net/url"

	"github.com/geniusrabbit/adcorelib/eventtraking/events"
)

type EventType interface {
	Pack() events.Code
}

// PixelGenerator object
type PixelGenerator[EventT EventType, LeadT fmt.Stringer] struct {
	hostname string
}

// NewPixelGenerator object
func NewPixelGenerator[EventT EventType, LeadT fmt.Stringer](hostname string) PixelGenerator[EventT, LeadT] {
	return PixelGenerator[EventT, LeadT]{
		hostname: hostname,
	}
}

// Event generates pixel URL with event registration
func (g PixelGenerator[EventT, LeadT]) Event(ev EventT, js bool) (a string, err error) {
	var (
		code = ev.Pack().Compress().URLEncode()
		u    = url.Values{"i": []string{code.String()}}
	)
	if err = code.ErrorObj(); err != nil {
		return a, err
	}
	if js {
		a = fmt.Sprintf("//%s/t/px.js?%s", g.hostname, u.Encode())
	} else {
		a = fmt.Sprintf("//%s/t/px.gif?%s", g.hostname, u.Encode())
	}
	return a, err
}

// EventDirect can be used in case of traking `direct` or `no-traking` ad type.
// Pixel must automaticaly redirect to `u` param after pixel will be registered
func (g PixelGenerator[EventT, LeadT]) EventDirect(ev EventT, direct string) (a string, err error) {
	var (
		code = ev.Pack().Compress().URLEncode()
		u    = url.Values{
			"i": []string{code.String()},
			"u": []string{direct},
		}
	)
	if err = code.ErrorObj(); err != nil {
		return a, err
	}
	return fmt.Sprintf("//%s/go/m?%s", g.hostname, u.Encode()), nil
}

// Lead URL traking for lead type of event
func (g PixelGenerator[EventT, LeadT]) Lead(lead LeadT) (string, error) {
	return fmt.Sprintf("//%s/lead?l=%s", g.hostname, url.QueryEscape(lead.String())), nil
}
