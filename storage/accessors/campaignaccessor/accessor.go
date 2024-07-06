package campaignaccessor

import (
	"sort"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/admodels/types"
	"github.com/geniusrabbit/adcorelib/models"
	"github.com/geniusrabbit/adcorelib/storage/loader"
)

// CampaignAccessor provides accessor to the admodel company type
type CampaignAccessor struct {
	loader.DataAccessor
	formatAccessor types.FormatsAccessor

	campaignList []*admodels.Campaign
}

// NewCampaignAccessor from dataAccessor
func NewCampaignAccessor(dataAccessor loader.DataAccessor, formatAccessor types.FormatsAccessor) *CampaignAccessor {
	return &CampaignAccessor{DataAccessor: dataAccessor, formatAccessor: formatAccessor}
}

// CampaignList returns list of prepared data
func (acc *CampaignAccessor) CampaignList() ([]*admodels.Campaign, error) {
	if acc.campaignList == nil || acc.NeedUpdate() {
		data, err := acc.Data()
		if err != nil {
			return nil, err
		}
		campaignList := make([]*admodels.Campaign, 0, len(data))
		for _, it := range data {
			campaignList = append(campaignList,
				admodels.CampaignFromModel(it.(*models.Campaign), acc.formatAccessor))
		}
		sort.Slice(campaignList, func(i, j int) bool { return campaignList[i].ID < campaignList[j].ID })
		acc.campaignList = campaignList
	}
	return acc.campaignList, nil
}

// CampaignByID returns campaign object with specific ID
func (acc *CampaignAccessor) CampaignByID(id uint64) (*admodels.Campaign, error) {
	list, err := acc.CampaignList()
	if err != nil {
		return nil, err
	}
	idx := sort.Search(len(list), func(i int) bool { return list[i].ID >= id })
	if idx >= 0 && idx < len(list) && list[idx].ID == id {
		return list[idx], nil
	}
	return nil, nil
}
