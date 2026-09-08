package types

import (
	"testing"
	"time"

	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/searchtypes"
	"github.com/geniusrabbit/gosql/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetExtZonesExclude(t *testing.T) {
	var fl BaseFilter
	fl.SetExtZones(gosql.NullableStringArray{"-tspop-zone-2", "-other"})
	assert.Equal(t, gosql.StringArray{"tspop-zone-2", "other"}, fl.ExtZones)
	assert.NotZero(t, fl.excludeMask&(1<<FieldExtZones))
}

func TestSetExtAppsInclude(t *testing.T) {
	var fl BaseFilter
	fl.SetExtApps(gosql.NullableStringArray{"34002", "app-b"})
	assert.Equal(t, gosql.StringArray{"34002", "app-b"}, fl.ExtApps)
	assert.Zero(t, fl.excludeMask&(1<<FieldExtApps))
}

func TestSetExtZonesEmpty(t *testing.T) {
	var fl BaseFilter
	fl.SetExtZones(nil)
	assert.Empty(t, fl.ExtZones)
	assert.Zero(t, fl.excludeMask&(1<<FieldExtZones))
}

func TestResetClearsExtFilters(t *testing.T) {
	var fl BaseFilter
	fl.SetExtZones(gosql.NullableStringArray{"-z1"})
	fl.SetExtApps(gosql.NullableStringArray{"-a1"})
	fl.Reset()
	assert.Empty(t, fl.ExtZones)
	assert.Empty(t, fl.ExtApps)
	assert.Zero(t, fl.excludeMask)
}

func TestExtZonesExcludeBlocksExternalTargetID(t *testing.T) {
	var fl BaseFilter
	fl.SetExtZones(gosql.NullableStringArray{"-tspop-zone-2"})
	err := fl.Test(stubTarget{zoneExt: "tspop-zone-2"})
	require.ErrorIs(t, err, ErrTargetNotAllowed)
	err = fl.Test(stubTarget{zoneExt: "other-zone"})
	require.NoError(t, err)
}

func TestExtAppsIncludeMatchesAppThenSiteExtID(t *testing.T) {
	var fl BaseFilter
	fl.SetExtApps(gosql.NullableStringArray{"34002"})
	err := fl.Test(stubTarget{appExt: "34002"})
	require.NoError(t, err)
	err = fl.Test(stubTarget{siteExt: "34002"})
	require.NoError(t, err)
	err = fl.Test(stubTarget{appExt: "other"})
	require.ErrorIs(t, err, ErrAppNotAllowed)
}

type stubTarget struct {
	appExt  string
	siteExt string
	zoneExt string
}

type emptyFormats struct{}

func (emptyFormats) List() []*Format                         { return nil }
func (emptyFormats) Bitset() *searchtypes.NumberBitset[uint] { return nil }
func (emptyFormats) TypeMask() FormatTypeBitset              { return 0 }

func (s stubTarget) Formats() BidFormater    { return emptyFormats{} }
func (stubTarget) Size() (int, int)          { return 0, 0 }
func (stubTarget) IsDebug() bool             { return false }
func (stubTarget) IsSecure() bool            { return false }
func (stubTarget) IsAdBlock() bool           { return false }
func (stubTarget) IsPrivateBrowsing() bool   { return false }
func (stubTarget) IsRobot() bool             { return false }
func (stubTarget) IsProxy() bool             { return false }
func (stubTarget) IsIPv6() bool              { return false }
func (stubTarget) IsInterstitial() bool      { return false }
func (stubTarget) DeviceInfo() *DeviceInfo   { return &DeviceInfo{} }
func (stubTarget) BrowserInfo() *BrowserInfo { return &BrowserInfo{} }
func (stubTarget) OSInfo() *OSInfo           { return &OSInfo{} }
func (stubTarget) TrafficSourceID() uint64   { return 0 }
func (stubTarget) AppID() uint64             { return 0 }
func (s stubTarget) AppInfo() *AppInfo {
	if s.appExt == "" {
		return nil
	}
	return &AppInfo{ExtID: s.appExt}
}
func (s stubTarget) SiteInfo() *SiteInfo {
	if s.siteExt == "" {
		return nil
	}
	return &SiteInfo{ExtID: s.siteExt}
}
func (stubTarget) Domain() []string           { return nil }
func (stubTarget) DomainName() string         { return "" }
func (stubTarget) GeoID() uint64              { return 0 }
func (stubTarget) GeoInfo() *GeoInfo          { return &GeoInfo{} }
func (stubTarget) CountryCode() string        { return "" }
func (stubTarget) CarrierInfo() *CarrierInfo  { return nil }
func (stubTarget) LanguageID() uint64         { return 0 }
func (stubTarget) LanguageCode() string       { return "" }
func (stubTarget) TargetID() uint64           { return 0 }
func (s stubTarget) ExtarnalTargetID() string { return s.zoneExt }
func (stubTarget) AuctionType() AuctionType   { return 0 }
func (stubTarget) Sex() uint                  { return 0 }
func (stubTarget) Age() uint                  { return 0 }
func (stubTarget) Tags() []string             { return nil }
func (stubTarget) Categories() []uint64       { return nil }
func (stubTarget) MinECPM() billing.Money     { return 0 }
func (stubTarget) Time() time.Time            { return time.Time{} }
func (stubTarget) CurrentGeoTime() time.Time  { return time.Time{} }
