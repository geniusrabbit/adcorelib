package adtype

// SourceEmpty abstraction
type SourceEmpty struct {
	PriceCorrectionReduce float64 // from 0 to 1
}

// Bid request for standart system filter
func (*SourceEmpty) Bid(request *BidRequest) Responser { return nil }

// ProcessResponseItem result or error
func (*SourceEmpty) ProcessResponseItem(Responser, ResponserItem) {}

// ID of the source driver
func (*SourceEmpty) ID() uint64 { return 0 }

// ObjectKey of the source driver
func (*SourceEmpty) ObjectKey() uint64 { return 0 }

// Protocol of the source driver
func (*SourceEmpty) Protocol() string { return "undefined" }

// Test request before processing
func (*SourceEmpty) Test(request *BidRequest) bool { return false }

// PriceCorrectionReduceFactor which is a potential
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (s *SourceEmpty) PriceCorrectionReduceFactor() float64 { return s.PriceCorrectionReduce }

// RequestStrategy description
func (*SourceEmpty) RequestStrategy() RequestStrategy { return 0 }
