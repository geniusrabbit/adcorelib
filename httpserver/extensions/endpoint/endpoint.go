package endpoint

import (
	"github.com/geniusrabbit/adcorelib/adtype"
)

// Source of advertisement
type Source interface {
	// Bid request for standart system filter
	Bid(request *adtype.BidRequest) adtype.Responser

	// ProcessResponse when need to fix the result and process all counters
	ProcessResponse(response adtype.Responser)
}

// Endpoint implementation
type Endpoint interface {
	// Codename of the enpoint
	Codename() string

	// Handle request and process response
	Handle(source Source, bidRequest *adtype.BidRequest) adtype.Responser
}
