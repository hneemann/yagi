##Generics for Go [![Build Status](https://travis-ci.org/hneemann/yagi.svg?branch=master)](https://travis-ci.org/hneemann/yagi)

This is a tool to add a simple template functionality to the Go language.

###Motivation

I am an experienced programmer with a 10+ years Java background and I use Go as my "first" language since about 
three years now. One of the Java features I used a lot are the generics available in Java. 
But the Java generics are complex and hard to understand. And they seem to overgrow: If you 
start to use them you get compiler warning very quickly. And to fix them you have to add generics to more and more 
classes. After a while there are corner cases which a very hard to fix. So you start to add `@SuppressWarnings` 
annotations to the code. I don't like the `@SuppressWarnings` annotation.
I remember situations where I had a really hard fight against the Java type system to get the generic 
types working without a compiler warning. 
And at the next day a small code modification introduces new compiler warnings.
Therefore, I can understand the go authors who do not want to add something like that to the Go language.

But sometimes you write a complex piece of code and then you realize that you can extract a data structure 
which some methods working on them. And that you could reuse the same structure on some other types. How to 
deal with such a situation without the usage of generics? Copy and Paste is not a good idea because the result is 
hard to maintain. Using the empty interface and throwing away all the compilers type checking is also not a good idea.
Or you can rewrite the code with the empty interface and implement a wrapper for each type which does all the 
boxing and unboxing from and to the empty interface. Also an error prone work.
  
What I would like to do is to keep my code nearly unchanged and generate implementations for other types based 
on the existing code. 
If this file is generated I can inspect it in detail if something goes wrong. If there is no other 
way I can also modify the file to fix a problem. This is not a good idea: The code is hard to maintain in the future! 
But it is possible if neccesary. 

###Former work

There are a lot of different implementations out there and all of them are usable and are working.
Here are some of them and my understanding of how they work. I have to apologise if I did not understand 
them correctly:

1. [gen](https://clipperhouse.github.io/gen/)
   is a tool which helps you to generate code from a template. You have to implement this template in Go.
   To add a new template you have to implement the `TypeWriter`-interface. This interface has a method with takes
   a `io.Writer` and the information about the concrete type as an argument. 
   To this `io.Writer` you have to write the generated concrete code.
   So creating a template is expensive and the templates are hard to test.

2. [genny](https://github.com/cheekybits/genny/)
   is also a tool to handle templates. With `genny` a template is a simple go file. So it is very cheap 
   to create and maintain. This go file has no extension's which break the code. 
   So the template can be compiled and tested in a idiomatic way.
   But it uses a simple text based search and replace technique to create the concrete types. So the template author
   has to take care about the names of the types and functions, to get the generated code working as expected.
   At the end the code looks somewhat strange.

3. [gonerics](http://bouk.co/blog/idiomatic-generics-in-go/)
   uses the go packages `go/parser` and `go/ast` to parse a go file. 
   The created ast is traversed and the generic types are renamed to the concrete types. 
   The templates are using simple names like 'T' or 'U' which a renamed by traversing the ast. 
   But the names of structs and functions are not touched. So if you want to generate code for more then one type
   the code has to live in different packages.

4. [goast](https://github.com/go-goast/goast/)
   uses a similar approach to adress the generics problem. So if you like `yagi` you should also take a look at
   [goast](https://github.com/go-goast/goast/)

###The Idea

I found it a good idea not only to rename the types, but also the structs and the functions which use this types if neccesary.
So I can generate various structs and methods and all the code can live in the same package or even in the same 
file. So I implemented a generic rename tool which parses the ast, looks for the generic types, looks which structs 
and functions are effected by this types and rename also the affected structs and functions in a propper way. Then 
the renamed ast is written to a file. And this can be done for every type I need and at the end I get a generated 
file which contains all the neccesary declarations.
     
###Example

Let us start with a simple list:

```go
package temp

type ITEM int

// List stores the elements
type List struct {
	items []ITEM
}

// Items returns the stored items
func (l List) Items() []ITEM {
	return l.items
}

// Add adds an element to the list
func (l *List) Add(item ITEM) {
	l.items = append(l.items, item)
}

// Len returns the number of elements in the list
func (l *List) Len() int {
	return len(l.items)
}
```

To allow the code generator to identify the generic type, the type needs to be annotated with a special comment:

```go
package temp

//generic
type ITEM int

// List stores the elements
type List struct {
	items []ITEM
}

// Items returns the stored items
func (l List) Items() []ITEM {
	return l.items
}

// Add adds an element to the list
func (l *List) Add(item ITEM) {
	l.items = append(l.items, item)
}

// Len returns the number of elements in the list
func (l *List) Len() int {
	return len(l.items)
}
```

That's it. This is ideomatic Go code which can be tested and used without any modification. It does not depend on 
any packages.

If we want to generate some other types we can invoke `yagi` by adding this go:generate statement to a 
file which lives in the parent directory of our list: 

    //go:generate yagi -tem=./temp/list.go -gen=int64;int32

The `-tem` flag points to the template, and 
the `-gen` flag says that I want to generate a list for the types `int64` and `int32`.
The package name of the generated file is set to the directory name of the generated 
file, so in most cases it will be ok. If you need an other name you can set it by `-pac=main`.
Before a file is written, it is checked if it already exists. If it exists, it is checked whether 
it was created by yagi. If not, you will get an error. 
So you can not overwrite a manualy created file by mistake.

Running `go generate` from the command line we get:
 
```go
// generated by yagi. Don't modify this file!
// Any changes will be lost if this file is regenerated.

package list

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

As you can see not only the type `ITEM` is renamed, but also the name of the structs are modified in a propper way.
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

The struct `KeyMagic` depends only on `KEY`, the struct `Map` depends on `KEY` and `VALUE`. Again we can generate some 
concrete other types:

    //go:generate yagi -tem=./temp/mmap.go -gen=string,int64;string,string
    
We want to create a `<string,int64>` and a `<string,string>` Map. And this is what we get:

```go
// generated by yagi. Don't modify this file!
// Any changes will be lost if this file is regenerated.

package mmap

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
And we get two factory methods (`NewStringInt64()` and `NewStringString()`) which will create the new types.

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

After that the tests and the example comming with the list, are still running without any modification. That's nice!

Now I can write the following code which is based on the list example given by the Go authors:

```go
package container

import "fmt"

//go:generate yagi -tem=./list/list.go -gen=int64;string

func ExampleList() {
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

This code works as expected: The content of the list (1,2,3,4) is printed twice but 
the first part uses `int64` the second `string`. So `PushBack`, `PushFront`, 
`InsertBefore` and `InsertAfter` are typed methods now. As you can see there are also 
two factory methods created: `NewInt64()` and `NewString()`.
You can find the generated code [here](https://github.com/hneemann/yagi/blob/master/example/container/list.go).

## Create Wrappers

If you don't like the code bloat that comes with the generation of such complete 
typed copys of the original code you can also create a template implementation of a 
wrapper for the original type. Then you only have to create typed wrappers. 
Imagine you have written a very complex implementation of a list which consists of a large amount of code.
This list (`largecode.List`) uses the empty interface to store the list items.
If you create a lot of typed copys of such a list you generate a large amount of 
mostly identical code. Maybe you do not want to do that.

So you can write a generic wrapper for the list:

```go
package wrap

import "github.com/hneemann/yagi/example/wrapper/largecode"

//generic
type ITEM int

type Wrapper struct {
	delegate largecode.List
}

// Add adds an element to the list
func (l *Wrapper) Add(item ITEM) {
	l.delegate.Add(item)
}

// Get returns an element from the list
func (l *Wrapper) Get(index int) ITEM {
	item, ok := l.delegate.Get(index).(ITEM)
	if ok {
		return item
	}
	panic("wrong type in list")
}

// Remove an element from the list
func (l *Wrapper) Remove(index int) {
	l.delegate.Remove(index)
}

// Len returns the number of elements in the list
func (l *Wrapper) Len() int {
	return l.delegate.Len()
}
```
The generation of wrappers is somewhat tricky because the type `Wrapper` does
not depend on a generic type. There are only some methods witch have `Wrapper` 
as an receiver which depend on the generic type.

Now you can generate type save wrappers for `largecode.List` and use the
type save wrappers instead of `largecode.List` itself:

```go
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
}
```

Again the methods `Add` and `Get` are typed now.
You can find the generated code [here](https://github.com/hneemann/yagi/blob/master/example/wrapper/wrapper.go).
  
### State of the Work

Here you can find a first implementation. Feel free to play around with the code. 
Up to now it's not tested on really complex code, so don't blame me if it does not 
work as expected. But I am happy about comments. 

One open issue is the handling of the comments in the template. It seems to me that they are 
bound to a fixed code position when they are parsed and stored in the ast. So if the type 
names became longer, the comments move around in the generated code. So at the moment the 
comments are simply removed from the generated code, which makes the code harder to read.

Other open issues are the special properties of the types: You can write a template which compares
two values to check whichever is greater. If you replace the template type by a struct you will
get compile time errors because structs are not comparable in that way.

---
I have to apologise my poor english.