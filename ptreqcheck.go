package ptreqcheck

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"log"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "ptreqcheck",
	Doc:      "check for two pointers compared with the '==' operator",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeTypes := []ast.Node{
		(*ast.BinaryExpr)(nil),
	}

	inspect.Preorder(nodeTypes, check(pass))

	return nil, nil
}

// check contains the logic for checking that '==' is used correctly in the code being analysed
func check(pass *analysis.Pass) func(ast.Node) {
	return func(node ast.Node) {
		expr := node.(*ast.BinaryExpr)
		// we are only interested in comparison
		if expr.Op != token.EQL {
			return
		}

		// get the types of the two operands
		x, xOK := pass.TypesInfo.Types[expr.X]
		y, yOK := pass.TypesInfo.Types[expr.Y]

		if !xOK || !yOK {
			return
		}

		if isPointer(x.Type) && isPointer(y.Type) {
			pass.Reportf(expr.Pos(), "Comparison of pointer types: `%s`", formatNode(expr))
		}
	}
}

func isPointer(x types.Type) bool {
	return x.String()[0] == '*'
}

func formatNode(node ast.Node) string {
	buf := new(bytes.Buffer)
	if err := format.Node(buf, token.NewFileSet(), node); err != nil {
		log.Printf("Error formatting expression: %v", err)
		return ""
	}

	return buf.String()
}
