package stdhttpclient

import "net/http"

type Request struct {
	HTTP *http.Request
}

func (r *Request) SetHeader(key, value string) {
	r.HTTP.Header.Set(key, value)
}
