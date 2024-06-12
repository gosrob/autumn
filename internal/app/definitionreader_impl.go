package app

import (
	"fmt"
	"go/ast"
	"strings"

	annotation "github.com/YReshetko/go-annotation/pkg"
	"github.com/gosrob/autumn/internal/logger"
	"github.com/gosrob/autumn/internal/util"
	"github.com/gosrob/autumn/internal/util/astutil"
	"github.com/gosrob/autumn/internal/util/cast"
	"github.com/gosrob/autumn/internal/util/nodeutil"
	"github.com/gosrob/autumn/internal/util/parser"
	pkg "github.com/gosrob/autumn/pkg/annotation"
	"github.com/samber/lo"
)

type goAnnotationBeanDefinitionReader struct{}

// ReadBeanFactoryDefinition implements BeanDefinitionReader.
func (g *goAnnotationBeanDefinitionReader) ReadBeanFactoryDefinition(n any) (FactoryFuncDefinition, error) {
	panic("unimplemented")
}

var GoAnnotationBeanDefinitionReader goAnnotationBeanDefinitionReader

// ReadBeanDefinition implements BeanDefinitionReader.
func (g *goAnnotationBeanDefinitionReader) ReadBeanDefinition(n any) (BeanDefinition, error) {
	var b BeanDefinition
	node, err := cast.Cast[annotation.Node](n)
	if err != nil {
		return b, err
	}

	beanClass, isBasicType, dependsOn, isPrimary, isLazy, alias, bean, err := extractNode(node)
	if err != nil {
		return b, err
	}
	b = BeanDefinition{
		IsBuiltinType: isBasicType,
		BeanClass:     beanClass,
		DependsOn:     dependsOn,
		IsPrimary:     isPrimary,
		IsLazy:        isLazy,
		Alias:         alias,
		Bean:          bean,
	}
	return b, nil
}

var _ BeanDefinitionReader = (*goAnnotationBeanDefinitionReader)(nil)

func extractNode(n annotation.Node) (beanClass BeanClass, isBasicType bool, dependsOn []string, isPrimary bool, isLazy bool, alias string, bean StructDefinition, err error) {
	var node annotation.Node
	node, err = cast.Cast[annotation.Node](n)
	if err != nil {
		logger.Logger.Warnf("extract node faild err:%s", err)
	}
	typeSpec := node.ASTNode().(*ast.TypeSpec)
	logger.Logger.Debugf("type spec is %+v", typeSpec)

	// isBasicType
	isBasicType = astutil.IsBasicType(node.ASTNode().(ast.Expr))
	// beanClass
	beanClass = BeanClass(node.Meta().PackageName() + "." + nodeutil.GetType(typeSpec.Name.Name).PureType)

	// dependsOn
	dependsOn, err = extractDepends(n.ASTNode(), n.Imports())

	annotations := annotation.FindAnnotations[pkg.Bean](node.Annotations())
	if len(annotations) > 0 {
		ano := annotations[0]
		isPrimary = ano.IsPrimary == "true"
		isLazy = ano.IsLazy == "true"
		alias = ano.Alias
	}
	bean.Annotations = annotations
	tp, err := extractStructType(node.ASTNode())
	if err == nil {
		bean.Name = astutil.BuildFullpathPackage(tp, node.Meta().PackageName())
	}
	fds, err := extractFields(n.ASTNode(), n.Imports())
	if err == nil {
		bean.Fields = fds
	}
	return
}

func extractStructType(spec ast.Node) (tp string, err error) {
	typeSpec, ok := spec.(*ast.TypeSpec)
	if !ok {
		return "", fmt.Errorf("params is not type spec")
	}

	if _, ok := typeSpec.Type.(*ast.StructType); ok {
		return typeSpec.Name.Name, nil
	}

	return "", fmt.Errorf("no struct type found")
}

func extractFields(n ast.Node, ims []*ast.ImportSpec) (fields []Field, err error) {
	visitor := astutil.NewVisitor(func(ns ast.Node) {
		logger.Logger.Debugf("start extractDepend %+v", n)
		if f, ok := ns.(*ast.Field); ok {
			field, err := extractField(f, ims)
			if err != nil {
				return
			}
			fields = append(fields, field)
		}
	})
	ast.Walk(&visitor, n)
	return
}

func extractField(f *ast.Field, ims []*ast.ImportSpec) (field Field, err error) {
	if len(f.Names) > 0 {
		field.Name = f.Names[0].Name
	}
	field.Type = extractBeanClass(f, ims)

	cmt, ok := astutil.Comment(f)
	if ok {
		as, errs := parser.Parse(cmt)
		if errs != nil {
			err = errs
			logger.Logger.Warnf("parse comment error: %s", cmt)
			return
		}
		field.Annotations = as
	}
	return
}

func extractDepends(n ast.Node, ims []*ast.ImportSpec) (dependsOn []string, err error) {
	visitor := astutil.NewVisitor(func(ns ast.Node) {
		logger.Logger.Debugf("start extractDepend %+v", n)
		if f, ok := ns.(*ast.Field); ok {
			bclz := extractBeanClass(f, ims)
			if bclz == "" {
				return
			}
			dependsOn = append(dependsOn, bclz)
		}
	})
	ast.Walk(&visitor, n)
	return
}

func extractBeanClass(f *ast.Field, ims []*ast.ImportSpec) string {
	pkg, name := astutil.GetFieldPackageAndTypeName(f)
	pkgArr := lo.Map(ims, func(item *ast.ImportSpec, index int) string {
		return strings.Trim(item.Path.Value, "\"")
	})
	find := lo.Filter(pkgArr, func(item string, index int) bool {
		return util.EndWith(item, pkg)
	})
	if len(find) > 0 {
		pkg = find[0]
	}

	return pkg + "." + name
}
