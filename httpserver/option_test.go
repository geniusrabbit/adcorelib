package httpserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Options(t *testing.T) {
	server, err := NewServer(
		WithServiceName("test"),
		WithDebugMode(true),
	)
	if assert.NoError(t, err) {
		assert.Equal(t, "test", server.serviceName, "invalid service name initialisation")
		assert.True(t, server.debug, "invalid debug mode setup")
		assert.True(t, server.logger != nil, "invalid logger initialisation")
	}
}
