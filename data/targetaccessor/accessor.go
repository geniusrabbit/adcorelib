package targetaccessor

import "geniusrabbit.dev/corelib/admodels"

// Accessor interface
type Accessor interface {
	TargetByID(id uint64) admodels.Target
}
