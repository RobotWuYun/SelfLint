package main

import (
	"github.com/jirfag/go-printf-func-name/pkg/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	//example.Run()
	singlechecker.Main(analyzer.Analyzer)
}
