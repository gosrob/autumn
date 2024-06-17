package app

import (
	"sort"
	"strings"

	"github.com/gosrob/autumn/internal/errorcode"
	"github.com/gosrob/autumn/internal/logger"
	"github.com/gosrob/autumn/internal/util/astutil"
	"github.com/gosrob/autumn/internal/util/stream"
	treemap "github.com/liyue201/gostl/ds/map"
	"github.com/liyue201/gostl/utils/comparator"
)

type BeanFactory interface {
	GetBean(className string, params ...string) (beanResolver, error)
	GetPrimaryBean(className string, params ...string) (beanResolver, error)
}

type ListableBeanFactory interface {
	BeanFactory
	GetBeans(className string) ([]beanResolver, error)
	GetAllResolvedBeans() []beanResolver
}

var ListableBeanFactoryer ListableBeanFactory

type DefaultBeanFactoryer struct {
	// created  map[string][]beanResolver
	created    *treemap.Map[string, []beanResolver]
	interfaces *treemap.Map[string, []beanResolver]
	registry   BeanRegistryer
	orders     []string
	init       bool
}

// GetPrimaryBean implements ListableBeanFactory.
func (d *DefaultBeanFactoryer) GetPrimaryBean(className string, params ...string) (beanResolver, error) {
	rbs, err := d.GetBeans(className)
	if err != nil {
		return nil, err
	}

	if len(rbs) != 1 {
		rbs = stream.OfSlice(rbs).Filter(func(br beanResolver) bool { return br.GetDefinitionBase().IsPrimary }).ToSlice()
	}
	if len(rbs) != 1 {
		return nil, &errorcode.BeanPrimaryError
	}
	return rbs[0], nil
}

// GetAllResolvedBeans implements ListableBeanFactory.
func (d *DefaultBeanFactoryer) GetAllResolvedBeans() []beanResolver {
	all := []beanResolver{}

	d.created.Traversal(func(key string, value []beanResolver) bool {
		for _, v := range value {
			logger.Logger.Debugf("iter %s", v.GetType())
		}
		all = append(all, value...)
		return true
	})
	all = sortBeansByTypeOrder(all, d.orders)
	return all
}

func sortBeansByTypeOrder(beans []beanResolver, typeOrder []string) []beanResolver {
	// 创建一个类型到顺序的映射
	typeIndexMap := make(map[string]int)
	for i, t := range typeOrder {
		typeIndexMap[t] = i
	}

	// 对 beans 进行排序
	sort.Slice(beans, func(i, j int) bool {
		typeI := beans[i].GetDefinitionBase().BeanClass
		typeJ := beans[j].GetDefinitionBase().BeanClass

		indexI, okI := typeIndexMap[string(typeI)]
		indexJ, okJ := typeIndexMap[string(typeJ)]

		if !okI || !okJ { // 如果某个类型不在排序列表中出现，将它们的 index 置为最大并按照字符串初始顺序。
			if !okI && !okJ {
				return strings.Compare(string(typeI), string(typeJ)) < 0
			}
			return okI
		}
		return indexI < indexJ
	})

	return beans
}

func (d *DefaultBeanFactoryer) push(className string, b beanResolver) {
	if stream.OfSlice(d.orders).Filter(func(s string) bool { return s == className }).Length() > 0 {
		return
	}
	logger.Logger.Debugf("create bean %s", className)
	insert, _ := d.created.Get(className)
	d.created.Insert(className, append(insert, b))
	d.orders = append(d.orders, className)
}

func (d *DefaultBeanFactoryer) pushInterface(interfaceName string, b beanResolver) {
	if stream.OfSlice(d.orders).Filter(func(s string) bool { return s == interfaceName }).Length() > 0 {
		return
	}
	logger.Logger.Debugf("create bean %s", interfaceName)
	insert, _ := d.interfaces.Get(interfaceName)
	d.interfaces.Insert(interfaceName, append(insert, b))
}

func NewDefaultBeanFactory(registry BeanRegistryer) *DefaultBeanFactoryer {
	return &DefaultBeanFactoryer{
		created:    treemap.New[string, []beanResolver](comparator.StringComparator, treemap.WithGoroutineSafe()),
		interfaces: treemap.New[string, []beanResolver](comparator.StringComparator, treemap.WithGoroutineSafe()),
		registry:   registry,
		orders:     []string{},
		init:       false,
	}
}

// GetBean implements ListableBeanFactory.
func (d *DefaultBeanFactoryer) GetBean(className string, params ...string) (beanResolver, error) {
	if b, err := d.created.Get(className); err == nil || len(b) > 0 {
		return b[0], nil
	}

	fds := d.registry.GetBeanFactoryDefinition(className)
	for _, fd := range fds {
		b, err := d.makeFactoryBean(fd, params...)
		if err != nil {
			return b, err
		}
		d.push(className, b)
		return b, nil
	}

	if d.init {
		return nil, errorcode.CreateZeroBeanError.Instance().Printf("need bean class: %s", className)
	}
	for _, bd := range d.registry.GetAllBeans() {
		if bd.IsInterface {
			continue
		}
		b, err := d.makeBean(bd)
		if err != nil {
			return b, err
		}
		d.push(string(bd.BeanClass), b)

	}
	// bds := d.registry.GetBeanDefinition(className)
	// for _, bd := range bds {
	// 	if bd.IsInterface {
	// 		continue
	// 	}
	// 	b, err := d.makeBean(bd)
	// 	if err != nil {
	// 		return b, err
	// 	}
	// 	d.push(className, b)
	// }

	allCreatedBean := d.GetAllResolvedBeans()
	// determine if any bean implements interface
	for _, bd := range d.registry.GetAllBeans() {
		bdsImplement := []beanResolver{}
		if !bd.IsInterface {
			continue
		}
		for _, rb := range allCreatedBean {
			if ok, err := astutil.CheckIfTypeImplementsInterfaceWithCache(string(rb.GetDefinitionBase().BeanClass), string(bd.BeanClass)); err == nil && ok {
				bdsImplement = append(bdsImplement, rb)
			}
		}
		if len(bdsImplement) > 0 {
			for _, v := range bdsImplement {
				d.pushInterface(string(bd.BeanClass), v)
			}
		}
	}

	d.init = true
	if b, err := d.created.Get(className); err == nil && len(b) > 0 {
		return b[0], nil
	}

	return nil, errorcode.CreateZeroBeanError.Instance().Printf("need bean class: %s", className)
}

// GetBeans implements ListableBeanFactory.
func (d *DefaultBeanFactoryer) GetBeans(className string) ([]beanResolver, error) {
	// run this to init created
	d.GetBean(className)

	if b, ok := d.created.Get(className); ok == nil && len(b) > 0 {
		return b, nil
	}

	if b, ok := d.interfaces.Get(className); ok == nil && len(b) > 0 {
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

func (d *DefaultBeanFactoryer) makeFactoryBean(bd FactoryFuncDefinition, params ...string) (bean beanResolver, err error) {
	bean = NewBeanResolveFactory("var1", bd.Pachage.CurrentPackage+"."+bd.Bean.Name, bd.DefinitionBase, bd.Bean.Results[0].TypeInfo.IsPointer, params...)
	return
}
