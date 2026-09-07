package leadtracker

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
	"github.com/geniusrabbit/adcorelib/fasttime"
)

// CompactLeadEvent is the S2S /lead JSON payload (HTTP acc/clk/lt/p).
//
// Eventstream already maps acv→adv_account_id and lpr→lead_price; lt is
// latitude on full events, so lead type is also emitted as ltid.
type CompactLeadEvent struct {
	Event        string `json:"e"`
	AccountID    uint64 `json:"acc"`
	AdvAccountID uint64 `json:"acv"`
	ImpressionID string `json:"imp"`
	LeadTypeID   uint32 `json:"lt"`
	LeadTypeCH   uint32 `json:"ltid"`
	Price        int64  `json:"p"`
	LeadPrice    int64  `json:"lpr"`
	Time         int64  `json:"tm"`
}

var errCompactRequired = errors.New("acc and clk are required")

func peekLeadArg(rctx *fasthttp.RequestCtx, key string) []byte {
	if v := rctx.QueryArgs().Peek(key); len(v) > 0 {
		return v
	}
	return rctx.PostArgs().Peek(key)
}

func parseCompactLead(rctx *fasthttp.RequestCtx) (*CompactLeadEvent, error) {
	accRaw := bytes.TrimSpace(peekLeadArg(rctx, "acc"))
	clkRaw := bytes.TrimSpace(peekLeadArg(rctx, "clk"))
	if len(accRaw) == 0 || len(clkRaw) == 0 {
		return nil, errCompactRequired
	}

	acc, err := strconv.ParseUint(string(accRaw), 10, 64)
	if err != nil {
		return nil, err
	}

	var lt uint32
	if ltRaw := bytes.TrimSpace(peekLeadArg(rctx, "lt")); len(ltRaw) > 0 {
		v, perr := strconv.ParseUint(string(ltRaw), 10, 32)
		if perr != nil {
			return nil, perr
		}
		lt = uint32(v)
	}

	var price billing.Money
	if pRaw := bytes.TrimSpace(peekLeadArg(rctx, "p")); len(pRaw) > 0 {
		f, perr := strconv.ParseFloat(string(pRaw), 64)
		if perr != nil {
			return nil, perr
		}
		price = billing.MoneyFloat(f)
	}

	clk := strings.TrimSpace(string(clkRaw))
	ev := &CompactLeadEvent{
		Event:        string(events.Lead),
		AccountID:    acc,
		AdvAccountID: acc,
		ImpressionID: clk,
		LeadTypeID:   lt,
		LeadTypeCH:   lt,
		Price:        int64(price),
		LeadPrice:    int64(price),
		Time:         int64(fasttime.UnixTimestampNano()),
	}
	return ev, nil
}
