package adtype

import "github.com/geniusrabbit/adcorelib/rand"

// NewRequestID generates new unique request ID
//
//go:inline
func NewRequestID() string {
	return rand.UUID()
}

// NewImpressionID generates new unique impression ID
//
//go:inline
func NewImpressionID() string {
	return rand.UUID()
}

// NewAdResponseItemID generates new unique ad response item ID
//
//go:inline
func NewAdResponseItemID() string {
	return rand.UUID()
}
