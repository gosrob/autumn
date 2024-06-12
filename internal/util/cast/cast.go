package cast

import "github.com/gosrob/autumn/internal/errorcode"

func Cast[T any](node any) (T, error) {
	var zero T
	if v, ok := node.(T); ok {
		return v, nil
	}

	return zero, errorcode.CastError.Instance()
}
