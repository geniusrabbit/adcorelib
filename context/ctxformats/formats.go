package ctxformats

import (
	"context"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

type ctxKey struct{}

// WithContext returns new context with format accessor.
func WithContext(ctx context.Context, accessor types.FormatsAccessor) context.Context {
	return context.WithValue(ctx, ctxKey{}, accessor)
}

// FromContext returns format accessor from context.
func FromContext(ctx context.Context) types.FormatsAccessor {
	if v := ctx.Value(ctxKey{}); v != nil {
		return v.(types.FormatsAccessor)
	}
	return nil
}
