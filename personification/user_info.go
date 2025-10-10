package personification

import (
	"github.com/geniusrabbit/udetect"
	"github.com/google/uuid"
)

type (
	User    = udetect.User
	Device  = udetect.Device
	Geo     = udetect.Geo
	Carrier = udetect.Carrier
)

var (
	GeoDefault    = udetect.GeoDefault
	DeviceDefault = udetect.DeviceDefault
)

// UserInfo value
type UserInfo struct {
	User   *User
	Device *Device
	Geo    *Geo
}

// UUID of the user
func (i *UserInfo) UUID() string {
	if i == nil || i.User == nil || isEmptyUUIDPtr(&i.User.UUID) {
		return ""
	}
	return i.User.UUID.String()
}

// SessionID of the user
func (i *UserInfo) SessionID() string {
	if i == nil || i.User == nil {
		return ""
	}
	return i.User.SessionID
}

// Fingerprint of the iser
func (i *UserInfo) Fingerprint() string {
	if i == nil || i.User == nil {
		return ""
	}
	return i.User.FingerPrintID
}

// Country info
func (i *UserInfo) Country() *Geo {
	if i == nil || i.Geo == nil {
		return &udetect.GeoDefault
	}
	return i.Geo
}

// Ages of the user
func (i *UserInfo) Ages() (from, to int) {
	if i == nil || i.User == nil {
		return 0, 0
	}
	return i.User.AgeStart, i.User.AgeEnd
}

// ETag of the user
func (i *UserInfo) ETag() string {
	if i == nil || i.User == nil {
		return ""
	}
	return i.User.ETag
}

// Keywords of the user
func (i *UserInfo) Keywords() string {
	if i == nil || i.User == nil {
		return ""
	}
	return i.User.Keywords
}

// MostPossibleSex of the user
func (i *UserInfo) MostPossibleSex() int {
	if i == nil || i.User == nil {
		return 0
	}
	return i.User.MostPossibleSex()
}

// DeviceInfo get method
func (i *UserInfo) DeviceInfo() *Device {
	if i == nil || i.Device == nil {
		return &DeviceDefault
	}
	return i.Device
}

// GeoInfo get method
func (i *UserInfo) GeoInfo() *Geo {
	if i == nil || i.Geo == nil {
		return &GeoDefault
	}
	return i.Geo
}

// GeoInfo get method
func (i *UserInfo) CarrierInfo() *Carrier {
	return i.GeoInfo().Carrier
}

func isEmptyUUIDPtr(uuid *uuid.UUID) bool {
	if uuid != nil {
		for i := range len(*uuid) {
			if (*uuid)[i] != 0 {
				return false
			}
		}
	}
	return true
}
