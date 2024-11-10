package cluster

import (
	"sync"

	"github.com/demdxx/gocast/v2"
)

type SyncValue[T any] struct {
	mx  sync.RWMutex
	val T
}

func NewSyncValue[T any](val T) *SyncValue[T] {
	return &SyncValue[T]{val: val}
}

func (v *SyncValue[T]) Value() T {
	v.mx.RLock()
	defer v.mx.RUnlock()
	return v.val
}

func (v *SyncValue[T]) SetValue(val any) {
	v.mx.Lock()
	defer v.mx.Unlock()
	v.val = gocast.Cast[T](val)
}
