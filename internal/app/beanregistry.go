package app

import (
	"github.com/gosrob/autumn/internal/errorcode"
	"github.com/gosrob/autumn/internal/util/stream"
)

type BeanRegistryer interface {
	RegisterBeanDefinition(beanDefinition BeanDefinition)
	RegisterBeanDefinitionByName(name string, beanDefinition BeanDefinition)
	RegisterBeanFactoryDefinition(factoryFuncDefinition FactoryFuncDefinition)
	ResolveAlias() error
	GetAllBeans() []BeanDefinition
	GetAllFactoryBeans() []FactoryFuncDefinition
	GetBeanDefinition(className string) []BeanDefinition
	GetBeanFactoryDefinition(className string) []FactoryFuncDefinition
	GetBeanClassStrs() []string
	Reset()
}

type beanRegistry struct {
	beans        map[string][]BeanDefinition
	beansFactory map[string][]FactoryFuncDefinition
}

// GetBeanClassStrs implements BeanRegistryer.
func (b *beanRegistry) GetBeanClassStrs() []string {
	bdclass := []string{}
	for _, bd := range b.GetAllBeans() {
		bdclass = append(bdclass, string(bd.BeanClass))
	}
	for _, bd := range b.GetAllFactoryBeans() {
		bdclass = append(bdclass, string(bd.BeanClass))
	}
	return stream.OfSlice(bdclass).Filter(stream.Distinct[string]()).ToSlice()
}

// Reset implements BeanRegistryer.
func (b *beanRegistry) Reset() {
	b.beans = make(map[string][]BeanDefinition)
	b.beansFactory = make(map[string][]FactoryFuncDefinition)
}

// GetBeanDefinition implements BeanRegistryer.
func (b *beanRegistry) GetBeanDefinition(className string) []BeanDefinition {
	return b.beans[className]
}

// GetBeanFactoryDefinition implements BeanRegistryer.
func (b *beanRegistry) GetBeanFactoryDefinition(className string) []FactoryFuncDefinition {
	return b.beansFactory[className]
}

// GetAllBeans implements BeanRegistryer.
func (b *beanRegistry) GetAllBeans() []BeanDefinition {
	return stream.FlatMap(stream.OfMap(b.beans), func(t stream.Pair[string, []BeanDefinition]) []BeanDefinition {
		return t.Val2
	}).
		Filter(stream.DistinctBy(func(t BeanDefinition) string { return string(t.BeanClass) })).
		ToSlice()
}

// GetAllFactoryBeans implements BeanRegistryer.
func (b *beanRegistry) GetAllFactoryBeans() []FactoryFuncDefinition {
	return stream.FlatMap(stream.OfMap(b.beansFactory), func(t stream.Pair[string, []FactoryFuncDefinition]) []FactoryFuncDefinition {
		return t.Val2
	}).ToSlice()
}

// ResolveAlias implements BeanRegistryer.
func (b *beanRegistry) ResolveAlias() error {
	stream.OfMap(b.beans).RangeErr(func(p stream.Pair[string, []BeanDefinition]) error {
		bds := p.Val2
		for _, v := range bds {
			alias := v.Alias
			if alias == "" {
				continue
			}

			if _, ok := b.beans[alias]; ok {
				return errorcode.DepulicatedBeanAliasError.Instance().Printf("%s has set", alias)
			}
			b.beans[alias] = append(b.beans[alias], v)
		}
		return nil
	})

	stream.OfMap(b.beansFactory).RangeErr(func(p stream.Pair[string, []FactoryFuncDefinition]) error {
		bds := p.Val2
		for _, v := range bds {
			alias := v.Alias
			if alias == "" {
				continue
			}

			if _, ok := b.beansFactory[alias]; ok {
				return errorcode.DepulicatedBeanAliasError.Instance().Printf("%s has set", alias)
			}
			b.beansFactory[alias] = append(b.beansFactory[alias], v)
		}
		return nil
	})
	return nil
}

// RegisterBeanFactoryDefinition implements BeanRegistryer.
func (b *beanRegistry) RegisterBeanFactoryDefinition(factoryFuncDefinition FactoryFuncDefinition) {
	b.beansFactory[string(factoryFuncDefinition.BeanClass)] = append(b.beansFactory[string(factoryFuncDefinition.BeanClass)], factoryFuncDefinition)
}

// RegisterBeanDefinitionByName implements BeanRegistryer.
func (b *beanRegistry) RegisterBeanDefinitionByName(name string, beanDefinition BeanDefinition) {
	panic("unimplemented")
}

// RegisterBeanDefinition implements BeanRegistryer.
func (b *beanRegistry) RegisterBeanDefinition(beanDefinition BeanDefinition) {
	b.beans[string(beanDefinition.BeanClass)] = append(b.beans[string(beanDefinition.BeanClass)], beanDefinition)
}

var _ BeanRegistryer = (*beanRegistry)(nil)

var BeanRegistry beanRegistry
