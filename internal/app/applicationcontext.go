package app

import (
	"context"

	"github.com/gosrob/autumn/pkg/annotation"
)

type ApplicationContext struct {
	Metainfo annotation.MetaInfo
	BeanDefinitionReader
	BeanRegistryer
}

var ApplicationContexter ApplicationContext

func (a *ApplicationContext) ScanCollect(ctx context.Context) {
}

func (a *ApplicationContext) Run(ctx context.Context) map[string][]byte {
	// NOTE: Run to here means that all annotations are already collected, and parsed, but alias name is not set.

	// TODO: So first set alias name in registry
	a.SetAlias()

	// TODO: begin to create bean with zero value(if not has factory func), if has factory func, then initialize it with factory func, and push to resolvedRegistry
	a.CreateZeroBean()

	// TODO: before inject attributes, check attributes, if one attributes has gt 1 beans && these beans do not have primary flag set, then FatalLog this error
	a.Check()

	// TODO: when create all value, then inject value, if value is not constructed, then try to find factory func, if no factory func found, then FatalLog(this error)
	a.Wire()

	// TODO: before register to container, we need to check variable name conflick, if name conflicks, generate new names one by one.
	a.ConferName()

	// TODO: run to here all beans and all attributes are set.  then register this to sambar/do container, which is a generic type container
	a.Inject()

	// Future TODO: in future canbe  cp all files to ./tmp directory, after that we can rewrite Struct to Struct_base, and our Struct Proxy canbe Struct, so code in our project do not need to change type from struct to structProxy
	return nil
}

func (a *ApplicationContext) SetAlias() {
}

func (a *ApplicationContext) CreateZeroBean() {
}

func (a *ApplicationContext) Check() {
}

func (a *ApplicationContext) Wire() {
}

func (a *ApplicationContext) ConferName() {
}

func (a *ApplicationContext) Inject() {
}
