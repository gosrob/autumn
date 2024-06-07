package app

type BeanRegistryer interface {
	RegisterBeanDefinition(beanDefinition BeanDefinition)
	RegisterBeanDefinitionByName(name string, beanDefinition BeanDefinition)
}

type beanRegistry struct{}

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
