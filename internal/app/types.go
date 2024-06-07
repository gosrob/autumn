package app

type AnnotationList []any

// # represents a struct
//
// this struct contains struct field
type StructDefinition struct {
	Annotations AnnotationList

	Name string

	// struct field
	Fields []Field
}

// # represent a struct field
type Field struct {
	// Field name
	Name string

	// # Field Type
	//
	// which is full package type, if one field is in another package and is a pointer, type is *github.com/xxx/xxx/pkg.{TypeName}.
	//
	// if field type is basic type, it just act like a normal type.
	Type string

	Annotations AnnotationList
}

// FuncDefinition represents the structure of a function definition, including its name,
// parameters, and results.
type FuncDefinition struct {
	Name    string  // Name is the name of the function.
	Params  []Param // Params is a slice of parameters the function accepts.
	Results []Param // Results is a slice of parameters the function returns.

	Annotations AnnotationList
}

// Param represents a parameter within a function definition.
type Param struct {
	Name string // Name is the name of the parameter.

	Type string // Type is the type of the parameter. It can be a full package type
	// (e.g., *github.com/xxx/xxx/pkg.{TypeName} for a pointer type in another package)
	// or a basic type.
}
