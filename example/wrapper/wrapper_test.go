package wrapper

import "fmt"

//go:generate yagi -tem=./wrap/wrapper.go -gen=int64;string

func ExampleWrapper() {
	{
		m := WrapperInt64{}
		m.Add(1)
		m.Add(2)
		fmt.Println(m.Get(1))
	}
	{
		m := WrapperString{}
		m.Add("1")
		m.Add("2")
		fmt.Println(m.Get(1))
	}
	// Output:
	// 2
	// 2
}
