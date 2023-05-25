package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"selfLint/example"
)

func main() {
	exampleAst()
}

// to be able to run this like "go run main.go -- input.go"
func exampleAst() {
	v := example.Visitor{Fset: token.NewFileSet()}

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
