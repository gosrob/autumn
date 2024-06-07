package processor

import (
	annotation "github.com/YReshetko/go-annotation/pkg"
	. "github.com/gosrob/autumn/internal/app"
	. "github.com/gosrob/autumn/internal/logger"
)

type processor struct{}

var Processor = processor{}

// Name implements annotation.AnnotationProcessor.
func (p *processor) Name() string {
	return "autum framework"
}

// Output implements annotation.AnnotationProcessor.
func (p *processor) Output() map[string][]byte {
	return ApplicationContexter.Run()
}

// Process implements annotation.AnnotationProcessor.
func (p *processor) Process(node annotation.Node) error {
	bd, err := ApplicationContexter.ReadBeanDefinition(node)
	if err != nil {
		Logger.Warnf("failed parse beanDefinition, %+v", node)
		return err
	}
	ApplicationContexter.RegisterBeanDefinition(bd)
	// TODO: register factory func
	return nil
}

// Version implements annotation.AnnotationProcessor.
func (p *processor) Version() string {
	return "0.1"
}

var _ annotation.AnnotationProcessor = (*processor)(nil)
