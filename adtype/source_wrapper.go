package adtype

type minimalSourceWrapper struct {
	source SourceMinimal
}

// ID of the source driver
func (wp minimalSourceWrapper) ID() uint64 { return 0 }

// Protocol of the source driver
func (wp minimalSourceWrapper) Protocol() string { return "undefined" }

// Test request before processing
func (wp minimalSourceWrapper) Test(request *BidRequest) bool { return true }

// Bid request for standart system filter
func (wp minimalSourceWrapper) Bid(request *BidRequest) Responser { return wp.source.Bid(request) }

// ProcessResponseItem result or error
func (wp minimalSourceWrapper) ProcessResponseItem(response Responser, item ResponserItem) {
	wp.source.ProcessResponseItem(response, item)
}

// PriceCorrectionReduceFactor which is a potential
// Returns percent from 0 to 1 for reducing of the value
// If there is 10% of price correction, it means that 10% of the final price must be ignored
func (wp minimalSourceWrapper) PriceCorrectionReduceFactor() float64 { return 0 }

// RequestStrategy description
func (wp minimalSourceWrapper) RequestStrategy() RequestStrategy {
	return AsynchronousRequestStrategy
}

// ToSource interface from different types of interfaces with the implementation of unsupported methods
func ToSource(val SourceMinimal) Source {
	switch v := val.(type) {
	case Source:
		return v
	case SourceMinimal:
		return minimalSourceWrapper{source: v}
	}
	return nil
}
