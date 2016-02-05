package concrete

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	c, err := New("int")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c.Instance))
	assert.Equal(t, 1, len(c.Instance[0]))
	assert.Equal(t, "int", c.Instance[0][0])
}

func TestTwoInstances(t *testing.T) {
	c, err := New("int32;int64")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(c.Instance))
	assert.Equal(t, 1, len(c.Instance[0]))
	assert.Equal(t, "int32", c.Instance[0][0])
	assert.Equal(t, 1, len(c.Instance[1]))
	assert.Equal(t, "int64", c.Instance[1][0])
}

func TestEmptyInstances(t *testing.T) {
	_, err := New("int32;")
	assert.Error(t, err)
}

func TestEmptyInstances2(t *testing.T) {
	_, err := New("string,int32;string, ")
	assert.Error(t, err)
}

func TestNotSameSize(t *testing.T) {
	_, err := New("string,int32;string")
	assert.Error(t, err)
}

func TestPointer(t *testing.T) {
	_, err := New("string,*int32")
	assert.NoError(t, err)
}

func TestEmpty(t *testing.T) {
	_, err := New("")
	assert.Error(t, err)
}
