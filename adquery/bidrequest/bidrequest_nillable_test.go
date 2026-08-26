package bidrequest

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/geniusrabbit/adcorelib/adtype"
)

func TestBidRequestIsNil(t *testing.T) {
	var typedNil *BidRequest
	assert.True(t, typedNil.IsNil())
	assert.True(t, adtype.IsNil(typedNil))
	assert.False(t, (&BidRequest{}).IsNil())
	assert.False(t, adtype.IsNil(&BidRequest{}))
}

func TestBidRequestAccessPointNilReceiver(t *testing.T) {
	var typedNil *BidRequest
	assert.Nil(t, typedNil.AccessPoint())
	assert.Nil(t, (&BidRequest{}).AccessPoint())
}

func TestNewErrorResponseTypedNilRequest(t *testing.T) {
	err := errors.New("no supported impressions")
	var typedNil *BidRequest
	var resp *adtype.ResponseError
	require.NotPanics(t, func() {
		resp = adtype.NewErrorResponse(typedNil, err)
	})
	require.NotNil(t, resp)
	assert.Nil(t, resp.Request())
	assert.True(t, adtype.IsNil(resp.Request()))
	assert.Equal(t, err, resp.Error())
}

func TestNewErrorResponseUntypedNilRequest(t *testing.T) {
	err := errors.New("bad request")
	resp := adtype.NewErrorResponse(nil, err)
	assert.Nil(t, resp.Request())
	assert.Equal(t, err, resp.Error())
}

func TestNewErrorResponseKeepsLiveRequest(t *testing.T) {
	req := &BidRequest{IDVal: "req-1"}
	err := errors.New("skipped")
	resp := adtype.NewErrorResponse(req, err)
	assert.Same(t, req, resp.Request())
	assert.False(t, adtype.IsNil(resp.Request()))
	assert.Equal(t, err, resp.Error())
}
