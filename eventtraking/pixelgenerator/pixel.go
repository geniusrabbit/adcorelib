//
// @project geniusrabbit::archivarious 2017 - 2018, 2021, 2024 - 2026
//

package pixelgenerator

import (
	"fmt"
	"strings"

	"github.com/geniusrabbit/adcorelib/eventtraking/events"
)

const queryHex = "0123456789ABCDEF"

type EventType interface {
	Pack() events.Code
}

// PixelGenerator object
type PixelGenerator[EventT EventType, LeadT fmt.Stringer] struct {
	schemaHost string
}

// NewPixelGenerator object
func NewPixelGenerator[EventT EventType, LeadT fmt.Stringer](schemaHost string) PixelGenerator[EventT, LeadT] {
	if !strings.HasPrefix(schemaHost, "http://") && !strings.HasPrefix(schemaHost, "https://") {
		schemaHost = "https://" + strings.TrimLeft(schemaHost, "/")
	}
	return PixelGenerator[EventT, LeadT]{
		schemaHost: schemaHost,
	}
}

// Event generates pixel URL with event registration
func (g PixelGenerator[EventT, LeadT]) Event(ev EventT, js bool) (string, error) {
	code := ev.Pack().Compress().URLEncode()
	if err := code.ErrorObj(); err != nil {
		return "", err
	}
	path := "/t/px.gif?i="
	if js {
		path = "/t/px.js?i="
	}
	return g.writePixel(path, code.Data(), false, ""), nil
}

// EventDirect can be used in case of traking `direct` or `no-traking` ad type.
// Pixel must automaticaly redirect to `u` param after pixel will be registered
func (g PixelGenerator[EventT, LeadT]) EventDirect(ev EventT, direct string) (string, error) {
	code := ev.Pack().Compress().URLEncode()
	if err := code.ErrorObj(); err != nil {
		return "", err
	}
	return g.writePixel("/go/m?i=", code.Data(), true, direct), nil
}

// Lead URL traking for lead type of event
func (g PixelGenerator[EventT, LeadT]) Lead(lead LeadT) (string, error) {
	value := lead.String()
	var b strings.Builder
	b.Grow(len(g.schemaHost) + len("/lead?l=") + queryEscapedLen(value))
	b.WriteString(g.schemaHost)
	b.WriteString("/lead?l=")
	appendQueryEscape(&b, value)
	return b.String(), nil
}

func (g PixelGenerator[EventT, LeadT]) writePixel(path string, code []byte, withDirect bool, direct string) string {
	n := len(g.schemaHost) + len(path) + base64QueryLen(code)
	if withDirect {
		n += len("&u=") + queryEscapedLen(direct)
	}
	var b strings.Builder
	b.Grow(n)
	b.WriteString(g.schemaHost)
	b.WriteString(path)
	appendBase64Query(&b, code)
	if withDirect {
		b.WriteString("&u=")
		appendQueryEscape(&b, direct)
	}
	return b.String()
}

// appendBase64Query writes a base64.URLEncoding token. That alphabet is
// query-safe except '=', which QueryEscape turns into %3D.
func appendBase64Query(b *strings.Builder, code []byte) {
	for _, c := range code {
		if c == '=' {
			b.WriteString("%3D")
			continue
		}
		b.WriteByte(c)
	}
}

func base64QueryLen(code []byte) int {
	n := len(code)
	for _, c := range code {
		if c == '=' {
			n += 2
		}
	}
	return n
}

// appendQueryEscape matches net/url.QueryEscape: space is '+', and every
// byte outside [A-Za-z0-9-_.~] is %HH.
func appendQueryEscape(b *strings.Builder, s string) {
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c == ' ':
			b.WriteByte('+')
		case queryShouldEscape(c):
			b.WriteByte('%')
			b.WriteByte(queryHex[c>>4])
			b.WriteByte(queryHex[c&15])
		default:
			b.WriteByte(c)
		}
	}
}

func queryEscapedLen(s string) int {
	n := len(s)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c != ' ' && queryShouldEscape(c) {
			n += 2
		}
	}
	return n
}

func queryShouldEscape(c byte) bool {
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}
	switch c {
	case '-', '_', '.', '~':
		return false
	default:
		return true
	}
}
