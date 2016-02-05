package intset

import "strconv"

type IntSet map[int]struct{}

func (i IntSet) String() string {
	n := ""
	for i, _ := range i {
		if len(n) > 0 {
			n += " "
		}
		n += strconv.FormatInt(int64(i), 10)
	}
	return n
}

func (i *IntSet) Add(a int) {
	(*i)[a] = struct{}{}
}

func (i *IntSet) AddAll(a IntSet) {
	for el, _ := range a {
		i.Add(el)
	}
}

func (i IntSet) Has(a int) bool {
	_, ok := i[a]
	return ok
}
