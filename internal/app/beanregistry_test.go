package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBeanRegistry(t *testing.T) {
	BeanRegistry.RegisterBeanDefinition(BeanDefinition{})
	BeanRegistry.RegisterBeanDefinition(BeanDefinition{})
	BeanRegistry.RegisterBeanFactoryDefinition(FactoryFuncDefinition{})

	assert.Greater(t, len(BeanRegistry.GetAllBeans()), 1)
	assert.Equal(t, len(BeanRegistry.GetAllFactoryBeans()), 1)
}
