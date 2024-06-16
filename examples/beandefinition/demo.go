package beandefinition

import (
	"github.com/gosrob/autumn/examples/beandefinition/subdirectory"
)

// @MetaInfo(wirePath="wire/wire.go")
// @Test(key="there", key2="here")
// @Bean(isPrimary="true", isLazy="false", alias="demo")
type DefinitionDemo struct { // is come
	// @TestField(key="there")

	// @Autowired(key="there")
	B subdirectory.B // this is comment

	// @Autowired()
	Inter []DemoInterface
}

// @Bean(isLazy="false")
type DemoInterface interface {
	Hello()
}

// @Bean(isPrimary="true", isLazy="false", alias="demo")
func ProduceDefinitionDemo(b *subdirectory.C) *subdirectory.B {
	return &subdirectory.B{}
}
