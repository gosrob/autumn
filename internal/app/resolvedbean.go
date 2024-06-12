package app

type beanResolver interface {
	GetDecl() string
	GetType() string
	GetAlias() string
	RandomName()
	GetConstructText() string
}

type beanResolve struct {
	decl string
	tp   string
}

// GetAlias implements beanResolver.
func (b *beanResolve) GetAlias() string {
	panic("unimplemented")
}

// GetConstructText implements beanResolver.
func (b *beanResolve) GetConstructText() string {
	panic("unimplemented")
}

// RandomName implements beanResolver.
func (b *beanResolve) RandomName() {
	panic("unimplemented")
}

// GetDecl implements beanResolver.
func (b *beanResolve) GetDecl() string {
	return b.decl
}

// GetType implements beanResolver.
func (b *beanResolve) GetType() string {
	return b.tp
}

var (
	_            beanResolver = (*beanResolve)(nil)
	BeanResolver beanResolve
)
