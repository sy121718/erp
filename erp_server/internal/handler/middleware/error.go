package middleware

import (
	"fmt"

	"erp-server/pkg/errors"
	"erp-server/pkg/log"
	"erp-server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler 全局错误处理中间件
// 使用 Gin 的 c.Error() 方式处理错误
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			handleGinErrors(c)
		}
	}
}

// handleGinErrors 处理 Gin 的错误列表
func handleGinErrors(c *gin.Context) {
	// 获取最后一个错误（通常是业务错误）
	err := c.Errors.Last().Err

	// 检查是否为 AppError
	if appErr, ok := errors.IsAppError(err); ok {
		handleAppError(c, appErr)
		return
	}

	// 其他未知错误
	log.Error("未知错误", zap.Error(err))
	c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Response{
		Code: 500001,
		Msg:  "服务器内部错误",
	})
}

// handleAppError 处理应用错误
func handleAppError(c *gin.Context, err *errors.AppError) {
	switch err.Code {
	case 400:
		log.Warn("参数错误", zap.Int("err_code", err.ErrCode), zap.String("msg", err.Message))
		c.AbortWithStatusJSON(http.StatusBadRequest, utils.Response{
			Code: err.ErrCode,
			Msg:  err.Message,
		})

	case 401:
		log.Warn("未授权", zap.Int("err_code", err.ErrCode), zap.String("msg", err.Message))
		c.AbortWithStatusJSON(http.StatusUnauthorized, utils.Response{
			Code: err.ErrCode,
			Msg:  err.Message,
		})

	case 403:
		log.Warn("禁止访问", zap.Int("err_code", err.ErrCode), zap.String("msg", err.Message))
		c.AbortWithStatusJSON(http.StatusForbidden, utils.Response{
			Code: err.ErrCode,
			Msg:  err.Message,
		})

	case 404:
		log.Warn("资源不存在", zap.Int("err_code", err.ErrCode), zap.String("msg", err.Message))
		c.AbortWithStatusJSON(http.StatusNotFound, utils.Response{
			Code: err.ErrCode,
			Msg:  err.Message,
		})

	default:
		log.Error("服务器错误", zap.Int("err_code", err.ErrCode), zap.String("msg", err.Message), zap.Error(err.Err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, utils.Response{
			Code: err.ErrCode,
			Msg:  err.Message,
		})
	}
}

// BindJSON 绑定 JSON 参数，失败时返回错误
func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.Error(errors.ErrInvalidParams.WithError(err))
		return false
	}
	return true
}

// BindQuery 绑定 Query 参数，失败时返回错误
func BindQuery(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		c.Error(errors.ErrInvalidParams.WithError(err))
		return false
	}
	return true
}

// ParseID 解析路径参数 ID，失败时返回错误
func ParseID(c *gin.Context) (int64, bool) {
	idStr := c.Param("id")
	var id int64
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil || id <= 0 {
		c.Error(errors.ErrInvalidField.WithError(fmt.Errorf("ID格式错误: %s", idStr)))
		return 0, false
	}
	return id, true
}