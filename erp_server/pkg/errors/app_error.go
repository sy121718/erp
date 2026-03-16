package errors

import (
	"fmt"
)

// AppError 应用错误类型
type AppError struct {
	Code    int    // HTTP 状态码
	ErrCode int    // 业务错误码
	Message string // 错误消息
	Err     error  // 原始错误
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.ErrCode, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.ErrCode, e.Message)
}

// Unwrap 返回原始错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithError 添加原始错误
func (e *AppError) WithError(err error) *AppError {
	e.Err = err
	return e
}

// New 创建新错误
func New(httpCode, errCode int, message string) *AppError {
	return &AppError{
		Code:    httpCode,
		ErrCode: errCode,
		Message: message,
	}
}

// NewBadRequest 400 错误
func NewBadRequest(errCode int, message string) *AppError {
	return New(400, errCode, message)
}

// NewUnauthorized 401 错误
func NewUnauthorized(errCode int, message string) *AppError {
	return New(401, errCode, message)
}

// NewForbidden 403 错误
func NewForbidden(errCode int, message string) *AppError {
	return New(403, errCode, message)
}

// NewNotFound 404 错误
func NewNotFound(errCode int, message string) *AppError {
	return New(404, errCode, message)
}

// NewInternalError 500 错误
func NewInternalError(errCode int, message string) *AppError {
	return New(500, errCode, message)
}

// 预定义通用错误
var (
	// 400 Bad Request
	ErrInvalidParams    = NewBadRequest(400001, "请求参数错误")
	ErrInvalidJSON      = NewBadRequest(400002, "JSON格式错误")
	ErrMissingField     = NewBadRequest(400003, "缺少必填字段")
	ErrInvalidField     = NewBadRequest(400004, "字段格式错误")

	// 401 Unauthorized
	ErrUnauthorized     = NewUnauthorized(401001, "未授权，请先登录")
	ErrTokenExpired     = NewUnauthorized(401002, "Token已过期")
	ErrTokenInvalid     = NewUnauthorized(401003, "Token无效")
	ErrSignatureInvalid = NewUnauthorized(401004, "签名验证失败")

	// 403 Forbidden
	ErrForbidden        = NewForbidden(403001, "无权访问")
	ErrNotAdmin         = NewForbidden(403002, "需要管理员权限")

	// 404 Not Found
	ErrNotFound         = NewNotFound(404001, "资源不存在")
	ErrUserNotFound     = NewNotFound(404002, "用户不存在")

	// 500 Internal Server Error
	ErrInternal         = NewInternalError(500001, "服务器内部错误")
	ErrDatabase         = NewInternalError(500002, "数据库错误")
	ErrRedis            = NewInternalError(500003, "Redis错误")
)

// IsAppError 检查是否为 AppError
func IsAppError(err error) (*AppError, bool) {
	if e, ok := err.(*AppError); ok {
		return e, true
	}
	return nil, false
}