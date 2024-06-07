package app

type beanResolver interface {
	Resolve(n any)
}

type beanResolve struct{}

// Resolve implements beanResolver.
func (b *beanResolve) Resolve(n any) {
	panic("unimplemented")
}

var (
	_            beanResolver = (*beanResolve)(nil)
	BeanResolver beanResolve
)
