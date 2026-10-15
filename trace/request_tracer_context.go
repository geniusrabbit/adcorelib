package trace

import "context"

var (
	ctxTracerKey = struct{ string }{"tracer"}
)

func ContextWithRequestTracer(ctx context.Context, tracer RequestTracer) context.Context {
	if tracer.IsNoopTracer() {
		return ctx
	}
	return context.WithValue(ctx, ctxTracerKey, tracer)
}

func RequestTracerFromContext(ctx context.Context) RequestTracer {
	tracer, ok := ctx.Value(ctxTracerKey).(RequestTracer)
	if !ok {
		return Noop
	}
	return tracer
}

func RequestTracerBegin(ctx context.Context, key string) (RequestTracer, context.Context) {
	tracer := RequestTracerFromContext(ctx)
	if tracer.IsNoopTracer() {
		return tracer, ctx
	}
	tracer = tracer.StartSubTrace(key)
	return tracer, ContextWithRequestTracer(ctx, tracer)
}
