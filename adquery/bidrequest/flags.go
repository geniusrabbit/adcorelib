package bidrequest

// BidRequestFlags defines flags for bid requests.
type BidRequestFlags uint8

const (
	// BidRequestFlagAdBlock indicates if adblock is enabled
	BidRequestFlagAdBlock BidRequestFlags = 1 << iota
	// BidRequestFlagPrivateBrowsing indicates if private browsing is enabled
	BidRequestFlagPrivateBrowsing
	// BidRequestFlagSecure indicates if the request is secure
	BidRequestFlagSecure
	// BidRequestFlagBot indicates if the request is from a bot
	BidRequestFlagBot
	// BidRequestFlagProxy indicates if the request is from a proxy
	BidRequestFlagProxy
)

// IsSet checks if a specific flag is set in the BidRequestFlags.
func (f BidRequestFlags) IsSet(flag BidRequestFlags) bool {
	return f&flag == flag
}

// Set sets or unsets a specific flag in the BidRequestFlags.
func (f *BidRequestFlags) Set(flag BidRequestFlags, set bool) {
	if set {
		*f |= flag
	} else {
		*f &^= flag
	}
}

// If returns the flag if the condition is true, otherwise returns 0.
func (f BidRequestFlags) If(b bool) BidRequestFlags {
	if b {
		return f
	}
	return 0
}
