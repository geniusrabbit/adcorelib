package accesspoint

import (
	"context"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/platform/info"
	"github.com/valyala/fasthttp"
)

// Factory of the plaform object
type Factory interface {
	New(ctx context.Context, accessPoint *admodels.RTBAccessPoint, opts ...any) (Platformer, error)
	Info() info.Platform
}

// Platformer interface of the platform executor
type Platformer interface {
	// ID of source
	ID() uint64

	// Codename of the platform
	Codename() string

	// Protocol name of the platform
	Protocol() string

	// HTTPHandler of the raw HTTP request
	HTTPHandler(ctx context.Context, rctx *fasthttp.RequestCtx) adtype.Responser
}
