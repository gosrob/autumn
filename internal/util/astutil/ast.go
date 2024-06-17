package astutil

import (
	"go/ast"
	"path"

	annotation "github.com/YReshetko/go-annotation/pkg"
	"github.com/gosrob/autumn/internal/errorcode"
	"github.com/gosrob/autumn/internal/util/nodeutil"
	"github.com/gosrob/autumn/internal/util/parser"
	"github.com/gosrob/autumn/internal/util/stream"
)

func AstCast[T any](node ast.Node) (T, error) {
	var zero T
	if v, ok := node.(T); ok {
		return v, nil
	}

	return zero, errorcode.CastError.DeepCopy().APrintf("cannot cast to ast type").Printf(nodeutil.GetTypeFromAny(node), nodeutil.GetTypeFromAny(zero))
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

func IsInterface(expr ast.Expr) bool {
	_, ok := expr.(*ast.InterfaceType)
	return ok
}

// check if an ast.Field is from another package
func IsAnotherPackage(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.SelectorExpr:
		if ident, ok := e.X.(*ast.Ident); ok {
			return ident.Obj == nil
		}
	case *ast.StarExpr:
		return IsAnotherPackage(e.X)
	}
	return false
}

func IsAnotherPackageOrGetPkgPath(expr ast.Expr, ims []string) (bool, string) {
	if !IsAnotherPackage(expr) {
		return false, ""
	}
	ims = stream.OfSlice(ims).Map(stream.TrimQuote()).ToSlice()

	switch v := expr.(type) {
	case *ast.SelectorExpr:
		if ident, ok := v.X.(*ast.Ident); ok {
			for _, im := range ims {
				if path.Base(im) == "." {
					return false, ""
				}
				if ident.Name == path.Base(im) {
					return true, im
				}
			}
		}
	case *ast.StarExpr:
		if sel, ok := v.X.(*ast.SelectorExpr); ok {
			if ident, ok := sel.X.(*ast.Ident); ok {
				for _, im := range ims {
					if path.Base(im) == "." {
						return false, ""
					}
					if ident.Name == path.Base(im) {
						return true, im
					}
				}
			}
		}
	}
	return false, ""
}

func GetFieldPackageAndTypeName(f *ast.Field) (pkg string, name string) {
	ident, ok := f.Type.(*ast.Ident)
	if ok {
		name = ident.Name
	} else if sel, ok := f.Type.(*ast.SelectorExpr); ok {
		pkg = sel.X.(*ast.Ident).Name
		name = sel.Sel.Name
	} else if arrayType, ok := f.Type.(*ast.ArrayType); ok {
		if ident, ok := arrayType.Elt.(*ast.Ident); ok {
			name = ident.Name
		} else if sel, ok := arrayType.Elt.(*ast.SelectorExpr); ok {
			pkg = sel.X.(*ast.Ident).Name
			name = sel.Sel.Name
		}
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

func GetExprType(tp ast.Expr) string {
	switch v := tp.(type) {
	case *ast.Ident:
		return v.Name
	case *ast.StarExpr:
		return "*" + GetExprType(v.X)
	case *ast.SelectorExpr:
		return GetExprType(v.X) + "." + v.Sel.Name
	case *ast.ArrayType:
		return "[]" + GetExprType(v.Elt)
	case *ast.MapType:
		return "map[" + GetExprType(v.Key) + "]" + GetExprType(v.Value)
	case *ast.StructType:
		if len(v.Fields.List) > 0 {
			if ident, ok := v.Fields.List[0].Type.(*ast.Ident); ok {
				return ident.Name
			}
		}
		return "struct"
	default:
		return ""
	}
}

func GetImportsStrings(n annotation.Node) []string {
	return stream.Map(
		stream.OfSlice(n.Imports()),
		func(t *ast.ImportSpec) string { return t.Path.Value }).
		Filter(stream.Distinct[string]()).
		ToSlice()
}
