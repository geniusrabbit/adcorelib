package accessors

import (
	"time"

	"github.com/geniusrabbit/adcorelib/adtype"
)

type AccessorWithMain struct {
	mainSource       adtype.Source
	mainSourceWeight float32
	other            adtype.SourceAccessor
}

func NewAccessorWithMainSource(main adtype.Source, mainSourceWeight float32, other adtype.SourceAccessor) adtype.SourceAccessor {
	if main == nil {
		return other
	}
	return &AccessorWithMain{
		mainSource:       main,
		mainSourceWeight: mainSourceWeight,
		other:            other,
	}
}

func (a *AccessorWithMain) Iterator(request *adtype.BidRequest) adtype.SourceIterator {
	return func(yield func(float32, adtype.Source) bool) {
		if a.mainSource != nil {
			if !yield(a.mainSourceWeight, a.mainSource) {
				return
			}
		}

		for w, src := range a.other.Iterator(request) {
			if src == nil {
				break
			}
			if !yield(w, src) {
				return
			}
		}
	}
}

func (a *AccessorWithMain) SourceByID(id uint64) (adtype.Source, error) {
	if a.mainSource != nil && a.mainSource.ID() == id {
		return a.mainSource, nil
	}
	return a.other.SourceByID(id)
}

func (a *AccessorWithMain) SetTimeout(timeout time.Duration) {
	if a.mainSource != nil {
		if srcSetTM, _ := a.mainSource.(adtype.SourceTimeoutSetter); srcSetTM != nil {
			srcSetTM.SetTimeout(timeout)
		}
	}
	a.other.SetTimeout(timeout)
}
