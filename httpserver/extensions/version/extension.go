package version

import (
	"context"
	"encoding/json"

	"github.com/fasthttp/router"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"

	"github.com/geniusrabbit/adcorelib/context/version"
)

// Version type alias
type Version = version.Version

// Add version to context
func WithContext(ctx context.Context, ver *Version) context.Context {
	return version.WithContext(ctx, ver)
}

// Extension for version printing
type Extension struct {
	Path string
}

func (e *Extension) InitRouter(ctx context.Context, router *router.Router, tracer opentracing.Tracer) {
	path := e.Path
	if path == "" {
		path = "/version"
	}
	router.GET(path, func(rctx *fasthttp.RequestCtx) {
		ver := version.Get(ctx)
		rctx.SetContentType("application/json")
		rctx.SetStatusCode(fasthttp.StatusOK)
		_ = json.NewEncoder(rctx).Encode(map[string]string{
			"version": ver.Public(),
			"date":    ver.Date,
		})
	})
}
