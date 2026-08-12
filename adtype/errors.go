//
// @project GeniusRabbit corelib 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package adtype

import "github.com/geniusrabbit/adcorelib/errtype"

// Set of errors
var (
	// For bidding validation
	ErrInvalidCur                = errtype.Error("BID currency is not valid")
	ErrInvalidCreativeSize       = errtype.Error("creative size is invalid")
	ErrInvalidViewType           = errtype.Error("view type is invalid")
	ErrLowPrice                  = errtype.Error("BID price is lower than floor price")
	ErrResponseEmpty             = errtype.Error("response is empty")
	ErrResponseSkipped           = errtype.Error("response is skipped")
	ErrResponseNoBid             = errtype.Error("response no bid")
	ErrResponseInvalidRequest    = errtype.Error("response invalid request")
	ErrResponseItemEmpty         = errtype.Error("response item is empty")
	ErrResponseItemSkipped       = errtype.Error("response item is skipped")
	ErrResponseInvalidType       = errtype.Error("invalid response type")
	ErrResponseInvalidGroup      = errtype.Error("system not support group winners")
	ErrInvalidItemInitialisation = errtype.Error("invalid item initialisation")
	ErrSourceEmptyRejected       = errtype.Error("empty source rejected")
)

// NoSupportError object
type NoSupportError struct {
	NSField string
}

// Error text
func (e NoSupportError) Error() string {
	return e.NSField + " is not supported"
}
