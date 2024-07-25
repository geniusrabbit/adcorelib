//
// @project GeniusRabbit corelib 2016, 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016, 2018
//

// @TODO add format types

package adtype

import (
	"slices"
	"sort"

	"github.com/geniusrabbit/udetect"
)

// BaseFilter object
type BaseFilter struct {
	Secure               int8     // 0 - any, 1 - only secure, 2 - no secure
	Adblock              int8     // 0 - any, 1 - only adblock, 2 - no adblock
	PrivateBrowsing      int8     // 0 - any, 1 - only private, 2 - no private
	IP                   int8     // 0 - any, 1 - IPv4, 2 - IPv6
	Devices              []uint64 // Devices type
	OS                   []uint64
	OSExclude            []uint64
	Browsers             []uint64
	BrowsersExclude      []uint64
	Categories           []uint64
	Countries            []string
	Applications         []uint64
	ApplicationsExclude  []uint64
	Domains              []string
	DomainsExclude       []string
	Zones                []uint64
	ZonesExclude         []uint64
	ExternalZones        []string
	ExternalZonesExclude []string
}

// Normalise params
func (f *BaseFilter) Normalise() {
	slices.Sort(f.Devices)
	slices.Sort(f.Categories)
	slices.Sort(f.Countries)
	slices.Sort(f.Applications)
	slices.Sort(f.ApplicationsExclude)
	slices.Sort(f.Domains)
	slices.Sort(f.DomainsExclude)
	slices.Sort(f.Zones)
	slices.Sort(f.ZonesExclude)
	slices.Sort(f.ExternalZones)
	slices.Sort(f.ExternalZonesExclude)
}

// Test base from search request
func (f *BaseFilter) Test(request *BidRequest) bool {
	switch {
	case (request.IsSecure() && f.Secure == 2) || (!request.IsSecure() && f.Secure == 1):
		return false
	case (request.IsAdblock() && f.Adblock == 2) || (!request.IsAdblock() && f.Adblock == 1):
		return false
	case (request.IsPrivateBrowsing() && f.PrivateBrowsing == 2) || (!request.IsPrivateBrowsing() && f.PrivateBrowsing == 1):
		return false
	}

	var (
		deviceType  udetect.DeviceType
		countryCode string
	)

	if request.Device != nil {
		deviceType = request.Device.DeviceType
	}
	if request.User != nil {
		if request.User.Geo != nil {
			countryCode = request.User.Geo.Country
		}
	}

	return true &&
		(len(f.Devices) /* ***** */ < 1 || hasInTArr(uint64(deviceType), f.Devices)) &&
		(len(f.Countries) /* *** */ < 1 || hasInTArr(countryCode, f.Countries)) &&
		(len(f.Categories) /* ** */ < 1 || intersecTArr(request.Categories(), f.Categories)) &&
		(len(f.ApplicationsExclude) < 1 || !hasInTArr(request.AppID(), f.ApplicationsExclude)) &&
		(len(f.Applications) /*  */ < 1 || hasInTArr(request.AppID(), f.Applications)) &&
		(len(f.DomainsExclude) /**/ < 1 || !hasOneInStringArr(request.Domain(), f.DomainsExclude)) &&
		(len(f.Domains) /* ***** */ < 1 || hasOneInStringArr(request.Domain(), f.Domains)) &&
		(len(f.ZonesExclude) /*  */ < 1 || !intersecTArr(request.TargetIDs(), f.ZonesExclude)) &&
		(len(f.Zones) /* ******* */ < 1 || intersecTArr(request.TargetIDs(), f.Zones)) &&
		(len(f.ExternalZonesExclude) < 1 || !hasOneInStringArr(request.ExtTargetIDs(), f.ExternalZonesExclude)) &&
		(len(f.ExternalZones) /* */ < 1 || hasOneInStringArr(request.ExtTargetIDs(), f.ExternalZones))
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func hasInTArr[T ~string | ~int | ~int64 | ~uint | ~uint64](v T, arr []T) bool {
	i := sort.Search(len(arr), func(i int) bool { return arr[i] >= v })
	return i >= 0 && i < len(arr) && v == arr[i]
}

func hasOneInStringArr(arr1, arr2 []string) bool {
	for _, v := range arr1 {
		if hasInTArr(v, arr2) {
			return true
		}
	}
	return false
}

func intersecTArr[T ~string | ~int | ~int64 | ~uint | ~uint64](cat1, cat2 []T) bool {
	if len(cat1) < 1 && len(cat2) < 1 {
		return true
	}
	if len(cat1) < 1 || len(cat2) < 1 {
		return false
	}

	for _, c1 := range cat1 {
		if hasInTArr(c1, cat2) {
			return true
		}
	}
	return false
}
