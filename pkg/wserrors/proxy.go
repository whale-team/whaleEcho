package wserrors

import (
	stderrors "errors"
	"fmt"

	"github.com/pkg/errors"
)

var Errorf = errors.Errorf

// Wrap is as the proxy for github.com/pkg/errors.Wrap func.
// func Wrap(err error, message string) error {
// 	return errors.Wrap(err, message)
// }
var Wrap = errors.Wrap

// Wrapf is as the proxy for github.com/pkg/errors.Wrapf func.
// func Wrapf(err error, format string, args ...interface{}) error {
// 	return errors.Wrapf(err, format, args...)
// }
var Wrapf = errors.Wrapf

// // WithMessage is as the proxy for github.com/pkg/errors.WithMessage func.
// // func WithMessage(err error, message string) error {
// // 	return errors.WithMessage(err, message)
// // }
// var WithMessage = errors.WithMessage

// // WithMessagef is as the proxy for github.com/pkg/errors.WithMessagef func.
// // func WithMessagef(err error, format string, args ...interface{}) error {
// // 	return errors.WithMessagef(err, format, args...)
// // }
// var WithMessagef = errors.WithMessagef

func WithMessagef(err error, message string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*WsError)
	if !ok {
		return New(Internal, fmt.Sprintf(message, args...))
	}
	return New(_err.Status, fmt.Sprintf(message, args...))
}

// Cause is as the proxy for github.com/pkg/errors.Cause func.
// func Cause(err error) error {
// 	return errors.Cause(err)
// }
var Cause = errors.Cause

// WithStack is as the proxy for github.com/pkg/errors.WithStack func.
// func WithStack(err error) error {
// 	return errors.WithStack(err)
// }
var WithStack = errors.WithStack

// Is reports whether any error in err's chain matches target.
// The chain consists of err itself followed by the sequence of errors obtained by repeatedly calling Unwrap.
// An error is considered to match a target if it is equal to that target or if it implements a method Is(error) bool such that Is(target) returns true.
// var Is = errors.Is

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func Is(err, target error) bool { return stderrors.Is(err, target) }

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target interface{}) bool { return stderrors.As(err, target) }

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}
