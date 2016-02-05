package concrete

import (
	"errors"
	"strings"
)

type Types []string

type Instances struct {
	Instance []Types
}

func New(types string) (*Instances, error) {
	con := Instances{}
	inst := strings.Split(types, ";")
	for _, i := range inst {
		t := strings.Split(i, ",")
		var ty Types
		for _, t := range t {
			t = strings.TrimSpace(t)
			if t == "" {
				return nil, errors.New("empty concrete type")
			}
			ty = append(ty, t)
		}
		if len(ty) == 0 {
			return nil, errors.New("no concrete type given")
		}
		if len(con.Instance) > 0 {
			if len(ty) != len(con.Instance[0]) {
				return nil, errors.New("not all instances have same number of types")
			}
		}
		con.Instance = append(con.Instance, ty)
	}
	return &con, nil
}
