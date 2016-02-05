package mmap

//generic
type KEY int

//generic
type VALUE int

type KeyMagic struct {
	counter int
	lastKey KEY
}

func (km *KeyMagic) doMagicOnKey(key KEY) {
	km.counter++
	km.lastKey = key
}

type Map struct {
	items    map[KEY]VALUE
	keyMagic KeyMagic
}

func New() *Map {
	return &Map{make(map[KEY]VALUE), KeyMagic{}}
}

func (m *Map) Put(key KEY, value VALUE) {
	m.items[key] = value
	m.keyMagic.doMagicOnKey(key)
}

func (m Map) Get(key KEY) VALUE {
	return m.items[key]
}
