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

var defaultUserdata = User{Geo: &udetect.GeoDefault}

// Native asset IDs
const (
	NativeAssetUndefined = iota
	NativeAssetTitle
	NativeAssetLegend
	NativeAssetMainImage
	NativeAssetIcon
	NativeAssetRating
	NativeAssetSponsored
)

// BidRequest for internal using
type BidRequest struct {
	Ctx context.Context `json:"-"`

	ID              string                 `json:"id,omitempty"`    // Auction ID
	ExtID           string                 `json:"bidid,omitempty"` // External Auction ID
	AccessPoint     AccessPoint            `json:"-"`
	Debug           bool                   `json:"debug,omitempty"`
	AuctionType     types.AuctionType      `json:"auction_type,omitempty"`
	RequestCtx      *fasthttp.RequestCtx   `json:"-"` // HTTP request context
	Request         any                    `json:"-"` // Contains original request from RTB or another protocol
	Person          personification.Person `json:"-"`
	Imps            []Impression           `json:"imps,omitempty"`
	AppTarget       *admodels.Application  `json:"app_target,omitempty"`
	Device          *udetect.Device        `json:"device,omitempty"`
	App             *udetect.App           `json:"app,omitempty"`
	Site            *udetect.Site          `json:"site,omitempty"`
	User            *User                  `json:"user,omitempty"`
	Secure          int                    `json:"secure,omitempty"`
	Adblock         int                    `json:"adb,omitempty"`
	PrivateBrowsing int                    `json:"pb,omitempty"`
	Ext             map[string]any         `json:"ext,omitempty"`
	Timemark        time.Time              `json:"timemark,omitempty"`
	Tracer          any                    `json:"-"`

	targetIDs         []uint64
	externalTargetIDs []string
	categoryArray     []uint64
	domain            []string
	tags              []string
	formats           []*types.Format
	formatBitset      searchtypes.NumberBitset[uint]
	formatTypeMask    types.FormatTypeBitset
	sourceIDs         []uint64
}

// String implements of fmt.Stringer interface
func (r *BidRequest) String() (res string) {
	if data, err := json.MarshalIndent(r, "", "  "); err != nil {
		res = `{"error":"` + err.Error() + `"}`
	} else {
		res = string(data)
	}
	return
}

// ProjectID value
func (r *BidRequest) ProjectID() uint64 {
	return 0
}

// Init basic information
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

// HTTPRequest object
func (r *BidRequest) HTTPRequest() *fasthttp.RequestCtx {
	return r.RequestCtx
}

// ServiceDomain name
func (r *BidRequest) ServiceDomain() string {
	return string(r.RequestCtx.URI().Host())
}

// SetSourceFilter by IDs
func (r *BidRequest) SetSourceFilter(ids ...uint64) {
	if len(r.sourceIDs) > 0 {
		r.sourceIDs = r.sourceIDs[:0]
	}
	if len(ids) > 0 {
		r.sourceIDs = append(r.sourceIDs, ids...)
	}
}

// SourceFilterCheck returns the list of available sources
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

// Formats list
func (r *BidRequest) Formats() []*types.Format {
	if len(r.formats) < 1 {
		for _, imp := range r.Imps {
			r.formats = append(r.formats, imp.Formats()...)
		}
	}
	return r.formats
}

// FormatBitset of IDs
func (r *BidRequest) FormatBitset() *searchtypes.NumberBitset[uint] {
	if r.formatBitset.Len() < 1 {
		for _, f := range r.Formats() {
			r.formatBitset.Set(uint(f.ID))
		}
	}
	return &r.formatBitset
}

// FormatTypeMask of formats
func (r *BidRequest) FormatTypeMask() types.FormatTypeBitset {
	if r.formatTypeMask.IsEmpty() {
		r.formatTypeMask.SetFromFormats(r.Formats()...)
	}
	return r.formatTypeMask
}

// Size of the area of visibility
func (r *BidRequest) Size() (width, height int) {
	return r.Width(), r.Height()
}

// Width size
func (r *BidRequest) Width() int {
	if r.Device == nil || r.Device.Browser == nil {
		return 0
	}
	return r.Device.Browser.Width
}

// Height size
func (r *BidRequest) Height() int {
	if r.Device == nil || r.Device.Browser == nil {
		return 0
	}
	return r.Device.Browser.Height
}

// Tags list
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

// TargetID value
func (r *BidRequest) TargetID() uint64 {
	if len(r.Imps) == 1 && r.Imps[0].Target != nil {
		return r.Imps[0].Target.ID()
	}
	return 0
}

// TargetIDs by request
func (r *BidRequest) TargetIDs() []uint64 {
	targets, _ := r.getTargetIDs()
	return targets
}

// ExtTargetIDs by request
func (r *BidRequest) ExtTargetIDs() []string {
	_, extTargets := r.getTargetIDs()
	return extTargets
}

func (r *BidRequest) getTargetIDs() (ids []uint64, externalIDs []string) {
	if r.targetIDs == nil && r.externalTargetIDs == nil && len(r.Imps) > 0 {
		for _, imp := range r.Imps {
			if imp.Target != nil {
				r.targetIDs = append(r.targetIDs, imp.Target.ID())
			}
			if imp.ExtTargetID != "" {
				r.externalTargetIDs = append(r.externalTargetIDs, imp.ExtTargetID)
			}
		}
		if r.targetIDs == nil {
			r.targetIDs = []uint64{}
		}
	}
	return r.targetIDs, r.externalTargetIDs
}

// Domain of site or bundle name
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

// DomainName of site or bundle name
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

// Sex by request
func (r *BidRequest) Sex() uint {
	if r == nil || r.User == nil {
		return 0
	}
	return uint(r.User.Sex())
}

// AppID by request
func (r *BidRequest) AppID() uint64 {
	if r == nil || r.AppTarget == nil {
		return 0
	}
	return r.AppTarget.ID
}

// GeoID by request
func (r *BidRequest) GeoID() uint64 {
	if r == nil || r.User == nil || r.User.Geo == nil {
		return 0
	}
	return uint64(r.User.Geo.ID)
}

// GeoCode by request
func (r *BidRequest) GeoCode() string {
	if r == nil || r.User == nil || r.User.Geo == nil {
		return "**"
	}
	return r.User.Geo.Country
}

// City by request
func (r *BidRequest) City() string {
	if r == nil || r.User == nil || r.User.Geo == nil {
		return ""
	}
	return r.User.Geo.City
}

// LanguageID value
func (r *BidRequest) LanguageID() uint64 {
	return uint64(languages.GetLanguageIdByCodeString(
		r.BrowserInfo().PrimaryLanguage,
	))
}

// BrowserID by request
func (r *BidRequest) BrowserID() uint64 {
	if r.Device == nil || r.Device.Browser == nil {
		return 0
	}
	return r.Device.Browser.ID
}

// OSID by request
func (r *BidRequest) OSID() uint64 {
	if r.Device == nil || r.Device.OS == nil {
		return 0
	}
	return uint64(r.Device.OS.ID)
}

// Gender which the most relevant
func (r *BidRequest) Gender() byte {
	if r.User == nil || len(r.User.Gender) != 1 {
		return '?'
	}
	return r.User.Gender[0]
}

// Age which the most relevant
func (r *BidRequest) Age() uint {
	if r.User == nil {
		return 0
	}
	if r.User.AgeStart <= r.User.AgeEnd {
		return uint(r.User.AgeStart)
	}
	return uint(r.User.AgeStart)
}

// Ages which the most relevant
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
		uint(r.User.AgeStart),
		uint(r.User.AgeEnd),
	}
}

// Keywords for request
func (r *BidRequest) Keywords() []string {
	if r == nil || r.User == nil {
		return nil
	}
	return strings.Split(r.User.Keywords, ",")
}

// Categories for request
func (r *BidRequest) Categories() []uint64 {
	// if r.categoryArray == nil {
	// 	if r.App != nil {
	// 	}
	// 	if r.Site != nil {
	// 	}
	// }
	return r.categoryArray
}

// IsSecure request
func (r *BidRequest) IsSecure() bool {
	return r.Secure == 1
}

// IsAdblock request
func (r *BidRequest) IsAdblock() bool {
	return r.Adblock == 1
}

// IsPrivateBrowsing request
func (r *BidRequest) IsPrivateBrowsing() bool {
	return r.PrivateBrowsing == 1
}

// SiteInfo object
func (r *BidRequest) SiteInfo() *udetect.Site {
	if r.Site != nil {
		return r.Site
	}
	if r.App == nil {
		return &udetect.SiteDefault
	}
	return nil
}

// AppInfo object
func (r *BidRequest) AppInfo() *udetect.App {
	return r.App
}

// UserInfo data
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

// DeviceInfo data
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

// DeviceID value
func (r *BidRequest) DeviceID() uint64 {
	if r != nil && r.Device != nil {
		return uint64(r.Device.ID)
	}
	return 0
}

// DeviceType item
func (r *BidRequest) DeviceType() uint64 {
	if r == nil {
		return 0
	}
	return uint64(r.DeviceInfo().DeviceType)
}

// OSInfo data
func (r *BidRequest) OSInfo() *udetect.OS {
	if r == nil {
		return nil
	}
	return r.DeviceInfo().OS
}

// BrowserInfo data
func (r *BidRequest) BrowserInfo() *udetect.Browser {
	if r == nil {
		return nil
	}
	return r.DeviceInfo().Browser
}

// MinECPM value of request acceptable
func (r *BidRequest) MinECPM() (minBid billing.Money) {
	for _, imp := range r.Imps {
		if minBid == 0 {
			minBid = imp.BidFloor
		} else if imp.BidFloor > 0 && minBid < imp.BidFloor {
			minBid = imp.BidFloor
		}
	}
	return
}

// GeoInfo data
func (r *BidRequest) GeoInfo() *udetect.Geo {
	if r == nil {
		return nil
	}
	return r.UserInfo().Geo
}

// CarrierInfo data
func (r *BidRequest) CarrierInfo() *udetect.Carrier {
	if geo := r.GeoInfo(); geo != nil {
		return geo.Carrier
	}
	return nil
}

// IsIPv6 address
func (r *BidRequest) IsIPv6() bool {
	return r != nil && r.User != nil && r.User.Geo != nil && r.User.Geo.IsIPv6()
}

// Get context item by key
func (r *BidRequest) Get(key string) any {
	if r.Ext == nil {
		return nil
	}
	return r.Ext[key]
}

// Set context item with key
func (r *BidRequest) Set(key string, val any) {
	if r.Ext == nil {
		r.Ext = map[string]any{}
	}
	r.Ext[key] = val
}

// Unset context item with keys
func (r *BidRequest) Unset(keys ...string) {
	if r.Ext == nil {
		return
	}

	for _, key := range keys {
		delete(r.Ext, key)
	}
}

// ImpressionUpdate each
func (r *BidRequest) ImpressionUpdate(fn func(imp *Impression) bool) {
	for i, imp := range r.Imps {
		if fn(&imp) {
			r.Imps[i] = imp
		}
	}
}

// ImpressionByID object
func (r *BidRequest) ImpressionByID(id string) *Impression {
	for _, im := range r.Imps {
		if im.ID == id {
			return &im
		}
	}
	return nil
}

// ImpressionByIDvariation returns impression by ID which can contains any postfix
func (r *BidRequest) ImpressionByIDvariation(id string) *Impression {
	for _, im := range r.Imps {
		if strings.HasPrefix(id, im.ID) {
			return &im
		}
	}
	return nil
}

// Time of request
func (r *BidRequest) Time() time.Time {
	return r.Timemark
}

// func (r *BidRequest) reset() {
// 	r.targetIDs = r.targetIDs[:0]
// 	r.externalTargetIDs = r.externalTargetIDs[:0]
// 	r.categoryArray = r.categoryArray[:0]
// 	r.domain = r.domain[:0]
// 	r.tags = r.tags[:0]
// 	r.formats = r.formats[:0]
// 	r.sourceIDs = r.sourceIDs[:0]
// 	r.Imps = r.Imps[:0]
// 	r.formatBitset.Reset()
// 	r.Tracer = nil
// 	r.Ext = nil
// }

///////////////////////////////////////////////////////////////////////////////
/// Validation
///////////////////////////////////////////////////////////////////////////////

// Validate request by currency
func (r *BidRequest) Validate() error {
	return nil
}
