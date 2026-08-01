package bidrequest

import (
	"time"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
)

// BidTargetWrapper wraps a [BidRequest] and an [adtype.Impression] together
// to implement [types.TargetPointer]. Methods that are impression-specific
// (Size, TargetID, IsInterstitial) are resolved from Imp; all others delegate
// to BidReq.
type BidTargetWrapper struct {
	BidReq *BidRequest
	Imp    *adtype.Impression
}

// CountryCode implements [types.TargetPointer].
func (b *BidTargetWrapper) CountryCode() string {
	return b.BidReq.GeoInfo().Country
}

// ExtarnalTargetID implements [types.TargetPointer].
func (b *BidTargetWrapper) ExtarnalTargetID() string {
	return b.Imp.ExternalTargetID
}

// LanguageCode implements [types.TargetPointer].
func (b *BidTargetWrapper) LanguageCode() string {
	if langs := b.BidReq.BrowserInfo().Languages; len(langs) > 0 {
		return langs[0]
	}
	return ""
}

// Age returns the estimated age of the user in years.
func (b *BidTargetWrapper) Age() uint {
	return b.BidReq.Age()
}

// AppID returns the target application ID.
func (b *BidTargetWrapper) AppID() uint64 {
	return b.BidReq.AppID()
}

// AppInfo returns metadata about the target application.
func (b *BidTargetWrapper) AppInfo() *types.AppInfo {
	return b.BidReq.AppInfo()
}

// BrowserInfo returns the detected browser metadata.
func (b *BidTargetWrapper) BrowserInfo() *types.BrowserInfo {
	return b.BidReq.BrowserInfo()
}

// CarrierInfo returns the mobile carrier metadata.
func (b *BidTargetWrapper) CarrierInfo() *types.CarrierInfo {
	return b.BidReq.CarrierInfo()
}

// Categories returns the IAB content category IDs of the request.
func (b *BidTargetWrapper) Categories() []uint64 {
	return b.BidReq.Categories()
}

// CurrentGeoTime returns the current wall-clock time adjusted to the user's geo timezone.
func (b *BidTargetWrapper) CurrentGeoTime() time.Time {
	return b.BidReq.CurrentGeoTime()
}

// DeviceInfo returns the full device metadata.
func (b *BidTargetWrapper) DeviceInfo() *types.DeviceInfo {
	return b.BidReq.DeviceInfo()
}

// Domain returns the list of domains associated with the request.
func (b *BidTargetWrapper) Domain() []string {
	return b.BidReq.Domain()
}

// DomainName returns the primary domain or app bundle name.
func (b *BidTargetWrapper) DomainName() string {
	return b.BidReq.DomainName()
}

// Formats returns the ad formats available in the bid request.
func (b *BidTargetWrapper) Formats() types.BidFormater {
	return b.BidReq.Formats()
}

// GeoID returns the geographic region ID of the user.
func (b *BidTargetWrapper) GeoID() uint64 {
	return b.BidReq.GeoID()
}

// GeoInfo returns the full geographic metadata of the user.
func (b *BidTargetWrapper) GeoInfo() *types.GeoInfo {
	return b.BidReq.GeoInfo()
}

// IsAdBlock returns true if an ad-blocker was detected on the client.
func (b *BidTargetWrapper) IsAdBlock() bool {
	return b.BidReq.IsAdBlock()
}

// IsDebug returns true if the request is in debug mode.
func (b *BidTargetWrapper) IsDebug() bool {
	return b.BidReq.IsDebug()
}

// IsIPv6 returns true if the client is connected over IPv6.
func (b *BidTargetWrapper) IsIPv6() bool {
	return b.BidReq.IsIPv6()
}

// IsInterstitial returns true if the impression slot is interstitial (full-screen).
// The value is read from Imp, not from BidReq.
func (b *BidTargetWrapper) IsInterstitial() bool {
	return b.Imp.IsInterstitial()
}

// IsPrivateBrowsing returns true if the request originates from a private / incognito session.
func (b *BidTargetWrapper) IsPrivateBrowsing() bool {
	return b.BidReq.IsPrivateBrowsing()
}

// IsProxy returns true if a proxy or VPN was detected.
func (b *BidTargetWrapper) IsProxy() bool {
	return b.BidReq.IsProxy()
}

// IsRobot returns true if the client was identified as a bot or crawler.
func (b *BidTargetWrapper) IsRobot() bool {
	return b.BidReq.IsRobot()
}

// IsSecure returns true if the request was made over HTTPS.
func (b *BidTargetWrapper) IsSecure() bool {
	return b.BidReq.IsSecure()
}

// LanguageID returns the browser language ID.
func (b *BidTargetWrapper) LanguageID() uint64 {
	return b.BidReq.LanguageID()
}

// MinECPM returns the minimum effective CPM floor for this request.
func (b *BidTargetWrapper) MinECPM() billing.Money {
	if b == nil || b.Imp == nil {
		return 0
	}
	return b.Imp.BidFloorCPM
}

// OSInfo returns the detected operating system metadata.
func (b *BidTargetWrapper) OSInfo() *types.OSInfo {
	return b.BidReq.OSInfo()
}

// Sex returns the inferred sex of the user.
func (b *BidTargetWrapper) Sex() uint {
	return b.BidReq.Sex()
}

// SiteInfo returns the site metadata associated with the request.
func (b *BidTargetWrapper) SiteInfo() *types.SiteInfo {
	return b.BidReq.SiteInfo()
}

// Size returns the width and height of the impression slot.
// The value is read from Imp, not from BidReq.
func (b *BidTargetWrapper) Size() (w, h int) {
	return b.Imp.Size()
}

// Tags returns the targeting tag list associated with the request.
func (b *BidTargetWrapper) Tags() []string {
	return b.BidReq.Tags()
}

// TargetID returns the ID of the specific ad zone / placement target.
// Returns 0 if the wrapper, impression, or target is nil.
// The value is read from Imp.Target, not from BidReq.
func (b *BidTargetWrapper) TargetID() uint64 {
	if b == nil || b.Imp == nil || b.Imp.Target == nil {
		return 0
	}
	return b.Imp.Target.ID()
}

// Time returns the wall-clock time at which the bid request started.
func (b *BidTargetWrapper) Time() time.Time {
	return b.BidReq.Time()
}

// TrafficSourceID returns the ID of the traffic source that delivered the request.
func (b *BidTargetWrapper) TrafficSourceID() uint64 {
	return b.BidReq.TrafficSourceID()
}

var _ types.TargetPointer = (*BidTargetWrapper)(nil)
