package simple

import (
	"context"
	"strings"

	"github.com/google/uuid"
	useragent "github.com/mileusna/useragent"

	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/udetect"
)

// VersionItem is one catalog version row. Min/Max come from match_ver_min_exp
// and match_ver_max_exp. There is no further nesting.
type VersionItem struct {
	ID  uint
	Min types.Version
	Max types.Version
}

// Item is a parentless OS or browser family. Versions are the child rows.
type Item struct {
	ID       uint
	Name     string
	Versions []*VersionItem
}

// SimpleClient represents a simple implementation of the Client interface
// that uses the mileusna/useragent package for user-agent parsing.
type SimpleClient struct {
	BrowserList []*Item
	OSList      []*Item
}

// New returns a SimpleClient with the given catalog lists.
func New(browsers, os []*Item) *SimpleClient {
	return &SimpleClient{
		BrowserList: browsers,
		OSList:      os,
	}
}

func (s *SimpleClient) Detect(ctx context.Context, req *udetect.Request) (*udetect.Response, error) {
	if isEmpty(req.UID) {
		req.UID = uuid.New()
	}
	if isEmpty(req.SessID) {
		req.SessID = uuid.New()
	}
	ua := useragent.Parse(req.UA)
	osID, osVerID := s.osGet(ua.OS, ua.OSVersion)
	brID, brVerID := s.browserGet(ua.Name, ua.Version)
	return &udetect.Response{
		User: &udetect.User{
			UUID:      req.UID,
			SessionID: req.SessID.String(),
		},
		Device: &udetect.Device{
			DeviceType: deviceType(&ua),
			OS: &udetect.OS{
				ID:        osID,
				VersionID: osVerID,
				Name:      ua.OS,
				Version:   ua.OSVersion,
			},
			Browser: &udetect.Browser{
				ID:              uint64(brID),
				VersionID:       uint64(brVerID),
				Name:            ua.Name,
				Version:         ua.Version,
				DNT:             req.DNT,
				LMT:             req.LMT,
				AdBlock:         req.AdBlock,
				IsRobot:         b2i[int8](ua.Bot || ua.Name == "curl"),
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

func (s *SimpleClient) browserGet(name, version string) (rootID, versionID uint) {
	return lookup(s.BrowserList, name, version)
}

func (s *SimpleClient) osGet(name, version string) (rootID, versionID uint) {
	return lookup(s.OSList, name, version)
}

func lookup(list []*Item, name, version string) (rootID, versionID uint) {
	for _, item := range list {
		if item == nil || !strings.EqualFold(item.Name, name) {
			continue
		}
		return item.ID, matchVersion(item.Versions, types.IgnoreParseVersion(version))
	}
	return 0, 0
}

// matchVersion picks the child whose range contains ua: Min <= ua < Max.
// An empty Max is unbounded. An empty Min is not a range. Overlaps keep the greater Min.
func matchVersion(items []*VersionItem, ua types.Version) uint {
	var best *VersionItem
	for _, item := range items {
		if item == nil || item.Min.IsEmpty() {
			continue
		}
		if ua.Less(item.Min) {
			continue
		}
		if !item.Max.IsEmpty() && !ua.Less(item.Max) {
			continue
		}
		if best == nil || best.Min.Less(item.Min) {
			best = item
		}
	}
	if best == nil {
		return 0
	}
	return best.ID
}
