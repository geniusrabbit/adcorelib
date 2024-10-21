package endpoint

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httphandler"
)

type dummySource struct{}

func (dummySource) Bid(request *adtype.BidRequest) adtype.Responser { return nil }
func (dummySource) ProcessResponse(response adtype.Responser)       {}

type dummyZoneAccessor struct{}

func (dummyZoneAccessor) TargetByID(uint64) (admodels.Target, error)       { return nil, nil }
func (dummyZoneAccessor) TargetByCodename(string) (admodels.Target, error) { return nil, nil }

func Test_Options(t *testing.T) {
	server := NewExtension(
		WithHTTPHandlerWrapper(&httphandler.HTTPHandlerWrapper{}),
		WithAdvertisementSource(dummySource{}),
		WithZoneAccessor(&dummyZoneAccessor{}),
	)
	assert.True(t, server.source != nil, "invalid SSP server initialisation")
	assert.True(t, server.handlerWrapper != nil, "invalid handlerWrapper initialisation")
	assert.True(t, server.zoneAccessor != nil, "invalid zoneAccessor initialisation")
}
