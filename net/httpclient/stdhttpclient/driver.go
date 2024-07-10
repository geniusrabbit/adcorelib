package stdhttpclient

import (
	"io"
	"net/http"
)

type Driver struct {
	HTTPClient *http.Client
}

func NewDriver() *Driver {
	return NewDriverWithHTTPClient(&http.Client{})
}

func NewDriverWithHTTPClient(client *http.Client) *Driver {
	return &Driver{HTTPClient: client}
}

func (d *Driver) Request(method, url string, body io.Reader) (*Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	return &Request{HTTP: req}, nil
}

func (d *Driver) Do(req *Request) (*Response, error) {
	resp, err := d.HTTPClient.Do(req.HTTP)
	if err != nil {
		return nil, err
	}
	return &Response{HTTP: resp}, nil
}

func (d *Driver) NoopRequest() *Request {
	return nil
}
