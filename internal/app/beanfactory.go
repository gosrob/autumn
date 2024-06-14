package app

import (
	"github.com/gosrob/autumn/internal/errorcode"
)

type BeanFactory interface {
	GetBean(className string) (beanResolver, error)
}

type ListableBeanFactory interface {
	BeanFactory
	GetBeans(className string) ([]beanResolver, error)
	GetAllResolvedBeans() []beanResolver
}

var ListableBeanFactoryer ListableBeanFactory

type DefaultBeanFactoryer struct {
	created  map[string][]beanResolver
	registry BeanRegistryer
}

// GetAllResolvedBeans implements ListableBeanFactory.
func (d *DefaultBeanFactoryer) GetAllResolvedBeans() []beanResolver {
	all := []beanResolver{}

	for _, v := range d.created {
		all = append(all, v...)
	}

	return all
}

func NewDefaultBeanFactory(registry BeanRegistryer) *DefaultBeanFactoryer {
	return &DefaultBeanFactoryer{
		created:  map[string][]beanResolver{},
		registry: registry,
	}
}

// GetBean implements ListableBeanFactory.
func (d *DefaultBeanFactoryer) GetBean(className string) (beanResolver, error) {
	if b, ok := d.created[className]; ok && len(b) > 0 {
		return b[0], nil
	}
	bds := d.registry.GetBeanDefinition(className)
	for _, bd := range bds {
		b, err := d.makeBean(bd)
		if err != nil {
			return b, err
		}
		d.created[className] = append(d.created[className], b)
	}

	fds := d.registry.GetBeanFactoryDefinition(className)
	for _, fd := range fds {
		b, err := d.makeFactoryBean(fd)
		if err != nil {
			return b, err
		}
		d.created[className] = append(d.created[className], b)
	}

	if b, ok := d.created[className]; ok && len(b) > 0 {
		return b[0], nil
	}

	return nil, &errorcode.CreateZeroBeanError
}

// GetBeans implements ListableBeanFactory.
func (d *DefaultBeanFactoryer) GetBeans(className string) ([]beanResolver, error) {
	// run this to init created
	d.GetBean(className)

	if b, ok := d.created[className]; ok && len(b) > 0 {
		return b, nil
	}

	return nil, &errorcode.CreateZeroBeanError
}

var (
	_                  ListableBeanFactory = (*DefaultBeanFactoryer)(nil)
	DefaultBeanFactory ListableBeanFactory
)

func (d *DefaultBeanFactoryer) makeBean(bd BeanDefinition) (bean beanResolver, err error) {
	bean = NewBeanResolveConstructor("var1", string(bd.BeanClass), bd.DefinitionBase)
	return
}

func (d *DefaultBeanFactoryer) makeFactoryBean(bd FactoryFuncDefinition) (bean beanResolver, err error) {
	bean = NewBeanResolveFactory("var1", bd.Pachage.CurrentPackage+"."+bd.Bean.Name, bd.DefinitionBase, bd.Bean.Results[0].TypeInfo.IsPointer)
	return
}
