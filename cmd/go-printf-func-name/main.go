package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"
	"selfLint/pkg/analyzer"
)

func main() {
	//example.Run()
	singlechecker.Main(analyzer.Analyzer)
}
