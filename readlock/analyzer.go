// Package readlock --
package readlock

import (
	"errors"
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Doc explaining the tool.
const Doc = "Tool to recursive read locks in Go"

var errUnsafePackage = errors.New(
	"recursive read lock mutex detected",
)

// Analyzer runs static analysis.
var Analyzer = &analysis.Analyzer{
	Name:     "readlock",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, errors.New("analyzer is not type *inspector.Inspector")
	}

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	hasPendingMutex := false
	inspect.Preorder(nodeFilter, func(node ast.Node) {
		switch stmt := node.(type) {
		case *ast.CallExpr:
			if hasPendingMutex {
				ast.Inspect(node, func(node2 ast.Node) bool {
					switch stmt2 := node2.(type) {
					case *ast.CallExpr:
						if sel2, ok2 := stmt2.Fun.(*ast.SelectorExpr); ok2 && sel2.Sel.Name == "RLock" {
							fmt.Println("Found another call to RLock while not unlocked!", stmt2)
							pass.Reportf(
								node.Pos(),
								fmt.Sprintf(
									"%v",
									errUnsafePackage,
								),
							)
							return false
						}
					}
					return true
				})
			}
			if sel, ok := stmt.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "RLock" {
				fmt.Println("Found, setting pending mutex")
				hasPendingMutex = true
				return
			}
		}
	})

	return nil, nil
}

// func isPkgDot(expr ast.Expr, pkg, name string) bool {
// 	sel, ok := expr.(*ast.SelectorExpr)
// 	res := ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
// 	return res
// }

// func isIdent(expr ast.Expr, ident string) bool {
// 	id, ok := expr.(*ast.Ident)
// 	return ok && id.Name == ident
// }
