package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActiveStatus(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		assert.Equal(t, "active", StatusActive.Name())
		assert.True(t, StatusActive.IsActive())
		assert.Equal(t, "pause", StatusPause.Name())
		assert.False(t, StatusPause.IsActive())
	})
	t.Run("decode", func(t *testing.T) {
		tests := []struct {
			name     string
			expected ActiveStatus
		}{
			{
				name:     "active",
				expected: StatusActive,
			},
			{
				name:     "pause",
				expected: StatusPause,
			},
			{
				name:     "other",
				expected: StatusPause,
			},
		}
		for _, test := range tests {
			assert.Equal(t, ActiveNameToStatus(test.name), test.expected)
			data, err := json.Marshal(test.expected)
			assert.NoError(t, err)
			var st ActiveStatus
			err = json.Unmarshal(data, &st)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, st)
		}
	})
}

func TestApproveStatus(t *testing.T) {
	t.Run("encode", func(t *testing.T) {
		assert.Equal(t, "pending", StatusPending.Name())
		assert.False(t, StatusPending.IsApproved())
		assert.Equal(t, "approved", StatusApproved.Name())
		assert.True(t, StatusApproved.IsApproved())
		assert.Equal(t, "rejected", StatusRejected.Name())
		assert.False(t, StatusRejected.IsApproved())
	})
	t.Run("decode", func(t *testing.T) {
		tests := []struct {
			name     string
			expected ApproveStatus
		}{
			{
				name:     "pending",
				expected: StatusPending,
			},
			{
				name:     "approved",
				expected: StatusApproved,
			},
			{
				name:     "rejected",
				expected: StatusRejected,
			},
		}
		for _, test := range tests {
			assert.Equal(t, ApproveNameToStatus(test.name), test.expected)
			data, err := json.Marshal(test.expected)
			assert.NoError(t, err)
			var st ApproveStatus
			err = json.Unmarshal(data, &st)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, st)
		}
	})
}
