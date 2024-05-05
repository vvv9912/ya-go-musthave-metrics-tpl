package linterOsExit

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"log"
)

var ErrLinterOsExit = &analysis.Analyzer{
	Name: "osCheckErr",
	Doc:  "osCheckErr",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	log.Println("starting osCheckErr")
	selector := func(xx *ast.SelectorExpr) {
		if ident, ok := xx.X.(*ast.Ident); ok {
			if ident.Name == "os" && xx.Sel.Name == "Exit" {
				pass.Reportf(xx.Pos(), "os.Exit() found in main")
			}
		}
	}

	funcdecl := func(x *ast.FuncDecl) {
		if x.Name.Name != "main" {
			return
		}

		for _, stmt := range x.Body.List {
			switch stmt := stmt.(type) {
			case *ast.ExprStmt:
				ast.Inspect(stmt, func(node ast.Node) bool {
					switch xx := node.(type) {
					case *ast.SelectorExpr:
						selector(xx)

					}
					return true
				})

			}
		}
		return
	}

	for _, f := range pass.Files {
		if f.Name.Name != "main" {
			continue
		}

		ast.Inspect(f, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.FuncDecl:
				funcdecl(x)
			}
			return true
		})
	}
	return nil, nil
}
