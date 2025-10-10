package simple

import (
	"strings"

	"github.com/geniusrabbit/udetect"
	"github.com/google/uuid"
	useragent "github.com/mileusna/useragent"
	"golang.org/x/exp/constraints"
)

func deviceType(ua *useragent.UserAgent) udetect.DeviceType {
	if ua.Mobile {
		return udetect.DeviceTypeMobile
	}
	if ua.Tablet {
		return udetect.DeviceTypeTablet
	}
	if ua.Desktop {
		return udetect.DeviceTypePC
	}
	if strings.Contains(ua.Name, "AppleTV") {
		return udetect.DeviceTypeTV
	}
	if strings.Contains(ua.String, "PlayStation") || strings.Contains(ua.String, "Xbox") {
		return udetect.DeviceTypeSetTopBox
	}
	return udetect.DeviceTypeUnknown
}

func isEmpty(uid uuid.UUID) bool {
	for i := range len(uid) {
		if uid[i] != 0 {
			return false
		}
	}
	return true
}

func b2i[R constraints.Integer | constraints.Float](b bool) R {
	if b {
		return 1
	}
	return 0
}
