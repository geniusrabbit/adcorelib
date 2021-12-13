package admodels

import "strings"

// AdLink URL to the target
type AdLink struct {
	ID   uint64
	Link string // link to target
}

// IsInsecure link type
func (l *AdLink) IsInsecure() bool {
	return strings.HasPrefix(l.Link, "http://")
}
