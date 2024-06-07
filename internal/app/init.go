package app

func init() {
	GoAnnotationBeanDefinitionReader = goAnnotationBeanDefinitionReader{}
	BeanRegistry = beanRegistry{}
	ApplicationContexter = ApplicationContext{
		BeanDefinitionReader: &GoAnnotationBeanDefinitionReader,
		BeanRegistryer:       &BeanRegistry,
	}
}
