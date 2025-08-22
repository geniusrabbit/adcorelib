package actiontracker

import (
	"testing"

	"github.com/geniusrabbit/notificationcenter/v2/dummy"
	"github.com/stretchr/testify/assert"

	"github.com/geniusrabbit/adcorelib/eventtraking/eventgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
	"github.com/geniusrabbit/adcorelib/eventtraking/urlgenerator"
	"github.com/geniusrabbit/adcorelib/httpserver/wrappers/httphandler"
)

type (
	TestEvent    = eventgenerator.TestEvent
	TestLead     = eventgenerator.TestLead
	TestUserInfo = eventgenerator.TestUserInfo
)

func Test_Options(t *testing.T) {
	eventGenerator := eventgenerator.New(
		"test",
		func() *TestEvent { return &TestEvent{} },
		func() *TestUserInfo { return &TestUserInfo{} },
	)
	eventStream := eventstream.New(
		&dummy.Publisher{},
		&dummy.Publisher{},
		eventGenerator,
	)
	server := NewExtension(
		WithURLGenerator[*TestEvent](&urlgenerator.Generator[*TestEvent, *TestLead, *TestUserInfo]{}),
		WithHTTPHandlerWrapper[*TestEvent](&httphandler.HTTPHandlerWrapper{}),
		WithEventStream[*TestEvent](eventStream),
		WithEventAllocator(func() *TestEvent { return &TestEvent{} }),
	)
	assert.True(t, server.eventStream != nil, "invalid eventstream server initialisation")
	assert.True(t, server.handlerWrapper != nil, "invalid handlerWrapper initialisation")
	assert.True(t, server.urlGenerator != nil, "invalid URLGenerator initialisation")
}
