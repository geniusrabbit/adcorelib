package personification

import (
	"github.com/geniusrabbit/udetect"
)

// PersonType is a personification type
type PersonType struct {
	Request       *Request
	UserInfoValue UserInfo
}

// User info data
func (p *PersonType) UserInfo() *UserInfo {
	return &p.UserInfoValue
}

// IsInited person in database
func (p *PersonType) IsInited() bool { return false }

// Properties for domain
func (p *PersonType) Properties(name string) Properties { return nil }

// Predict what does he likes?
func (p *PersonType) Predict(req *PredictRequest) (*PredictResponse, error) {
	return nil, nil
}

// PredictPrice what minimal
func (p *PersonType) PredictPrice(req *PredictPriceRequest) (*PredictPriceResponse, error) {
	return nil, nil
}

var EmptyPerson = &PersonType{
	UserInfoValue: UserInfo{
		User:   &udetect.User{},
		Device: &udetect.DeviceDefault,
		Geo:    &udetect.GeoDefault,
	},
}
