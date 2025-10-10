package simple

import (
	"context"
	"strings"

	"github.com/google/uuid"
	useragent "github.com/mileusna/useragent"

	"github.com/geniusrabbit/udetect"
)

type Item struct {
	ID   uint
	Name string
}

// SimpleClient represents a simple implementation of the Client interface
// that uses the mileusna/useragent package for user-agent parsing.
type SimpleClient struct {
	BrowserList []*Item
	OSList      []*Item
}

func (s *SimpleClient) Detect(ctx context.Context, req *udetect.Request) (*udetect.Response, error) {
	if isEmpty(req.UID) {
		req.UID = uuid.New()
	}
	if isEmpty(req.SessID) {
		req.SessID = uuid.New()
	}
	ua := useragent.Parse(req.UA)
	return &udetect.Response{
		User: &udetect.User{
			UUID:      req.UID,
			SessionID: req.SessID.String(),
		},
		Device: &udetect.Device{
			DeviceType: deviceType(&ua),
			OS: &udetect.OS{
				ID:      s.osGet(ua.OS),
				Name:    ua.OS,
				Version: ua.OSVersion,
			},
			Browser: &udetect.Browser{
				ID:              uint64(s.browserGet(ua.Name)),
				Name:            ua.Name,
				Version:         ua.Version,
				DNT:             req.DNT,
				LMT:             req.LMT,
				AdBlock:         req.AdBlock,
				IsRobot:         b2i[int8](ua.Bot),
				Languages:       req.Languages,
				PrimaryLanguage: req.PrimaryLanguage,
				JS:              req.JS,
				UA:              req.UA,
				Ref:             req.Ref,
				Width:           req.Width,
				Height:          req.Height,
				FlashVer:        req.FlashVer,
			},
		},
	}, nil
}

func (s *SimpleClient) browserGet(name string) uint {
	for _, brw := range s.BrowserList {
		if strings.EqualFold(brw.Name, name) {
			return brw.ID
		}
	}
	return 0
}

func (s *SimpleClient) osGet(name string) uint {
	for _, os := range s.OSList {
		if strings.EqualFold(os.Name, name) {
			return os.ID
		}
	}
	return 0
}
