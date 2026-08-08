package accessors

import (
	"context"
	"time"

	"github.com/geniusrabbit/adcorelib/adtype"
)

type MainAccessor struct {
	mainSource adtype.Source
}

func NewMainAccessor(mainSource adtype.Source) adtype.SourceAccessor {
	return &MainAccessor{mainSource: mainSource}
}

func (a *MainAccessor) Iterator(request adtype.BidRequester) adtype.SourceIterator {
	return func(yield func(float32, adtype.Source) bool) {
		if a.mainSource != nil {
			if !yield(1, a.mainSource) {
				return
			}
		}
	}
}

func (a *MainAccessor) SourceByID(ctx context.Context, id uint64) (adtype.Source, error) {
	if a.mainSource != nil && a.mainSource.ID() == id {
		return a.mainSource, nil
	}
	return nil, nil
}

func (a *MainAccessor) SetTimeout(ctx context.Context, timeout time.Duration) {
	if srcSetTM, _ := a.mainSource.(adtype.SourceTimeoutSetter); srcSetTM != nil {
		srcSetTM.SetTimeout(timeout)
	}
}
