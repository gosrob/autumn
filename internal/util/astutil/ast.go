package astutil

import (
	"go/ast"

	"github.com/gosrob/autumn/internal/errorcode"
	"github.com/gosrob/autumn/internal/util/parser"
)

func AstCast[T any](node ast.Node) (T, error) {
	var zero T
	if v, ok := node.(T); ok {
		return v, nil
	}

	return zero, errorcode.CastError.Instance().Printf("cannot cast to ast type")
}

func IsBasicType(expr ast.Expr) bool {
	// Set of basic Go types
	basicTypes := map[string]struct{}{
		"int":        {},
		"int8":       {},
		"int16":      {},
		"int32":      {},
		"int64":      {},
		"uint":       {},
		"uint8":      {},
		"uint16":     {},
		"uint32":     {},
		"uint64":     {},
		"uintptr":    {},
		"float32":    {},
		"float64":    {},
		"complex64":  {},
		"complex128": {},
		"bool":       {},
		"string":     {},
		"rune":       {},
	}

	// Type assertion to *ast.Ident
	if ident, ok := expr.(*ast.Ident); ok {
		// Check if the identifier name matches any basic type
		_, exists := basicTypes[ident.Name]
		return exists
	}

	// Return false if expr is not an identifier or it is not a basic type
	return false
}

// check if an ast.Field is from another package
func IsAnotherPackage(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		if ident, ok := e.X.(*ast.Ident); ok {
			return ident.Obj == nil
		}
	}
	return false
}

func GetFieldPackageAndTypeName(f *ast.Field) (pkg string, name string) {
	ident, ok := f.Type.(*ast.Ident)
	if ok {
		name = ident.Name
	} else if sel, ok := f.Type.(*ast.SelectorExpr); ok {
		pkg = sel.X.(*ast.Ident).Name
		name = sel.Sel.Name
	}
	return
}

func BuildFullpathPackage(tp string, pkg string) string {
	return pkg + "." + tp
}

func AnnotationsByNode(n ast.Node) ([]parser.Annotation, bool) {
	c, ok := Comment(n)
	if !ok {
		return nil, false
	}

	a := panicOnError(parser.Parse(c))
	if len(a) == 0 {
		return nil, false
	}
	return a, true
}

func panicOnError[T any](v T, e error) T {
	if e != nil {
		panic(e)
	}
	return v
}
