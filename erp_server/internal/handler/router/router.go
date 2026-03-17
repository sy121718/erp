package router

import (
	adminHandler "erp-server/internal/handler/admin"
	"erp-server/internal/handler/middleware"
	userHandler "erp-server/internal/handler/user"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	api := r.Group("/api")
	{
		// 管理员模块
		admin := api.Group("/admin")
		adminHandler.RegisterRoutes(admin)

		// 管理员-用户管理模块
		adminUser := admin.Group("/user")
		userHandler.RegisterAdminRoutes(adminUser)

		// 用户模块
		user := api.Group("/user")
		userHandler.RegisterRoutes(user)
	}

	return r
}