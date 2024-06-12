package app

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"testing"

	"github.com/gosrob/autumn/internal/logger"
)

func TestExtractBeanDepends(t *testing.T) {
	defer logger.Logger.CatchPanic()
	filePath := "../../testdata/beandefinition/demo.go"
	src, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	imports := []*ast.ImportSpec{}
	ast.Inspect(node, func(n ast.Node) bool {
		if importSpec, ok := n.(*ast.ImportSpec); ok {
			imports = append(imports, importSpec)
		}
		return true
	})

	extractDepends(node, imports)
	extractFields(node, imports)
}
