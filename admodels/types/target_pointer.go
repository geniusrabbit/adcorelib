//
// @project GeniusRabbit corelib 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package types

import (
	"time"

	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/searchtypes"
	"github.com/geniusrabbit/udetect"
)

type (
	DeviceInfo  = udetect.Device
	BrowserInfo = udetect.Browser
	OSInfo      = udetect.OS
	AppInfo     = udetect.App
	SiteInfo    = udetect.Site
	GeoInfo     = udetect.Geo
	CarrierInfo = udetect.Carrier
)

// BidFormater defines the interface for managing ad formats in a bid request.
type BidFormater interface {
	// List returns the list of formats
	List() []*Format

	// Bitset returns the bitset of format IDs
	Bitset() *searchtypes.NumberBitset[uint]

	// TypeMask returns the format type mask
	TypeMask() FormatTypeBitset
}

// TargetPointer describer of base target params
type TargetPointer interface {
	// BidFormater of the request
	Formats() BidFormater

	// Size of the area of visibility
	Size() (width, height int)

	// Request state and flags
	IsDebug() bool           // True if debug mode is enabled
	IsSecure() bool          // True if request is HTTPS
	IsAdBlock() bool         // True if adblock detected
	IsPrivateBrowsing() bool // True if in incognito
	IsRobot() bool           // True if bot detected
	IsProxy() bool           // True if proxy detected
	IsIPv6() bool            // True if IP is IPv6

	// Device and environment
	DeviceInfo() *DeviceInfo   // Full device info
	BrowserInfo() *BrowserInfo // Browser info
	OSInfo() *OSInfo           // OS info

	// App, site, geo
	TrafficSourceID() uint64   // Traffic source ID
	AppID() uint64             // Target app ID
	AppInfo() *AppInfo         // App info
	SiteInfo() *SiteInfo       // Site info or default
	Domain() []string          // Domain list
	DomainName() string        // Main domain or bundle name
	GeoID() uint64             // Geo ID
	GeoInfo() *GeoInfo         // Geo info
	CarrierInfo() *CarrierInfo // Carrier info
	LanguageID() uint64        // Browser language

	// TargetID of the specific point
	TargetID() uint64

	Sex() uint // Sex of the user
	Age() uint // Age in years

	// Tags list
	Tags() []string

	// Categories of the current request
	Categories() []uint64

	// MinECPM value
	MinECPM() billing.Money

	// Time of the request start
	Time() time.Time
	CurrentGeoTime() time.Time
}

// MultiTargetPointer extends standart target pointer untile multi zone targetting
type MultiTargetPointer interface {
	TargetPointer
	TargetIDs() []uint64
}
