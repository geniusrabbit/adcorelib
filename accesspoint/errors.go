package accesspoint

import "errors"

var (
	// ErrSkipBidRequest according to unternal reasons
	ErrSkipBidRequest = errors.New(`skip bid request`)
	// ErrNoBidRequest if no any selected ADS
	ErrNoBidRequest = errors.New(`nobid request`)
)
