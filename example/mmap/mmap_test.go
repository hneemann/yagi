package mmap

import "fmt"

//go:generate yagi -tem=./temp/mmap.go -gen=string,int64;string,string

func ExampleMMap() {
	{
		m := NewStringInt64()
		m.Put("first", 1)
		m.Put("second", 2)
		fmt.Println(m.Get("first"))
	}
	{
		m := NewStringString()
		m.Put("first", "1")
		m.Put("second", "2")
		fmt.Println(m.Get("first"))
	}
	// Output:
	// 1
	// 1
}
