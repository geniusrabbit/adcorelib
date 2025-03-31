package admodels

import "github.com/geniusrabbit/adcorelib/admodels/types"

// TrafficRouter represents a traffic router in the system.
type TrafficRouter struct {
	ID uint64

	RTBSourceIDs []uint64

	// Percentage of traffic to be routed
	Percent float32

	Filter types.BaseFilter
}

// Test checks if the target matches the filter criteria.
func (d *TrafficRouter) Test(target types.TargetPointer) bool {
	return d.Filter.Test(target)
}

// ObjectKey of the router
func (d *TrafficRouter) ObjectKey() uint64 {
	return d.ID
}
