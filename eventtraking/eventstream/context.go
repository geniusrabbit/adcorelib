package eventstream

import (
	"context"
)

var (
	CtxStreamObject = struct{ s string }{"eventstream"}
	CtxWinsObject   = struct{ s string }{"eventstream.wins"}
)

// StreamFromContext object
func StreamFromContext(ctx context.Context) Stream {
	return ctx.Value(CtxStreamObject).(Stream)
}

// WithStream puts stream to context
func WithStream(ctx context.Context, stream Stream) context.Context {
	return context.WithValue(ctx, CtxStreamObject, stream)
}

// WinsFromContext object
func WinsFromContext(ctx context.Context) *WinNotifier {
	return ctx.Value(CtxWinsObject).(*WinNotifier)
}

// WithWins puts stream to context
func WithWins(ctx context.Context, wins *WinNotifier) context.Context {
	return context.WithValue(ctx, CtxWinsObject, wins)
}
