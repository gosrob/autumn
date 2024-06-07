package astutil

import (
	"go/ast"

	"github.com/gosrob/autumn/internal/errorcode"
)

func AstCast[T any](node ast.Node) (T, error) {
	var zero T
	if v, ok := node.(T); ok {
		return v, nil
	}

	return zero, errorcode.CastError.Instance().Printf("cannot cast to ast type")
}

func Cast[T any](node any) (T, error) {
	var zero T
	if v, ok := node.(T); ok {
		return v, nil
	}

	return zero, errorcode.CastError.Instance()
}
