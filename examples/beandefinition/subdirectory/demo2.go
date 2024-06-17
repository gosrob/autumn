package subdirectory

// @Bean(isPrimary="true")
type B struct {
	// @Autowired
	Cc *C
}

// @Bean(isPrimary="true", isLazy="false")
type C struct{}

func (c *C) Hello() {
	panic("not implemented") // TODO: Implement
}
