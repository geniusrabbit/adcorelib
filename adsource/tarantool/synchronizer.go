package tarantool

import (
	"time"

	libtarantool "github.com/tarantool/go-tarantool"

	"geniusrabbit.dev/adcorelib/storage"
)

type synchronizer struct {
	conn *libtarantool.Connection

	namespace string

	lastSyncTime time.Time
}

// NewSynchronizer object
func NewSynchronizer(opts ...OptionFnk) (*synchronizer, error) {
	conn, opt, err := newConnection(opts...)
	if err != nil {
		return nil, err
	}
	return &synchronizer{namespace: opt.namespace, conn: conn}, nil
}

func (sync *synchronizer) Sync(reader storage.Reader) error {
	return nil
}
