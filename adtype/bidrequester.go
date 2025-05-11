package adtype

import (
	"fmt"
	"time"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/searchtypes"
	"github.com/geniusrabbit/udetect"
	"github.com/valyala/fasthttp"
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
	List() []*types.Format

	// Bitset returns the bitset of format IDs
	Bitset() *searchtypes.NumberBitset[uint]

	// TypeMask returns the format type mask
	TypeMask() types.FormatTypeBitset
}

// BidRequester defines the interface for interacting with a bid request in the ad system.
type BidRequester interface {
	fmt.Stringer

	AuctionID() string // Unique ID of the auction
	ProjectID() uint64 // Project ID (placeholder)

	// HTTP and domain info
	HTTPRequest() *fasthttp.RequestCtx // Raw HTTP context
	ServiceDomain() string             // Host of the request

	// Request state and flags
	IsSecure() bool          // True if request is HTTPS
	IsAdblock() bool         // True if adblock detected
	IsPrivateBrowsing() bool // True if in incognito
	IsRobot() bool           // True if bot detected
	IsProxy() bool           // True if proxy detected
	IsIPv6() bool            // True if IP is IPv6

	// Source filtering
	SourceFilterCheck(id uint64) bool // Check if source ID is allowed

	// Device and environment
	DeviceInfo() *DeviceInfo   // Full device info
	BrowserInfo() *BrowserInfo // Browser info
	OSInfo() *OSInfo           // OS info

	// App, site, geo
	AppID() uint64             // Target app ID
	AppInfo() *AppInfo         // App info
	SiteInfo() *SiteInfo       // Site info or default
	Domain() []string          // Domain list
	DomainName() string        // Main domain or bundle name
	GeoInfo() *GeoInfo         // Geo info
	CarrierInfo() *CarrierInfo // Carrier info

	// User info
	UserInfo() *User      // Full user info
	LanguageID() uint64   // Browser language
	Keywords() []string   // User keywords
	Tags() []string       // Combined tags
	Categories() []uint64 // Categories (currently cached)

	// Formats accessor
	Formats() BidFormater

	// Financial
	MinECPM() billing.Money // Max of all bid floors

	// Impressions
	Impressions() []*Impression                 // List of impressions
	ImpressionUpdate(fn func(*Impression) bool) // Modify all impressions
	ImpressionByID(id string) *Impression       // Exact match

	// Ext map manipulation
	Get(key string) any      // Get ext value
	Set(key string, val any) // Set ext value
	Unset(keys ...string)    // Unset keys

	Time() time.Time       // Timestamp of the request
	Done() <-chan struct{} // Channel that closes when context is done
	Validate() error       // Validate request
	Release()              // Release resources
}
