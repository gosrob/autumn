package annotation

type Bean struct {
	IsPrimary string `annotation:"name=isPrimary"`
	IsLazy    string `annotation:"name=isLazy"`
	Alias     string `annotation:"name=alias"`
}
