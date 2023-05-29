package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "selflint",
	Doc:      "Checks func format.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		funcDecl := node.(*ast.FuncDecl)
		checkFuncFormat(funcDecl)

		if isLeafFunc(funcDecl) {
			checkLeafFuncFormat(funcDecl)
		}
	})
	return nil, nil
}

// 检查是否是叶子函数
func isLeafFunc(funcDecl *ast.FuncDecl) bool {
	for _, v := range funcDecl.Body.List {
		if es, ok := v.(*ast.ExprStmt); ok {
			if _, isFunc := es.X.(*ast.CallExpr); isFunc {
				return false
			}
		}
	}
	return true
}

// 检查函数通用结构
func checkLeafFuncFormat(funcDecl *ast.FuncDecl) {}

// 检查叶子函数结构
func checkFuncFormat(funcDecl *ast.FuncDecl) {}
