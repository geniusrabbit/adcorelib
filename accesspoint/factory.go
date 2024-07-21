package accesspoint

import (
	"context"

	"github.com/valyala/fasthttp"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/platform/info"
)

// Factory of the plaform object
type Factory interface {
	New(ctx context.Context, accessPoint *admodels.RTBAccessPoint, opts ...any) (Platformer, error)
	Info() info.Platform
}

// Platformer interface of the platform executor
type Platformer interface {
	// ID of the access point
	ID() uint64

	// Codename of the access point
	Codename() string

	// Protocol name of the access point
	Protocol() string

	// HTTPHandler of the raw HTTP request
	HTTPHandler(ctx context.Context, rctx *fasthttp.RequestCtx) adtype.Responser
}
