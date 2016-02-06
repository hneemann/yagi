package temp

//generic
type ITEM int

// List stores the elements
type List struct {
	items []ITEM
}

// Items returns the stored items
func (l List) Items() []ITEM {
	return l.items
}

// Add adds an element to the list
func (l *List) Add(item ITEM) {
	l.items = append(l.items, item)
}

// Len returns the number of elements in the list
func (l *List) Len() int {
	return len(l.items)
}
