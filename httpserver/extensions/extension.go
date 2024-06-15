package extensions

import (
	"context"

	"github.com/fasthttp/router"
	"github.com/opentracing/opentracing-go"
)

// ServerExtension provides abstraction server extension
type ServerExtension interface {
	InitRouter(ctx context.Context, router *router.Router, tracer opentracing.Tracer)
}
