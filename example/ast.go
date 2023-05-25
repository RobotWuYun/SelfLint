package example

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"
)

// to be able to run this like "go run main.go -- input.go"
func Run() {
	v := Visitor{Fset: token.NewFileSet()}

	for _, filePath := range os.Args[1:] {
		if filePath == "--" {
			continue
		}

		f, err := parser.ParseFile(v.Fset, filePath, nil, 0)
		if err != nil {
			log.Fatalf("Failed to parse file %s: %s", filePath, err)
		}

		ast.Walk(&v, f)
	}
}

type Visitor struct {
	Fset *token.FileSet
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return v
	}

	params := funcDecl.Type.Params.List
	if len(params) != 2 { // [0] must be format (string), [1] must be args (...interface{})
		return v
	}

	firstParamType, ok := params[0].Type.(*ast.Ident)
	if !ok { // first param type isn't identificator so it can't be of type "string"
		return v
	}

	if firstParamType.Name != "string" { // first param (format) type is not string
		return v
	}

	secondParamType, ok := params[1].Type.(*ast.Ellipsis)
	if !ok { // args are not ellipsis (...args)
		return v
	}

	elementType, ok := secondParamType.Elt.(*ast.InterfaceType)
	if !ok { // args are not of interface type, but we need interface{}
		return v
	}

	if elementType.Methods != nil && len(elementType.Methods.List) != 0 {
		return v // has >= 1 method in interface, but we need an empty interface "interface{}"
	}

	if strings.HasSuffix(funcDecl.Name.Name, "f") {
		return v
	}

	fmt.Printf("%s: printf-like formatting function '%s' should be named '%sf'\n",
		v.Fset.Position(node.Pos()), funcDecl.Name.Name, funcDecl.Name.Name)
	return v
}
