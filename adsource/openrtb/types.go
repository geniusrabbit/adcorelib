package openrtb

import (
	"errors"

	"geniusrabbit.dev/adcorelib/admodels"
)

// Request type enum
const (
	RequestTypeJSON      = admodels.RTBRequestTypeJSON
	RequestTypeXML       = admodels.RTBRequestTypeXML
	RequestTypeProtobuff = admodels.RTBRequestTypeProtoBUFF
)

// Errors set
var (
	ErrResponseAreNotSecure  = errors.New("response are not secure")
	ErrInvalidResponseStatus = errors.New("invalid response status")
	ErrNoCampaignsStatus     = errors.New("no campaigns response")
)
