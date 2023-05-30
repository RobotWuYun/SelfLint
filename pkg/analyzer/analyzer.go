package analyzer

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
	"strings"
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

		var hasTrace bool

		hasTrace = hasTrace || checkFuncDoc(funcDecl)

		hasTrace = hasTrace || checkFuncFormat(funcDecl)

		if isLeafFunc(funcDecl) {
			checkLeafFuncFormat(funcDecl)
		}

		if !hasTrace {
			pass.Reportf(funcDecl.Pos(), "func %s has not trace", funcDecl.Name.Name)
		}

	})
	return nil, nil
}

// isLeafFunc 检查是否是叶子函数
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

// checkLeafFuncFormat 检查叶子函数结构
func checkLeafFuncFormat(funcDecl *ast.FuncDecl) {}

// checkFuncFormat 检查函数通用结构
func checkFuncFormat(funcDecl *ast.FuncDecl) (hasTrace bool) {
	// check trace
	for _, v := range funcDecl.Body.List {
		if as, ok := v.(*ast.AssignStmt); ok {
			for _, asChild := range as.Rhs {
				if asCF, ok := asChild.(*ast.CallExpr); ok {
					if _, isFunc := asCF.Fun.(*ast.SelectorExpr); isFunc {
						hasTrace = hasTrace || checkHasTrace(asCF)
					}
				}
			}
		} else if es, ok := v.(*ast.ExprStmt); ok {
			if ce, isFunc := es.X.(*ast.CallExpr); isFunc {
				hasTrace = hasTrace || checkHasTrace(ce)
			}
		}
	}
	return
}

func checkFuncDoc(funcDecl *ast.FuncDecl) (hasTrace bool) {
	return strings.Contains(funcDecl.Doc.Text(), "ignore trace")
}

// checkHasTrace 检查是否有trace
func checkHasTrace(funcExpr *ast.CallExpr) bool {
	if _, isFunc := funcExpr.Fun.(*ast.SelectorExpr); isFunc {
		return strings.Contains(funcExpr.Fun.(*ast.SelectorExpr).X.(*ast.Ident).Name, "trace")
	}
	return false
}
