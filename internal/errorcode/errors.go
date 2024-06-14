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

var DepulicatedBeanAliasError Error = Error{
	Msg:   "bean alias has already set",
	Code:  1,
	inner: nil,
}

var CreateZeroBeanError Error = Error{
	Msg:   "bean can not create, because there is no definitio to create it",
	Code:  2,
	inner: nil,
}

var BeanPrimaryError Error = Error{
	Msg:   "if bean has more then one definition, there must be only one primary bean",
	Code:  3,
	inner: nil,
}

var BeanNotFindError Error = Error{
	Msg:   "all beans has be set, but can not found wire bean",
	Code:  4,
	inner: nil,
}

var MetaInfoNotDefined Error = Error{
	Msg:   "must define one metainfo",
	Code:  5,
	inner: nil,
}
