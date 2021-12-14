//
// @project GeniusRabbit rotator 2017 - 2019, 2021
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019, 2021
//

package admodels

import (
	"strings"

	"github.com/geniusrabbit/gosql/pgtype"

	"geniusrabbit.dev/corelib/admodels/types"
	"geniusrabbit.dev/corelib/models"
)

// RTBAccessPoint for DSP connect.
// It means that this is entry point which contains
// information for access and search data
type RTBAccessPoint struct {
	ID      uint64
	Title   string
	Company *Company

	Codename string // Unical name of the access point
	Headers  pgtype.Hstore

	AuctionType types.AuctionType // default: 0 – first price type, 1 – second price type

	// RevenueShareReduce represents extra reduce factor to nevilate AdExchange and SSP discrepancy.
	// It means that the final bid respose will be decresed by RevenueShareReduce %
	// Example:
	//   1. Found advertisement with `bid=1.0$`
	//   2. Final `bid = bid - $AdSourceComission{%} - $AdExchangeComission{%} - $RevenueShareReduce{%}`
	RevenueShareReduce float64 // % 100_00, 10000 -> 100%, 6550 -> 65.5%

	Protocol string // rtb as default
	RPS      int    // 0 – unlimit
	Timeout  int    // In milliseconds

	Filter types.BaseFilter

	Flags pgtype.Hstore
}

// RTBAccessPoint create new DSP connect.
func RTBAccessPointFromModel(cl *models.RTBAccessPoint, comp *Company) (src *RTBAccessPoint) {
	if comp == nil {
		return nil
	}

	var (
		filter = types.BaseFilter{
			Secure:          cl.Secure,
			Adblock:         cl.AdBlock,
			PrivateBrowsing: cl.PrivateBrowsing,
			IP:              cl.IP,
		}
	)

	filter.Set(types.FieldFormat, cl.Formats)
	filter.Set(types.FieldDeviceTypes, cl.DeviceTypes)
	filter.Set(types.FieldDevices, cl.Devices)
	filter.Set(types.FieldOS, cl.OS)
	filter.Set(types.FieldBrowsers, cl.Browsers)
	filter.Set(types.FieldCategories, cl.Categories)
	filter.Set(types.FieldCountries, cl.Countries)
	filter.Set(types.FieldLanguages, cl.Languages)
	filter.Set(types.FieldZones, cl.Zones)
	filter.Set(types.FieldDomains, cl.Domains)

	return &RTBAccessPoint{
		ID:                 cl.ID,
		Company:            comp,
		Codename:           cl.Codename,
		Protocol:           strings.ToLower(cl.Protocol),
		Headers:            cl.Headers,
		AuctionType:        cl.AuctionType,
		RPS:                cl.RPS,
		Timeout:            cl.Timeout,
		Filter:             filter,
		RevenueShareReduce: cl.RevenueShareReduce,
		Flags:              cl.Flags,
	}
}

// Test RTB source
func (s *RTBAccessPoint) Test(t types.TargetPointer) bool {
	return s.Filter.Test(t)
}

// TestFormat available in filter
func (s *RTBAccessPoint) TestFormat(f *types.Format) bool {
	return s.Filter.TestFormat(f)
}

// RevenueShareReduceFactor from 0. to 1.
func (s *RTBAccessPoint) RevenueShareReduceFactor() float64 {
	return s.RevenueShareReduce / 10000.
}
