package cluster

type Valuer[T any] interface {
	Value() T
	SetValue(any)
}
