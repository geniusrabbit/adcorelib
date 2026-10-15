package trace

import (
	"encoding/json"
	"fmt"
)

type RequestTracerOpt func() (key string, val any)

func Error(err error) RequestTracerOpt {
	return func() (string, any) {
		return "error", err.Error()
	}
}

func Message(format string, args ...any) RequestTracerOpt {
	return func() (string, any) {
		return "message", fmt.Sprintf(format, args...)
	}
}

func Value(key string, val any) RequestTracerOpt {
	return func() (string, any) {
		return key, val
	}
}

func ValueFnk(key string, fn func() any) RequestTracerOpt {
	return func() (string, any) {
		return key, fn()
	}
}

type RequestTracer interface {
	Trace(msg string, opts ...RequestTracerOpt)
	StartSubTrace(key string) RequestTracer
	StartSubTraceFmt(format string, args ...any) RequestTracer
	_isNoopTracer() bool
}

type RequestTracerEntry struct {
	Name  string
	trace []any
}

func NewRequestTracer(name string) RequestTracer {
	return &RequestTracerEntry{Name: name}
}

func (t *RequestTracerEntry) Trace(msg string, opts ...RequestTracerOpt) {
	fields := make(map[string]any, len(opts)+1)
	if msg != "" {
		fields["message"] = msg
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		key, val := opt()
		fields[key] = val
	}
	if len(fields) == 0 {
		return
	}
	t.trace = append(t.trace, fields)
}

func (t *RequestTracerEntry) StartSubTrace(key string) RequestTracer {
	subTrace := &RequestTracerEntry{Name: key}
	t.trace = append(t.trace, subTrace)
	return subTrace
}

func (t *RequestTracerEntry) StartSubTraceFmt(format string, args ...any) RequestTracer {
	return t.StartSubTrace(fmt.Sprintf(format, args...))
}

func (t *RequestTracerEntry) _isNoopTracer() bool {
	return false
}

// MarshalJSON encodes the log as {"name":"...","traces":[...]}.
// A nil receiver and an unnamed empty log encode as null. Errors are written as their Error() string.
func (t *RequestTracerEntry) MarshalJSON() ([]byte, error) {
	if t == nil || (t.Name == "" && len(t.trace) == 0) {
		return []byte("null"), nil
	}
	var traces any
	if len(t.trace) > 0 {
		traces = traceValues(t.trace)
	}
	return json.Marshal(struct {
		Name   string `json:"name,omitempty"`
		Traces any    `json:"traces"`
	}{
		Name:   t.Name,
		Traces: traces,
	})
}

func traceValues(values []any) []any {
	if len(values) == 0 {
		return nil
	}
	out := make([]any, len(values))
	for i, value := range values {
		out[i] = jsonTraceValue(value)
	}
	return out
}

func jsonTraceValue(value any) any {
	if err, ok := value.(error); ok && err != nil {
		return err.Error()
	}
	return value
}
