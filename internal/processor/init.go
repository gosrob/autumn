package processor

import (
	annotation "github.com/YReshetko/go-annotation/pkg"
	pkg "github.com/gosrob/autumn/pkg/annotation"
)

func init() {
	annotation.Register[pkg.Autowired](&Processor)
	annotation.Register[pkg.Bean](&Processor)
}
