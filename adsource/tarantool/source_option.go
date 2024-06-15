package tarantool

import tarantool "github.com/tarantool/go-tarantool"

type option struct {
	addr      string
	opts      tarantool.Opts
	namespace string
}

// OptionFnk type
type OptionFnk func(opt *option)

// OptionAddr of source option
func OptionAddr(addr string) OptionFnk {
	return func(opt *option) {
		opt.addr = addr
	}
}

// OptionConnectOpts of source option
func OptionConnectOpts(opts tarantool.Opts) OptionFnk {
	return func(opt *option) {
		opt.opts = opts
	}
}

// OptionLogin of source option
func OptionLogin(user, password string) OptionFnk {
	return func(opt *option) {
		opt.opts.User = user
		opt.opts.Pass = password
	}
}

// OptionNamespace of source option
func OptionNamespace(namespace string) OptionFnk {
	return func(opt *option) {
		opt.namespace = namespace
	}
}
