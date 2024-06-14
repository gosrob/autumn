package app_test

import (
	"testing"

	"github.com/gosrob/autumn/internal/app"
	"github.com/gosrob/autumn/internal/logger"
	"github.com/stretchr/testify/suite"
)

type registerSuite struct {
	suite.Suite
	r app.BeanRegistryer
}

func (suite *registerSuite) SetupTest() {
	logger.Logger.SetIsDebug(true)
}

func (suite *registerSuite) TearDownTest() {
	suite.r.Reset()
	logger.Logger.SetIsDebug(false)
}

func (s *registerSuite) TestBasicBeans() {
	s.r.RegisterBeanDefinition(app.BeanDefinition{
		DefinitionBase: app.DefinitionBase{
			BeanClass: "test1",
		},
	})
	s.r.RegisterBeanDefinition(app.BeanDefinition{
		DefinitionBase: app.DefinitionBase{
			BeanClass: "test2",
		},
	})

	s.Assert().Equal(len(s.r.GetAllBeans()), 2, "beans num must gt 0")
}

func (s *registerSuite) TestBasicBeansGetBeans() {
	s.r.RegisterBeanDefinition(app.BeanDefinition{
		DefinitionBase: app.DefinitionBase{
			BeanClass: "test1",
		},
	})
	s.r.RegisterBeanDefinition(app.BeanDefinition{
		DefinitionBase: app.DefinitionBase{
			BeanClass: "test2",
		},
	})

	b := s.r.GetBeanDefinition("test2")
	s.Assert().NotNil(b)
}

func TestBeanRegistry(t *testing.T) {
	suite.Run(t, &registerSuite{
		r: &app.BeanRegistry,
	})
}
