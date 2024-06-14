package processor

import (
	annotation "github.com/YReshetko/go-annotation/pkg"
	"github.com/gosrob/autumn/internal/logger"
	pkg "github.com/gosrob/autumn/pkg/annotation"
)

func init() {
	logger.Logger.Info("inited")
	annotation.Register[pkg.Autowired](&Processor)
	annotation.Register[pkg.Bean](&Processor)
	annotation.Register[pkg.MetaInfo](&Processor)
}
