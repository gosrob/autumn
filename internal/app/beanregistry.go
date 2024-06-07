package app

type BeanRegistryer interface {
	RegisterBeanDefinition(beanDefinition BeanDefinition)
	RegisterBeanDefinitionByName(name string, beanDefinition BeanDefinition)
	RegisterBeanFactoryDefinition(factoryFuncDefinition FactoryFuncDefinition)
}

type beanRegistry struct{}

// RegisterBeanFactoryDefinition implements BeanRegistryer.
func (b *beanRegistry) RegisterBeanFactoryDefinition(factoryFuncDefinition FactoryFuncDefinition) {
	panic("unimplemented")
}

// RegisterBeanDefinitionByName implements BeanRegistryer.
func (b *beanRegistry) RegisterBeanDefinitionByName(name string, beanDefinition BeanDefinition) {
	panic("unimplemented")
}

// RegisterBeanDefinition implements BeanRegistryer.
func (b *beanRegistry) RegisterBeanDefinition(beanDefinition BeanDefinition) {
	panic("unimplemented")
}

var _ BeanRegistryer = (*beanRegistry)(nil)

var BeanRegistry beanRegistry
