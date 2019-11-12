package funcdecl

import (
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "funcdecl",
	Doc:      `find function declarations`,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	inspect.Nodes([]ast.Node{
		(*ast.File)(nil), // 関数定義だけでなくファイルにも関心が出た
		(*ast.FuncDecl)(nil),
	}, func(n ast.Node, push bool) bool {
		if !push {
			return false
		}

		switch n := n.(type) {
		case *ast.File:
			f := pass.Fset.File(n.Pos())
			if strings.HasSuffix(f.Name(), "_test.go") {
				return false
			}

			return !generated(n) // ジェネレータで生成されたファイルのサブツリーには関心がない
		case *ast.FuncDecl:
			pass.Reportf(n.Pos(), `found %s`, n.Name)
			return false
		default:
			panic(n)
		}
	})

	return nil, nil
}

// https://github.com/golang/go/issues/13560#issuecomment-288457920
var pattern = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)

// ファイルのどこかに生成されたことを表すコメントがある
func generated(f *ast.File) bool {
	for _, c := range f.Comments {
		for _, l := range c.List {
			if pattern.MatchString(l.Text) {
				return true
			}
		}
	}
	return false
}
