// Package generify holds the code to parse the ast of the given template
package generify

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"io"
	"strings"

	"github.com/hneemann/yagi/concrete"
	"github.com/hneemann/yagi/generify/intset"
)

var newline = []byte("\n\n")

func checkNodeIsOffType(n ast.Node, kind ast.ObjKind) (*ast.Ident, bool) {
	if ident, ok := n.(*ast.Ident); ok {
		return ident, ident.Obj != nil && ident.Obj.Kind == kind
	}
	return nil, false
}

type declWithDependency struct {
	decl             ast.Decl
	usedTypes        intset.IntSet
	writtenInstances []concrete.Types
}

func (dwd declWithDependency) String() string {
	buffer := new(bytes.Buffer)
	fset := token.NewFileSet()
	printer.Fprint(buffer, fset, dwd.decl)
	return fmt.Sprintf("used: %v\n%v", dwd.usedTypes, buffer)
}

func (dwd *declWithDependency) isAllreadyWritten(types concrete.Types) bool {
	for _, wi := range dwd.writtenInstances {
		same := 0
		for i := range dwd.usedTypes {
			if wi[i] == types[i] {
				same++
			}
		}
		if same == len(dwd.usedTypes) {
			return true
		}
	}
	dwd.writtenInstances = append(dwd.writtenInstances, types)
	return false
}

// Generify holds the data used to work on the ast
type Generify struct {
	// the parsed original template
	file *ast.File
	// the concrete types for which the code is generated
	concreteTypes *concrete.Instances
	// the name of the generic types
	genTypes []string
	// all the declarations from the template
	genericDecls []*declWithDependency
	// list of rename actions which are to perform on the ast to get a concrete type
	renameActions []renameAction
}

type renameAction interface {
	rename(concrete.Types)
}

func (g *Generify) addRenameAction(renameAction renameAction) {
	g.renameActions = append(g.renameActions, renameAction)
}

// New creates a new Generify instance
func New(file *ast.File, concreteTypes *concrete.Instances) *Generify {
	return &Generify{file: file, concreteTypes: concreteTypes}
}

// Do creates a concrete ast from the generic one and writes it to the given io.Writer
func (g *Generify) Do(packageName string, w io.Writer) error {
	var decls []ast.Decl
	g.genTypes, decls = findGenerics(g.file)
	if len(g.genTypes) == 0 {
		return fmt.Errorf("no generic types found")
	}
	if len(g.genTypes) != len(g.concreteTypes.Instance[0]) {
		return fmt.Errorf("there are %d generic types but %d concrete types", len(g.genTypes), len(g.concreteTypes.Instance[0]))
	}

	removeCommentsFrom(decls)

	decls = splitDeclsToUngroupedDecls(decls)

	g.genericDecls = g.inspectAllDeclsForDependencies(decls)

	g.renameStructsAndVars()

	g.renameFunctions()

	file := ast.File{Name: g.file.Name, Decls: g.staticDecls(), Scope: g.file.Scope, Imports: g.file.Imports}
	if packageName != "" {
		file.Name.Name = packageName
	}

	fset := token.NewFileSet()
	err := printer.Fprint(w, fset, &file)
	if err != nil {
		return err
	}
	w.Write(newline)

	for _, types := range g.concreteTypes.Instance {

		// rename all identifiers
		for _, ra := range g.renameActions {
			ra.rename(types)
		}

		// write the renamed ast
		for _, decl := range g.genericDecls {
			if len(decl.usedTypes) > 0 && !decl.isAllreadyWritten(types) {
				err := printer.Fprint(w, fset, decl.decl)
				if err != nil {
					return err
				}
				w.Write(newline)
			}
		}
	}

	return nil
}

func (g *Generify) staticDecls() []ast.Decl {
	decls := []ast.Decl{}
	for _, d := range g.genericDecls {
		if len(d.usedTypes) == 0 {
			decls = append(decls, d.decl)
		}
	}
	return decls
}

type declVisitor interface {
	ast.Visitor
	finalize(*declWithDependency)
}

func (g *Generify) walk(v declVisitor) {
	for _, decl := range g.genericDecls {
		ast.Walk(v, decl.decl)
		v.finalize(decl)
	}
}

// find type declarations which are marked with the "generic" comment
// returns this type names and the remaining declarations
func findGenerics(file *ast.File) ([]string, []ast.Decl) {
	var genTypes []string

	newDecls := []ast.Decl{}
	for _, d := range file.Decls {
		remove := false
		if gd, ok := d.(*ast.GenDecl); ok {
			if gd.Tok == token.TYPE {
				if strings.TrimSpace(gd.Doc.Text()) == "generic" {
					if len(gd.Specs) == 1 {
						if ts, ok := gd.Specs[0].(*ast.TypeSpec); ok {
							// pick the generic type
							genTypes = append(genTypes, ts.Name.Name)
							remove = true
						}
					}
				}
			}
		}
		if !remove {
			newDecls = append(newDecls, d)
		}
	}

	return genTypes, newDecls
}

func removeCommentsFrom(decls []ast.Decl) {
	for _, decl := range decls {
		ast.Inspect(decl, func(n ast.Node) bool {
			switch d := n.(type) {
			case *ast.GenDecl:
				d.Doc = nil
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.ImportSpec:
						s.Doc = nil
						s.Comment = nil
					case *ast.ValueSpec:
						s.Doc = nil
						s.Comment = nil
					case *ast.TypeSpec:
						s.Doc = nil
						s.Comment = nil
					}
				}
			case *ast.FuncDecl:
				d.Doc = nil
			case *ast.StructType:
				for _, field := range d.Fields.List {
					field.Doc = nil
					field.Comment = nil
				}
			}
			return true
		})
	}
}

func splitDeclsToUngroupedDecls(decls []ast.Decl) []ast.Decl {
	newDecls := []ast.Decl{}
	for _, decl := range decls {
		copyDecl := true
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			if len(genDecl.Specs) > 1 {
				copyDecl = false
				for _, spec := range genDecl.Specs {
					gd := ast.GenDecl{Tok: genDecl.Tok, Specs: []ast.Spec{spec}}
					newDecls = append(newDecls, &gd)
				}
			}
		}
		if copyDecl {
			newDecls = append(newDecls, decl)
		}
	}
	return newDecls
}

type simpleVisitor struct {
	g          *Generify
	foundTypes intset.IntSet
}

func newSimpleVisitor(g *Generify) *simpleVisitor {
	return &simpleVisitor{g, make(intset.IntSet)}
}

type simpleRename struct {
	ident    *ast.Ident
	genIndex int
}

func (ir simpleRename) rename(t concrete.Types) {
	ir.ident.Name = t[ir.genIndex]
}

func (sv *simpleVisitor) Visit(n ast.Node) ast.Visitor {
	if id, ok := checkNodeIsOffType(n, ast.Typ); ok {
		for i, gen := range sv.g.genTypes {
			if id.Name == gen {
				sv.g.addRenameAction(simpleRename{id, i})
				sv.foundTypes[i] = struct{}{}
			}
		}
	}
	return sv
}

func (g *Generify) inspectAllDeclsForDependencies(decls []ast.Decl) []*declWithDependency {
	newDecls := []*declWithDependency{}
	for _, decl := range decls {
		sv := newSimpleVisitor(g)
		ast.Walk(sv, decl)
		newDecls = append(newDecls, &declWithDependency{decl, sv.foundTypes, nil})
	}
	return newDecls
}

type renameVisitor struct {
	g           *Generify
	kind        ast.ObjKind
	origName    string
	usedIndices intset.IntSet
	wasActive   bool
}

func (rv *renameVisitor) Visit(n ast.Node) ast.Visitor {
	if id, ok := checkNodeIsOffType(n, rv.kind); ok {
		if id.Name == rv.origName {
			rv.g.addRenameAction(multiRename{rv.origName, id, rv.usedIndices})
			rv.wasActive = true
		}
	}
	return rv
}

func (rv *renameVisitor) finalize(d *declWithDependency) {
	if rv.wasActive {
		d.usedTypes.AddAll(rv.usedIndices)
	}
	rv.wasActive = false
}

type multiRename struct {
	origName    string
	ident       *ast.Ident
	usedIndices intset.IntSet
}

func (mr multiRename) rename(ct concrete.Types) {
	// ToDo: this operation is done over and over again!
	n := mr.origName
	for i, conName := range ct {
		if _, ok := mr.usedIndices[i]; ok {
			conName = strings.Replace(conName, "*", "P", -1)
			conName = strings.Replace(conName, ".", "", -1)
			n += strings.Title(conName)
		}
	}
	mr.ident.Name = n
}

func (g *Generify) renameStructsAndVars() {
	for _, decl := range g.genericDecls {
		if genDecl, ok := decl.decl.(*ast.GenDecl); ok {
			switch spec := genDecl.Specs[0].(type) {
			case *ast.TypeSpec:
				g.walk(&renameVisitor{g, ast.Typ, spec.Name.Name, decl.usedTypes, false})
			case *ast.ValueSpec:
				for _, name := range spec.Names {
					g.walk(&renameVisitor{g, ast.Var, name.Name, decl.usedTypes, false})
				}
			}
		}
	}
}

func (g *Generify) renameFunctions() {
	for _, decl := range g.genericDecls {
		if funcDecl, ok := decl.decl.(*ast.FuncDecl); ok {
			g.walk(&renameVisitor{g, ast.Fun, funcDecl.Name.Name, decl.usedTypes, false})
		}
	}
}
