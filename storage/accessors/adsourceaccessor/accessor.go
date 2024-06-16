package adsourceaccessor

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/demdxx/xtypes"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"geniusrabbit.dev/adcorelib/admodels"
	"geniusrabbit.dev/adcorelib/adtype"
	"geniusrabbit.dev/adcorelib/context/ctxlogger"
	"geniusrabbit.dev/adcorelib/models"
	"geniusrabbit.dev/adcorelib/platform/info"
	"geniusrabbit.dev/adcorelib/storage/accessors/companyaccessor"
	"geniusrabbit.dev/adcorelib/storage/loader"
)

var errUnsupportedSourceProtocol = errors.New("unsupported source protocol")

type sourceFactory interface {
	New(ctx context.Context, source *admodels.RTBSource, opts ...any) (adtype.SourceTester, error)
	Info() info.Platform
	Protocols() []string
}

// Accessor object ad reloader
type Accessor struct {
	mx sync.Mutex

	loader.DataAccessor
	companyAccessor *companyaccessor.CompanyAccessor

	mainContext context.Context

	factories  map[string]sourceFactory
	sourceList []adtype.Source
}

// NewAccessor object
func NewAccessor(ctx context.Context, dataAccessor loader.DataAccessor, companyAccessor *companyaccessor.CompanyAccessor, factories ...sourceFactory) (*Accessor, error) {
	if dataAccessor == nil {
		return nil, errors.New("data accessor is required")
	}
	if companyAccessor == nil {
		return nil, errors.New("company accessor is required")
	}
	mapFactory := map[string]sourceFactory{}
	for _, fact := range factories {
		for _, protoName := range fact.Protocols() {
			mapFactory[protoName] = fact
		}
	}
	return &Accessor{
		mainContext:     ctx,
		DataAccessor:    dataAccessor,
		companyAccessor: companyAccessor,
		factories:       mapFactory,
	}, nil
}

// SourceList returns list of sources
func (acc *Accessor) SourceList() ([]adtype.Source, error) {
	if acc.sourceList != nil && !acc.NeedUpdate() {
		return acc.sourceList, nil
	}

	acc.mx.Lock()
	defer acc.mx.Unlock()
	if !acc.NeedUpdate() {
		return acc.sourceList, nil
	}

	sources, err := acc.Data()
	if err != nil {
		return nil, err
	}

	acc.sourceList = xtypes.SliceApply(sources, func(src any) adtype.Source {
		company, err := acc.companyAccessor.CompanyByID(src.(*models.RTBSource).CompanyID)
		if err != nil {
			ctxlogger.Get(acc.mainContext).Error("get company by ID", zap.Error(err))
			return nil
		}
		rtbModel := admodels.RTBSourceFromModel(src.(*models.RTBSource), company)
		if src, err := acc.newSource(acc.mainContext, rtbModel); err == nil {
			return src
		} else {
			ctxlogger.Get(acc.mainContext).Error("create RTB source",
				zap.Uint64("source_id", rtbModel.ID),
				zap.String("source_protocol", rtbModel.Protocol),
				zap.Error(err))
		}
		return nil
	}).Sort(func(i, j adtype.Source) bool { return i.ID() < j.ID() })

	return acc.sourceList, nil
}

// Iterator returns the configured queue acc
func (acc *Accessor) Iterator(request *adtype.BidRequest) adtype.SourceIterator {
	list, _ := acc.SourceList()
	return NewPriorityIterator(request, list)
}

// SourceByID returns source instance
func (acc *Accessor) SourceByID(id uint64) (adtype.Source, error) {
	list, err := acc.SourceList()
	if err != nil {
		return nil, err
	}
	idx := sort.Search(len(list), func(i int) bool { return list[i].ID() >= id })
	if idx >= 0 && idx < len(list) && list[idx].ID() == id {
		return list[idx], nil
	}
	return nil, nil
}

func (acc *Accessor) newSource(ctx context.Context, src *admodels.RTBSource) (adtype.Source, error) {
	if acc.factories == nil {
		return nil, errors.Wrap(errUnsupportedSourceProtocol, src.Protocol)
	}
	fact := acc.factories[src.Protocol]
	if fact == nil {
		return nil, errors.Wrap(errUnsupportedSourceProtocol, src.Protocol)
	}
	return fact.New(ctx, src)
}

// SetTimeout for sourcer
func (acc *Accessor) SetTimeout(timeout time.Duration) {
	list, _ := acc.SourceList()
	for _, src := range list {
		if srcSetTM, _ := src.(adtype.SourceTimeoutSetter); srcSetTM != nil {
			srcSetTM.SetTimeout(timeout)
		}
	}
}

var _ adtype.SourceAccessor = &Accessor{}
