// this file defines a wrapper for largecode.List
// which is type save an does the boxing and unboxing
package wrap

import "github.com/hneemann/yagi/example/wrapper/largecode"

//generic
type ITEM int

type Wrapper struct {
	delegate largecode.List
}

// Add adds an element to the list
func (l *Wrapper) Add(item ITEM) {
	l.delegate.Add(item)
}

// Get returns an element from the list
func (l *Wrapper) Get(index int) ITEM {
	item, ok := l.delegate.Get(index).(ITEM)
	if ok {
		return item
	}
	panic("wrong type in list")
}

// Remove an element from the list
func (l *Wrapper) Remove(index int) {
	l.delegate.Remove(index)
}

// Len returns the number of elements in the list
func (l *Wrapper) Len() int {
	return l.delegate.Len()
}
