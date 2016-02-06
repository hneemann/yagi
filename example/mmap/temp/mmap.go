package temp

//generic
type KEY int

//generic
type VALUE int

// KeyMagic does some magic on the keys.
// In this example the insertions are counted and
// the last inserted key is stored.
type KeyMagic struct {
	counter int
	lastKey KEY
}

func (km *KeyMagic) doMagicOnKey(key KEY) {
	km.counter++
	km.lastKey = key
}

// Map holds the map and a KeyMagic struct
type Map struct {
	items    map[KEY]VALUE
	keyMagic KeyMagic
}

// New creates a new map
func New() *Map {
	return &Map{make(map[KEY]VALUE), KeyMagic{}}
}

// Put adds a key,value pair to the map
func (m *Map) Put(key KEY, value VALUE) {
	m.items[key] = value
	m.keyMagic.doMagicOnKey(key)
}

// Get a value from the map
func (m Map) Get(key KEY) VALUE {
	return m.items[key]
}
