package zoneaccessor

import (
	"sort"

	"geniusrabbit.dev/adcorelib/admodels"
	"geniusrabbit.dev/adcorelib/models"
	"geniusrabbit.dev/adcorelib/storage/accessors/companyaccessor"
	"geniusrabbit.dev/adcorelib/storage/loader"
)

// ZoneAccessor provides accessor to the admodel company type
type ZoneAccessor struct {
	loader.DataAccessor
	companyAccessor *companyaccessor.CompanyAccessor

	zoneList []admodels.Target
}

// NewZoneAccessor from dataAccessor
func NewZoneAccessor(dataAccessor loader.DataAccessor, companyAccessor *companyaccessor.CompanyAccessor) *ZoneAccessor {
	return &ZoneAccessor{DataAccessor: dataAccessor, companyAccessor: companyAccessor}
}

// ZoneList returns list of prepared data
func (acc *ZoneAccessor) ZoneList() ([]admodels.Target, error) {
	if acc.zoneList == nil || acc.NeedUpdate() {
		data, err := acc.Data()
		if err != nil {
			return nil, err
		}
		zoneList := make([]admodels.Target, 0, len(data))
		for _, it := range data {
			zone := *it.(*models.Zone)
			if zone.Type.IsSmartlink() {
				target := admodels.SmartlinkFromModel(zone)
				target.Comp, _ = acc.companyAccessor.CompanyByID(it.(*models.Zone).CompanyID)
				zoneList = append(zoneList, target)
			} else {
				target := admodels.ZoneFromModel(zone)
				target.Comp, _ = acc.companyAccessor.CompanyByID(it.(*models.Zone).CompanyID)
				zoneList = append(zoneList, target)
			}
		}
		sort.Slice(zoneList, func(i, j int) bool { return zoneList[i].ID() < zoneList[j].ID() })
		acc.zoneList = zoneList
	}
	return acc.zoneList, nil
}

// TargetByID returns campaign object with specific ID
func (acc *ZoneAccessor) TargetByID(id uint64) (admodels.Target, error) {
	list, err := acc.ZoneList()
	if err != nil {
		return nil, err
	}
	idx := sort.Search(len(list), func(i int) bool { return list[i].ID() >= id })
	if idx >= 0 && idx < len(list) && list[idx].ID() == id {
		return list[idx], nil
	}
	return nil, nil
}
