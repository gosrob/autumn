package main

import (
	"fmt"

	_ "github.com/YReshetko/go-annotation/annotations/constructor"
	_ "github.com/YReshetko/go-annotation/annotations/mapper"
	_ "github.com/YReshetko/go-annotation/annotations/validator"
	"github.com/gosrob/autumn/internal/logger"
	_ "github.com/gosrob/autumn/internal/processor"
	"github.com/gosrob/autumn/pkg/container"
	"github.com/gosrob/autumn/testdata/beandefinition/subdirectory"
	_ "github.com/gosrob/autumn/wire"
	"github.com/samber/do/v2"
)

func main() {
	defer logger.Logger.CatchPanic()
	a := do.MustInvoke[subdirectory.B](container.Container)
	fmt.Println(a)
	// logger.Logger.SetIsDebug(true)
	// annotation.Process()
}
