package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	s := Set{}
	assert.False(t, s.Has(2))
	s.Add(1)
	assert.False(t, s.Has(2))
	assert.True(t, s.Has(1))
}

func TestAddAll(t *testing.T) {
	s := Set{}
	s.Add(1)
	s2 := Set{}
	s2.Add(5)
	s2.Add(9)
	s.AddAll(s2)

	assert.True(t, s.Has(1))
	assert.True(t, s.Has(5))
	assert.True(t, s.Has(9))
}

func TestItems(t *testing.T) {
	s := Set{}
	s.Add(1)
	s.Add(2)
	items := s.Items()
	assert.Equal(t, 2, len(items))
	assert.True(t, (items[0] == 1 && items[1] == 2) ||
		(items[0] == 2 && items[1] == 1))
}

func TestString(t *testing.T) {
	s := Set{}
	s.Add(1)
	s.Add(2)

	str := s.String()
	assert.True(t, str == "1 2" || str == "2 1")
}
