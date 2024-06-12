package beandefinition

import (
	"github.com/gosrob/autumn/testdata/beandefinition/subdirectory"
	"github.com/zrb/bufio"
)

// @Test(key="there", key2="here")
// @Bean(isPrimary="true", isLazy="false", alias="demo")
type DefinitionDemo struct { // is come
	// @TestField(key="there")
	Reader bufio.ReadWriter // this is comment
}

// @Bean(isPrimary="true", isLazy="false", alias="demo")
func ProduceDefinitionDemo() *subdirectory.B {
	return nil
}
