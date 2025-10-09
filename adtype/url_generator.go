//
// @project GeniusRabbit corelib 2017 - 2019, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019, 2024
//

package adtype

import "github.com/geniusrabbit/adcorelib/eventtraking/events"

// URLGenerator of advertisement
type URLGenerator interface {
	// CDNURL returns full URL to path
	CDNURL(path string) string

	// LibURL returns full URL to lib file path
	LibURL(path string) string

	// PixelURL generator from response of item
	// @param js generates the JavaScript pixel type
	PixelURL(event events.Type, status uint8, item ResponseItem, response Response, js bool) (string, error)

	// PixelDirectURL generator from response of item
	PixelDirectURL(event events.Type, status uint8, item ResponseItem, response Response, direct string) (string, error)

	// PixelLead URL
	PixelLead(item ResponseItem, response Response, js bool) (string, error)

	// MustClickURL generator from response of item
	MustClickURL(item ResponseItem, response Response) string

	// ClickURL generator from response of item
	ClickURL(item ResponseItem, response Response) (string, error)

	// ClickRouterURL returns router pattern
	ClickRouterURL() string

	// DirectURL generator from response of item
	DirectURL(event events.Type, item ResponseItem, response Response) (string, error)

	// DirectRouterURL returns router pattern
	DirectRouterURL() string

	// WinURL generator from response of item
	WinURL(event events.Type, status uint8, item ResponseItem, response Response) (string, error)

	// BillingNoticeURL generator from response of item
	BillingNoticeURL(event events.Type, status uint8, item ResponseItem, response Response) (string, error)

	// WinRouterURL returns router pattern
	WinRouterURL() string
}
