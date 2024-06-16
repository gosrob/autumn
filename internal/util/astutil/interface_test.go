package astutil_test

import (
	"testing"

	"github.com/gosrob/autumn/internal/logger"
	"github.com/gosrob/autumn/internal/util/astutil"
)

func TestStructImpl(t *testing.T) {
	ok, err := astutil.CheckIfTypeImplementsInterfaceWithCache("github.com/gosrob/autumn/internal/app.beanRegistry", "github.com/gosrob/autumn/internal/app.BeanRegistryer")
	logger.Logger.Info(ok, err)
}

func TestStructImplEmptyInterface(t *testing.T) {
	// ok, err := astutil.CheckIfTypeImplementsInterfaceWithCache("github.com/gosrob/autumn/internal/app.beanRegistry", "github.com/gosrob/autumn/internal/app.BeanRegistryer")
	ok, err := astutil.CheckIfTypeImplementsInterfaceWithCache("github.com/gosrob/autumn/internal/app.beanRegistry", "github.com/gosrob/autumn/examples/beandefinition.DemoInterface")
	logger.Logger.Info(ok, err)
}
