package set

import "fmt"

//go:generate yagi -tem=set.go -gen=int

//generic
type ITEM interface{}

// Set represents a simple set
type Set map[ITEM]struct{}

func (i Set) String() string {
	n := ""
	for i := range i {
		if len(n) > 0 {
			n += " "
		}
		n += fmt.Sprintf("%v", i)
	}
	return n
}

// Add adds a item to the set
func (i *Set) Add(a ITEM) {
	(*i)[a] = struct{}{}
}

// AddAll adds a other set to this set
func (i *Set) AddAll(a Set) {
	for el := range a {
		i.Add(el)
	}
}

// Has checks if this set contains the item
func (i Set) Has(a ITEM) bool {
	_, ok := i[a]
	return ok
}

// Has checks if this set contains the item
func (i Set) Items() (res []ITEM) {
	for item := range i {
		res = append(res, item)
	}
	return
}
