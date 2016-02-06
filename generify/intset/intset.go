package intset

import "strconv"

// IntSet represents a simple set of int
type IntSet map[int]struct{}

func (i IntSet) String() string {
	n := ""
	for i := range i {
		if len(n) > 0 {
			n += " "
		}
		n += strconv.FormatInt(int64(i), 10)
	}
	return n
}

// Add adds a int to the set
func (i *IntSet) Add(a int) {
	(*i)[a] = struct{}{}
}

// AddAll adds a other set to this set
func (i *IntSet) AddAll(a IntSet) {
	for el := range a {
		i.Add(el)
	}
}

// Has checks if this set contains a
func (i IntSet) Has(a int) bool {
	_, ok := i[a]
	return ok
}
