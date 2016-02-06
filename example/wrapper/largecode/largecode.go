// Imagine this is a large codebase dealing with ITEMS
// and you need a lot of concrete types and you find
// it to expensive to generate a lot of typed copies
// of all the code in this file.
package largecode

// List stores the elements
type List struct {
	items []interface{}
}

// Add adds an element to the list
func (l *List) Add(item interface{}) {
	l.items = append(l.items, item)
}

// Get gets an element to the list
func (l *List) Get(index int) interface{} {
	return l.items[index]
}

// Remove an element from the list
func (l *List) Remove(index int) {
	l.items = append(l.items[:index], l.items[index+1:]...)
}

// Len returns the number of elements in the list
func (l *List) Len() int {
	return len(l.items)
}
