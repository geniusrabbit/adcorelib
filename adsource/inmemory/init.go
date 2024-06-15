package inmemory

import (
	"context"

	"geniusrabbit.dev/adcorelib/adtype"
	"geniusrabbit.dev/adcorelib/platform/info"
	"geniusrabbit.dev/adcorelib/storage/accessors/campaignaccessor"
)

const protocol = "inmemory"

type factory struct {
	adCampaigns *campaignaccessor.CampaignAccessor
}

// NewFactory from adsource
func NewFactory(ctx context.Context, adCampaigns *campaignaccessor.CampaignAccessor) *factory {
	return &factory{adCampaigns: adCampaigns}
}

func (f *factory) New(ctx context.Context) (adtype.SourceTester, error) {
	d := &driver{adCampaigns: f.adCampaigns}
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
