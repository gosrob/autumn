package app

type BeanDefinitionReader interface {
	ReadBeanDefinition(node any) (BeanDefinition, error)
	ReadBeanFactoryDefinition(node any) (FactoryFuncDefinition, error)
}
