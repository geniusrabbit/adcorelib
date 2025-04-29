package endpoint

import (
	"net/url"
	"strings"

	"github.com/demdxx/gocast/v2"
	"github.com/valyala/fasthttp"

	"github.com/geniusrabbit/adcorelib/admodels/types"
)

func peekOneFromQuery(query *fasthttp.Args, keys ...string) string {
	for _, key := range keys {
		if v := query.Peek(key); len(v) > 0 {
			return string(v)
		}
	}
	return ""
}

func directTypeMask(is bool) types.FormatTypeBitset {
	if is {
		return *types.NewFormatTypeBitset(types.FormatDirectType)
	}
	return types.FormatTypeBitsetEmpty
}

func domain(surl string) (name string) {
	if len(surl) < 1 {
		name = ""
	} else if len(surl) < 7 {
		name = strings.Split(surl, ",")[0]
	} else {
		switch strings.ToLower(surl[:7]) {
		case "http://", "https:/":
			if u, err := url.Parse(surl); nil == err {
				name = u.Host
			}
		}
	}
	return name
}

func sexFrom(v int) string {
	switch v {
	case 1:
		return "M"
	case 2:
		return "F"
	}
	return "?"
}

func getSizeByCtx(ctx *fasthttp.RequestCtx) (sw, sh, minSW, minSH int) {
	var (
		queryArgs = ctx.QueryArgs()
		w         = string(queryArgs.Peek("w"))
		h         = string(queryArgs.Peek("h"))
		minW      = string(queryArgs.Peek("mw"))
		minH      = string(queryArgs.Peek("mh"))
	)

	if isEmptyNumString(w) && isEmptyNumString(h) {
		if s := strings.Split(string(queryArgs.Peek("fmt")), "x"); len(s) > 0 {
			if len(s) == 2 {
				w, h = s[0], s[1]
			} else {
				w = s[0]
			}
		}
	}

	sw, sh, minSW, minSH = gocast.Int(w), gocast.Int(h),
		gocast.Int(minW), gocast.Int(minH)

	if sw < minSW {
		sw, minSW = minSW, sw
	}
	if sh < minSH {
		sh, minSH = minSH, sh
	}
	return sw, sh, minSW, minSH
}

func minInt(v1, v2 int) int {
	if v1 > 0 {
		return v1
	}
	return v2
}

func ifPositiveNumber(v1, v2 int) int {
	if v1 > 0 {
		return v1
	}
	return v2
}

func isEmptyNumString(s string) bool {
	return s == "" || s == "0"
}
