package ctxperson

import (
	"context"
	"testing"

	"geniusrabbit.dev/corelib/personification"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	ctx := context.Background()
	ctx = WithPersonClient(ctx, &personification.DummyClient{})
	assert.NotNil(t, Get(ctx))
}
