package endpoint

import (
	"github.com/geniusrabbit/adcorelib/adtype"
)

// Source of advertisement
type Source interface {
	// Bid request for standart system filter
	Bid(request adtype.BidRequester) adtype.Response

	// ProcessResponse when need to fix the result and process all counters
	ProcessResponse(response adtype.Response)
}

// Endpoint implementation
type Endpoint interface {
	// Codename of the endpoint
	Codename() string

	// Handle request and process response
	Handle(source Source, bidRequest adtype.BidRequester) adtype.Response
}
