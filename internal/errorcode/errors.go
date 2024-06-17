package errorcode

import (
	"fmt"
	"strings"
)

// Error structure
type Error struct {
	Msg   string
	Code  int64
	inner error
}

// Constructor for Error
func NewError(msg string, code int64, inner error) *Error {
	return &Error{
		Msg:   msg,
		Code:  code,
		inner: inner,
	}
}

// Error implements error interface
func (e *Error) Error() string {
	parts := []string{fmt.Sprintf("Error Code: %d, Message: %s", e.Code, e.Msg)}

	// Collect inner error messages
	for inner := e.inner; inner != nil; {
		if innerErr, ok := inner.(*Error); ok && innerErr != nil {
			parts = append(parts, fmt.Sprintf("Caused by - Code: %d, Message: %s", innerErr.Code, innerErr.Msg))
			inner = innerErr.inner
		} else {
			parts = append(parts, fmt.Sprintf("Caused by - %s", inner.Error()))
			break
		}
	}

	// Reverse the collected parts to display innermost (root) error on top
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}

	return strings.Join(parts, "\n")
}

// DeepCopy creates a deep copy of the error instance
func (e *Error) DeepCopy() *Error {
	var copiedInner error
	if e.inner != nil {
		inner, ok := e.inner.(*Error)
		if ok {
			copiedInner = inner.DeepCopy()
		} else {
			copiedInner = e.inner // non-Error inner error, not copying
		}
	}

	return &Error{
		Msg:   e.Msg,
		Code:  e.Code,
		inner: copiedInner,
	}
}

// APrintf adds formatted message to the error
func (e *Error) APrintf(format string, params ...any) *Error {
	newMsg := fmt.Sprintf(format, params...)
	return NewError(fmt.Sprintf("%s: %s", e.Msg, newMsg), e.Code, e.inner)
}

func (e *Error) Printf(params ...any) *Error {
	e.Msg = fmt.Sprintf(e.Msg, params...)
	return e
}

// Variable to verify that Error implements the error interface
var _ error = (*Error)(nil)

var CastError Error = Error{
	Msg:   "cast fail, from %s to %s",
	Code:  0,
	inner: nil,
}

var FactoryCanOnlyProduceOneError Error = Error{
	Msg:   "bean factory [%s] method must produce 1 bean, but there produce ",
	Code:  1,
	inner: nil,
}

var DepulicatedBeanAliasError Error = Error{
	Msg:   "bean alias [%s] has already set",
	Code:  1,
	inner: nil,
}

var CreateZeroBeanError Error = Error{
	Msg:   "[%s] bean can not create, because there is no definitio to create it",
	Code:  2,
	inner: nil,
}

var BeanPrimaryError Error = Error{
	Msg:   "[%s] if bean has more then one definition, there must be only one primary bean",
	Code:  3,
	inner: nil,
}

var BeanNotFindError Error = Error{
	Msg:   "[%s] all beans has be set, but can not found wire bean",
	Code:  4,
	inner: nil,
}

var MetaInfoNotDefined Error = Error{
	Msg:   "must define one metainfo",
	Code:  5,
	inner: nil,
}

var CannotInjectToArray Error = Error{
	Msg:   "attributes is not array, but beans is array",
	Code:  6,
	inner: nil,
}
