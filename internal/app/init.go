package app

func init() {
	GoAnnotationBeanDefinitionReader = goAnnotationBeanDefinitionReader{}
	BeanRegistry = beanRegistry{
		beans:        map[string][]BeanDefinition{},
		beansFactory: map[string][]FactoryFuncDefinition{},
	}
	DefaultBeanFactory = NewDefaultBeanFactory(&BeanRegistry)
	ApplicationContexter = ApplicationContext{
		BeanDefinitionReader: &GoAnnotationBeanDefinitionReader,
		BeanRegistryer:       &BeanRegistry,
		ListableBeanFactory:  DefaultBeanFactory,
		ProjectScanner:       ProjectScanner{},
	}
}
