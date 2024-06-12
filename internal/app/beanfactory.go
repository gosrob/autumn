package app

type BeanFactory interface {
	GetBean(className string) (beanResolver, error)
}

type ListableBeanFactory interface {
	BeanFactory
	GetBeans(className string) ([]beanResolver, error)
}

var ListableBeanFactoryer ListableBeanFactory

type DefaultBeanFactoryer struct {
	created  map[string][]beanResolver
	registry BeanRegistryer
}

func NewDefaultBeanFactory(registry BeanRegistryer) *DefaultBeanFactoryer {
	return &DefaultBeanFactoryer{
		created:  map[string][]beanResolver{},
		registry: registry,
	}
}

// GetBean implements ListableBeanFactory.
func (d *DefaultBeanFactoryer) GetBean(className string) (beanResolver, error) {
	panic("unimplemented")
}

// GetBeans implements ListableBeanFactory.
func (d *DefaultBeanFactoryer) GetBeans(className string) ([]beanResolver, error) {
	panic("unimplemented")
}

var (
	_                  ListableBeanFactory = (*DefaultBeanFactoryer)(nil)
	DefaultBeanFactory ListableBeanFactory
)
