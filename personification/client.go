package personification

import (
	"context"

	"github.com/geniusrabbit/udetect"
	"github.com/google/uuid"
)

// Client interface
type Client interface {
	Detect(ctx context.Context, req *udetect.Request) (*udetect.Response, error)
}

// Connect to the udetect server
func Connect(tr udetect.Transport) Client {
	return udetect.NewClient(tr)
}

type DummyClient struct {
}

func (DummyClient) Detect(ctx context.Context, req *udetect.Request) (*udetect.Response, error) {
	if isEmpty(req.UID) {
		req.UID = uuid.New()
	}
	if isEmpty(req.SessID) {
		req.SessID = uuid.New()
	}
	return &udetect.Response{
		User: &udetect.User{
			UUID:      req.UID,
			SessionID: req.SessID.String(),
		},
		Device: &udetect.Device{
			Browser: &udetect.Browser{
				DNT:             int(req.DNT),
				LMT:             int(req.LMT),
				Adblock:         int(req.Adblock),
				Languages:       req.Languages,
				PrimaryLanguage: req.PrimaryLanguage,
				UA:              req.UA,
				Ref:             req.Ref,
				JS:              int(req.JS),
				Width:           req.Width,
				Height:          req.Height,
				FlashVer:        req.FlashVer,
			},
		},
	}, nil
}

func isEmpty(uid uuid.UUID) bool {
	for i := 0; i < len(uid); i++ {
		if uid[i] != 0 {
			return false
		}
	}
	return true
}
