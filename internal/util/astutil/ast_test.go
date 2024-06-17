package astutil_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"testing"

	"github.com/gosrob/autumn/internal/logger"
	"github.com/gosrob/autumn/internal/util/astutil"
)

func TestParseCommentToAnnotation(t *testing.T) {
	logger.Logger.SetIsDebug(true)
	defer logger.Logger.CatchPanic()
	filePath := "../../../examples/beandefinition/demo.go"
	src, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse source: %v", err)
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.GenDecl:
			for _, spec := range n.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					// docs := n.Doc // Get the attached comments/docs.
					// logger.Logger.Debugf("annotations for %s with comments %+v", ts.Name.Name, docs)
					// Assuming astutil.AnnotationsByNode can interpret the docs
					as, err := astutil.AnnotationsByNode(n)
					logger.Logger.Debugf("annotations for %s is %+v, %+v", ts.Name.Name, as, err)
					for _, v := range ts.Type.(*ast.StructType).Fields.List {
						as, err := astutil.AnnotationsByNode(v)
						logger.Logger.Debugf("annotations for %s is %+v, %+v", ts.Name.Name, as, err)
					}
				}
			}
		}
		return true
	})
}
