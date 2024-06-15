package accesspointaccessor

import (
	"context"
	"sort"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"geniusrabbit.dev/adcorelib/accesspoint"
	"geniusrabbit.dev/adcorelib/admodels"
	"geniusrabbit.dev/adcorelib/context/ctxlogger"
	"geniusrabbit.dev/adcorelib/models"
	"geniusrabbit.dev/adcorelib/storage/accessors/companyaccessor"
	"geniusrabbit.dev/adcorelib/storage/loader"
)

var (
	errUnsupportedAccessPointProtocol = errors.New("unsupported DSP protocol")
	errUndefinedDSPPlatform           = errors.New("undefined DSP platform")
)

// Accessor object ad reloader
type Accessor struct {
	loader.DataAccessor
	mx sync.RWMutex

	mainContext context.Context

	companyAccessor *companyaccessor.CompanyAccessor

	factories       map[string]accesspoint.Factory
	factoryList     []accesspoint.Factory
	accesspointList []accesspoint.Platformer
	accesspointMap  map[string]accesspoint.Platformer
}

// NewAccessor object
func NewAccessor(
	ctx context.Context,
	dataAccessor loader.DataAccessor,
	companyAccessor *companyaccessor.CompanyAccessor,
	factoryList ...accesspoint.Factory,
) (*Accessor, error) {
	if dataAccessor == nil {
		return nil, errors.New("data accessor is required")
	}
	if companyAccessor == nil {
		return nil, errors.New("company accessor is required")
	}
	factories := make(map[string]accesspoint.Factory, len(factoryList))
	for _, fact := range factoryList {
		factories[fact.Info().Protocol] = fact
	}
	return &Accessor{
		mainContext:     ctx,
		DataAccessor:    dataAccessor,
		companyAccessor: companyAccessor,
		factories:       factories,
		factoryList:     factoryList,
		accesspointList: nil,
		accesspointMap:  nil,
	}, nil
}

// AccesspointList returns list of accesspoints
func (acc *Accessor) AccesspointList() ([]accesspoint.Platformer, error) {
	if acc.accesspointList == nil || acc.NeedUpdate() {
		acc.mx.Lock()
		defer acc.mx.Unlock()
		if acc.accesspointList == nil || acc.NeedUpdate() {
			accesspoints, err := acc.Data()
			if err != nil {
				return nil, err
			}
			list := make([]accesspoint.Platformer, 0, len(acc.accesspointList))
			for _, src := range accesspoints {
				company, err := acc.companyAccessor.CompanyByID(src.(*models.RTBAccessPoint).CompanyID)
				if err != nil {
					ctxlogger.Get(acc.mainContext).Error("get company by ID", zap.Error(err))
				} else {
					accessPointModel := admodels.RTBAccessPointFromModel(src.(*models.RTBAccessPoint), company)
					accessPoint, err := acc.newAccessPoint(acc.mainContext, accessPointModel)
					if err == nil {
						list = append(list, accessPoint)
					} else {
						ctxlogger.Get(acc.mainContext).Error("create DSP accessor",
							zap.Uint64("access_point_id", accessPointModel.ID),
							zap.String("access_point_protocol", accessPointModel.Protocol),
							zap.Error(err))
					}
				}
			}
			sort.Slice(list, func(i, j int) bool { return list[i].ID() < list[j].ID() })
			acc.accesspointList = list
			// Update access point mapping
			mapped := map[string]accesspoint.Platformer{}
			for _, plt := range list {
				mapped[plt.Codename()] = plt
			}
			acc.accesspointMap = mapped
		}
	}
	return acc.accesspointList, nil
}

// AccessPointByID returns source instance
func (acc *Accessor) AccessPointByID(id uint64) (accesspoint.Platformer, error) {
	list, err := acc.AccesspointList()
	if err != nil {
		return nil, err
	}
	idx := sort.Search(len(list), func(i int) bool { return list[i].ID() >= id })
	if idx >= 0 && idx < len(list) && list[idx].ID() == id {
		return list[idx], nil
	}
	return nil, nil
}

// ListFactories of platforms
func (acc *Accessor) ListFactories() []accesspoint.Factory {
	return acc.factoryList
}

// PlatformByProtocol returns platform by codename and protocol
func (acc *Accessor) PlatformByProtocol(protocol, codename string) (accesspoint.Platformer, error) {
	_, err := acc.AccesspointList()
	if err != nil {
		return nil, err
	}
	acc.mx.RLock()
	defer acc.mx.RUnlock()
	plt := acc.accesspointMap[codename]
	if plt == nil {
		return nil, errors.Wrap(errUndefinedDSPPlatform, codename)
	}
	if protocol != "" && plt.Protocol() != protocol {
		return nil, errors.Wrap(errUnsupportedAccessPointProtocol, protocol)
	}
	return plt, nil
}

func (acc *Accessor) newAccessPoint(ctx context.Context, accessPoint *admodels.RTBAccessPoint) (accesspoint.Platformer, error) {
	if acc.factories == nil {
		return nil, errors.Wrap(errUnsupportedAccessPointProtocol, accessPoint.Protocol)
	}
	fact := acc.factories[accessPoint.Protocol]
	if fact == nil {
		return nil, errors.Wrap(errUnsupportedAccessPointProtocol, accessPoint.Protocol)
	}
	return fact.New(ctx, accessPoint)
}

var _ accesspoint.DSPPlatformAccessor = &Accessor{}
