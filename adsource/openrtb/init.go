package openrtb

import (
	"context"

	"geniusrabbit.dev/adcorelib/admodels"
	"geniusrabbit.dev/adcorelib/adsource/srctestwrapper"
	"geniusrabbit.dev/adcorelib/adtype"
	"geniusrabbit.dev/adcorelib/platform/info"
)

const protocol = "openrtb"

type factory struct{}

func NewFactory() *factory {
	return &factory{}
}

func (*factory) New(ctx context.Context, source *admodels.RTBSource, opts ...any) (adtype.SourceTester, error) {
	dr, err := newDriver(ctx, source, opts...)
	if err != nil {
		return nil, err
	}
	return srctestwrapper.Wrap(source, dr), nil
}

func (*factory) Info() info.Platform {
	return info.Platform{
		Name:        "OpenRTB",
		Protocol:    protocol,
		Versions:    []string{"2.3", "2.4", "2.5", "2.6", "3.0"},
		Description: "",
		Docs: []info.Documentation{
			{
				Title: "OpenRTB (Real-Time Bidding)",
				Link:  "https://www.iab.com/guidelines/real-time-bidding-rtb-project/",
			},
			{
				Title: "Digital Video Ad Serving Template (VAST)",
				Link:  "https://www.iab.com/guidelines/vast/",
			},
		},
		Subprotocols: []info.Subprotocol{
			{
				Name:     "VAST",
				Protocol: "vast",
			},
			{
				Name:     "OpenNative",
				Protocol: "opennative",
			},
		},
	}
}

func (*factory) Protocols() []string {
	return []string{"openrtb", "openrtb2", "openrtb3"}
}
