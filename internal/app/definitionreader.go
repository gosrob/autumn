package app

type BeanDefinitionReader interface {
	ReadBeanDefinition(node any) (BeanDefinition, error)
}

type goAnnotationBeanDefinitionReader struct{}

var GoAnnotationBeanDefinitionReader goAnnotationBeanDefinitionReader

// ReadBeanDefinition implements BeanDefinitionReader.
func (g *goAnnotationBeanDefinitionReader) ReadBeanDefinition(node any) (BeanDefinition, error) {
	panic("unimplemented")
}

var _ BeanDefinitionReader = (*goAnnotationBeanDefinitionReader)(nil)
