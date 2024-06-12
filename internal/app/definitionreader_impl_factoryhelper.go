package app

import (
	"go/ast"

	annotation "github.com/YReshetko/go-annotation/pkg"
	"github.com/gosrob/autumn/internal/errorcode"
	"github.com/gosrob/autumn/internal/logger"
	"github.com/gosrob/autumn/internal/util/astutil"
	"github.com/gosrob/autumn/internal/util/logic"
	"github.com/gosrob/autumn/internal/util/nodeutil"
	"github.com/gosrob/autumn/internal/util/pkginfo"
	"github.com/gosrob/autumn/internal/util/stream"
	pkg "github.com/gosrob/autumn/pkg/annotation"
)

func extractFunc(n annotation.Node) (isBuiltinType bool, beanClass string, dependsOn []string, isPrimary bool, isLazy bool, alias string, bean FuncDefinition, pkgInfo Package, err error) {
	pkgInfo.FileAbsolutePath = n.Meta().Dir()
	pkgInfo.CurrentPackage = n.Meta().PackageName()
	pkgInfo.CurrentFullPackage = pkginfo.GetFullPackage(pkgInfo.FileAbsolutePath).ImportPath

	funcDecl := n.ASTNode().(*ast.FuncDecl)
	if len(funcDecl.Type.Results.List) > 1 {
		err = errorcode.FactoryCanOnlyProduceOneError.Instance().Printf(" produced %d", len(funcDecl.Type.Results.List))
		return
	}
	resultsList := funcDecl.Type.Results.List
	isBuiltinType = astutil.IsBasicType(resultsList[0].Type)

	ims := stream.Map(stream.OfSlice(n.Imports()), func(t *ast.ImportSpec) string { return t.Path.Value }).Map(stream.TrimQuote()).ToSlice()
	isAnotherPkg, path := astutil.IsAnotherPackageOrGetPkgPath(resultsList[0].Type, ims)

	pkgPath := logic.OrGet(isAnotherPkg, path, pkgInfo.CurrentFullPackage)
	tp := astutil.GetExprType(resultsList[0].Type)
	beanClass = astutil.BuildFullpathPackage(nodeutil.GetType(tp).TypeName, pkgPath)
	logger.Logger.Debugf("func pkgPath is %s", pkgPath)

	annotations := annotation.FindAnnotations[pkg.Bean](n.Annotations())
	if len(annotations) > 0 {
		ano := annotations[0]
		isPrimary = ano.IsPrimary == "true"
		isLazy = ano.IsLazy == "true"
		alias = ano.Alias
	}
	bean.Annotations = annotations
	bean.Name = extractFuncName(n)
	bean.Params = extractParam(n)
	bean.Results = extractResult(n)
	return
}

func extractFuncName(n annotation.Node) string {
	node := n.ASTNode()
	if fn, ok := node.(*ast.FuncDecl); ok {
		return fn.Name.Name
	}
	return ""
}

func extractParam(n annotation.Node) (params []Param) {
	node := n.ASTNode()
	ims := astutil.GetImportsStrings(n)

	// Assuming node represents a function declaration
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		for _, field := range funcDecl.Type.Params.List {
			tp := nodeutil.GetType(astutil.GetExprType(field.Type))
			isAnotherPkg, path := astutil.IsAnotherPackageOrGetPkgPath(field.Type, ims)

			pkgPath := logic.OrGet(isAnotherPkg, path, pkginfo.GetFullPackage(n.Meta().Dir()).ImportPath)
			tps := astutil.BuildFullpathPackage(nodeutil.GetType(tp.PureType).TypeName, pkgPath)

			for _, name := range field.Names {
				params = append(params, Param{
					Name: name.Name,
					Type: tps,
				})
			}
		}
	}

	return params
}

func extractResult(n annotation.Node) (results []Param) {
	node := n.ASTNode()
	ims := astutil.GetImportsStrings(n)

	// Assuming node represents a function declaration
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		for _, field := range funcDecl.Type.Results.List {
			tp := nodeutil.GetType(astutil.GetExprType(field.Type))
			isAnotherPkg, path := astutil.IsAnotherPackageOrGetPkgPath(field.Type, ims)

			pkgPath := logic.OrGet(isAnotherPkg, path, pkginfo.GetFullPackage(n.Meta().Dir()).ImportPath)
			tps := astutil.BuildFullpathPackage(nodeutil.GetType(tp.PureType).TypeName, pkgPath)

			results = append(results, Param{
				Name: "",
				Type: tps,
			})
		}
	}
	return
}
