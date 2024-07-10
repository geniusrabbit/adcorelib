// Package openrtb provides implementation of the OpenRTB protocol for the adsource package.
// Supported versions: 2.3, 2.4, 2.5, 2.6, 3.0+
package openrtb

import (
	"context"
	"time"

	"github.com/demdxx/gocast/v2"
	"github.com/geniusrabbit/adcorelib/admodels"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/net/httpclient"
	"github.com/geniusrabbit/adcorelib/platform/info"
)

const (
	protocol       = "openrtb"
	defaultTimeout = 150 * time.Millisecond
)

type NewClientFnk[NetDriver httpclient.Driver[Rq, Rs], Rq httpclient.Request, Rs httpclient.Response] func(context.Context, time.Duration) (NetDriver, error)

type factory[NetDriver httpclient.Driver[Rq, Rs], Rq httpclient.Request, Rs httpclient.Response] struct {
	newClientFnk NewClientFnk[NetDriver, Rq, Rs]
}

func NewFactory[ND httpclient.Driver[Rq, Rs], Rq httpclient.Request, Rs httpclient.Response](newClient NewClientFnk[ND, Rq, Rs]) *factory[ND, Rq, Rs] {
	return &factory[ND, Rq, Rs]{
		newClientFnk: newClient,
	}
}

func (fc *factory[ND, Rq, Rs]) New(ctx context.Context, source *admodels.RTBSource, opts ...any) (adtype.SourceTester, error) {
	ncli, err := fc.newClientFnk(ctx, gocast.IfThen(
		source.Timeout > 0,
		time.Duration(source.Timeout)*time.Millisecond,
		defaultTimeout,
	))
	if err != nil {
		return nil, err
	}
	dr, err := newDriver(ctx, source, ncli, opts...)
	if err != nil {
		return nil, err
	}
	return dr, nil
}

func (*factory[ND, Rq, Rs]) Info() info.Platform {
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

func (*factory[ND, Rq, Rs]) Protocols() []string {
	return []string{"openrtb", "openrtb2", "openrtb3"}
}
