//
// @project GeniusRabbit corelib 2016 – 2019, 2024
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019, 2024
//

package adtype

import (
	openrtbnreq "github.com/bsm/openrtb/native/request"
	openrtbnreq3 "github.com/bsm/openrtb/v3/native/request"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/searchtypes"
)

// Impression target object represents impression place
// which can be used for the ad placement
type Impression struct {
	ID               string `json:"id,omitempty"`    // Internal impression ID
	ExternalID       string `json:"extid,omitempty"` // External impression ID (ImpID)
	ExternalTargetID string `json:"exttrgid"`        // External zone/unit/spot ID (tagid)

	Request any    `json:"request,omitempty"` // Contains subrequest from RTB or another protocol
	Target  Target `json:"target,omitempty"`

	BidFloorCPM       billing.Money `json:"bid_floor,omitempty"`
	PurchaseViewPrice billing.Money `json:"purchase_view_price,omitempty"`

	// Ad position
	Pos   int `json:"pos,omitempty"` // 5.4 Ad Position
	Count int `json:"cnt,omitempty"` // Count of places for multiple banners

	// Sizes and position on the screen
	X         int `json:"x,omitempty"`
	Y         int `json:"y,omitempty"`
	Width     int `json:"w,omitempty"`
	Height    int `json:"h,omitempty"`
	WidthMax  int `json:"wm,omitempty"`
	HeightMax int `json:"hm,omitempty"`

	// Additional identifiers
	SubID1 string `json:"subid1,omitempty"`
	SubID2 string `json:"subid2,omitempty"`
	SubID3 string `json:"subid3,omitempty"`
	SubID4 string `json:"subid4,omitempty"`
	SubID5 string `json:"subid5,omitempty"`

	// Format types for impression
	FormatTypes  types.FormatTypeBitset `json:"-"`
	FormatCodes  []string               `json:"-"`
	formats      []*types.Format
	formatBitset searchtypes.NumberBitset[uint]

	Ext map[string]any `json:"ext,omitempty"`
}

// InitFormats initialize formats by size and types or by codes
func (i *Impression) InitFormats(formats types.FormatsAccessor) {
	if i == nil {
		return
	}
	if len(i.FormatCodes) > 0 {
		i.InitFormatsByCodes(i.FormatCodes, formats)
	} else {
		i.InitFormatsBySizeAndTypes(i.WidthMax, i.HeightMax, i.Width, i.Height, i.FormatTypes, formats)
	}
}

// Init internal information
func (i *Impression) InitFormatsBySizeAndTypes(w, h, minw, minh int, formatTypes types.FormatTypeBitset, formats types.FormatsAccessor) {
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

	// Extract formats by size
	i.formats = formats.FormatsBySize(w+10, h+10, minw, minh, formatTypes)

	// Update format types
	i.formatBitset.Reset()
	if formatTypes.IsEmpty() {
		for _, f := range i.formats {
			i.FormatTypes.SetOneBitset(f.Types)
			i.formatBitset.Set(uint(f.ID))
		}
	} else {
		i.FormatTypes = formatTypes
		for _, f := range i.formats {
			i.formatBitset.Set(uint(f.ID))
		}
	}
}

// InitFormatsByCodes initialize formats by codes
func (i *Impression) InitFormatsByCodes(codes []string, formats types.FormatsAccessor) {
	if len(codes) < 1 {
		return
	}
	i.formats = i.formats[:0]
	i.formatBitset.Reset()
	i.FormatTypes.Reset()

	for _, code := range codes {
		if f := formats.FormatByCode(code); f != nil {
			i.FormatTypes.SetOneBitset(f.Types)
			i.formats = append(i.formats, f)
			i.formatBitset.Set(uint(f.ID))
		}
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
	return &i.formatBitset
}

// IDByFormat return specific ID to link format
func (i *Impression) IDByFormat(format *types.Format) string {
	return i.ID + "_" + format.Codename
}

// TargetID returns the target ID
func (i *Impression) TargetID() uint {
	if i == nil || i.Target == nil {
		return 0
	}
	return uint(i.Target.ID())
}

// AccountID returns the account ID
func (i *Impression) AccountID() uint64 {
	if i != nil && i.Target != nil {
		return i.Target.Account().ID()
	}
	return 0
}

// IsDirect returns true if the impression is direct
func (i *Impression) IsDirect() bool {
	return i.FormatTypes.Is(types.FormatDirectType)
}

// IsNative returns true if the impression is native
func (i *Impression) IsNative() bool {
	return i.FormatTypes.Is(types.FormatNativeType)
}

// IsStandart returns true if the impression is standart
func (i *Impression) IsStandart() bool {
	return false ||
		i.FormatTypes.Is(types.FormatBannerType) ||
		i.FormatTypes.Is(types.FormatBannerHTML5Type)
}

// CommissionShareFactor which system get from publisher from 0 to 1
//
//go:inline
func (i *Impression) CommissionShareFactor() float64 {
	if i == nil || i.Target == nil {
		return 0
	}
	return i.Target.CommissionShareFactor()
}

// PurchasePrice return the price of need to pay for the action
// to the connected network or application if price is fixed
func (i *Impression) PurchasePrice(action Action) billing.Money {
	if i == nil {
		return 0
	}
	if action.IsView() && i.PurchaseViewPrice > 0 {
		return i.PurchaseViewPrice
	}
	if i.Target != nil {
		return i.Target.PurchasePrice(action)
	}
	return 0
}

// RequestObj return the request object related to the impression
func (i *Impression) RequestObj() any {
	if i == nil {
		return nil
	}
	return i.Request
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
func (i *Impression) RTBNativeRequest() *openrtbnreq.Request {
	r, ok := i.Request.(*openrtbnreq.Request)
	if !ok {
		return nil
	}
	return r
}

// RTBNativeRequestV3 object
func (i *Impression) RTBNativeRequestV3() *openrtbnreq3.Request {
	r, ok := i.Request.(*openrtbnreq3.Request)
	if !ok {
		return nil
	}
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
