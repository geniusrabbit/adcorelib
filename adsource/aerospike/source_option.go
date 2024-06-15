package aerospike

type option struct {
	hostname  string
	port      int
	namespace string
	udfPath   string
}

// OptionFnk type
type OptionFnk func(opt *option)

// OptionHostname of source option
func OptionHostname(hostname string, port int) OptionFnk {
	return func(opt *option) {
		opt.hostname = hostname
		opt.port = port
	}
}

// OptionNamespace of source option
func OptionNamespace(namespace string) OptionFnk {
	return func(opt *option) {
		opt.namespace = namespace
	}
}

// OptionUDFPath of source option
func OptionUDFPath(path string) OptionFnk {
	return func(opt *option) {
		opt.udfPath = path
	}
}
