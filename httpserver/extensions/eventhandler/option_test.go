package eventhandler

import (
	"testing"

	"github.com/geniusrabbit/notificationcenter/v2/dummy"
	"github.com/stretchr/testify/assert"

	"geniusrabbit.dev/adcorelib/eventtraking/eventgenerator"
	"geniusrabbit.dev/adcorelib/eventtraking/eventstream"
	"geniusrabbit.dev/adcorelib/httpserver/wrappers/httphandler"
)

func Test_Options(t *testing.T) {
	eventGenerator := eventgenerator.New("test")
	eventStream := eventstream.New(
		&dummy.Publisher{},
		&dummy.Publisher{},
		eventGenerator,
	)
	server := NewExtension(
		WithHTTPHandlerWrapper(&httphandler.HTTPHandlerWrapper{}),
		WithEventStream(eventStream),
	)
	assert.True(t, server.eventStream != nil, "invalid eventstream server initialisation")
	assert.True(t, server.handlerWrapper != nil, "invalid handlerWrapper initialisation")
}
