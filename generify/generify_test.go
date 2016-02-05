package generify

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"testing"

	"github.com/hneemann/yagi/concrete"

	"github.com/stretchr/testify/assert"
)

func getFile(t *testing.T, code string) *ast.File {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	assert.NoError(t, err)
	return file
}

func getSource(t *testing.T, n ast.Node) string {
	fset := token.NewFileSet()
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, n)
	assert.NoError(t, err)
	return buf.String()
}

func TestSplitDeclsToUngroupedDecls1(t *testing.T) {
	f := getFile(t, `package a
type (
	a int
	b int
)		
`)
	decls := splitDeclsToUngroupedDecls(f.Decls)
	assert.Equal(t, `package a

type a int
type b int
`, getSource(t, &ast.File{Name: &ast.Ident{Name: "a"}, Decls: decls}))
}

func TestSplitDeclsToUngroupedDecls2(t *testing.T) {
	f := getFile(t, `package a
var (
	a int
	b int
)		
`)
	decls := splitDeclsToUngroupedDecls(f.Decls)
	assert.Equal(t, `package a

var a int
var b int
`, getSource(t, &ast.File{Name: &ast.Ident{Name: "a"}, Decls: decls}))
}

func gen(t *testing.T, code string, types string) string {
	file := getFile(t, code)

	c, err := concrete.New(types)
	assert.NoError(t, err)

	gen := New(file, c)
	assert.NoError(t, err)
	var buf bytes.Buffer
	err = gen.Do("", &buf)
	assert.NoError(t, err)

	return buf.String()
}

func TestDoubleGeneration(t *testing.T) {
	out := gen(t, `package test

//generic
type KEY int
//generic
type VALUE int

var a KEY

type str struct {
	key KEY
	value VALUE
}`, "string,int32;string,int64")

	assert.Equal(t, `package test


var aString string

type strStringInt32 struct {
	key	string
	value	int32
}

type strStringInt64 struct {
	key	string
	value	int64
}

`, out)
}

func TestSimpleFunction(t *testing.T) {
	out := gen(t, `package test

//generic
type VALUE int

func min(a,b VALUE) bool {
	if a<b {
		return a
	}
	return b
}`, "int32;int64")

	assert.Equal(t, `package test


func minInt32(a, b int32) bool {
	if a < b {
		return a
	}
	return b
}

func minInt64(a, b int64) bool {
	if a < b {
		return a
	}
	return b
}

`, out)
}

func TestFunction(t *testing.T) {
	out := gen(t, `package test

//generic
type NUMBER int

func DoSomething(fn func(a, b NUMBER) bool, a, b NUMBER) NUMBER {
	if fn(a, b) {
		return a
	}
	return b
}`, "int32;int64")

	assert.Equal(t, `package test


func DoSomethingInt32(fn func(a, b int32) bool, a, b int32) int32 {
	if fn(a, b) {
		return a
	}
	return b
}

func DoSomethingInt64(fn func(a, b int64) bool, a, b int64) int64 {
	if fn(a, b) {
		return a
	}
	return b
}

`, out)
}

func TestFunctionExpand(t *testing.T) {
	out := gen(t, `package test

//generic
type ITEM int

type List struct {
	items []ITEM
}

func (l *List) len() int {
	return len(l.items)
}
`, "int32;int64")

	assert.Equal(t, `package test


type ListInt32 struct{ items []int32 }

func (l *ListInt32) len() int {
	return len(l.items)
}

type ListInt64 struct{ items []int64 }

func (l *ListInt64) len() int {
	return len(l.items)
}

`, out)
}
