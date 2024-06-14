package astutil_test

import (
	"testing"

	"github.com/gosrob/autumn/internal/logger"
	"github.com/gosrob/autumn/internal/util/astutil"
	"github.com/gosrob/autumn/internal/util/pkginfo"
)

func TestStructImpl(t *testing.T) {
	ok, err := astutil.CheckIfTypeImplementsInterfaceWithCache("github.com/gosrob/autumn/internal/app.beanRegistry", "github.com/gosrob/autumn/internal/app.BeanRegistryer", pkginfo.GetFullPackage("").Module.Dir)
	logger.Logger.Info(ok, err)
}
