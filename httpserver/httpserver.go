// package httpserver provides basic HTTP server with DSP handlers
package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"sync/atomic"

	"github.com/fasthttp/router"
	fastp "github.com/flf2ko/fasthttp-prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/geniusrabbit/adcorelib/context/ctxlogger"
	"github.com/geniusrabbit/adcorelib/gtracing"
	"github.com/geniusrabbit/adcorelib/httpserver/extensions"
)

type (
	// personalizedHandler func(person personification.Person, ctx *fasthttp.RequestCtx)
	customRouterFnk func(router *router.Router)
)

// Server implements basic HTTP infostructure and routing
type Server struct {
	// Debug mode of the server
	debug bool

	// service name
	serviceName string

	// net connection to the listen port
	httpConnection net.Listener

	// httpServer object
	httpServer *fasthttp.Server

	// Extensions of the server
	extensions []extensions.ServerExtension

	// Custom router function
	customRouter customRouterFnk

	// Metrics prepared methods
	// metrics Metrics

	// tracer interface helping to trace application
	tracer opentracing.Tracer

	// Setup into the 1 when server is shuting down
	shutdownMode uint32

	// Logger base object
	logger *zap.Logger
}

// NewServer http server object
func NewServer(options ...Option) (*Server, error) {
	srv := &Server{}
	for _, opt := range options {
		opt(srv)
	}
	if srv.logger == nil {
		srv.logger = zap.L().With(zap.String("module", "httpserver"))
	}
	if err := srv.initTracer(); err != nil {
		return nil, err
	}
	return srv, nil
}

// Listen server address
func (srv *Server) Listen(ctx context.Context, address string) (err error) {
	if srv.httpServer == nil {
		srv.httpServer = &fasthttp.Server{ReadBufferSize: 1 << 20}
	}

	p := fastp.NewPrometheus("fasthttp")
	srv.httpServer.Handler = srv.corsHandler(
		p.WrapHandler(srv.newRouter(ctx)),
	)

	srv.httpConnection, err = net.Listen("tcp4", address)
	if err != nil {
		return err
	}

	return srv.httpServer.Serve(srv.httpConnection)
}

// Shutdown server gracefully
func (srv *Server) Shutdown() {
	srv.logger.Debug("Shutdown the HTTP server", zap.String("method", "Shutdown"))

	if !atomic.CompareAndSwapUint32(&srv.shutdownMode, 0, 1) {
		return
	}

	if srv.httpConnection != nil {
		srv.httpConnection.Close()
		srv.httpConnection = nil
	}
}

// IsShutdownMode on or off
func (srv *Server) IsShutdownMode() bool {
	return atomic.LoadUint32(&srv.shutdownMode) == 1
}

func (srv *Server) newRouter(ctx context.Context) *router.Router {
	ctxlogger.Get(ctx).Debug("Initialise the router", zap.String("method", "newRouter"))

	nrt := router.New()
	nrt.PanicHandler = srv.panicCallback

	// Prepare routing by extensions
	for _, ext := range srv.extensions {
		ext.InitRouter(ctx, nrt, srv.tracer)
	}

	// Utility part
	nrt.GET("/healthcheck", srv.healthCheck)
	nrt.GET("/check", srv.check)

	if srv.customRouter != nil {
		srv.customRouter(nrt)
	}

	return nrt
}

///////////////////////////////////////////////////////////////////////////////
/// Handlers
///////////////////////////////////////////////////////////////////////////////

func (srv *Server) healthCheck(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetStatusCode(http.StatusOK)
	headers := strings.TrimSpace(ctx.Request.Header.String())

	_ = json.NewEncoder(ctx.Response.BodyWriter()).Encode(&struct {
		Status  string `json:"status"`
		Headers any    `json:"headers"`
	}{
		Status:  "ok",
		Headers: strings.Split(headers, "\r\n"),
	})
}

func (srv *Server) check(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetStatusCode(http.StatusOK)
	_, _ = fmt.Fprint(ctx.Response.BodyWriter(), `{"status":"ok"}`)
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

// CORS handler adds CORS headers to the response for all requests including fetch
func (srv *Server) corsHandler(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "HEAD,GET,POST,PUT,PATCH,OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		if string(ctx.Method()) == "OPTIONS" {
			ctx.SetStatusCode(http.StatusNoContent)
			return
		}
		next(ctx)
	}
}

func (srv *Server) panicCallback(ctx *fasthttp.RequestCtx, rcv any) {
	_ = srv.logError(fmt.Errorf("server panic: %+v\n%s", rcv, debug.Stack()))
	if srv.debug {
		msg := fmt.Sprintf("server painc: %v\n", rcv)
		_, _ = ctx.Write([]byte(msg))
		_, _ = ctx.Write(debug.Stack())
	}
	ctx.SetStatusCode(http.StatusInternalServerError)
}

func (srv *Server) initTracer() (err error) {
	if srv.tracer, err = gtracing.InitTracer(srv.serviceName, srv.logger); err != nil {
		_ = srv.logError(err)
	}
	return nil
}

func (srv *Server) logError(err error) error {
	if err != nil {
		srv.logger.Error("", zap.Error(err))
	}
	return err
}
