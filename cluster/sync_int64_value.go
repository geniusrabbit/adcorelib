package cluster

import (
	"sync/atomic"

	"github.com/demdxx/gocast/v2"
)

type SyncInt64Value struct {
	val int64
}

func NewSyncInt64Value(val int64) *SyncInt64Value {
	return &SyncInt64Value{val: val}
}

func (v *SyncInt64Value) Value() int64 {
	return atomic.LoadInt64(&v.val)
}

func (v *SyncInt64Value) SetValue(val any) {
	atomic.StoreInt64(&v.val, gocast.Int64(val))
}
