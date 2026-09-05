package bidresponse

import (
	"testing"

	"github.com/geniusrabbit/adcorelib/adtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseItemSkipsNils(t *testing.T) {
	live := &ResponseItemBlank{
		ItemID: "ad1",
		Imp:    &adtype.Impression{ID: "imp1"},
	}
	var typedNil *ResponseItemBlank
	r := &Response{
		items: []adtype.ResponseItemCommon{nil, typedNil, live},
	}

	got := r.Item("imp1")
	require.NotNil(t, got)
	assert.Equal(t, "ad1", got.ID())
	assert.Nil(t, r.Item("missing"))
}

func TestResponseItemNilReceiverAndEmpty(t *testing.T) {
	assert.Nil(t, (*Response)(nil).Item("imp1"))
	assert.Nil(t, (&Response{}).Item("imp1"))

	var typedNil *ResponseItemBlank
	r := &Response{items: []adtype.ResponseItemCommon{nil, typedNil}}
	assert.Nil(t, r.Item("imp1"))
}

func TestResponseValidateSkipsNils(t *testing.T) {
	live := &ResponseItemBlank{
		ItemID: "ad1",
		Imp:    &adtype.Impression{ID: "imp1"},
	}
	var typedNil *ResponseItemBlank
	r := &Response{
		items: []adtype.ResponseItemCommon{nil, typedNil, live},
	}
	assert.NoError(t, r.Validate())
}

func TestFilterNilItems(t *testing.T) {
	live := &ResponseItemBlank{
		ItemID: "ad1",
		Imp:    &adtype.Impression{ID: "imp1"},
	}
	var typedNil *ResponseItemBlank
	got := filterNilItems([]adtype.ResponseItemCommon{nil, typedNil, live})
	require.Len(t, got, 1)
	assert.Equal(t, "ad1", got[0].ID())
	assert.Empty(t, filterNilItems([]adtype.ResponseItemCommon{nil, typedNil}))
	assert.Nil(t, filterNilItems(nil))
}

func TestAddItemIgnoresNil(t *testing.T) {
	live := &ResponseItemBlank{
		ItemID: "ad1",
		Imp:    &adtype.Impression{ID: "imp1"},
	}
	var typedNil *ResponseItemBlank
	r := &Response{}
	r.AddItem(nil)
	r.AddItem(typedNil)
	r.AddItem(live)
	require.Len(t, r.items, 1)
	assert.Equal(t, "ad1", r.items[0].ID())
}
