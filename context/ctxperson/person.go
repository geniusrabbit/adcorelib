package ctxperson

import (
	"context"

	"github.com/geniusrabbit/adcorelib/personification"
	"github.com/geniusrabbit/adcorelib/personification/dummy"
)

var (
	// CtxPersonObject reference to the person accessor
	CtxPersonObject = struct{ s string }{"person"}
	dummyClient     dummy.DummyClient
)

// Get logger object
func Get(ctx context.Context) personification.Client {
	if cli := ctx.Value(CtxPersonObject); cli != nil {
		return cli.(personification.Client)
	}
	return &dummyClient
}

// WithPersonClient puts person to context
func WithPersonClient(ctx context.Context, person personification.Client) context.Context {
	return context.WithValue(ctx, CtxPersonObject, person)
}
