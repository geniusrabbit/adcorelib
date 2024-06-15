package httpserver

import "geniusrabbit.dev/adcorelib/adtype"

// Source of advertisement
type Source interface {
	// Bid request for standart system filter
	Bid(request *adtype.BidRequest) adtype.Responser

	// ProcessResponse when need to fix the result and process all counters
	ProcessResponse(response adtype.Responser)
}
