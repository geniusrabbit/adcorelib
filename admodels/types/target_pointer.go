//
// @project GeniusRabbit corelib 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package types

import (
	"time"

	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/searchtypes"
)

// TargetPointer describer of base target params
type TargetPointer interface {
	// Format list to target
	Formats() []*Format

	// FormatBitset of IDs
	FormatBitset() *searchtypes.NumberBitset[uint]

	// FormatTypeMask of formats
	FormatTypeMask() FormatTypeBitset

	// Size of the area of visibility
	Size() (width, height int)

	// Tags list
	Tags() []string

	// TargetID of the specific point
	TargetID() uint64

	// Domain prepared to targeting
	Domain() []string

	// Sex of the user
	Sex() uint

	// Age of the user
	Age() uint

	// Categories of the current request
	Categories() []uint64

	// GeoID of the target GEO
	GeoID() uint64

	// City of the target GEO
	City() string

	// LanguageID of targeting
	LanguageID() uint64

	// DeviceType value
	DeviceType() uint64

	// DeviceID of the target
	DeviceID() uint64

	// OSID of the user
	OSID() uint64

	// BrowserID of the user
	BrowserID() uint64

	// MinECPM value
	MinECPM() billing.Money

	// IsSecure connection of the request
	IsSecure() bool

	// IsAdblock used of the connection of the request
	IsAdblock() bool

	// IsPrivateBrowsing used of the connection of the request
	IsPrivateBrowsing() bool

	// IsIPv6 of the request
	IsIPv6() bool

	// Time of the request start
	Time() time.Time
}

// MultiTargetPointer extends standart target pointer untile multi zone targetting
type MultiTargetPointer interface {
	TargetPointer
	TargetIDs() []uint64
}
