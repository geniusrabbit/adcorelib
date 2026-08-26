package adtype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type nillablePtr struct{}

func (p *nillablePtr) IsNil() bool { return p == nil }

func TestIsNil(t *testing.T) {
	assert.True(t, IsNil(nil))
	assert.True(t, IsNil((*nillablePtr)(nil)))
	assert.False(t, IsNil(&nillablePtr{}))
}
