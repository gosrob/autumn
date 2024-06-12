package app

import (
	annotation "github.com/YReshetko/go-annotation/pkg"
	"github.com/gosrob/autumn/internal/logger"
	"github.com/gosrob/autumn/internal/util/cast"
)

type goAnnotationBeanDefinitionReader struct{}

// ReadBeanFactoryDefinition implements BeanDefinitionReader.
func (g *goAnnotationBeanDefinitionReader) ReadBeanFactoryDefinition(n any) (FactoryFuncDefinition, error) {
	var factory FactoryFuncDefinition
	node, err := cast.Cast[annotation.Node](n)
	if err != nil {
		return factory, err
	}
	isBuiltinType, beanClass, dependsOn, isPrimary, isLazy, alias, bean, pkg, err := extractFunc(node)
	if err != nil {
		logger.Logger.Warnf("parse bean func factory fail %s", err)
		return factory, err
	}
	factory = FactoryFuncDefinition{
		IsBuiltinType: isBuiltinType,
		BeanClass:     BeanClass(beanClass),
		DependsOn:     dependsOn,
		IsPrimary:     isPrimary,
		IsLazy:        isLazy,
		Alias:         alias,
		Consturcts:    map[string]Constructor{},
		Bean:          bean,
		Pachage:       pkg,
	}

	return factory, nil
}

var GoAnnotationBeanDefinitionReader goAnnotationBeanDefinitionReader

// ReadBeanDefinition implements BeanDefinitionReader.
func (g *goAnnotationBeanDefinitionReader) ReadBeanDefinition(n any) (BeanDefinition, error) {
	var b BeanDefinition
	node, err := cast.Cast[annotation.Node](n)
	if err != nil {
		return b, err
	}

	beanClass, isBasicType, dependsOn, isPrimary, isLazy, alias, bean, pkg, err := extractNode(node)
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
		Pachage:       pkg,
	}
	return b, nil
}

var _ BeanDefinitionReader = (*goAnnotationBeanDefinitionReader)(nil)
