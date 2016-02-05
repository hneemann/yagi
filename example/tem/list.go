package tem

//generic
type T int

// struct to store the elements
type List struct {
	items []T
}

func (l List) Items() []T {
	return l.items
}

func (l *List) Add(item T) {
	l.items = append(l.items, item)
}

func (l *List) len() int {
	return len(l.items)
}
