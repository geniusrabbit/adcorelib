package cluster

type ValueSetter interface {
	SetValue(any)
}

type Client interface {
	SubscribeValue(name string, val ValueSetter) error
}
