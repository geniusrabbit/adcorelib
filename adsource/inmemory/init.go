package inmemory

import (
	"context"

	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/billing"
	"github.com/geniusrabbit/adcorelib/platform/info"
	"github.com/geniusrabbit/adstorage/accessors/campaignaccessor"
)

const protocol = "inmemory"

type balanceManager interface {
	MakeVirtualView(test bool, ad *admodels.Ad, bidPrice billing.Money) error
}

type factory struct {
	adCampaigns    *campaignaccessor.CampaignAccessor
	balanceManager balanceManager
}

// NewFactory from adsource
func NewFactory(ctx context.Context, adCampaigns *campaignaccessor.CampaignAccessor, balanceManager balanceManager) *factory {
	return &factory{adCampaigns: adCampaigns, balanceManager: balanceManager}
}

func (f *factory) New(ctx context.Context) (adtype.SourceTester, error) {
	d := &driver{
		adCampaigns:    f.adCampaigns,
		balanceManager: f.balanceManager,
	}
	d.init()
	return d, nil
}

func (*factory) Info() info.Platform {
	return info.Platform{
		Name:        "InMemory",
		Protocol:    protocol,
		Versions:    []string{"1.0"},
		Description: "In memory advertisement accessor for basic use cases",
	}
}

func (*factory) Protocols() []string {
	return []string{protocol}
}
