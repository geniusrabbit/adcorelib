//
// @project GeniusRabbit rotator 2016 – 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 – 2019
//

package events

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"geniusrabbit.dev/corelib/msgpack"

	// "geniusrabbit.dev/corelib/billing/billingthrift"
	// "geniusrabbit.dev/corelib/msgpack/thrift"
	msgjson "geniusrabbit.dev/corelib/msgpack/json"
)

// Event status
const (
	StatusUndefined   = 0
	StatusSuccess     = 1
	StatusFailed      = 2
	StatusCompromised = 3
	StatusCustom      = 4 // User code
)

// var streamCoder = thrift.NewAPI(
// 	billingthrift.MoneyExt{},
// 	typeThriftExy{},
// )

var (
	streamCodeEncoder = &msgjson.EncodeGenerator{}
	streamCodeDecoder = &msgjson.DecodeGenerator{}
)

// Event struct
type Event struct {
	Time     int64  `thrift:",1" json:"tm,omitempty"`  // Timestamp
	Delay    uint64 `thrift:",2" json:"dl,omitempty"`  // Delay of preparation of Ads in Nanosecinds
	Duration uint64 `thrift:",3" json:"d,omitempty"`   // Duration in Nanoseconds
	Service  string `thrift:",4" json:"srv,omitempty"` // Service sender
	Cluster  string `thrift:",5" json:"cl,omitempty"`  // Cluster code (eu, us, as)
	Event    Type   `thrift:",6" json:"e,omitempty"`   // Event code string
	Status   uint8  `thrift:",7" json:"st,omitempty"`  // Status: 0 - undefined, 1 - success, 2 - failed, 3 - compromised
	// Accounts link information
	Project           uint64 `thrift:",8"  json:"pr,omitempty"`  // Project network ID
	PublisherCompany  uint64 `thrift:",9"  json:"pcb,omitempty"` // -- // --
	AdvertiserCompany uint64 `thrift:",10" json:"acv,omitempty"` // -- // --
	// Source
	AuctionID    string `thrift:",11" json:"auc,omitempty"`     // Internal Auction ID
	AuctionType  uint8  `thrift:",12" json:"auctype,omitempty"` // Aution type 1 - First price, 2 - Second price
	ImpID        string `thrift:",13" json:"imp,omitempty"`     // Sub ID of request for paticular impression spot
	ImpAdID      string `thrift:",14" json:"impad,omitempty"`   // Specific ID for paticular ad impression
	ExtAuctionID string `thrift:",15" json:"eauc,omitempty"`    // RTB Request/Response ID
	ExtImpID     string `thrift:",16" json:"eimp,omitempty"`    // RTB Imp ID
	ExtTargetID  string `thrift:",17" json:"extz,omitempty"`    // RTB Zone ID (tagid)
	Source       uint64 `thrift:",18" json:"src,omitempty"`     // Advertisement Source ID
	Network      string `thrift:",19" json:"net,omitempty"`     // Source Network Name or Domain (Cross sails)
	AccessPoint  uint64 `thrift:",20" json:"acp,omitempty"`     // Access Point ID to own Advertisement
	// State Location
	Platform    int    `thrift:",21" json:"pl,omitempty"`  // Where displaid? 0 – undefined, 1 – web site, 2 – native app, 3 – game
	Domain      string `thrift:",22" json:"dm,omitempty"`  // If not web site then "bundle"
	Application uint64 `thrift:",23" json:"app,omitempty"` // application ID (registered in the system)
	Zone        uint64 `thrift:",24" json:"z,omitempty"`   // -- // --
	PixelID     uint64 `thrift:",25" json:"pxl,omitempty"` // -- // --
	Campaign    uint64 `thrift:",26" json:"cmp,omitempty"` // Campaign info
	FormatID    uint64 `thrift:",27" json:"fmt,omitempty"` // Format ID
	AdID        uint64 `thrift:",28" json:"ad,omitempty"`  // -- // --
	AdW         int    `thrift:",29" json:"aw,omitempty"`  // -- // --
	AdH         int    `thrift:",30" json:"ah,omitempty"`  // -- // --
	SourceURL   string `thrift:",31" json:"su,omitempty"`  // Advertisement source URL (iframe, image, video, direct)
	WinURL      string `thrift:",32" json:"wu,omitempty"`  // Win URL used for RTB confirmation
	URL         string `thrift:",33" json:"u,omitempty"`   // Non modified target URL
	Jumper      uint64 `thrift:",34" json:"j,omitempty"`   // Jumper Page ID
	// Money section
	// We have a different situations how we are pay for the traffic
	// For example
	// 	1) If we have external traffic resource:
	//		ViewPrice = 0.01 CPV - Have to pay advertiser
	//		LeadPrice = 40$ - How much will we pay for the lead
	//		PurchaseViewPrice = 0.001 CPV - how much we will pay for this traffic
	//	2) If external second price auction type
	//		ViewPrice = 0.02 CPV - Have to pay advertiser
	//		PurchaseViewPrice = 0.01 CPV - how much we will pay for this traffic
	//  3) If we are buying the traffic for CPA
	//		LeadPrice = 20$ - Have to pay advertiser
	//		PurchaseViewPrice = 10$ - Have to pay to the source
	PricingModel uint `thrift:",36" json:"pm,omitempty"` // Display As CPM/CPC/CPA/CPI
	// Prurchase price from SSP or other TRAFFIC sources (menetisation of income requests)
	// We are buying the place of advertisement display
	PurchaseViewPrice   int64   `thrift:",37" json:"pvpr,omitempty"`   // Price of the view of source traffic cost
	PurchaseClickPrice  int64   `thrift:",38" json:"pcpr,omitempty"`   // Price of the click of source traffic cost
	PurchaseLeadPrice   int64   `thrift:",39" json:"plpr,omitempty"`   // Price of the lead of source traffic cost
	PotentialViewPrice  int64   `thrift:",40" json:"ptvpr,omitempty"`  // Price of the view of source of ads oitential
	PotentialClickPrice int64   `thrift:",41" json:"ptcpr,omitempty"`  // Price of the click of source of ads oitential
	PotentialLeadPrice  int64   `thrift:",42" json:"ptlpr,omitempty"`  // Price of the lead of source of ads oitential
	ViewPrice           int64   `thrift:",43" json:"vpr,omitempty"`    // Price per view
	ClickPrice          int64   `thrift:",44" json:"cpr,omitempty"`    // Price per click
	LeadPrice           int64   `thrift:",45" json:"lpr,omitempty"`    // Price per lead
	Competitor          uint64  `thrift:",46" json:"cmid,omitempty"`   // Competitor compaign ID
	CompetitorSource    uint64  `thrift:",47" json:"cmsrc,omitempty"`  // Competitor source ID
	CompetitorECPM      float64 `thrift:",48" json:"cmecpm,omitempty"` // Competitor ECPM or auction
	// User IDENTITY
	UDID        string `thrift:",49" json:"udi,omitempty"` // Unique Device ID (IDFA)
	UUID        string `thrift:",50" json:"uui,omitempty"` // User
	SessionID   string `thrift:",51" json:"ses,omitempty"` // -- // --
	Fingerprint string `thrift:",52" json:"fpr,omitempty"` // -- // --
	ETag        string `thrift:",53" json:"etg,omitempty"` // -- // --
	// Targeting
	Carrier         uint   `thrift:",54" json:"car,omitempty"`  // -- // --
	Country         string `thrift:",55" json:"cc,omitempty"`   // Country Code ISO-2
	City            string `thrift:",56" json:"ct,omitempty"`   // City Code
	Latitude        string `thrift:",57" json:"lt,omitempty"`   // -- // --
	Longitude       string `thrift:",58" json:"lg,omitempty"`   // -- // --
	Language        string `thrift:",59" json:"lng,omitempty"`  // en-US
	IPString        string `thrift:",60" json:"ip,omitempty"`   // -- // --
	Referer         string `thrift:",61" json:"ref,omitempty"`  // Referal link
	Page            string `thrift:",62" json:"page,omitempty"` // Page link
	UserAgent       string `thrift:",63" json:"ua,omitempty"`   // -- // --
	DeviceType      uint   `thrift:",64" json:"dvt,omitempty"`  // Device type 0 - Undefined, 1 - Desktop, etc.
	Device          uint   `thrift:",65" json:"dv,omitempty"`   // Device ID
	OS              uint   `thrift:",66" json:"os,omitempty"`   // OS ID
	Browser         uint   `thrift:",67" json:"br,omitempty"`   // Browser ID
	Categories      string `thrift:",68" json:"c,omitempty"`    // -- // --
	Adblock         uint   `thrift:",69" json:"ab,omitempty"`   // Ext params
	PrivateBrowsing uint   `thrift:",70" json:"prv,omitempty"`  // Private Mode
	Robot           uint   `thrift:",71" json:"rt,omitempty"`   // Robot traffic
	Proxy           uint   `thrift:",72" json:"pt,omitempty"`   // Proxy traffic
	Backup          uint   `thrift:",73" json:"bt,omitempty"`   // Backup Display Type
	X               int    `thrift:",74" json:"x,omitempty"`    // X - coord of addisplay or click position
	Y               int    `thrift:",75" json:"y,omitempty"`    // Y - coord of addisplay or click position
	W               int    `thrift:",76" json:"w,omitempty"`    // W - available space
	H               int    `thrift:",77" json:"h,omitempty"`    // H - available space

	SubID1 string `thrift:",78" json:"sd1,omitempty"`
	SubID2 string `thrift:",79" json:"sd2,omitempty"`
	SubID3 string `thrift:",80" json:"sd3,omitempty"`
	SubID4 string `thrift:",81" json:"sd4,omitempty"`
	SubID5 string `thrift:",82" json:"sd5,omitempty"`
}

// String imprlementation of Stringer interface
func (e *Event) String() string {
	if json, err := json.MarshalIndent(e, "", "  "); err == nil {
		return string(json)
	}
	return ""
}

// Normalize some params
func (e *Event) Normalize() *Event {
	e.X = positiveNumber(e.X)
	e.Y = positiveNumber(e.Y)
	e.W = positiveNumber(e.W)
	e.H = positiveNumber(e.H)
	e.AdW = positiveNumber(e.AdW)
	e.AdH = positiveNumber(e.AdH)
	return e
}

// Validate event object
func (e *Event) Validate() error {
	return nil
}

// DateTime object
func (e *Event) DateTime() time.Time {
	return time.Unix(0, e.Time)
}

// SetDateTime timestamp
func (e *Event) SetDateTime(tmNano int64) {
	e.Time = tmNano
}

// UpdateDateTime of event
func (e *Event) UpdateDateTime() {
	e.Time = time.Now().UnixNano()
}

// TimeFormat to string
func (e *Event) TimeFormat(format string) string {
	return e.DateTime().Format(format)
}

// EncodeCodeOld URL object
func (e *Event) EncodeCodeOld() (string, error) {
	data, err := msgpack.StdPack(e)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// DecodeCodeOld URL object
func (e *Event) DecodeCodeOld(data string) error {
	var decoded, err = base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}

	if err = msgpack.StdUnpack(decoded, e); err == nil {
		err = e.Validate()
	}
	return err
}

// Pack object to event Code
func (e *Event) Pack() Code {
	return ObjectCode(e, streamCodeEncoder)
}

// Unpack event object
func (e *Event) Unpack(data []byte, unpuckFnc ...func(code Code) Code) error {
	code := CodeObj(data, nil)
	if len(unpuckFnc) > 0 && unpuckFnc[0] != nil {
		code = unpuckFnc[0](code)
	}
	return code.DecodeObject(e, streamCodeDecoder)
}

// PreparedURL string
func (e *Event) PreparedURL(https bool, args ...string) string {
	var url = e.URL
	if !strings.HasPrefix(e.URL, "http://") &&
		!strings.HasPrefix(e.URL, "https://") {
		if strings.HasPrefix(e.URL, "//") {
			if https {
				url = "https:" + url
			} else {
				url = "http:" + url
			}
		} else if https {
			url = "https://" + url
		} else {
			url = "http://" + url
		}
	}

	if len(args) > 1 {
		return strings.NewReplacer(args...).Replace(url)
	}
	return url
}

///////////////////////////////////////////////////////////////////////////////
/// helpers
///////////////////////////////////////////////////////////////////////////////

func fixNumberString(n, def string) string {
	n = strings.TrimSpace(n)
	if strings.Trim(n, "0123456789") != "" {
		n = def
	}
	return n
}

func positiveNumber(v int) int {
	if v < 0 {
		return 0
	}
	return v
}
