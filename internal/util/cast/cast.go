package cast

import (
	"github.com/gosrob/autumn/internal/errorcode"
	"github.com/gosrob/autumn/internal/util/nodeutil"
)

func Cast[T any](node any) (T, error) {
	var zero T
	if v, ok := node.(T); ok {
		return v, nil
	}

	return zero, errorcode.CastError.DeepCopy().Printf(nodeutil.GetTypeFromAny(node), nodeutil.GetTypeFromAny(zero))
}
