package stdhttpclient

import (
	"io"
	"net/http"
)

type Response struct {
	HTTP *http.Response
}

func (r *Response) StatusCode() int {
	if r == nil {
		return 0
	}
	return r.HTTP.StatusCode
}

func (r *Response) Body() io.Reader {
	return r.HTTP.Body
}

func (r *Response) IsNoop() bool {
	return r == nil || r.HTTP == nil
}

func (r *Response) Close() error {
	if r.HTTP == nil {
		return nil
	}
	return r.HTTP.Body.Close()
}
