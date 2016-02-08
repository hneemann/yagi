package set

import (
	"bytes"
	"fmt"
)

//go:generate yagi -tem=set.go -gen=int

//generic
type ITEM interface{}

// Set represents a simple set
type Set map[ITEM]struct{}

// String satisfies the fmt.Stringer interface
func (i Set) String() string {
	var buffer bytes.Buffer
	for i := range i {
		if buffer.Len() > 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(fmt.Sprint(i))
	}
	return buffer.String()
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

// Items returns the set items as a slice
func (i Set) Items() (res []ITEM) {
	for item := range i {
		res = append(res, item)
	}
	return
}
