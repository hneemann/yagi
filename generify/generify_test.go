package generify

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
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
	f := getFile(t, `package test
type (
	a int
	b int
)`)
	decls := splitDeclsToUngroupedDecls(f.Decls)
	s := getSource(t, &ast.File{Name: &ast.Ident{Name: "a"}, Decls: decls})
	assert.Equal(t, 1, strings.Count(s, "package a\n"))
	assert.Equal(t, 1, strings.Count(s, "type a int\n"))
	assert.Equal(t, 1, strings.Count(s, "type b int\n"))
}

func TestSplitDeclsToUngroupedDecls2(t *testing.T) {
	f := getFile(t, `package test
var (
	a int
	b int
)`)
	decls := splitDeclsToUngroupedDecls(f.Decls)
	s := getSource(t, &ast.File{Name: &ast.Ident{Name: "a"}, Decls: decls})
	assert.Equal(t, 1, strings.Count(s, "package a\n"))
	assert.Equal(t, 1, strings.Count(s, "var a int\n"))
	assert.Equal(t, 1, strings.Count(s, "var b int\n"))
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

	assert.Equal(t, 0, strings.Count(out, "KEY"), out)
	assert.Equal(t, 0, strings.Count(out, "VALUE"), out)
	assert.Equal(t, 1, strings.Count(out, "type strStringInt32 struct {\n"), out)
	assert.Equal(t, 1, strings.Count(out, "type strStringInt64 struct {\n"), out)
	assert.Equal(t, 1, strings.Count(out, "var aString string\n"), out)
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

	assert.Equal(t, 0, strings.Count(out, "VALUE"), out)
	assert.Equal(t, 1, strings.Count(out, "func minInt32(a, b int32) bool {\n"), out)
	assert.Equal(t, 1, strings.Count(out, "func minInt64(a, b int64) bool {\n"), out)
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

	assert.Equal(t, 0, strings.Count(out, "NUMBER"), out)
	assert.Equal(t, 1, strings.Count(out, "func DoSomethingInt32(fn func(a, b int32) bool, a, b int32) int32 {\n"), out)
	assert.Equal(t, 1, strings.Count(out, "func DoSomethingInt64(fn func(a, b int64) bool, a, b int64) int64 {\n"), out)
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

	assert.Equal(t, 0, strings.Count(out, "func (l *List) len() int"), out)
	assert.Equal(t, 1, strings.Count(out, "type ListInt32 struct"), out)
	assert.Equal(t, 1, strings.Count(out, "func (l *ListInt32) len() int {\n"), out)
	assert.Equal(t, 1, strings.Count(out, "type ListInt64 struct"), out)
	assert.Equal(t, 1, strings.Count(out, "func (l *ListInt64) len() int {\n"), out)
}
