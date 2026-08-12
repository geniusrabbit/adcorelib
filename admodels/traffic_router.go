package admodels

import "github.com/geniusrabbit/adcorelib/admodels/types"

// TrafficRouter represents a traffic router in the system.
type TrafficRouter struct {
	ID uint64

	// IDs of RTB sources to which the traffic should be routed
	RTBSourceIDs []uint64

	// Percentage of traffic to be routed
	Percent float32

	// Filter criteria for routing traffic
	Filter types.BaseFilter
}

// Test checks if the target matches the filter criteria.
func (d *TrafficRouter) Test(target types.TargetPointer) error {
	return d.Filter.Test(target)
}

// ObjectKey of the router
func (d *TrafficRouter) ObjectKey() uint64 {
	return d.ID
}
