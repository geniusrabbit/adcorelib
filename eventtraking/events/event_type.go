package events

// Type of event
type Type string

func (t Type) String() string { return string(t) }

// Event types
const (
	Undefined  Type = ""
	Request    Type = "request"
	Impression Type = "impression"
	View       Type = "view"
	Direct     Type = "direct"
	Click      Type = "click"
	Lead       Type = "lead"
	// Source types
	SourceNoBid Type = "src.nobid"
	SourceBid   Type = "src.bid"
	SourceWin   Type = "src.win"
	SourceFail  Type = "src.fail"
	SourceSkip  Type = "src.skip"
	// Access Point types
	AccessPointNoBid Type = "ap.nobid"
	AccessPointBid   Type = "ap.bid"
	AccessPointWin   Type = "ap.win"
	AccessPointFail  Type = "ap.fail"
	AccessPointSkip  Type = "ap.skip"
)
