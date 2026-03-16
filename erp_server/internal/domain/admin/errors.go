package admin

// 预定义管理员错误
var (
	ErrAdminNotFound         = NewAdminError(404, "管理员不存在")
	ErrUsernameExists        = NewAdminError(10001, "用户名已存在")
	ErrPasswordIncorrect     = NewAdminError(10002, "用户名或密码错误")
	ErrAccountLocked         = NewAdminError(10003, "账户已锁定")
	ErrAccountDisabled       = NewAdminError(10004, "账户已禁用")
	ErrCaptchaInvalid        = NewAdminError(10005, "验证码错误")
	ErrOldPasswordIncorrect  = NewAdminError(10006, "原密码错误")
	ErrCannotDeleteSelf      = NewAdminError(10007, "不能删除自己")
	ErrCannotDisableSelf     = NewAdminError(10008, "不能禁用自己")
	ErrCannotModifyAdmin     = NewAdminError(10009, "无权修改超级管理员")
	ErrTokenInvalid          = NewAdminError(10010, "Token无效或已过期")
)

// AdminError 管理员错误
type AdminError struct {
	code    int
	message string
}

// NewAdminError 创建管理员错误
func NewAdminError(code int, message string) *AdminError {
	return &AdminError{code: code, message: message}
}

func (e *AdminError) Error() string {
	return e.message
}

// Code 返回错误码
func (e *AdminError) Code() int {
	return e.code
}

// Message 返回错误消息
func (e *AdminError) Message() string {
	return e.message
}