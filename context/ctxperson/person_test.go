package ctxperson

import (
	"context"
	"testing"

	"github.com/geniusrabbit/adcorelib/personification/dummy"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	ctx := context.Background()
	ctx = WithPersonClient(ctx, &dummy.DummyClient{})
	assert.NotNil(t, Get(ctx))
}
