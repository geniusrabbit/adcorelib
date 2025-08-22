package actiontracker

import (
	"context"

	"github.com/valyala/fasthttp"
)

type PriceExtractor func(ctx context.Context, rctx *fasthttp.RequestCtx) (float64, error)

func DefaultPriceExtractor(paramName string) PriceExtractor {
	return func(ctx context.Context, rctx *fasthttp.RequestCtx) (float64, error) {
		return rctx.QueryArgs().GetUfloat(paramName)
	}
}
