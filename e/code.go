package e

const (
	SuccessCode       = 10000
	InvalidReqCode    = 10010
	InternalErrorCode = 10020

	UnauthorizedCode = 401
	ForbiddenCode    = 403
)

var codeMsg map[int]string

func init() {
	codeMsg = make(map[int]string)
	codeMsg[SuccessCode] = "成功"
	codeMsg[InvalidReqCode] = "无效参数"
	codeMsg[InternalErrorCode] = "服务内部错误"
	codeMsg[UnauthorizedCode] = "未认证"
	codeMsg[ForbiddenCode] = "无权限"
}

//SetCodeMsg 设置错误码消息
func SetCodeMsg(code int, msg string) {
	codeMsg[code] = msg
}

//GetCodeMsg 获取错误码消息
func GetCodeMsg(code int) string {
	if msg, ok := codeMsg[code]; ok {
		return msg
	}
	return "未知错误"
}
