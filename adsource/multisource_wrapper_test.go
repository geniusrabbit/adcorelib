// Package adsource provides the implementation of ad source drivers for the AdEngine.
// This package offers a unified interface for interacting with various ad sources
// and includes a collection of standard ad sources such as in-memory, database,
// OpenRTB (Real-Time Bidding) sources, and more.
//
// The primary components of the package are:
// - MultisourceWrapper: An abstraction that manages multiple ad sources and controls
//   the request distribution and response collection from these sources.
// - Error definitions: Standardized error messages used across the package.
// - Internal methods: Helper functions and internal logic to support the primary
//   operations of the MultisourceWrapper and other components.
//
// The package ensures efficient and parallel processing of ad requests by utilizing
// a worker pool for executing bid requests. It also integrates with tracing and
// logging systems to provide detailed insights into the bidding process and performance.
//
// Key Features:
// - Unified interface for multiple ad sources
// - Support for parallel bid requests
// - Integration with tracing (using OpenTracing) and logging (using zap)
// - Metrics collection and reporting for monitoring performance
//
// Example usage:
//   wrapper, err := adsource.NewMultisourceWrapper(options...)
//   if err != nil {
//       log.Fatal(err)
//   }
//
//   response := wrapper.Bid(request)
//   if response.Error() != nil {
//       log.Println("Bid request failed:", response.Error())
//   } else {
//       log.Println("Bid request succeeded:", response.Ads())
//   }

package adsource

import (
	"context"
	"testing"
	"time"

	"github.com/demdxx/rpool/v2"
	"github.com/geniusrabbit/adcorelib/adquery/bidrequest"
	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/geniusrabbit/adcorelib/errtype"
	"github.com/geniusrabbit/adcorelib/eventtraking/events"
	"github.com/geniusrabbit/adcorelib/eventtraking/eventstream"
)

func TestMultisourceWrapper_sourceResponseLog(t *testing.T) {
	type fields struct {
		sources            adtype.SourceAccessor
		execpool           *rpool.Pool
		requestTimeout     time.Duration
		maxParallelRequest int
		metrics            Metrics
	}
	type args struct {
		in0      *bidrequest.BidRequest
		response adtype.Response
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrp := &MultisourceWrapper{
				sources:            tt.fields.sources,
				execpool:           tt.fields.execpool,
				requestTimeout:     tt.fields.requestTimeout,
				maxParallelRequest: tt.fields.maxParallelRequest,
				metrics:            tt.fields.metrics,
			}
			wrp.sourceResponseLog(tt.args.in0, tt.args.response)
		})
	}
}

type recordingStream struct {
	fraud int
	fail  int
}

func (recordingStream) SendEvent(context.Context, any) error { return nil }
func (recordingStream) Send(events.Type, uint8, adtype.Response, adtype.ResponseItem) error {
	return nil
}
func (recordingStream) SendLeadEvent(context.Context, any) error { return nil }
func (recordingStream) SendSourceSkip(adtype.Response) error     { return nil }
func (recordingStream) SendSourceNoBid(adtype.Response) error    { return nil }
func (s *recordingStream) SendSourceFail(adtype.Response) error  { s.fail++; return nil }
func (s *recordingStream) SendSourceFraud(adtype.Response) error { s.fraud++; return nil }
func (recordingStream) SendAccessPointBid(adtype.Response, ...adtype.ResponseItem) error {
	return nil
}
func (recordingStream) SendAccessPointSkip(adtype.Response) error  { return nil }
func (recordingStream) SendAccessPointNoBid(adtype.Response) error { return nil }
func (recordingStream) SendAccessPointFail(adtype.Response) error  { return nil }
func (recordingStream) SendAccessPointFraud(adtype.Response) error { return nil }

func TestSourceResponseLogFraudIsNotFail(t *testing.T) {
	stream := &recordingStream{}
	req := &bidrequest.BidRequest{}
	req.SetContext(eventstream.WithStream(context.Background(), stream))
	response := adtype.NewErrorResponse(req, errtype.Error("response fraud").WithMessage("blacklist"))

	(&MultisourceWrapper{}).sourceResponseLog(req, response)

	if stream.fraud != 1 || stream.fail != 0 {
		t.Fatalf("fraud=%d fail=%d, want fraud event", stream.fraud, stream.fail)
	}
}
