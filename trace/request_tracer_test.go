package trace

import (
	"context"
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestNoopTracerJSON(t *testing.T) {
	var _ json.Marshaler = RequestNoopTracer{}
	var _ json.Marshaler = (*RequestNoopTracer)(nil)

	cases := []struct {
		name string
		in   any
		want string
	}{
		{name: "value", in: RequestNoopTracer{}, want: "null"},
		{name: "pointer", in: &RequestNoopTracer{}, want: "null"},
		{name: "nil pointer", in: (*RequestNoopTracer)(nil), want: "null"},
		{name: "interface", in: RequestTracer(&RequestNoopTracer{}), want: "null"},
		{name: "nil interface value", in: RequestTracer((*RequestNoopTracer)(nil)), want: "null"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			raw, err := json.Marshal(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, string(raw))

			raw, err = jsonv2.Marshal(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, string(raw))
		})
	}

	type payload struct {
		Tracer RequestTracer `json:"tracer"`
	}
	raw, err := json.Marshal(payload{Tracer: &RequestNoopTracer{}})
	require.NoError(t, err)
	assert.JSONEq(t, `{"tracer":null}`, string(raw))

	raw, err = json.Marshal([]any{RequestNoopTracer{}, &RequestNoopTracer{}})
	require.NoError(t, err)
	assert.JSONEq(t, `[null,null]`, string(raw))
}

func TestRequestNoopTracerMethods(t *testing.T) {
	tracer := &RequestNoopTracer{}
	tracer.Trace("ignored", Value("n", 1))
	assert.Same(t, tracer, tracer.StartSubTrace("child"))
	assert.Same(t, tracer, tracer.StartSubTraceFmt("child %d", 1))

	var nilTracer *RequestNoopTracer
	assert.NotPanics(t, func() {
		nilTracer.Trace("k", Message("%s", "v"))
	})
	assert.Nil(t, nilTracer.StartSubTrace("child"))
	assert.Nil(t, nilTracer.StartSubTraceFmt("child %d", 1))
}

func TestRequestTracerEntry(t *testing.T) {
	root := &RequestTracerEntry{Name: "root"}
	root.Trace("hello one")
	root.Trace("", Value("n", 2))
	root.Trace("failed", Error(errors.New("boom")))

	require.Equal(t, []any{
		map[string]any{"message": "hello one"},
		map[string]any{"n": 2},
		map[string]any{"message": "failed", "error": "boom"},
	}, root.trace)

	child := root.StartSubTrace("child").(*RequestTracerEntry)
	child.Trace("", Value("n", 7))
	again := root.StartSubTraceFmt("child %d", 2).(*RequestTracerEntry)
	assert.Equal(t, "child 2", again.Name)
	assert.Equal(t, "child", child.Name)
	assert.NotSame(t, child, again)
	assert.Equal(t, []any{map[string]any{"n": 7}}, child.trace)
	assert.Same(t, child, root.trace[3])
	assert.Same(t, again, root.trace[4])
	assert.Empty(t, again.trace)

	assertRequestTracerEntryJSON(t, root, `{
		"name":"root",
		"traces":[
			{"message":"hello one"},
			{"n":2},
			{"message":"failed","error":"boom"},
			{"name":"child","traces":[{"n":7}]},
			{"name":"child 2","traces":null}
		]
	}`)
}

func TestRequestTracerContext(t *testing.T) {
	ctx := context.Background()
	assert.Same(t, Noop, RequestTracerFromContext(ctx))

	ctx = context.WithValue(ctx, ctxTracerKey, "not a tracer")
	assert.Same(t, Noop, RequestTracerFromContext(ctx))

	root := &RequestTracerEntry{}
	ctx = ContextWithRequestTracer(context.Background(), root)
	assert.Same(t, root, RequestTracerFromContext(ctx))
	assertRequestTracerEntryJSON(t, root, `null`)

	child, childCtx := RequestTracerBegin(ctx, "span")
	assert.Same(t, child, RequestTracerFromContext(childCtx))
	assert.Same(t, child, root.trace[0])
	assert.Same(t, root, RequestTracerFromContext(ctx))
	assertRequestTracerEntryJSON(t, root, `{"traces":[{"name":"span","traces":null}]}`)
	assertRequestTracerEntryJSON(t, RequestTracer(root), `{"traces":[{"name":"span","traces":null}]}`)

	noop, noopCtx := RequestTracerBegin(context.Background(), "span")
	assert.Same(t, Noop, noop)
	assert.Same(t, Noop, RequestTracerFromContext(noopCtx))
}

func assertRequestTracerEntryJSON(t *testing.T, in any, want string) {
	t.Helper()
	raw, err := json.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, want, string(raw))

	raw, err = jsonv2.Marshal(in)
	require.NoError(t, err)
	assert.JSONEq(t, want, string(raw))
}
