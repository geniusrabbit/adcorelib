package cluster

import (
	"sync/atomic"

	"github.com/demdxx/gocast/v2"
)

type SyncUInt64Value struct {
	val uint64
}

func NewSyncUInt64Value(val uint64) *SyncUInt64Value {
	return &SyncUInt64Value{val: val}
}

func (v *SyncUInt64Value) Value() uint64 {
	return atomic.LoadUint64(&v.val)
}

func (v *SyncUInt64Value) SetValue(val any) {
	atomic.StoreUint64(&v.val, gocast.Uint64(val))
}
