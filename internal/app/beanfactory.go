package app

type BeanFactory interface{}

type ListableBeanFactory interface {
	BeanFactory
}

var ListableBeanFactoryer ListableBeanFactory
