//
// @project GeniusRabbit corelib 2016
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016
//

package adtype

import "errors"

// Set of errors
var (
	// For bidding validation
	ErrInvalidCur                = errors.New("BID currency is not valid")
	ErrInvalidCreativeSize       = errors.New("creative size is invalid")
	ErrInvalidViewType           = errors.New("view type is invalid")
	ErrLowPrice                  = errors.New("BID price is lower than floor price")
	ErrResponseEmpty             = errors.New("response is empty")
	ErrResponseSkipped           = errors.New("response is skipped")
	ErrResponseItemEmpty         = errors.New("response item is empty")
	ErrResponseItemSkipped       = errors.New("response item is skipped")
	ErrResponseInvalidType       = errors.New("invalid response type")
	ErrResponseInvalidGroup      = errors.New("system not support group winners")
	ErrInvalidItemInitialisation = errors.New("invalid item initialisation")
)

// NoSupportError object
type NoSupportError struct {
	NSField string
}

// Error text
func (e NoSupportError) Error() string {
	return e.NSField + " is not supported"
}
