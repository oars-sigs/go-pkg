package perr

type Error struct {
	code   int
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

func (e *Error) Code() int {
	if e.code == 0 {
		return -1
	}
	return e.code
}

func (e *Error) Set(err error) *Error {
	e.detail = err.Error()
	return e
}

func (e *Error) SetCode(code int) *Error {
	e.code = code
	return e
}

func New(msg string) *Error {
	return &Error{msg: msg}
}

var (
	ErrInvalidReq   = New("无效请求参数").SetCode(10010)
	ErrInternal     = New("服务内部错误").SetCode(10020)
	ErrUnkown       = New("未知错误").SetCode(10030)
	ErrUnauthorized = New("未认证").SetCode(401)
	ErrForbidden    = New("无权限").SetCode(403)
	ErrUn           = New("未实名").SetCode(10040)
)
