package tem

//generic
type ITEM int

// struct to store the elements
type List struct {
	items []ITEM
}

func (l List) Items() []ITEM {
	return l.items
}

func (l *List) Add(item ITEM) {
	l.items = append(l.items, item)
}

func (l *List) len() int {
	return len(l.items)
}
