package ctxperson

import (
	"context"

	"geniusrabbit.dev/corelib/personification"
)

var (
	// CtxPersonObject reference to the person accessor
	CtxPersonObject = struct{ s string }{"person"}
	dummyClient     personification.DummyClient
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
