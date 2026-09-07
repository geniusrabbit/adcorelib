package leadtracker

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/geniusrabbit/notificationcenter/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"

	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventgenerator"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
)

type unpackLead struct {
	TestLead
	unpacked []byte
	err      error
}

func (l *unpackLead) Unpack(data []byte) error {
	l.unpacked = append([]byte(nil), data...)
	return l.err
}

func newTestExt(pub notificationcenter.Publisher, lead eventgenerator.Allocator[*unpackLead]) *Extension[*unpackLead] {
	gen := eventgenerator.New(
		"test",
		func() *TestEvent { return &TestEvent{} },
		func() *TestUserInfo { return &TestUserInfo{} },
	)
	stream := eventstream.New(pub, &dummyDiscard{}, gen)
	if lead == nil {
		lead = func() *unpackLead { return &unpackLead{err: errors.New("no packed lead")} }
	}
	return NewExtension(
		WithEventStream[*unpackLead](stream),
		WithLeadAllocator(lead),
	)
}

type dummyDiscard struct{}

func (dummyDiscard) Publish(context.Context, ...any) error { return nil }

func serveLead(t *testing.T, ext *Extension[*unpackLead], uri string, postForm string) *fasthttp.RequestCtx {
	t.Helper()
	rctx := &fasthttp.RequestCtx{}
	if postForm != "" {
		rctx.Request.Header.SetMethod(http.MethodPost)
		rctx.Request.Header.SetContentType("application/x-www-form-urlencoded")
		rctx.Request.SetBodyString(postForm)
		rctx.Request.SetRequestURI(uri)
	} else {
		rctx.Request.SetRequestURI(uri)
	}
	ext.eventLeadHandler("plain")(context.Background(), rctx)
	return rctx
}

func TestCompactLeadJSONTags(t *testing.T) {
	ev := CompactLeadEvent{
		Event:        "lead",
		AccountID:    42,
		AdvAccountID: 42,
		ImpressionID: "clk-hex",
		LeadTypeID:   1,
		LeadTypeCH:   1,
		Price:        int64(billing.MoneyFloat(10.0)),
		LeadPrice:    int64(billing.MoneyFloat(10.0)),
		Time:         123,
	}
	raw, err := json.Marshal(ev)
	require.NoError(t, err)
	var m map[string]any
	require.NoError(t, json.Unmarshal(raw, &m))
	assert.Equal(t, "lead", m["e"])
	assert.Equal(t, float64(42), m["acc"])
	assert.Equal(t, float64(42), m["acv"])
	assert.Equal(t, "clk-hex", m["imp"])
	assert.Equal(t, float64(1), m["lt"])
	assert.Equal(t, float64(1), m["ltid"])
	assert.Equal(t, float64(billing.MoneyFloat(10.0)), m["p"])
	assert.Equal(t, float64(billing.MoneyFloat(10.0)), m["lpr"])
	assert.Equal(t, float64(123), m["tm"])
}

func TestCompactPostbackQuery(t *testing.T) {
	var got []any
	pub := notificationcenter.FuncPublisher(func(_ context.Context, messages ...any) error {
		got = append(got, messages...)
		return nil
	})
	ext := newTestExt(pub, nil)
	rctx := serveLead(t, ext, "/lead?acc=42&clk=uud220&lt=1&p=10", "")
	assert.Equal(t, http.StatusOK, rctx.Response.StatusCode())
	require.Len(t, got, 1)
	ev, ok := got[0].(*CompactLeadEvent)
	require.True(t, ok)
	assert.Equal(t, "lead", ev.Event)
	assert.Equal(t, uint64(42), ev.AccountID)
	assert.Equal(t, "uud220", ev.ImpressionID)
	assert.Equal(t, uint32(1), ev.LeadTypeID)
	assert.Equal(t, int64(billing.MoneyFloat(10.0)), ev.Price)
	assert.Equal(t, ev.AccountID, ev.AdvAccountID)
	assert.Equal(t, ev.LeadTypeID, ev.LeadTypeCH)
	assert.Equal(t, ev.Price, ev.LeadPrice)
	assert.NotZero(t, ev.Time)
}

func TestCompactPostbackPOSTForm(t *testing.T) {
	var got []any
	pub := notificationcenter.FuncPublisher(func(_ context.Context, messages ...any) error {
		got = append(got, messages...)
		return nil
	})
	ext := newTestExt(pub, nil)
	rctx := serveLead(t, ext, "/lead", "acc=7&clk=abc&lt=0")
	assert.Equal(t, http.StatusOK, rctx.Response.StatusCode())
	require.Len(t, got, 1)
	ev := got[0].(*CompactLeadEvent)
	assert.Equal(t, uint64(7), ev.AccountID)
	assert.Equal(t, "abc", ev.ImpressionID)
	assert.Equal(t, uint32(0), ev.LeadTypeID)
	assert.Zero(t, ev.Price)
}

func TestCompactMissingAccOrClk(t *testing.T) {
	ext := newTestExt(&dummyDiscard{}, nil)
	assert.Equal(t, http.StatusBadRequest, serveLead(t, ext, "/lead?clk=x", "").Response.StatusCode())
	assert.Equal(t, http.StatusBadRequest, serveLead(t, ext, "/lead?acc=1", "").Response.StatusCode())
}

func TestCompactOmitsLtDefaultsZero(t *testing.T) {
	var got []any
	pub := notificationcenter.FuncPublisher(func(_ context.Context, messages ...any) error {
		got = append(got, messages...)
		return nil
	})
	ext := newTestExt(pub, nil)
	serveLead(t, ext, "/lead?acc=1&clk=c", "")
	require.Len(t, got, 1)
	assert.Equal(t, uint32(0), got[0].(*CompactLeadEvent).LeadTypeID)
}

func TestClkTakesPrecedenceOverPackedL(t *testing.T) {
	var got []any
	pub := notificationcenter.FuncPublisher(func(_ context.Context, messages ...any) error {
		got = append(got, messages...)
		return nil
	})
	unpacked := &unpackLead{err: errors.New("should not unpack")}
	ext := newTestExt(pub, func() *unpackLead { return unpacked })
	rctx := serveLead(t, ext, "/lead?acc=1&clk=imp1&l=packed", "")
	assert.Equal(t, http.StatusOK, rctx.Response.StatusCode())
	require.Len(t, got, 1)
	_, ok := got[0].(*CompactLeadEvent)
	assert.True(t, ok)
	assert.Nil(t, unpacked.unpacked)
}

func TestLegacyPackedLeadUnchanged(t *testing.T) {
	var got []any
	pub := notificationcenter.FuncPublisher(func(_ context.Context, messages ...any) error {
		got = append(got, messages...)
		return nil
	})
	lead := &unpackLead{}
	ext := newTestExt(pub, func() *unpackLead { return lead })
	rctx := serveLead(t, ext, "/lead?l=oktoken", "")
	assert.Equal(t, http.StatusOK, rctx.Response.StatusCode())
	assert.Equal(t, []byte("oktoken"), lead.unpacked)
	require.Len(t, got, 1)
	_, isCompact := got[0].(*CompactLeadEvent)
	assert.False(t, isCompact)
}

func TestLegacyUnpackError400(t *testing.T) {
	ext := newTestExt(&dummyDiscard{}, func() *unpackLead {
		return &unpackLead{err: errors.New("bad")}
	})
	assert.Equal(t, http.StatusBadRequest, serveLead(t, ext, "/lead?l=nope", "").Response.StatusCode())
}

func TestPublishErrorStill200(t *testing.T) {
	pub := notificationcenter.FuncPublisher(func(context.Context, ...any) error {
		return errors.New("bus down")
	})
	ext := newTestExt(pub, nil)
	rctx := serveLead(t, ext, "/lead?acc=1&clk=c", "")
	assert.Equal(t, http.StatusOK, rctx.Response.StatusCode())
}
