package app

import (
	"fmt"
	"strings"

	"github.com/gosrob/autumn/internal/util"
	"github.com/gosrob/autumn/internal/util/nodeutil"
)

type beanResolver interface {
	GetDecl() string
	GetType() string
	GetAlias() string
	RandomName()
	GetConstructText() string
	GetDefinitionBase() DefinitionBase
	SetResolved(fieldName string)
	IsFiledResolved(fieldName string) bool
}

type beanResolve struct {
	decl           string
	tp             string
	initialContent string
	bd             DefinitionBase
	resolved       map[string]any
}

// IsFiledResolved implements beanResolver.
func (b *beanResolve) IsFiledResolved(fieldName string) bool {
	_, ok := b.resolved[fieldName]

	return ok
}

// SetResolved implements beanResolver.
func (b *beanResolve) SetResolved(fieldName string) {
	b.resolved[fieldName] = nil
}

// GetDefinitionBase implements beanResolver.
func (b *beanResolve) GetDefinitionBase() DefinitionBase {
	return b.bd
}

func NewBeanResolveConstructor(decl string, tp string, bd DefinitionBase) *beanResolve {
	b := &beanResolve{
		decl:           decl,
		tp:             nodeutil.GetShortIdentityByFullPath(tp),
		initialContent: "%s := %s{}\n",
		bd:             bd,
		resolved:       map[string]any{},
	}
	b.RandomName()
	return b
}

func NewBeanResolveFactory(decl string, fn string, bd DefinitionBase, isPointer bool, params ...string) *beanResolve {
	p := strings.Join(params, ",")
	starExp := ""
	if isPointer {
		starExp = "*"
	}
	b := &beanResolve{
		decl:           decl,
		tp:             fn,
		initialContent: fmt.Sprintf("%s := %s%s(%s) \n", "%s", starExp, "%s", p),
		bd:             bd,
		resolved:       map[string]any{},
	}
	b.RandomName()
	return b
}

// GetAlias implements beanResolver.
func (b *beanResolve) GetAlias() string {
	panic("unimplemented")
}

// GetConstructText implements beanResolver.
func (b *beanResolve) GetConstructText() string {
	return fmt.Sprintf(b.initialContent, b.decl, b.tp)
}

// RandomName implements beanResolver.
func (b *beanResolve) RandomName() {
	b.decl = b.decl + "_" + util.RandomStr(5)
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
