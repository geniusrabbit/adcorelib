// thread-safeless

package adsourceaccessor

import (
	"github.com/geniusrabbit/adcorelib/adtype"
)

type priorityIterator struct {
	index   int
	request *adtype.BidRequest
	sources []adtype.Source
}

// NewPriorityIterator from request and source
func NewPriorityIterator(request *adtype.BidRequest, sources []adtype.Source) adtype.SourceIterator {
	return &priorityIterator{
		index:   0,
		request: request,
		sources: sources,
	}
}

func (iter *priorityIterator) Next() (src adtype.Source) {
	if iter.index < len(iter.sources) {
		src = iter.sources[iter.index]
		iter.index++
	}
	return src
}

var _ adtype.SourceIterator = &priorityIterator{}
