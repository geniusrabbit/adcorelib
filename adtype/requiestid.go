package adtype

import (
	"encoding/hex"

	"github.com/google/uuid"
)

// NewRequestID generates new unique request ID
//
//go:inline
func NewRequestID() string {
	return GenSomeID()
}

// NewImpressionID generates new unique impression ID
//
//go:inline
func NewImpressionID() string {
	return GenSomeID()
}

// NewAdResponseItemID generates new unique ad response item ID
//
//go:inline
func NewAdResponseItemID() string {
	return GenSomeID()
}

func GenSomeID() string {
	newUUID := uuid.New()
	return hex.EncodeToString(newUUID[:])
}
