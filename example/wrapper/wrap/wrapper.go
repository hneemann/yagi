// this file defines a wrapper for largecode.List
// which is type save an does the boxing and unboxing
package wrap

import "github.com/hneemann/yagi/example/wrapper/largecode"

//generic
type ITEM int

type Wrapper struct {
	parent largecode.List
}

// Add adds an element to the list
func (l *Wrapper) Add(item ITEM) {
	l.parent.Add(item)
}

// Get gets an element to the list
func (l *Wrapper) Get(index int) ITEM {
	item, ok := l.parent.Get(index).(ITEM)
	if ok {
		return item
	}
	panic("wrong type in list")
}

// Remove an element from the list
func (l *Wrapper) Remove(index int) {
	l.parent.Remove(index)
}

// Len returns the number of elements in the list
func (l *Wrapper) Len() int {
	return l.parent.Len()
}
