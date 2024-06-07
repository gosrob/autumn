package app

type ApplicationContext struct {
	BeanDefinitionReader
	BeanRegistryer
}

var ApplicationContexter ApplicationContext

func (a ApplicationContext) Run() map[string][]byte {
	// NOTE: Run to here means that all annotations are already collected, and parsed, but alias name is not set.

	// TODO: So first set alias name in registry

	// TODO: begin to create bean with zero value(if not has factory func), if has factory func, then initialize it with factory func, and push to resolvedRegistry

	// TODO: before inject attributes, check attributes, if one attributes has gt 1 beans && these beans do not have primary flag set, then FatalLog this error

	// TODO: when create all value, then inject value, if value is not constructed, then try to find factory func, if no factory func found, then FatalLog(this error)

	// TODO: before register to container, we need to check variable name conflick, if name conflicks, generate new names one by one.

	// TODO: run to here all beans and all attributes are set.  then register this to sambar/do container, which is a generic type container
	return nil
}
