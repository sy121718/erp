package admin

import (
	service "erp-server/internal/application/admin"
	"erp-server/internal/handler/middleware"
	adminRepo "erp-server/internal/infrastructure/repository/admin"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册管理员模块路由
func RegisterRoutes(r *gin.RouterGroup) {
	// 初始化管理员模块
	repo := adminRepo.NewAdminRepository()
	svc := service.NewAdminService(repo)
	handler := NewHandler(svc)

	// 公开接口（无需认证）
	r.POST("/login", handler.Login)
	r.GET("/captcha", handler.GetCaptcha)
	r.POST("/refresh-token", handler.RefreshToken)

	// 需要认证的接口
	auth := r.Group("")
	auth.Use(middleware.Auth())

	{
		// 当前管理员操作
		auth.POST("/logout", handler.Logout)
		auth.GET("/profile", handler.GetProfile)
		auth.POST("/password/change", handler.ChangePassword)

		// 管理员管理
		auth.GET("/list", handler.List)
		auth.GET("/:id", handler.GetProfile)
		auth.POST("/create", handler.Create)
		auth.POST("/update/:id", handler.Update)
		auth.POST("/delete/:id", handler.Delete)
		auth.POST("/ban/:id", handler.Ban)
		auth.POST("/unban/:id", handler.Unban)

		// 重置密码（仅超管）
		auth.POST("/password/reset/:id", handler.ResetPassword)

		// 强制下线（仅超管）
		auth.POST("/force-logout/:id", handler.ForceLogout)
	}
}