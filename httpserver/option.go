package httpserver

import (
	"geniusrabbit.dev/adcorelib/httpserver/extensions"
	"go.uber.org/zap"

	"github.com/valyala/fasthttp"
)

// Option type
type Option func(srv *Server)

// WithServiceName which represents the internal name of the service
func WithServiceName(name string) Option {
	return func(srv *Server) {
		srv.serviceName = name
	}
}

// WithDebugMode of the server
func WithDebugMode(debug bool) Option {
	return func(srv *Server) {
		srv.debug = debug
	}
}

// WithCustomHTTPServer setup customly configured server
func WithCustomHTTPServer(server *fasthttp.Server) Option {
	return func(srv *Server) {
		srv.httpServer = server
	}
}

// WithExtensions setup custom
func WithExtensions(exts ...extensions.ServerExtension) Option {
	return func(srv *Server) {
		srv.extensions = exts
	}
}

// WithCustomRouter registrator
func WithCustomRouter(fouterFnk customRouterFnk) Option {
	return func(srv *Server) {
		srv.customRouter = fouterFnk
	}
}

// WithLogger interface
func WithLogger(logger *zap.Logger) Option {
	return func(srv *Server) {
		srv.logger = logger
	}
}
