package app

func init() {
	GoAnnotationBeanDefinitionReader = goAnnotationBeanDefinitionReader{}
	BeanRegistry = beanRegistry{
		beans:        map[string][]BeanDefinition{},
		beansFactory: map[string][]FactoryFuncDefinition{},
	}
	ApplicationContexter = ApplicationContext{
		BeanDefinitionReader: &GoAnnotationBeanDefinitionReader,
		BeanRegistryer:       &BeanRegistry,
	}
}
