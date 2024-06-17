package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/gosrob/autumn/internal/errorcode"
	"github.com/gosrob/autumn/internal/logger"
	"github.com/gosrob/autumn/internal/util/logic"
	"github.com/gosrob/autumn/internal/util/pkginfo"
	"github.com/gosrob/autumn/internal/util/stream"
	"github.com/gosrob/autumn/pkg/annotation"
	pkg "github.com/gosrob/autumn/pkg/annotation"
)

type ApplicationContext struct {
	Metainfo *annotation.MetaInfo
	BeanDefinitionReader
	BeanRegistryer
	ListableBeanFactory

	containerBuilder strings.Builder
}

var ApplicationContexter ApplicationContext

func (a *ApplicationContext) Run(ctx context.Context) (files map[string][]byte) {
	var err error
	// NOTE: Run to here means that all annotations are already collected, and parsed, but alias name is not set.

	// TODO: begin to create bean with zero value(if not has factory func), if has factory func, then initialize it with factory func, and push to resolvedRegistry
	err = a.CreateZeroBean()
	if err != nil {
		logger.Logger.Fatalf("create bean fail %s", err)
		return nil
	}

	// TODO: before inject attributes, check attributes, if one attributes has gt 1 beans && these beans do not have primary flag set, then FatalLog this error
	err = a.Check()
	if err != nil {
		logger.Logger.Fatalf("check beans error %s", err)
		return nil
	}

	// TODO: before register to container, we need to check variable name conflick, if name conflicks, generate new names one by one.
	a.ConferName()

	deferedFn := a.initBeanBuilder()

	// TODO: when create all value, then inject value, if value is not constructed, then try to find factory func, if no factory func found, then FatalLog(this error)
	a.WriteCreateBean()
	err = a.Wire()
	if err != nil {
		logger.Logger.Fatalf("%s", err)
		return nil
	}

	// TODO: So first set alias name in registry
	a.SetAlias()

	// TODO: run to here all beans and all attributes are set.  then register this to sambar/do container, which is a generic type container
	a.Inject()

	deferedFn()
	// Future TODO: in future canbe  cp all files to ./tmp directory, after that we can rewrite Struct to Struct_base, and our Struct Proxy canbe Struct, so code in our project do not need to change type from struct to structProxy
	files = map[string][]byte{}
	files[a.Metainfo.WirePath] = []byte(a.containerBuilder.String())
	return
}

func (a *ApplicationContext) initBeanBuilder() func() {
	template := `
package %s

import (
	"github.com/gosrob/autumn/pkg/container"
	"github.com/samber/do/v2"
    %s
)

func init() {
    `

	ims := []string{}
	for _, rb := range a.GetAllResolvedBeans() {
		ims = append(ims, rb.GetDefinitionBase().Pachage.CurrentFullPackage)
	}
	ims = stream.OfSlice(ims).
		Map(func(s string) string { return fmt.Sprintf(`"%s"`, s) }).
		Filter(stream.Distinct[string]()).
		ToSlice()

	pkg := pkginfo.GetPackageFromFilePath(a.Metainfo.WirePath)
	if a.Metainfo.WirePkg != "" {
		pkg = a.Metainfo.WirePkg
	}
	a.containerBuilder.WriteString(
		fmt.Sprintf(template+"\n", pkg, strings.Join(ims, "\n")),
	)

	return func() {
		a.containerBuilder.WriteString("}")
	}
}

func (a *ApplicationContext) SetAlias() {
	a.ResolveAlias()
}

func (a *ApplicationContext) CreateZeroBean() error {
	// Get bean trigger factory create bean
	for _, bd := range a.GetAllBeans() {
		a.GetBean(string(bd.BeanClass))
	}
	for _, bd := range a.GetAllFactoryBeans() {
		if len(bd.Bean.Params) <= 0 {
			a.GetBean(string(bd.BeanClass))
		}
	}
	for _, bd := range a.GetAllFactoryBeans() {
		if len(bd.Bean.Params) > 0 {
			callParams := []string{}
			params := bd.Bean.Params
			for _, p := range params {
				b, err := a.GetBean(p.Type)
				if err != nil {
					return &errorcode.BeanNotFindError
				}

				callParams = append(callParams, logic.OrGet(p.TypeInfo.IsPointer, "&"+b.GetDecl(), b.GetDecl()))
			}
			_, err := a.GetBean(string(bd.BeanClass), callParams...)
			if err != nil {
				return err
			}
		} else {
			_, err := a.GetBean(string(bd.BeanClass))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *ApplicationContext) Check() error {
	if a.Metainfo == nil {
		return &errorcode.MetaInfoNotDefined
	}
	bdclass := a.GetBeanClassStrs()
	for _, class := range bdclass {
		beans, _ := a.GetBeans(class)
		if len(beans) <= 1 {
			continue
		}
		if stream.OfSlice(beans).Filter(func(br beanResolver) bool {
			return br.GetDefinitionBase().IsPrimary
		}).Length() != 1 {
			return &errorcode.BeanPrimaryError
		}
	}

	return nil
}

func (a *ApplicationContext) WriteCreateBean() error {
	for _, bd := range a.GetAllResolvedBeans() {
		if bd.GetDefinitionBase().IsInterface {
			continue
		}
		a.containerBuilder.WriteString(bd.GetConstructText())
	}
	return nil
}

func (a *ApplicationContext) Wire() error {
	for _, bd := range a.GetAllBeans() {
		if bd.IsInterface {
			continue
		}
		s, err := wireAttribute(bd, a.ListableBeanFactory, a.BeanRegistryer)
		if err != nil {
			return err
		}
		a.containerBuilder.WriteString(s)
	}

	return nil
}

func wireAttribute(bd BeanDefinition, lbf ListableBeanFactory, br BeanRegistryer) (string, error) {
	wireFunc := func(bdi BeanDefinition, rbi beanResolver) (string, error) {
		var builder strings.Builder
		wireTempl := `
    %s.%s = %s
        `
		wireArrayTempl := `
    %s.%s = append(%s.%s, %s)
        `
		for _, f := range bdi.Bean.Fields {
			if rbi.IsFiledResolved(f.Name) {
				continue
			}
			autowired := FindFieldAnnotation[pkg.Autowired](f)
			if len(autowired) <= 0 {
				continue
			}
			rb, err := lbf.GetPrimaryBean(f.Type)
			if err != nil {
				return "", errorcode.BeanNotFindError.DeepCopy().Printf("%s", err)
			}

			isInterface := false
			bd := br.GetBeanDefinition(f.Type)
			if len(bd) > 0 && bd[0].IsInterface {
				isInterface = true
			}
			if !f.TypeInfo.IsArray {
				builder.WriteString(
					fmt.Sprintf(wireTempl, rbi.GetDecl(), f.Name, logic.OrGet(isInterface, "&"+rb.GetDecl(), rb.GetDecl())),
				)
			} else {
				rbs, _ := lbf.GetBeans(f.Type)
				for _, rb := range rbs {
					builder.WriteString(
						fmt.Sprintf(wireArrayTempl, rbi.GetDecl(), f.Name, rbi.GetDecl(), f.Name, logic.OrGet(isInterface, "&"+rb.GetDecl(), rb.GetDecl())),
					)
				}
			}
			rbi.SetResolved(f.Name)
		}
		return builder.String(), nil
	}

	wireRb, _ := lbf.GetBean(string(bd.BeanClass))

	s, err := wireFunc(bd, wireRb)
	if err != nil {
		return "", errorcode.BeanNotFindError.DeepCopy().Printf("%s", err)
	}

	return s, nil
}

func (a *ApplicationContext) ConferName() {
	bdclass := a.GetBeanClassStrs()
	names := map[string]any{}
	for _, class := range bdclass {
		beans, _ := a.GetBeans(class)
		for _, rb := range beans {
			conferName(rb, names)
		}
	}
}

func conferName(rd beanResolver, names map[string]any) {
	if _, ok := names[string(rd.GetDecl())]; ok {
		rd.RandomName()
		conferName(rd, names)
	}
}

func (a *ApplicationContext) Inject() {
	temp := `
	do.ProvideNamedValue(container.Container, "%s", %s)
`
	for _, rb := range a.GetAllResolvedBeans() {
		inject := fmt.Sprintf(temp+"\n", rb.GetDefinitionBase().BeanClass, rb.GetDecl())
		a.containerBuilder.WriteString(inject)
	}
}
