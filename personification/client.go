package personification

import (
	"context"

	"github.com/geniusrabbit/udetect"
)

// Client interface
type Client interface {
	Detect(ctx context.Context, req *udetect.Request) (*udetect.Response, error)
}
