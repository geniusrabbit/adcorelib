package trakeraction

import (
	"testing"

	"github.com/geniusrabbit/notificationcenter/v2/dummy"
	"github.com/stretchr/testify/assert"

	"github.com/geniusrabbit/adcorelib/eventtraking/eventgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httphandler"
	"github.com/geniusrabbit/adcorelib/urlgenerator"
)

func Test_Options(t *testing.T) {
	eventGenerator := eventgenerator.New("test")
	eventStream := eventstream.New(
		&dummy.Publisher{},
		&dummy.Publisher{},
		eventGenerator,
	)
	server := NewExtension(
		WithURLGenerator(&urlgenerator.Generator{}),
		WithHTTPHandlerWrapper(&httphandler.HTTPHandlerWrapper{}),
		WithEventStream(eventStream),
	)
	assert.True(t, server.eventStream != nil, "invalid eventstream server initialisation")
	assert.True(t, server.handlerWrapper != nil, "invalid handlerWrapper initialisation")
	assert.True(t, server.urlGenerator != nil, "invalid URLGenerator initialisation")
}
