//
// @project GeniusRabbit corelib 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package adtype

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/geniusrabbit/udetect"
	"github.com/valyala/fasthttp"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/i18n/languages"
	"github.com/geniusrabbit/adcorelib/personification"
	"github.com/geniusrabbit/adcorelib/searchtypes"
)

// defaultUserdata initializes a default User with Geo set to GeoDefault
var defaultUserdata = User{Geo: &udetect.GeoDefault}

// Native asset IDs enumeration for different native ad components
const (
	NativeAssetUndefined = iota // Undefined asset
	NativeAssetTitle            // Title asset
	NativeAssetLegend           // Legend asset
	NativeAssetMainImage        // Main image asset
	NativeAssetIcon             // Icon asset
	NativeAssetRating           // Rating asset
	NativeAssetSponsored        // Sponsored asset
)

// BidRequest represents a bid request in the ad system.
// It contains all necessary information for processing an ad bid.
type BidRequest struct {
	Ctx context.Context `json:"-"` // Context for request handling

	ID          string                 `json:"id,omitempty"`           // Auction ID
	ExtID       string                 `json:"bidid,omitempty"`        // External Auction ID
	AccessPoint AccessPoint            `json:"-"`                      // Access point information
	Debug       bool                   `json:"debug,omitempty"`        // Debug mode flag
	AuctionType types.AuctionType      `json:"auction_type,omitempty"` // Type of auction
	RequestCtx  *fasthttp.RequestCtx   `json:"-"`                      // HTTP request context
	Request     any                    `json:"-"`                      // Original request from RTB or another protocol
	Person      personification.Person `json:"-"`                      // Personification data
	Imps        []Impression           `json:"imps,omitempty"`         // List of impressions

	AppTarget       *admodels.Application `json:"app_target,omitempty"` // Target application
	Device          *udetect.Device       `json:"device,omitempty"`     // Device information
	App             *udetect.App          `json:"app,omitempty"`        // App information
	Site            *udetect.Site         `json:"site,omitempty"`       // Site information
	User            *User                 `json:"user,omitempty"`       // User information
	Secure          int                   `json:"secure,omitempty"`     // Security flag (1 if secure)
	Adblock         int                   `json:"adb,omitempty"`        // Adblock flag (1 if adblock enabled)
	PrivateBrowsing int                   `json:"pb,omitempty"`         // Private browsing flag (1 if enabled)
	Ext             map[string]any        `json:"ext,omitempty"`        // Additional extensions
	Timemark        time.Time             `json:"timemark,omitempty"`   // Timestamp of the request
	Tracer          any                   `json:"-"`                    // Tracing information

	// Internal caches for efficient access
	targetIDs         []uint64                       // Cached target IDs
	externalTargetIDs []string                       // Cached external target IDs
	categoryArray     []uint64                       // Cached category IDs
	domain            []string                       // Cached domains
	tags              []string                       // Cached tags
	formats           []*types.Format                // Cached formats
	formatBitset      searchtypes.NumberBitset[uint] // Bitset for format IDs
	formatTypeMask    types.FormatTypeBitset         // Bitmask for format types
	sourceIDs         []uint64                       // Cached source IDs
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
	if r.formats != nil {
		r.formats = r.formats[:0]
	}
	r.formatBitset.Reset()

	r.ImpressionUpdate(func(imp *Impression) bool {
		imp.Init(formats)
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
		r.sourceIDs = append(r.sourceIDs, ids...)
	}
}

// SourceFilterCheck checks if a given source ID is allowed by the current filter.
// Returns true if no filter is set or if the ID is present in the filter.
func (r *BidRequest) SourceFilterCheck(id uint64) bool {
	if len(r.sourceIDs) < 1 {
		return true
	}
	for _, sid := range r.sourceIDs {
		if sid == id {
			return true
		}
	}
	return false
}

// Formats returns the list of formats associated with the BidRequest.
// If the formats slice is empty, it aggregates formats from all impressions.
func (r *BidRequest) Formats() []*types.Format {
	if len(r.formats) < 1 {
		for _, imp := range r.Imps {
			r.formats = append(r.formats, imp.Formats()...)
		}
	}
	return r.formats
}

// FormatBitset returns a bitset representing the format IDs in the BidRequest.
// It populates the bitset if it's currently empty.
func (r *BidRequest) FormatBitset() *searchtypes.NumberBitset[uint] {
	if r.formatBitset.Len() < 1 {
		for _, f := range r.Formats() {
			r.formatBitset.Set(uint(f.ID))
		}
	}
	return &r.formatBitset
}

// FormatTypeMask returns a bitmask representing the types of formats in the BidRequest.
// It populates the mask if it's currently empty.
func (r *BidRequest) FormatTypeMask() types.FormatTypeBitset {
	if r.formatTypeMask.IsEmpty() {
		r.formatTypeMask.SetFromFormats(r.Formats()...)
	}
	return r.formatTypeMask
}

// Size returns the width and height of the area of visibility for the ad.
func (r *BidRequest) Size() (w, h int) { return r.Width(), r.Height() }

// Width returns the width of the device's browser.
// Returns 0 if device or browser information is unavailable.
func (r *BidRequest) Width() int {
	if r.Device == nil || r.Device.Browser == nil {
		return 0
	}
	return r.Device.Browser.Width
}

// Height returns the height of the device's browser.
// Returns 0 if device or browser information is unavailable.
func (r *BidRequest) Height() int {
	if r.Device == nil || r.Device.Browser == nil {
		return 0
	}
	return r.Device.Browser.Height
}

// Tags returns a list of tags associated with the BidRequest.
// It aggregates keywords from the user and site information.
func (r *BidRequest) Tags() []string {
	if r.tags != nil {
		return r.tags
	}
	if r != nil {
		if r.User != nil && len(r.User.Keywords) > 0 {
			r.tags = strings.Split(r.User.Keywords, ",")
		}
		if r.Site != nil && len(r.Site.Keywords) > 0 {
			r.tags = append(r.tags, strings.Split(r.Site.Keywords, ",")...)
		}
	}
	return r.tags
}

// TargetID returns the target ID if there is exactly one impression with a target.
// Otherwise, returns 0.
func (r *BidRequest) TargetID() uint64 {
	if len(r.Imps) == 1 && r.Imps[0].Target != nil {
		return r.Imps[0].Target.ID()
	}
	return 0
}

// TargetIDs returns a slice of target IDs associated with the BidRequest.
func (r *BidRequest) TargetIDs() []uint64 {
	targets, _ := r.getTargetIDs()
	return targets
}

// ExtTargetIDs returns a slice of external target IDs associated with the BidRequest.
func (r *BidRequest) ExtTargetIDs() []string {
	_, extTargets := r.getTargetIDs()
	return extTargets
}

// getTargetIDs is a helper method that retrieves both target and external target IDs.
// It caches the results to optimize repeated access.
func (r *BidRequest) getTargetIDs() (ids []uint64, externalIDs []string) {
	if r.targetIDs == nil && r.externalTargetIDs == nil && len(r.Imps) > 0 {
		for _, imp := range r.Imps {
			if imp.Target != nil {
				r.targetIDs = append(r.targetIDs, imp.Target.ID())
			}
			if imp.ExternalTargetID != "" {
				r.externalTargetIDs = append(r.externalTargetIDs, imp.ExternalTargetID)
			}
		}
		if r.targetIDs == nil {
			r.targetIDs = []uint64{}
		}
	}
	return r.targetIDs, r.externalTargetIDs
}

// Domain returns a list of domains associated with the site or app.
// It prepares the domain list by aggregating from Site and App information.
func (r *BidRequest) Domain() []string {
	if r.domain == nil {
		if r.Site != nil {
			r.domain = r.Site.DomainPrepared()
		}
		if r.App != nil {
			r.domain = r.App.DomainPrepared()
		}
	}
	return r.domain
}

// DomainName returns the primary domain name of the site or the bundle name of the app.
func (r *BidRequest) DomainName() string {
	if r != nil {
		if r.Site != nil {
			return r.Site.Domain
		}
		if r.App != nil {
			return r.App.Bundle
		}
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

// GeoID returns the geographical ID associated with the user.
// Returns 0 if geographical information is unavailable.
func (r *BidRequest) GeoID() uint64 {
	if r == nil || r.User == nil || r.User.Geo == nil {
		return 0
	}
	return uint64(r.User.Geo.ID)
}

// GeoCode returns the country code associated with the user's geography.
// Returns "**" if geographical information is unavailable.
func (r *BidRequest) GeoCode() string {
	if r == nil || r.User == nil || r.User.Geo == nil {
		return "**"
	}
	return r.User.Geo.Country
}

// City returns the city associated with the user's geography.
// Returns an empty string if geographical information is unavailable.
func (r *BidRequest) City() string {
	if r == nil || r.User == nil || r.User.Geo == nil {
		return ""
	}
	return r.User.Geo.City
}

// LanguageID returns the language ID based on the primary language of the browser.
func (r *BidRequest) LanguageID() uint64 {
	return uint64(languages.GetLanguageIdByCodeString(
		r.BrowserInfo().PrimaryLanguage,
	))
}

// BrowserID returns the ID of the user's browser.
// Returns 0 if device or browser information is unavailable.
func (r *BidRequest) BrowserID() uint64 {
	if r.Device == nil || r.Device.Browser == nil {
		return 0
	}
	return r.Device.Browser.ID
}

// OSID returns the ID of the user's operating system.
// Returns 0 if device or OS information is unavailable.
func (r *BidRequest) OSID() uint64 {
	if r.Device == nil || r.Device.OS == nil {
		return 0
	}
	return uint64(r.Device.OS.ID)
}

// Gender returns the most relevant gender as a byte.
// Returns '?' if gender information is unavailable or invalid.
func (r *BidRequest) Gender() byte {
	if r.User == nil || len(r.User.Gender) != 1 {
		return '?'
	}
	return r.User.Gender[0]
}

// Age returns the most relevant age of the user.
// Returns the starting age if AgeStart <= AgeEnd, otherwise returns AgeStart.
func (r *BidRequest) Age() uint {
	if r.User == nil {
		return 0
	}
	if r.User.AgeStart <= r.User.AgeEnd {
		return uint(r.User.AgeStart)
	}
	return uint(r.User.AgeStart)
}

// Ages returns a range of ages [AgeStart, AgeEnd].
// If AgeStart > AgeEnd, it still returns [AgeStart, AgeEnd].
func (r *BidRequest) Ages() [2]uint {
	if r.User == nil {
		return [2]uint{0, 1000}
	}
	if r.User.AgeStart <= r.User.AgeEnd {
		return [2]uint{
			uint(r.User.AgeStart),
			uint(r.User.AgeEnd),
		}
	}
	return [2]uint{
		uint(r.User.AgeEnd),
		uint(r.User.AgeStart),
	}
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
	// Future implementation for aggregating categories from App and Site
	// if r.categoryArray == nil {
	// 	if r.App != nil {
	// 	}
	// 	if r.Site != nil {
	// 	}
	// }
	return r.categoryArray
}

// IsSecure checks if the request is made over a secure connection.
func (r *BidRequest) IsSecure() bool { return r.Secure == 1 }

// IsAdblock checks if the user has an ad blocker enabled.
func (r *BidRequest) IsAdblock() bool { return r.Adblock == 1 }

// IsPrivateBrowsing checks if the user is in private browsing mode.
func (r *BidRequest) IsPrivateBrowsing() bool { return r.PrivateBrowsing == 1 }

// IsRobot checks if the user is a robot.
func (r *BidRequest) IsRobot() bool { return false }

// IsProxy checks if the user is using a proxy.
func (r *BidRequest) IsProxy() bool { return false }

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
func (r *BidRequest) UserInfo() *User {
	if r == nil {
		return nil
	}
	if r.User == nil {
		r.User = &User{}
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

// DeviceID returns the ID of the device associated with the BidRequest.
// Returns 0 if device information is unavailable.
func (r *BidRequest) DeviceID() uint64 {
	if r != nil && r.Device != nil {
		return uint64(r.Device.ID)
	}
	return 0
}

// DeviceType returns the type of the device as an unsigned integer.
// Returns 0 if device information is unavailable.
func (r *BidRequest) DeviceType() uint64 {
	if r == nil {
		return 0
	}
	return uint64(r.DeviceInfo().DeviceType)
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

// IsIPv6 checks if the user's IP address is IPv6.
func (r *BidRequest) IsIPv6() bool {
	return r != nil && r.User != nil && r.User.Geo != nil && r.User.Geo.IsIPv6()
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

// ImpressionUpdate applies a provided function to each impression in the BidRequest.
// If the function returns true, the impression is updated.
func (r *BidRequest) ImpressionUpdate(fn func(imp *Impression) bool) {
	for i, imp := range r.Imps {
		if fn(&imp) {
			r.Imps[i] = imp
		}
	}
}

// ImpressionByID returns a pointer to the Impression with the specified ID.
// Returns nil if no matching impression is found.
func (r *BidRequest) ImpressionByID(id string) *Impression {
	for _, im := range r.Imps {
		if im.ID == id {
			return &im
		}
	}
	return nil
}

// ImpressionByIDvariation returns a pointer to the first Impression whose ID is a prefix of the provided ID.
// This allows matching impressions even if the ID contains additional postfixes.
// Returns nil if no matching impression is found.
func (r *BidRequest) ImpressionByIDvariation(id string) *Impression {
	for _, im := range r.Imps {
		if strings.HasPrefix(id, im.ID) {
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
