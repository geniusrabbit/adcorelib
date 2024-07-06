package tarantool

import (
	libtarantool "github.com/tarantool/go-tarantool"

	"github.com/geniusrabbit/adcorelib/adtype"
)

// Source provides access to the advertisement campaigns
type source struct {
	conn      *libtarantool.Connection
	namespace string
}

// NewSource accessor
func NewSource(opts ...OptionFnk) (*source, error) {
	conn, opt, err := newConnection(opts...)
	if err != nil {
		return nil, err
	}
	return &source{namespace: opt.namespace, conn: conn}, nil
}

// Bid request for standart system filter
func (src *source) Bid(request *adtype.BidRequest) adtype.Responser {
	var (
		list []any
		err  error
	)

	// if target := request.Target(); target != nil {
	// 	switch z := target.(type) {
	// 	case *models.Smartlink:
	// 		if len(z.Campaigns) > 0 {
	// 			err = src.conn.CallTyped("search_strict_ads", []any{}, &list)
	// 			return adtype.NewErrorResponse(request, err)
	// 		}
	// 	}
	// }

	if err = src.conn.CallTyped("search_ads", []any{}, &list); err != nil {
		return adtype.NewErrorResponse(request, err)
	}

	return adtype.NewErrorResponse(request, err)
}

// ProcessResponseItem result or error
func (src *source) ProcessResponseItem(resp adtype.Responser, item adtype.ResponserItem) {
	// TODO: increment counters
}
