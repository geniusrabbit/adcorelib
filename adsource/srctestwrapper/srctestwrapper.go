package srctestwrapper

import (
	"geniusrabbit.dev/adcorelib/admodels"
	"geniusrabbit.dev/adcorelib/adtype"
)

type sourceTester struct {
	adtype.SourceMinimal
	sourceInfo *admodels.RTBSource
}

func Wrap(sourceInfo *admodels.RTBSource, source adtype.SourceMinimal) adtype.SourceTester {
	return &sourceTester{SourceMinimal: source, sourceInfo: sourceInfo}
}

// ID of the source driver
func (w *sourceTester) ID() uint64 {
	return w.sourceInfo.ID
}

// Test current request for compatibility
func (w *sourceTester) Test(request *adtype.BidRequest) bool {
	return true
}

// PriceCorrectionReduceFactor which is a potential
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (w *sourceTester) PriceCorrectionReduceFactor() float64 {
	return w.sourceInfo.PriceCorrectionReduceFactor()
}

// RequestStrategy description
func (w *sourceTester) RequestStrategy() adtype.RequestStrategy {
	return adtype.AsynchronousRequestStrategy
}
