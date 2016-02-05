package tem

//generic
type KEY int

//generic
type VALUE int

type Map struct {
	items map[KEY]VALUE
}

func New() *Map {
	return &Map{make(map[KEY]VALUE)}
}

func (m *Map) Put(key KEY, value VALUE) {
	m.items[key] = value
}

func (m Map) Get(key KEY) VALUE {
	return m.items[key]
}
