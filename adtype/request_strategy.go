package adtype

// RequestStrategy defines politics of request sending
type RequestStrategy int

const (
	// AsynchronousRequestStrategy is default strategy implies
	// requesting all auction participants and choising the most
	// profitable variant of all
	AsynchronousRequestStrategy RequestStrategy = iota

	// SingleRequestStrategy tells that if response was
	// received it should be performed
	SingleRequestStrategy
)

func (rs RequestStrategy) IsSingle() bool {
	return rs == SingleRequestStrategy
}

func (rs RequestStrategy) IsAsynchronous() bool {
	return rs == AsynchronousRequestStrategy
}
