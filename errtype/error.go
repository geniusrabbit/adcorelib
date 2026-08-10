//
// @project GeniusRabbit corelib
//

package errtype

// Error is a reusable sentinel error that supports WithMessage while remaining
// compatible with errors.Is against the bare sentinel.
type Error string

// Error implements the error interface.
func (e Error) Error() string { return string(e) }

// WithMessage attaches detail text (typically cause.Error()).
// An empty message returns the sentinel itself.
func (e Error) WithMessage(msg string) error {
	if msg == "" {
		return e
	}
	return &withMessage{base: e, msg: msg}
}

type withMessage struct {
	base error
	msg  string
}

func (e *withMessage) Error() string   { return e.base.Error() + ": " + e.msg }
func (e *withMessage) Unwrap() error   { return e.base }
func (e *withMessage) Message() string { return e.msg }
