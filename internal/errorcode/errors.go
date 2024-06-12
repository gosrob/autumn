package errorcode

import "fmt"

type Error struct {
	Msg   string
	Code  int64
	inner error
}

// Error implements error.
func (e *Error) Error() string {
	return fmt.Sprintf("%s \n -----> inner error: %s", e.Msg, e.inner)
}

func (e *Error) Instance() *Error {
	ins := Error{
		Msg:   e.Msg,
		Code:  e.Code,
		inner: e.inner,
	}
	return &ins
}

func (e *Error) Printf(format string, params ...any) *Error {
	e.Msg += fmt.Sprintf(" : "+format, params)
	return e
}

var _ error = (*Error)(nil)

var CastError Error = Error{
	Msg:   "cast fail",
	Code:  0,
	inner: nil,
}

var FactoryCanOnlyProduceOneError Error = Error{
	Msg:   "bean factory can only produce one",
	Code:  1,
	inner: nil,
}
