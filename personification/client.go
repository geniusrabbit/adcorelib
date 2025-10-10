package personification

import (
	"context"

	"github.com/geniusrabbit/udetect"
)

type (
	Request  = udetect.Request
	Response = udetect.Response
)

// Client interface
type Client interface {
	Detect(ctx context.Context, req *Request) (*Response, error)
}
