package httpclient

import "io"

type Request interface {
	SetHeader(key, value string)
}

type Response interface {
	io.Closer
	StatusCode() int
	Body() io.Reader
	IsNoop() bool
}

type Driver[Rq Request, Rs Response] interface {
	Request(method, url string, body io.Reader) (Rq, error)
	NoopRequest() Rq
	Do(req Rq) (Rs, error)
}
