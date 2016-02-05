package intset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	s := IntSet{}
	assert.False(t, s.Has(2))
	s.Add(1)
	assert.False(t, s.Has(2))
	assert.True(t, s.Has(1))
}

func TestAddAll(t *testing.T) {
	s := IntSet{}
	s.Add(1)
	s2 := IntSet{}
	s2.Add(5)
	s2.Add(9)
	s.AddAll(s2)

	assert.True(t, s.Has(1))
	assert.True(t, s.Has(5))
	assert.True(t, s.Has(9))
}
