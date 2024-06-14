package app

type BeanClass string

type DefinitionBase struct {
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

	Alias string

	Pachage Package
}

// BeanDefinition represents the definition of a bean in the application.
type BeanDefinition struct {
	DefinitionBase
	Bean StructDefinition
}

type Constructor struct {
	ConstructBeanName string
	ConstructorFn     func() string
	ConstructorParam  []Param
	ConstructorResult []Param

	BeanIndex int
}
