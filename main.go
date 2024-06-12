package main

import (
	_ "github.com/YReshetko/go-annotation/annotations/constructor"
	_ "github.com/YReshetko/go-annotation/annotations/mapper"
	_ "github.com/YReshetko/go-annotation/annotations/validator"
	annotation "github.com/YReshetko/go-annotation/pkg"
	"github.com/gosrob/autumn/internal/logger"
	_ "github.com/gosrob/autumn/internal/processor"
)

func main() {
	defer logger.Logger.CatchPanic()
	annotation.Process()
}
