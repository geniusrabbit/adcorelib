// thread-safeless

package adsourceaccessor

import (
	"math/rand"

	"github.com/geniusrabbit/adcorelib/adtype"
)

type roundrobinIterator struct {
	index    int
	endIndex int
	request  *adtype.BidRequest
	sources  []adtype.Source
}

// NewRoundrobinIterator from request and source
func NewRoundrobinIterator(request *adtype.BidRequest, sources []adtype.Source) adtype.SourceIterator {
	startIndex := rand.Int() % len(sources)
	return &roundrobinIterator{
		index:    startIndex,
		endIndex: startIndex,
		request:  request,
		sources:  sources,
	}
}

func (iter *roundrobinIterator) Next() adtype.Source {
	if iter.index >= len(iter.sources) {
		return nil
	}
	src := iter.sources[iter.index]
	if iter.index++; iter.index > len(iter.sources) {
		iter.index = 0
	}
	if iter.index == iter.endIndex {
		return nil
	}
	return src
}

var _ adtype.SourceIterator = &roundrobinIterator{}
