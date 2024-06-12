package processor

import (
	"context"
	"go/ast"

	annotation "github.com/YReshetko/go-annotation/pkg"
	. "github.com/gosrob/autumn/internal/app"
	. "github.com/gosrob/autumn/internal/logger"
	"github.com/gosrob/autumn/internal/util/astutil"
)

type processor struct{}

var Processor = processor{}

// Name implements annotation.AnnotationProcessor.
func (p *processor) Name() string {
	return "autum framework"
}

// Output implements annotation.AnnotationProcessor.
func (p *processor) Output() map[string][]byte {
	return ApplicationContexter.Run(context.Background())
}

// Process implements annotation.AnnotationProcessor.
func (p *processor) Process(node annotation.Node) error {
	if _, error := astutil.AstCast[*ast.TypeSpec](node.ASTNode()); error == nil {
		bd, err := ApplicationContexter.ReadBeanDefinition(node)
		if err != nil {
			Logger.Warnf("failed parse beanDefinition, %+v", node)
			return err
		}
		ApplicationContexter.RegisterBeanDefinition(bd)
	}

	if _, error := astutil.AstCast[*ast.FuncDecl](node.ASTNode()); error == nil {
		fd, err := ApplicationContexter.ReadBeanFactoryDefinition(node)
		if err != nil {
			Logger.Warnf("failed parse beanDefinition, %+v", node)
			return err
		}
		ApplicationContexter.RegisterBeanFactoryDefinition(fd)
	}

	return nil
}

// Version implements annotation.AnnotationProcessor.
func (p *processor) Version() string {
	return "0.1"
}

var _ annotation.AnnotationProcessor = (*processor)(nil)
