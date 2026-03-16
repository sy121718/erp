package router

import (
	adminHandler "erp-server/internal/handler/admin"
	"erp-server/internal/handler/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 注册全局中间件
	// Gin 中间件执行顺序（洋葱模型，从外到内）：
	// 请求: Recovery -> Logger -> CORS -> ErrorHandler -> Handler
	// 响应: Handler -> ErrorHandler -> CORS -> Logger -> Recovery
	// 
	// ErrorHandler 必须在最内层（最后注册），这样才能在 Handler 之后执行
	// Logger 必须在 ErrorHandler 之外，这样才能在 ErrorHandler 之后读取状态码
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API 路由组
	api := r.Group("/api")
	{
		// 注册管理员模块路由
		admin := api.Group("/admin")
		adminHandler.RegisterRoutes(admin)
	}

	return r
}