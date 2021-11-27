package adtypeaccessor

import (
	"time"

	"geniusrabbit.dev/corelib/adtype"
)

// SourceReloaderFnk type
type SourceReloaderFnk func() ([]adtype.Source, error)

// Accessor object ad reloader
type Accessor struct {
	// reloader of objects
	reloader SourceReloaderFnk

	// List of sources
	sourceList []adtype.Source
}

// MustNewAccessor object
func MustNewAccessor(reloader SourceReloaderFnk) *Accessor {
	if reloader == nil {
		panic("reloader function is required")
	}
	accessor := &Accessor{reloader: reloader}
	accessor.Reload()
	return accessor
}

// Reload sources
func (accessor *Accessor) Reload() error {
	sources, err := accessor.reloader()
	if sources != nil && err == nil {
		accessor.sourceList = sources
	}
	return err
}

// Iterator returns the configured queue accessor
func (accessor *Accessor) Iterator(request *adtype.BidRequest) adtype.SourceIterator {
	return NewPriorityIterator(request, accessor.sourceList)
}

// SourceByID returns source instance
func (accessor *Accessor) SourceByID(id uint64) adtype.Source {
	for _, s := range accessor.sourceList {
		if s.ID() == id {
			return s
		}
	}
	return nil
}

// SetTimeout for sourcer
func (accessor *Accessor) SetTimeout(timeout time.Duration) {
	for _, src := range accessor.sourceList {
		if srcSetTM, _ := src.(adtype.SourceTimeoutSetter); srcSetTM != nil {
			srcSetTM.SetTimeout(timeout)
		}
	}
}

var _ adtype.SourceAccessor = &Accessor{}
