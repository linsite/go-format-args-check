package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name:     "goformatargs",
	Doc:      "Checks that printf-like functions have matched %s and args.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var checkMap map[string]bool = map[string]bool{
	"Printf":  true,
	"Sprintf": true,
	"Infof":   true,
	"Debugf":  true,
	"Warnf":   true,
	"Errorf":  true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspector.Preorder(nodeFilter, func(node ast.Node) {
		n, ok := node.(*ast.CallExpr)
		if !ok {
			return
		}
		fun, ok := n.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}
		if !checkMap[fun.Sel.Name] {
			return
		}
		if len(n.Args) == 0 {
			pass.Reportf(node.Pos(), "formatting function '%s' args shouldn't be 0", fun.Sel.Name)
			return
		}
		lit, ok := n.Args[0].(*ast.BasicLit)
		if !ok {
			return
		}
		left := strings.ReplaceAll(lit.Value, "%%", "")
		if strings.Count(left, "%") != len(n.Args)-1 {
			pass.Reportf(node.Pos(), "formatting function '%s' args should match %% count", fun.Sel.Name)
			return
		}
	})

	return nil, nil
}
