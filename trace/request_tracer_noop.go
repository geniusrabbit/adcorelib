package trace

type RequestNoopTracer struct {
}

// MarshalJSON encodes a noop tracer as JSON null.
// A value receiver covers both RequestNoopTracer and *RequestNoopTracer,
// including a non-addressable value stored in an interface.
func (RequestNoopTracer) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}

func (t *RequestNoopTracer) Trace(msg string, opts ...RequestTracerOpt) {
}

func (t *RequestNoopTracer) StartSubTrace(key string) RequestTracer {
	return t
}

func (t *RequestNoopTracer) StartSubTraceFmt(format string, args ...any) RequestTracer {
	return t
}

func (t *RequestNoopTracer) IsNoopTracer() bool {
	return true
}

var Noop = &RequestNoopTracer{}
