package companyaccessor

import (
	"sort"
	"sync"

	"geniusrabbit.dev/adcorelib/admodels"
	"geniusrabbit.dev/adcorelib/models"
	"geniusrabbit.dev/adcorelib/storage/loader"
)

// CompanyAccessor provides accessor to the admodel company type
type CompanyAccessor struct {
	loader.DataAccessor
	mx sync.RWMutex

	companyList []*admodels.Company
}

// NewCompanyAccessor from dataAccessor
func NewCompanyAccessor(dataAccessor loader.DataAccessor) *CompanyAccessor {
	return &CompanyAccessor{DataAccessor: dataAccessor}
}

// CompanyList returns list of prepared data
func (acc *CompanyAccessor) CompanyList() ([]*admodels.Company, error) {
	if acc.companyList == nil || acc.NeedUpdate() {
		acc.mx.Lock()
		defer acc.mx.Unlock()
		if acc.companyList == nil || acc.NeedUpdate() {
			data, err := acc.Data()
			if err != nil {
				return nil, err
			}
			companyList := make([]*admodels.Company, 0, len(data))
			for _, it := range data {
				companyList = append(companyList, admodels.CompanyFromModel(it.(*models.Company)))
			}
			sort.Slice(companyList, func(i, j int) bool { return companyList[i].ID < companyList[j].ID })
			acc.companyList = companyList
		}
	}
	return acc.companyList, nil
}

// CompanyByID returns company object with specific ID
func (acc *CompanyAccessor) CompanyByID(id uint64) (*admodels.Company, error) {
	list, err := acc.CompanyList()
	if err != nil {
		return nil, err
	}
	idx := sort.Search(len(list), func(i int) bool { return list[i].ID >= id })
	if idx >= 0 && idx < len(list) && list[idx].ID == id {
		return list[idx], nil
	}
	return nil, nil
}
