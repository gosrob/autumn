package main

import (
	_ "github.com/YReshetko/go-annotation/annotations/constructor"
	_ "github.com/YReshetko/go-annotation/annotations/mapper"
	_ "github.com/YReshetko/go-annotation/annotations/validator"
	annotation "github.com/YReshetko/go-annotation/pkg"
	"github.com/gosrob/autumn/internal/logger"
	_ "github.com/gosrob/autumn/internal/processor"
	_ "github.com/gosrob/autumn/wire"
)

func main() {
	defer logger.Logger.CatchPanic()
	logger.Logger.SetIsDebug(true)
	annotation.Process()
}
