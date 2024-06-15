package tarantool

import (
	"time"

	libtarantool "github.com/tarantool/go-tarantool"
)

func newConnection(opts ...OptionFnk) (*libtarantool.Connection, *option, error) {
	var (
		ttopts = libtarantool.Opts{
			Timeout:       500 * time.Millisecond,
			Reconnect:     1 * time.Second,
			MaxReconnects: 3,
		}
		opt  = &option{addr: "127.0.0.1:3013", opts: ttopts, namespace: "default"}
		conn *libtarantool.Connection
		err  error
	)

	for _, o := range opts {
		o(opt)
	}

	if conn, err = libtarantool.Connect(opt.addr, opt.opts); err != nil {
		return nil, nil, err
	}

	return conn, opt, nil
}
