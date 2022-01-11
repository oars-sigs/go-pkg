package e

type Error struct {
	msg    string
	detail string
}

func (e *Error) Error() string {
	if e.detail == "" {
		return e.msg
	}
	return e.detail
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Set(err error) error {
	e.detail = err.Error()
	return e
}

func NewError(msg string, err ...error) *Error {
	e := &Error{msg: msg}
	if len(err) > 0 && err[0] != nil {
		e.detail = err[0].Error()
	}
	return e
}
