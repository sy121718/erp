package user

// 预定义用户错误
var (
	ErrUserNotFound         = NewUserError(404, "用户不存在")
	ErrUsernameExists       = NewUserError(20001, "用户名已存在")
	ErrPasswordIncorrect    = NewUserError(20002, "用户名或密码错误")
	ErrAccountDisabled      = NewUserError(20003, "账户已禁用")
	ErrOldPasswordIncorrect = NewUserError(20004, "原密码错误")
	ErrCannotBanSelf        = NewUserError(20005, "不能操作自己")
	ErrCaptchaInvalid       = NewUserError(20006, "验证码错误")
)

// UserError 用户错误
type UserError struct {
	code    int
	message string
}

// NewUserError 创建用户错误
func NewUserError(code int, message string) *UserError {
	return &UserError{code: code, message: message}
}

func (e *UserError) Error() string {
	return e.message
}

// Code 返回错误码
func (e *UserError) Code() int {
	return e.code
}

// Message 返回错误消息
func (e *UserError) Message() string {
	return e.message
}
