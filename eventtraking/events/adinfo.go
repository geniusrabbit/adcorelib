package events

// AdInfo represents some additional information which aggregated during
// advertisement and ad-request processing
type AdInfo struct {
	Time     int64  `thrift:",1" json:"tm,omitempty"`   // Timestamp
	Type     string `thrift:",2" json:"type,omitempty"` // Type of message info
	Duration uint64 `thrift:",3" json:"d,omitempty"`    // Duration in Nanoseconds
	Service  string `thrift:",4" json:"srv,omitempty"`  // Service sender
	Cluster  string `thrift:",5" json:"cl,omitempty"`   // Cluster code (eu, us, as)
	Param1   int    `thrift:",6" json:"p1,omitempty"`   // Reserved
	Param2   int    `thrift:",7" json:"p2,omitempty"`   // Reserved
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
}
