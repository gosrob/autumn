package app

type BeanRegistryer interface {
	RegisterBeanDefinition(beanDefinition BeanDefinition)
	RegisterBeanDefinitionByName(name string, beanDefinition BeanDefinition)
	RegisterBeanFactoryDefinition(factoryFuncDefinition FactoryFuncDefinition)
}

type beanRegistry struct {
	beans        map[string][]BeanDefinition
	beansFactory map[string][]FactoryFuncDefinition
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
