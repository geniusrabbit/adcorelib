//
// @project GeniusRabbit corelib 2016 – 2019, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019, 2024
//

package adtype

import (
	openrtbnreq "github.com/bsm/openrtb/native/request"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/searchtypes"
)

// Impression target
type Impression struct {
	ID                string          `json:"id,omitempty"`                  // Internal impression ID
	ExtID             string          `json:"extid,omitempty"`               // External impression ID (ImpID)
	ExtTargetID       string          `json:"exttrgid"`                      // External zone ID (tagid)
	Request           any             `json:"request,omitempty"`             // Contains subrequest from RTB or another protocol
	Target            admodels.Target `json:"target,omitempty"`              //
	BidFloor          billing.Money   `json:"bid_floor,omitempty"`           //
	PurchaseViewPrice billing.Money   `json:"purchase_view_price,omitempty"` //
	Pos               int             `json:"pos,omitempty"`                 // 5.4 Ad Position
	Count             int             `json:"cnt,omitempty"`                 // Count of places for multiple banners

	// Sizes and position on the screen
	X    int `json:"x,omitempty"`  // Position on the site screen
	Y    int `json:"y,omitempty"`  //
	W    int `json:"w,omitempty"`  //
	H    int `json:"h,omitempty"`  //
	WMax int `json:"wm,omitempty"` //
	HMax int `json:"hm,omitempty"` //

	// Additional identifiers
	SubID1 string `json:"subid1,omitempty"`
	SubID2 string `json:"subid2,omitempty"`
	SubID3 string `json:"subid3,omitempty"`
	SubID4 string `json:"subid4,omitempty"`
	SubID5 string `json:"subid5,omitempty"`

	// Format types for impression
	FormatTypes  types.FormatTypeBitset `json:"-"`
	formats      []*types.Format
	formatBitset *searchtypes.NumberBitset[uint]

	Ext map[string]any `json:"ext,omitempty"`
}

// Init internal information
func (i *Impression) Init(formats types.FormatsAccessor) {
	var w, h, minw, minh = i.WMax, i.HMax, i.W, i.H
	if w <= 0 && h <= 0 {
		w, h = minw, minh
		minw, minh = minw-(minw/3), minh/3
	}
	if minw == 0 {
		minw = w - (w / 8)
	}
	if minh == 0 {
		minh = h - (h / 5)
	}

	i.formats = formats.FormatsBySize(w+10, h+10, minw, minh, i.FormatTypes)

	i.formatBitset = searchtypes.NewNumberBitset[uint]()
	for _, f := range i.formats {
		i.formatBitset.Set(uint(f.ID))
	}

	if i.FormatTypes.IsEmpty() {
		i.FormatTypes = *types.NewFormatTypeBitset().SetFromFormats(i.formats...)
	}
}

// Formats models
func (i *Impression) Formats() (f []*types.Format) {
	return i.formats
}

// FormatByType of formats
func (i *Impression) FormatByType(tp types.FormatType) *types.Format {
	for _, f := range i.formats {
		if f.Types.Is(tp) {
			return f
		}
	}
	return nil
}

// FormatBitset of IDs
func (i *Impression) FormatBitset() *searchtypes.NumberBitset[uint] {
	return i.formatBitset
}

// IDByFormat return specific ID to link format
func (i *Impression) IDByFormat(format *types.Format) string {
	return i.ID + "_" + format.Codename
}

// TargetID value
func (i *Impression) TargetID() uint {
	if i == nil || i.Target == nil {
		return 0
	}
	return uint(i.Target.ID())
}

// AccountID number
func (i *Impression) AccountID() uint64 {
	if i != nil && i.Target != nil {
		return i.Target.AccountID()
	}
	return 0
}

// IsDirect value
func (i *Impression) IsDirect() bool {
	return i.FormatTypes.Is(types.FormatDirectType)
}

// IsNative target support
func (i *Impression) IsNative() bool {
	return i.FormatTypes.Is(types.FormatNativeType)
}

// IsStandart target support
func (i *Impression) IsStandart() bool {
	return false ||
		i.FormatTypes.Is(types.FormatBannerType) ||
		i.FormatTypes.Is(types.FormatBannerHTML5Type)
}

// RevenueShareFactor value for the publisher company
func (i *Impression) RevenueShareFactor() float64 {
	if i == nil || i.Target == nil {
		return 0
	}
	return i.Target.RevenueShareFactor()
}

// ComissionShareFactor which system get from publisher from 0 to 1
func (i *Impression) ComissionShareFactor() float64 {
	if i == nil || i.Target == nil {
		return 0
	}
	return i.Target.ComissionShareFactor()
}

///////////////////////////////////////////////////////////////////////////////
/// OpenRTB methods
///////////////////////////////////////////////////////////////////////////////

// ContextType IDs 7.3
// @link https://www.iab.com/wp-content/uploads/2016/03/OpenRTB-Native-Ads-Specification-1-1_2016.pdf
func (i *Impression) ContextType() openrtbnreq.ContextTypeID {
	return openrtbnreq.ContextTypeContent
}

// ContextSubType IDs 7.4
// @link https://www.iab.com/wp-content/uploads/2016/03/OpenRTB-Native-Ads-Specification-1-1_2016.pdf
func (i *Impression) ContextSubType() openrtbnreq.ContextSubTypeID {
	return openrtbnreq.ContextSubTypeGeneral
}

// PlacementType IDs 7.5
// @link https://www.iab.com/wp-content/uploads/2016/03/OpenRTB-Native-Ads-Specification-1-1_2016.pdf
func (i *Impression) PlacementType() openrtbnreq.PlacementTypeID {
	return openrtbnreq.PlacementTypeRecommendation
}

// RTBNativeRequest object
func (i *Impression) RTBNativeRequest() (r *openrtbnreq.Request) {
	r, _ = i.Request.(*openrtbnreq.Request)
	return r
}

///////////////////////////////////////////////////////////////////////////////
/// Ext data methods
///////////////////////////////////////////////////////////////////////////////

// Get context item by key
func (i *Impression) Get(key string) any {
	if i.Ext == nil {
		return nil
	}
	return i.Ext[key]
}

// Set context item with key
func (i *Impression) Set(key string, val any) {
	if i.Ext == nil {
		i.Ext = map[string]any{}
	}
	i.Ext[key] = val
}

// Unset context item with keys
func (i *Impression) Unset(keys ...string) {
	if len(i.Ext) > 0 {
		for _, key := range keys {
			delete(i.Ext, key)
		}
	}
}
