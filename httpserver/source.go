package httpserver

import "github.com/geniusrabbit/adcorelib/adtype"

// Source of advertisement
type Source interface {
	// Bid request for standart system filter
	Bid(request adtype.BidRequester) adtype.Response

	// ProcessResponse when need to fix the result and process all counters
	ProcessResponse(response adtype.Response)
}
