package endpoint

import (
	"strconv"
	"strings"

	"github.com/demdxx/gocast/v2"
	"github.com/demdxx/xtypes"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/valyala/fasthttp"
)

// RequestOptions prepare
type RequestOptions struct {
	Debug       bool
	Request     *fasthttp.RequestCtx
	Count       int
	X, Y        int
	Width       int
	WidthMax    int
	Height      int
	HeightMax   int
	Page        string
	Keywords    string
	FormatCodes []string
	FormatTypes []string
	SubID1      string
	SubID2      string
	SubID3      string
	SubID4      string
	SubID5      string
}

// NewRequestOptions prepare
func NewRequestOptions(ctx *fasthttp.RequestCtx) *RequestOptions {
	var (
		queryArgs        = ctx.QueryArgs()
		w, h, minW, minH = getSizeByCtx(ctx)
		debug, _         = strconv.ParseBool(string(queryArgs.Peek("debug")))
		formats          = strings.Trim(string(queryArgs.Peek("adformat")), ", \t\n\r")
		formatCodes      []string
		formatTypes      = strings.Trim(string(queryArgs.Peek("type")), ", \t\n\r")
		formatTypeCodes  []string
	)
	if formats != "" && formats != "auto" && formats != "all" {
		formatCodes = xtypes.SliceApply(
			strings.Split(formats, ","),
			strings.TrimSpace,
		).Filter(isNotEmptyString)
	}
	if formatTypes != "" && formatTypes != "auto" && formatTypes != "all" {
		formatTypeCodes = xtypes.SliceApply(
			strings.Split(formatTypes, ","),
			strings.TrimSpace,
		).Filter(isNotEmptyString)
	}
	return &RequestOptions{
		Debug:       debug,
		Request:     ctx,
		X:           gocast.Int(string(queryArgs.Peek("x"))),
		Y:           gocast.Int(string(queryArgs.Peek("y"))),
		Width:       minW,
		WidthMax:    ifPositiveNumber(w, -1),
		Height:      minH,
		HeightMax:   ifPositiveNumber(h, -1),
		FormatCodes: formatCodes,
		FormatTypes: formatTypeCodes,
		Keywords:    strings.Trim(peekOneFromQuery(queryArgs, "keywords", "keyword", "kw"), ", \t\n\r"),
		SubID1:      peekOneFromQuery(queryArgs, "subid1", "subid", "s1"),
		SubID2:      peekOneFromQuery(queryArgs, "subid2", "s2"),
		SubID3:      peekOneFromQuery(queryArgs, "subid3", "s3"),
		SubID4:      peekOneFromQuery(queryArgs, "subid4", "s4"),
		SubID5:      peekOneFromQuery(queryArgs, "subid5", "s5"),
		Count:       gocast.Int(peekOneFromQuery(queryArgs, "count")),
	}
}

// NewDirectRequestOptions prepare
func NewDirectRequestOptions(ctx *fasthttp.RequestCtx) *RequestOptions {
	var (
		queryArgs = ctx.QueryArgs()
		debug, _  = strconv.ParseBool(string(queryArgs.Peek("debug")))
	)
	return &RequestOptions{
		Debug:       debug,
		Request:     ctx,
		Count:       0,
		X:           gocast.Int(string(queryArgs.Peek("x"))),
		Y:           gocast.Int(string(queryArgs.Peek("y"))),
		Width:       -1,
		Height:      -1,
		FormatCodes: nil,
		FormatTypes: nil,
		SubID1:      peekOneFromQuery(queryArgs, "subid1", "subid", "s1"),
		SubID2:      peekOneFromQuery(queryArgs, "subid2", "s2"),
		SubID3:      peekOneFromQuery(queryArgs, "subid3", "s3"),
		SubID4:      peekOneFromQuery(queryArgs, "subid4", "s4"),
		SubID5:      peekOneFromQuery(queryArgs, "subid5", "s5"),
	}
}

// GetFormatTypes prepare by request codes
func (opt *RequestOptions) GetFormatTypes() types.FormatTypeBitset {
	if opt == nil {
		return types.FormatTypeBitsetEmpty
	}
	if opt.Width == -1 && opt.Height == -1 {
		return types.FormatTypeBitsetDirect
	}
	formatTypes := types.FormatTypeBitsetEmpty
	for _, t := range opt.FormatTypes {
		formatTypes.SetOne(types.FormatTypeByName(t))
	}
	return formatTypes
}
