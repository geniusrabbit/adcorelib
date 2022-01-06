package adtype

// SourceEmpty abstraction
type SourceEmpty struct{}

// Bid request for standart system filter
func (*SourceEmpty) Bid(request *BidRequest) Responser { return nil }

// ProcessResponseItem result or error
func (*SourceEmpty) ProcessResponseItem(Responser, ResponserItem) {}

// ID of the source driver
func (*SourceEmpty) ID() uint64 { return 0 }

// Test request before processing
func (*SourceEmpty) Test(request *BidRequest) bool { return false }

// PriceCorrectionReduceFactor which is a potential
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (*SourceEmpty) PriceCorrectionReduceFactor() float64 { return 0 }

// RequestStrategy description
func (*SourceEmpty) RequestStrategy() RequestStrategy { return 0 }
