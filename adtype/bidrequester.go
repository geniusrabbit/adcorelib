package adtype

import (
	"context"
	"fmt"
	"time"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
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

// BidRequester defines the interface for interacting with a bid request in the ad system.
type BidRequester interface {
	fmt.Stringer

	ID() string         // Unique ID of the request
	ExternalID() string // External ID of the request (from OpenRTB)

	ProjectID() uint64              // Project ID (placeholder)
	AuctionID() string              // Unique ID of the auction
	ExternalAuctionID() string      // External ID of the auction (from OpenRTB)
	AuctionType() types.AuctionType // Type of auction (first price, second price, etc)

	WithFormats(types.FormatsAccessor) BidRequester // Set formats accessor

	// HTTP and domain info
	HTTPRequest() *fasthttp.RequestCtx // Raw HTTP context
	ServiceDomain() string             // Host of the request

	// Request state and flags
	IsDebug() bool           // True if debug mode is enabled
	IsSecure() bool          // True if request is HTTPS
	IsAdBlock() bool         // True if adblock detected
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
	TrafficSourceID() uint64   // Traffic source ID
	AppID() uint64             // Target app ID
	AppInfo() *AppInfo         // App info
	SiteInfo() *SiteInfo       // Site info or default
	Domain() []string          // Domain list
	DomainName() string        // Main domain or bundle name
	GeoID() uint64             // Geo ID
	GeoInfo() *GeoInfo         // Geo info
	CarrierInfo() *CarrierInfo // Carrier info

	// User info
	UserInfo() *User      // Full user info
	LanguageID() uint64   // Browser language
	Keywords() []string   // User keywords
	Tags() []string       // Combined tags
	Categories() []uint64 // Categories (currently cached)
	Sex() uint            // Sex of the user
	Age() uint            // Age of the user

	AccessPoint() AccessPoint   // Access point information
	Formats() types.BidFormater // Formats accessor

	// Financial
	MinECPM() billing.Money // Max of all bid floors

	// Size returns the width and height of the area of visibility for the ad.
	Size() (w, h int)

	// TargetID of the specific point
	TargetID() uint64
	TargetIDs() []uint64
	ExtTargetIDs() []string

	// Impressions
	Impressions() []*Impression                 // List of impressions
	ImpressionUpdate(fn func(*Impression) bool) // Modify all impressions
	ImpressionByID(id string) *Impression       // Exact match

	// Ext map manipulation
	Get(key string) any      // Get ext value
	Set(key string, val any) // Set ext value
	Unset(keys ...string)    // Unset keys

	Time() time.Time           // Timestamp of the request
	CurrentGeoTime() time.Time // Current time in the geo of the request
	Validate() error           // Validate request
	Release()                  // Release resources

	Context() context.Context // Request context
	SetContext(ctx context.Context)
	Done() <-chan struct{} // Channel that closes when context is done
}
