package dummy

import (
	"context"

	"github.com/geniusrabbit/udetect"
	"github.com/google/uuid"
)

type DummyClient struct {
}

func (DummyClient) Detect(ctx context.Context, req *udetect.Request) (*udetect.Response, error) {
	if isEmptyUUIDPtr(&req.UID) {
		req.UID = uuid.New()
	}
	if isEmptyUUIDPtr(&req.SessID) {
		req.SessID = uuid.New()
	}
	return &udetect.Response{
		User: &udetect.User{
			UUID:      req.UID,
			SessionID: req.SessID.String(),
		},
		Device: &udetect.Device{
			Browser: &udetect.Browser{
				DNT:             req.DNT,
				LMT:             req.LMT,
				AdBlock:         req.AdBlock,
				Languages:       req.Languages,
				PrimaryLanguage: req.PrimaryLanguage,
				UA:              req.UA,
				Ref:             req.Ref,
				JS:              req.JS,
				Width:           req.Width,
				Height:          req.Height,
				FlashVer:        req.FlashVer,
			},
		},
	}, nil
}
