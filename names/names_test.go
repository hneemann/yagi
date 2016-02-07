package names

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNames(t *testing.T) {
	data := []struct {
		out, tem, exp string
	}{
		{out: "", tem: "list.go", exp: "gen-list.go"},
		{out: "", tem: "./list/list.go", exp: "list.go"},
		{out: "z.go", tem: "./list/list.go", exp: "z.go"},
	}

	for _, d := range data {
		res := createOutNameInt(d.out, d.tem)
		assert.Equal(t, d.exp, res, "checked %v, expected '%v', got '%v'", fmt.Sprint(d), d.exp, res)
	}
}
