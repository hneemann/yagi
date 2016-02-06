package main

import "fmt"

//go:generate yagi -tem=./list/list.go -gen=int64;string -pkg=main

func main() {
	{
		l := NewInt64()
		e4 := l.PushBack(4)
		e1 := l.PushFront(1)
		l.InsertBefore(3, e4)
		l.InsertAfter(2, e1)

		// Iterate through list and print its contents.
		for e := l.Front(); e != nil; e = e.Next() {
			fmt.Println(e.Value)
		}
	}
	{
		l := NewString()
		e4 := l.PushBack("4")
		e1 := l.PushFront("1")
		l.InsertBefore("3", e4)
		l.InsertAfter("2", e1)

		// Iterate through list and print its contents.
		for e := l.Front(); e != nil; e = e.Next() {
			fmt.Println(e.Value)
		}
	}
}
