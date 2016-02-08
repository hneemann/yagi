package gmap

//go:generate yagi -tem=map.go -gen=string,int;string,int64;string,float64

//generic
type KEY interface{}

//generic
type VALUE interface{}

type Map map[KEY]VALUE

func (m Map) Get(key KEY) VALUE {
	return m[key]
}

func (m Map) Put(key KEY, value VALUE) {
	m[key] = value
}
