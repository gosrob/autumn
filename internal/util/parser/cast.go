package parser

import (
	"encoding/json"
	"reflect"

	"github.com/gosrob/autumn/internal/errorcode"
	"github.com/gosrob/autumn/internal/util/stream"
)

func Cast[T any](a Annotation) (T, error) {
	var t T
	name := reflect.TypeOf(t).Name()
	if a.Name() == name {
		err := mapToStruct(a.Params(), &t)
		if err == nil {
			return t, nil
		}
	}
	return t, &errorcode.CastError
}

func mapToStruct[T any](m any, dst *T) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func Castable[T any]() func(a annotation) bool {
	return func(a annotation) bool {
		if _, err := Cast[T](a); err != nil {
			return false
		}
		return true
	}
}

func FindAnnotations[T any](as []annotation) []T {
	return stream.Map(
		stream.OfSlice(as).
			Filter(Castable[T]()),
		func(t annotation) T {
			a, _ := Cast[T](t)
			return a
		}).ToSlice()
}
