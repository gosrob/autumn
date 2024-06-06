package app

type BeanClass string

// BeanDefinition represents the definition of a bean in the application.
type BeanDefinition struct {
	// IsBuiltinType indicates whether the type is a built-in Go type.
	IsBuiltinType bool

	// # beanClass is the fully qualified name of the bean's class.
	//
	// particularly {full-packageName}.{typeName}
	// TODO: current not use this field
	BeanClass BeanClass

	// TODO: Currently, this field is not used.
	DependsOn []string

	// IsPrimary indicates if the bean is the primary candidate for autowiring.
	// TODO: Currently, this field is not used.
	IsPrimary bool

	// IsLazy specifies whether the bean should be lazily initialized.
	// TODO: Currently, this field is not used.
	IsLazy bool

	// Consturcts is a map where each key is a string representing the name of a bean,
	// and each value is a Constructor function that creates an instance of that bean.
	// TODO: Currently, this field is not used.
	Consturcts map[string]Constructor

	Bean StructDefinition
}

type Constructor struct {
	ConstructBeanName string
	ConstructorFn     func() string
}
