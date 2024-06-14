package annotation

type MetaInfo struct {
	WirePath string `annotation:"name=wirePath"`
	WirePkg  string `annotation:"name=wirePkg"`
}
