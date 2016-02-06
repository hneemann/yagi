package main

import "fmt"

//go:generate yagi -tem=./temp/mmap.go -gen=string,int64;string,string -pkg=main

func main() {
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
}
