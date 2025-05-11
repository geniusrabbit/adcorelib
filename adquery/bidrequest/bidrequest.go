//
// @project GeniusRabbit corelib 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package bidrequest

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/demdxx/xtypes"
	"github.com/geniusrabbit/udetect"
	"github.com/valyala/fasthttp"
	"golang.org/x/exp/slices"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/i18n/languages"
	"github.com/geniusrabbit/adcorelib/personification"
)

// defaultUserdata initializes a default User with Geo set to GeoDefault
var defaultUserdata = adtype.User{Geo: &udetect.GeoDefault, AgeStart: 0, AgeEnd: 1000}

// BidRequestFlags defines flags for bid requests.
type BidRequestFlags uint8

const (
	// BidRequestFlagAdblock indicates if adblock is enabled
	BidRequestFlagAdblock BidRequestFlags = 1 << iota
	// BidRequestFlagPrivateBrowsing indicates if private browsing is enabled
	BidRequestFlagPrivateBrowsing
	// BidRequestFlagSecure indicates if the request is secure
	BidRequestFlagSecure
	// BidRequestFlagBot indicates if the request is from a bot
	BidRequestFlagBot
	// BidRequestFlagProxy indicates if the request is from a proxy
	BidRequestFlagProxy
)

// BidRequest represents a bid request in the ad system.
// It contains all necessary information for processing an ad bid.
type BidRequest struct {
	ID       string    `json:"id,omitempty"`       // Auction ID
	ExtID    string    `json:"bidid,omitempty"`    // External Auction ID
	Timemark time.Time `json:"timemark,omitempty"` // Timestamp of the request

	Ctx   context.Context `json:"-"`               // Context for request handling
	Debug bool            `json:"debug,omitempty"` // Debug mode flag

	// Source of the request
	AccessPoint adtype.AccessPoint `json:"-"` // Access point information

	AuctionType types.AuctionType      `json:"auction_type,omitempty"` // Type of auction
	RequestCtx  *fasthttp.RequestCtx   `json:"-"`                      // HTTP request context
	Request     any                    `json:"-"`                      // Original request from RTB or another protocol
	Person      personification.Person `json:"-"`                      // Personification data
	Imps        []adtype.Impression    `json:"imps,omitempty"`         // List of impressions

	AppTarget  *admodels.Application `json:"app_target,omitempty"` // Target application
	Device     *udetect.Device       `json:"device,omitempty"`     // Device information
	App        *udetect.App          `json:"app,omitempty"`        // App information
	Site       *udetect.Site         `json:"site,omitempty"`       // Site information
	User       *adtype.User          `json:"user,omitempty"`       // User information
	StateFlags BidRequestFlags       `json:"flags"`                // State flags for the request
	Ext        map[string]any        `json:"ext,omitempty"`        // Additional extensions
	Tracer     any                   `json:"-"`                    // Tracing information

	// Internal caches for efficient access
	categoryArray []uint64  // Cached category IDs
	domain        []string  // Cached domains
	tags          []string  // Cached tags
	formats       AdFormats // Formats interface for accessing formats
	sourceIDs     []uint64  // Cached source IDs
}

// String implements the fmt.Stringer interface for BidRequest.
// It returns a pretty-printed JSON representation of the BidRequest.
// If marshalling fails, it returns an error message in JSON format.
func (r *BidRequest) String() (res string) {
	if data, err := json.MarshalIndent(r, "", "  "); err != nil {
		res = `{"error":"` + err.Error() + `"}`
	} else {
		res = string(data)
	}
	return
}

// ProjectID returns the Project ID associated with the BidRequest.
// Currently returns 0 as a placeholder.
func (r *BidRequest) ProjectID() uint64 { return 0 }

// Init initializes the BidRequest with basic information.
// It resets the formats slice and format bitset,
// and initializes each impression using the provided formats accessor.
func (r *BidRequest) Init(formats types.FormatsAccessor) {
	// r.formats = r.formats[:0]
	// r.formatTypeMask.Reset()
	// r.formatBitset.Reset()
	r.formats.Reset()

	// Initialize each impression with the provided formats
	r.ImpressionUpdate(func(imp *adtype.Impression) bool {
		imp.InitFormats(formats)
		for _, f := range imp.Formats() {
			r.formats.Add(f)
		}
		return true
	})
}

// HTTPRequest returns the underlying HTTP request context.
func (r *BidRequest) HTTPRequest() *fasthttp.RequestCtx { return r.RequestCtx }

// ServiceDomain returns the domain of the service handling the request.
func (r *BidRequest) ServiceDomain() string { return string(r.RequestCtx.URI().Host()) }

// SetSourceFilter sets the source filter IDs for the BidRequest.
// It replaces any existing source IDs with the provided ones.
func (r *BidRequest) SetSourceFilter(ids ...uint64) {
	if len(r.sourceIDs) > 0 {
		r.sourceIDs = r.sourceIDs[:0]
	}
	if len(ids) > 0 {
		r.sourceIDs = xtypes.SliceUnique(append(r.sourceIDs, ids...))
		slices.Sort(r.sourceIDs)
	}
}

// SourceFilterCheck checks if a given source ID is allowed by the current filter.
// Returns true if no filter is set or if the ID is present in the filter.
func (r *BidRequest) SourceFilterCheck(id uint64) bool {
	if len(r.sourceIDs) < 1 {
		return true
	}
	_, found := slices.BinarySearch(r.sourceIDs, id)
	return found
}

// Formats returns the list of formats associated with the BidRequest.
// If the formats slice is empty, it aggregates formats from all impressions.
func (r *BidRequest) Formats() adtype.BidFormater {
	return &r.formats
}

// Tags returns a list of tags associated with the BidRequest.
// It aggregates keywords from the user and site information.
func (r *BidRequest) Tags() []string {
	if r == nil {
		return nil
	}
	if r.tags != nil {
		return r.tags
	}
	// Extract tags from user and site information
	if r.User != nil && len(r.User.Keywords) > 0 {
		r.tags = strings.Split(r.User.Keywords, ",")
	}
	// Extract tags from site information
	if r.Site != nil && len(r.Site.Keywords) > 0 {
		r.tags = append(r.tags, strings.Split(r.Site.Keywords, ",")...)
	}
	return r.tags
}

// Domain returns a list of domains associated with the site or app.
// It prepares the domain list by aggregating from Site and App information.
func (r *BidRequest) Domain() []string {
	if r.domain == nil {
		if r.Site != nil {
			r.domain = r.Site.DomainPrepared()
		} else if r.App != nil {
			r.domain = r.App.DomainPrepared()
		}
	}
	return r.domain
}

// DomainName returns the primary domain name of the site or the bundle name of the app.
func (r *BidRequest) DomainName() string {
	if r == nil {
		return ""
	}
	if r.Site != nil {
		return r.Site.Domain
	}
	if r.App != nil {
		return r.App.Bundle
	}
	return ""
}

// Sex returns the user's sex as an unsigned integer.
// Returns 0 if user information is unavailable.
func (r *BidRequest) Sex() uint {
	if r == nil || r.User == nil {
		return 0
	}
	return uint(r.User.Sex())
}

// AppID returns the ID of the target application.
// Returns 0 if the app target is unavailable.
func (r *BidRequest) AppID() uint64 {
	if r == nil || r.AppTarget == nil {
		return 0
	}
	return r.AppTarget.ID
}

// LanguageID returns the language ID based on the primary language of the browser.
func (r *BidRequest) LanguageID() uint64 {
	return uint64(languages.GetLanguageIdByCodeString(
		r.BrowserInfo().PrimaryLanguage,
	))
}

// Keywords returns a slice of keywords associated with the user.
// Returns nil if user information is unavailable.
func (r *BidRequest) Keywords() []string {
	if r == nil || r.User == nil {
		return nil
	}
	return strings.Split(r.User.Keywords, ",")
}

// Categories returns a slice of category IDs associated with the BidRequest.
// Currently, it returns the cached categoryArray.
// (Note: The implementation is incomplete and commented out for future development.)
func (r *BidRequest) Categories() []uint64 {
	return r.categoryArray
}

// IsSecure checks if the request is made over a secure connection.
func (r *BidRequest) IsSecure() bool { return r.StateFlags&BidRequestFlagSecure != 0 }

// IsAdblock checks if the user has an ad blocker enabled.
func (r *BidRequest) IsAdblock() bool { return r.StateFlags&BidRequestFlagAdblock != 0 }

// IsPrivateBrowsing checks if the user is in private browsing mode.
func (r *BidRequest) IsPrivateBrowsing() bool { return r.StateFlags&BidRequestFlagPrivateBrowsing != 0 }

// IsRobot checks if the user is a robot.
func (r *BidRequest) IsRobot() bool { return r.StateFlags&BidRequestFlagBot != 0 }

// IsProxy checks if the user is using a proxy.
func (r *BidRequest) IsProxy() bool { return r.StateFlags&BidRequestFlagProxy != 0 }

// IsIPv6 checks if the user's IP address is IPv6.
func (r *BidRequest) IsIPv6() bool {
	return r != nil && r.User != nil && r.User.Geo != nil && r.User.Geo.IsIPv6()
}

// SiteInfo returns the site information associated with the BidRequest.
// If the site is unavailable, it returns the default site information.
// Returns nil if neither site nor app information is available.
func (r *BidRequest) SiteInfo() *udetect.Site {
	if r.Site != nil {
		return r.Site
	}
	if r.App == nil {
		return &udetect.SiteDefault
	}
	return nil
}

// AppInfo returns the application information associated with the BidRequest.
func (r *BidRequest) AppInfo() *udetect.App { return r.App }

// UserInfo returns the user information associated with the BidRequest.
// It initializes default values if user or geographical information is missing.
func (r *BidRequest) UserInfo() *adtype.User {
	if r == nil {
		return nil
	}
	if r.User == nil {
		r.User = &adtype.User{}
		*r.User = defaultUserdata
	}
	if r.User.Geo == nil {
		r.User.Geo = &udetect.Geo{}
		*r.User.Geo = udetect.GeoDefault
	}
	if r.User.Geo.Carrier == nil {
		r.User.Geo.Carrier = &udetect.Carrier{}
		*r.User.Geo.Carrier = udetect.CarrierDefault
	}
	return r.User
}

// DeviceInfo returns the device information associated with the BidRequest.
// It initializes default values if device, browser, or OS information is missing.
func (r *BidRequest) DeviceInfo() *udetect.Device {
	if r == nil {
		return nil
	}
	if r.Device == nil {
		r.Device = &udetect.Device{}
		*r.Device = udetect.DeviceDefault
	}
	if r.Device.Browser == nil {
		r.Device.Browser = &udetect.Browser{}
		*r.Device.Browser = udetect.BrowserDefault
	}
	if r.Device.OS == nil {
		r.Device.OS = &udetect.OS{}
		*r.Device.OS = udetect.OSDefault
	}
	return r.Device
}

// OSInfo returns the operating system information associated with the BidRequest.
func (r *BidRequest) OSInfo() *udetect.OS {
	if r == nil {
		return nil
	}
	return r.DeviceInfo().OS
}

// BrowserInfo returns the browser information associated with the BidRequest.
func (r *BidRequest) BrowserInfo() *udetect.Browser {
	if r == nil {
		return nil
	}
	return r.DeviceInfo().Browser
}

// MinECPM calculates and returns the minimum ECPM (Effective Cost Per Mille) acceptable for the BidRequest.
// It iterates through all impressions and selects the highest bid floor.
func (r *BidRequest) MinECPM() (minBid billing.Money) {
	for _, imp := range r.Imps {
		if minBid == 0 {
			minBid = max(imp.BidFloorCPM, 0)
		} else if imp.BidFloorCPM > 0 && minBid < imp.BidFloorCPM {
			minBid = imp.BidFloorCPM
		}
	}
	return minBid
}

// GeoInfo returns the geographical information associated with the BidRequest.
func (r *BidRequest) GeoInfo() *udetect.Geo {
	if r == nil {
		return nil
	}
	return r.UserInfo().Geo
}

// CarrierInfo returns the carrier information associated with the user's geography.
func (r *BidRequest) CarrierInfo() *udetect.Carrier {
	if geo := r.GeoInfo(); geo != nil {
		return geo.Carrier
	}
	return nil
}

// Get retrieves a value from the BidRequest's extension map by key.
// Returns nil if the key does not exist.
func (r *BidRequest) Get(key string) any {
	if r.Ext == nil {
		return nil
	}
	return r.Ext[key]
}

// Set sets a key-value pair in the BidRequest's extension map.
func (r *BidRequest) Set(key string, val any) {
	if r.Ext == nil {
		r.Ext = map[string]any{}
	}
	r.Ext[key] = val
}

// Unset removes one or more keys from the BidRequest's extension map.
func (r *BidRequest) Unset(keys ...string) {
	if r.Ext == nil {
		return
	}
	for _, key := range keys {
		delete(r.Ext, key)
	}
}

// Impressions returns a slice of impressions associated with the BidRequest.
func (r *BidRequest) Impressions() []adtype.Impression {
	return r.Imps
}

// ImpressionUpdate applies a provided function to each impression in the BidRequest.
// If the function returns true, the impression is updated.
func (r *BidRequest) ImpressionUpdate(fn func(imp *adtype.Impression) bool) {
	for i, imp := range r.Imps {
		if fn(&imp) {
			r.Imps[i] = imp
		}
	}
}

// ImpressionByID returns a pointer to the Impression with the specified ID.
// Returns nil if no matching impression is found.
func (r *BidRequest) ImpressionByID(id string) *adtype.Impression {
	for _, im := range r.Imps {
		if im.ID == id {
			return &im
		}
	}
	return nil
}

// Time returns the timestamp of the BidRequest.
func (r *BidRequest) Time() time.Time { return r.Timemark }

// Validate performs validation on the BidRequest.
// Currently, it always returns nil, but can be extended to include validation logic.
func (r *BidRequest) Validate() error { return nil }

// Done returns a channel that is closed when the context of the BidRequest is done.
func (r *BidRequest) Done() <-chan struct{} {
	return r.Ctx.Done()
}

// Release releases any resources associated with the BidRequest.
// If the original request implements the releaser interface, it calls the Release method.
func (r *BidRequest) Release() {
	type releaser interface {
		Release()
	}
	if r != nil && r.Request != nil {
		if r, ok := r.Request.(releaser); ok {
			r.Release()
		}
	}
}
