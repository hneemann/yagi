##yagi

This a simple tool to add a template functionality to the go language.

###Motivation

I am a experienced programmer with a 10+ years Java background and I use Go as my "first" language since about 
three years. One of the Java features I used a lot are the generics available in Java. 
But the Java generics are complex and hard to understand. And they seem to overgrow: If you 
start to use them you get compiler warning very quickly. And to fix them you have to add generics to more and more 
classes. After a while there are corner cases which a very hard to fix, so you start to add `@SuppressWarnings` 
annotations to the code. I don't like the `@SuppressWarnings` annotation.
I remember situations where I had a really hard fight with the Java type system to get the generic types working 
without a compiler warning, and at the next day a small modification on the code starts the next hard fight.
Therefore, I can understand the go authors who do not want to add something like that to the Go language.

But sometimes you write a complex piece of code and then you realize that you can extract a data structure 
which some methods working on them. And that you could reuse the same structure on some other types. How to 
deal with such a situation without the usage of generics? Copy and Paste is not a good idea because the result is 
hard to maintain. Using the empty interface and throwing away all the compilers type checking is also not a good idea.
Or you can rewrite the code with the empty interface and implement a wrapper for each type which does all the 
boxing and unboxing from and to the empty interface?

Using templates like in C++ seems to be hard to debug if there goes something wrong because you can not see which
code is generated on the fly. (But I'm not a experienced C++ programmer so maybe I'm wrong!)
  
What I would like to do is to keep my code nearly unchanged and generate implementations for other types based 
on the existing code. 
If this file is generated I can inspect it in detail if something goes wrong. If there is no other 
way I can also modify the file to fix a problem. If I do so, the code is hard to maintain in the future, but it 
is possible. 

###Former work

There are a lot of different implementations out there and all are usable and are working.
Here are some of them and my understanding of how they work. I have to apologise if I did not understand 
them correctly: 

1. [gen](https://clipperhouse.github.io/gen/)
   `gen` is a tool which heps you generate code from a template. You have to implement this template in go.
   To add a new template you have to implement the `TypeWriter`-interface. this interface has a method with
   a `io.Writer` as a parameter. To this writer you hafe to write the template code.
   So creating a template is expensive and the templates are hard to test.

2. [genny](https://github.com/cheekybits/genny/)
   `genny` is also a tool to handle templates. With `genny` a template is a simple go file. So it is very cheap 
   to create and maintain. This go file has no extension's which break the code. 
   So the template can be compiled and tested in a idiomatic way.
   But it uses a simple text based find and replace technique to create the concrete types. So the template author
   has to take care about the names of the types, so that working code is generated.

3. [gonerics](http://bouk.co/blog/idiomatic-generics-in-go/)
   `goneric` uses the go packages to parse a go file and the the ast is traversed and the generic types are 
   renamed to the concrete types. The templates uses simple names like 'T' or 'U' which a renamed traversing 
   the ast. Also the author of the template has to take care about the names of structs and functions to generate 
   code which can live in one package.  

###The Idea

I found it a good idea not only to rename the typs, but also the structs and the functions which use this types.
So I can generate various structs and methods and all the code can live in the same package or even in the same 
file. So I implemented a generic rename tool which parses the ast looks for the generic types, looks which structs 
and functions use this types and rename also this struczs and functions in a propper way. The I can write the 
renamed ast to a file. And this can be done for every type I need and at the end I get a generated file which caontains 
all the needed declarations.
     
###Example

As always when talking about generics we start with a simple list:

```go
package tem

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

func (l *List) Len() int {
    return len(l.items)
}
```

To allow the code generator to identify the generic type the type is annotated with a special comment:

```go
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

func (l *List) Len() int {
    return len(l.items)
}
```

That's it. This is ideomatic go code which can be tested and used without any modification. It does not depend on 
any packages.

If we want to generate some other types we can invoke `yagi` by adding this go:generate statement to a 
file which lives in the parent directory of our list: 

    //go:generate yagi -tem=./tem/list.go -gen=int64;int32 -pkg=main

The `-gen` flag says that I need the list for the type `int64` and `int32` and the `-pkg` flag sets the package 
name to `main`.
Starting `go generate` from the command line we get:
 
```go
// generated by yagi. Don't modify this file!
// Any changes will be lost if this file is regenerated.

package main

type ListInt64 struct {
    items []int64
}

func (l ListInt64) Items() []int64 {
    return l.items
}

func (l *ListInt64) Add(item int64) {
    l.items = append(l.items, item)
}

func (l *ListInt64) Len() int {
    return len(l.items)
}

type ListInt32 struct {
    items []int32
}

func (l ListInt32) Items() []int32 {
    return l.items
}

func (l *ListInt32) Add(item int32) {
    l.items = append(l.items, item)
}

func (l *ListInt32) Len() int {
    return len(l.items)
}
```

As you can see not only the type `ITEM` is renamed, but also the method names are modified in a propper way.
But we can do something more complex:

Imagine a map which has a KEY and a VALUE type. And we want to do some magic on the keys:

```go
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
```

The type `KeyMagic` depends only on `KEY`, the type `Map` depends on `KEY` and `VALUE`. Again we can generate some 
concrete other types:

    //go:generate yagi -tem=./temp/mmap.go -gen=string,int64;string,string -pkg=main
    
We want to create a `<string,int64>` and a `<string,string>` Map. And this is what we get:

```go
// generated by yagi. Don't modify this file!
// Any changes will be lost if this file is regenerated.

package main

type KeyMagicString struct {
	counter int
	lastKey string
}

func (km *KeyMagicString) doMagicOnKey(key string) {
	km.counter++
	km.lastKey = key
}

type MapStringInt64 struct {
	items    map[string]int64
	keyMagic KeyMagicString
}

func NewStringInt64() *MapStringInt64 {
	return &MapStringInt64{make(map[string]int64), KeyMagicString{}}
}

func (m *MapStringInt64) Put(key string, value int64) {
	m.items[key] = value
	m.keyMagic.doMagicOnKey(key)
}

func (m MapStringInt64) Get(key string) int64 {
	return m.items[key]
}

type MapStringString struct {
	items    map[string]string
	keyMagic KeyMagicString
}

func NewStringString() *MapStringString {
	return &MapStringString{make(map[string]string), KeyMagicString{}}
}

func (m *MapStringString) Put(key string, value string) {
	m.items[key] = value
	m.keyMagic.doMagicOnKey(key)
}

func (m MapStringString) Get(key string) string {
	return m.items[key]
}
```

We get three types: `MapStringInt64` and `MapStringString` as expected and one type `KeyMagicString` which 
is shared by the other two types.

### A container/list example

Go comes with a implementation of a double linked list: [golang.org/pkg/container/list](https://golang.org/pkg/container/list/).
What is to do to generify this list implementation?

The implementation works directly on the empty interface. So I have to introduce a new type:

```go
//generic
type ITEM interface{}
```

After that I have to modify the code in a way that it uses this new type instead of the empty interface. 
There are seven usages of the empty interface and I can simply replace `interface{}` by `ITEM` in the file seven times.

After that the tests and the example comming with the list are still running without any modification. That's nice!

Now I can write the folowing code which is based on the list example given by the go authors:

```go
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
```

This code works as expected: The content of the list (1,2,3,4 ) is printed twice but 
the first part uses `int64` the second `string` so `PushBack`, `PushFront`, 
`InsertBefore` and `InsertAfter` are typed methods now. As you can see there are also 
two factory methods created: `NewInt64()` and `NewString()`.
  
### State of the Work

Here you can find a first implementation. Feel free to play around with the code. 
Up to now it's not tested on really complex code, so don't blame me if it does not 
work as expected. But I am happy about comments. 

One open issue is the handling of the comments in the template. It seems to me that they are 
bound to a fixed code position when they are parsed and stored in the ast. So if the type 
names became longer, the comments move around in the generated code. So at the moment the 
comments are simply removed from the code, which makes it harder to understand the
generated code.