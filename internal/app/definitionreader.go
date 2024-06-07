package app

import (
	"go/ast"

	"github.com/gosrob/autumn/internal/util/astutil"
)

type BeanDefinitionReader interface {
	ReadBeanDefinition(node any) (BeanDefinition, error)
	ReadBeanFactoryDefinition(node any) (FactoryFuncDefinition, error)
}

type goAnnotationBeanDefinitionReader struct{}

// ReadBeanFactoryDefinition implements BeanDefinitionReader.
func (g *goAnnotationBeanDefinitionReader) ReadBeanFactoryDefinition(n any) (FactoryFuncDefinition, error) {
	panic("unimplemented")
}

var GoAnnotationBeanDefinitionReader goAnnotationBeanDefinitionReader

// ReadBeanDefinition implements BeanDefinitionReader.
func (g *goAnnotationBeanDefinitionReader) ReadBeanDefinition(n any) (BeanDefinition, error) {
	var b BeanDefinition
	_, err := astutil.Cast[*ast.StructType](n)
	if err != nil {
		return b, err
	}

	return b, nil
}

var _ BeanDefinitionReader = (*goAnnotationBeanDefinitionReader)(nil)
