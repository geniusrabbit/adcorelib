//
// @project GeniusRabbit corelib
//

package errtype

import "fmt"

type WithMessageError interface {
	error
	WithMessage(msg string) WithMessageError
	WithMessageFmt(format string, args ...any) WithMessageError
}

// Error is a reusable sentinel error that supports WithMessage while remaining
// compatible with errors.Is against the bare sentinel.
type Error string

// Error implements the error interface.
func (e Error) Error() string { return string(e) }

// WithMessage attaches detail text (typically cause.Error()).
// An empty message returns the sentinel itself.
func (e Error) WithMessage(msg string) WithMessageError {
	if msg == "" {
		return e
	}
	return &withMessage{base: e, msg: msg}
}

// WithMessageFmt attaches detail text (typically cause.Error()).
// An empty message returns the sentinel itself.
func (e Error) WithMessageFmt(format string, args ...any) WithMessageError {
	return e.WithMessage(fmt.Sprintf(format, args...))
}

type withMessage struct {
	base error
	msg  string
}

func (e *withMessage) Error() string   { return e.base.Error() + ": " + e.msg }
func (e *withMessage) Unwrap() error   { return e.base }
func (e *withMessage) Message() string { return e.msg }

func (e *withMessage) WithMessage(msg string) WithMessageError {
	if msg == "" {
		return e
	}
	return &withMessage{base: e.base, msg: e.msg + ": " + msg}
}

func (e *withMessage) WithMessageFmt(format string, args ...any) WithMessageError {
	return e.WithMessage(fmt.Sprintf(format, args...))
}
