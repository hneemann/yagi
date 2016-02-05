package tem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	l := List{}
	l.Add(1)
	l.Add(7)
	l.Add(-3)
	assert.Equal(t, 3, l.len())
	assert.Equal(t, 1, l.Items()[0])
	assert.Equal(t, 7, l.Items()[1])
	assert.Equal(t, -3, l.Items()[2])
}
