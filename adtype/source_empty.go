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

// RevenueShareReduceFactor which is a potential
func (*SourceEmpty) RevenueShareReduceFactor() float64 { return 0 }

// RequestStrategy description
func (*SourceEmpty) RequestStrategy() RequestStrategy { return 0 }
