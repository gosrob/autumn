package beandefinition

import (
	"github.com/gosrob/autumn/testdata/beandefinition/subdirectory"
	"github.com/zrb/bufio"
)

// @MetaInfo(wirePath="wire/wire.go")
// @Test(key="there", key2="here")
// @Bean(isPrimary="false", isLazy="false", alias="demo")
type DefinitionDemo struct { // is come
	// @TestField(key="there")
	Reader bufio.ReadWriter // this is comment

	// @Autowired(key="there")
	B subdirectory.B // this is comment
}

// @Bean(isPrimary="true", isLazy="false", alias="demo")
func ProduceDefinitionDemo(b *subdirectory.B) *subdirectory.B {
	return nil
}

// @Bean(isPrimary="true", isLazy="false", alias="demo")
func ProduceDefinition(b *subdirectory.B) *DefinitionDemo {
	return nil
}
