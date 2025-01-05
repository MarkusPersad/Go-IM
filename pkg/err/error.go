package err

var (
	BadRequest   = NewError(1000, "参数错误")
	TokenInvalid = NewError(1001, "Token无效")
	Forbidden    = NewError(1002, "无权限")
	UserExists   = NewError(1003, "用户已存在")
	CheckCode    = NewError(1004, "验证码错误")
	Timeout      = NewError(1005, "登陆超时")
	NotFound     = NewError(1006, "未找到")
	PassError    = NewError(1007, "密码错误")
)

type PersonalError struct {
	Code    int
	Message string
}

func (e *PersonalError) Error() string {
	return e.Message
}

func NewError(code int, msg string) *PersonalError {
	return &PersonalError{
		Code:    code,
		Message: msg,
	}
}
