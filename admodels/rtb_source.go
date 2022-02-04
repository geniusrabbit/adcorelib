//
// @project GeniusRabbit rotator 2017 - 2019, 2021
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019, 2021
//

package admodels

import (
	"strings"

	"github.com/geniusrabbit/gosql"

	"geniusrabbit.dev/corelib/admodels/types"
	"geniusrabbit.dev/corelib/models"
)

// RTBRequestType contains type of representation of request information
type RTBRequestType = types.RTBRequestType

// Request types
const (
	RTBRequestTypeUndefined       = types.RTBRequestTypeUndefined
	RTBRequestTypeJSON            = types.RTBRequestTypeJSON
	RTBRequestTypeXML             = types.RTBRequestTypeXML
	RTBRequestTypeProtoBUFF       = types.RTBRequestTypeProtoBUFF
	RTBRequestTypePOSTFormEncoded = types.RTBRequestTypePOSTFormEncoded
	RTBRequestTypePLAINTEXT       = types.RTBRequestTypePLAINTEXT
)

// RTBSourceOptions flags
type RTBSourceOptions struct {
	ErrorsIgnore bool
	Trace        bool
}

// RTBSource describe the source of external DSP platform or similar exchange protocol.
// All that sources have similar options and very common prefilter configurations
type RTBSource struct {
	ID      uint64
	Company *Company

	Protocol    string         // rtb as default
	URL         string         // RTB client request URL
	Method      string         // HTTP method GET, POST, ect; Default POST
	RequestType RTBRequestType // 1 - json, 2 - xml, 3 - ProtoBUFF, 4 - MultipleFormaData, 5 - PLAINTEXT
	Headers     gosql.NullableJSON

	AuctionType types.AuctionType // default: 0 – first price type, 1 – second price type
	RPS         int               // 0 – unlimit
	Timeout     int               // In milliseconds
	Options     RTBSourceOptions  //
	Filter      types.BaseFilter  //

	Accuracy              float64 // Price accuracy for auction in percentages
	PriceCorrectionReduce float64 // % 100, 80%, 65.5%
	MinimalWeight         float64

	Flags  gosql.NullableJSON
	Config gosql.NullableJSON
}

// RTBSourceFromModel convert database model to specified model
func RTBSourceFromModel(cl *models.RTBSource, comp *Company) (src *RTBSource) {
	if comp == nil {
		return nil
	}

	var (
		opt = RTBSourceOptions{
			ErrorsIgnore: cl.Flag("errors_ignore") == 1,
			Trace:        cl.Flag("trace") == 1,
		}
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

	return &RTBSource{
		ID:                    cl.ID,
		Company:               comp,
		Protocol:              strings.ToLower(cl.Protocol),
		URL:                   cl.URL,
		Method:                strings.ToUpper(cl.Method),
		RequestType:           cl.RequestType,
		Headers:               cl.Headers,
		AuctionType:           cl.AuctionType,
		RPS:                   cl.RPS,
		Timeout:               cl.Timeout,
		Options:               opt,
		Filter:                filter,
		Accuracy:              cl.Accuracy,
		PriceCorrectionReduce: cl.PriceCorrectionReduce,
		MinimalWeight:         cl.MinimalWeight,
		Flags:                 cl.Flags,
		Config:                cl.Config,
	}
}

// Test RTB source
func (s *RTBSource) Test(t types.TargetPointer) bool {
	return s.Filter.Test(t)
}

// TestFormat available in filter
func (s *RTBSource) TestFormat(f *types.Format) bool {
	return s.Filter.TestFormat(f)
}

// PriceCorrectionReduceFactor which is a potential
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (s *RTBSource) PriceCorrectionReduceFactor() float64 {
	return s.PriceCorrectionReduce / 100.0
}
