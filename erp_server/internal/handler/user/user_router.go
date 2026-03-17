package user

import (
	service "erp-server/internal/application/user"
	"erp-server/internal/handler/middleware"
	userRepo "erp-server/internal/infrastructure/repository/user"

	"github.com/gin-gonic/gin"
)

var handler *Handler

func initHandler() {
	if handler != nil {
		return
	}
	repo := userRepo.NewUserRepository()
	svc := service.NewUserService(repo)
	handler = NewHandler(svc)
}

// RegisterRoutes 注册用户模块路由（用户侧）
func RegisterRoutes(r *gin.RouterGroup) {
	initHandler()

	// 公开接口（无需认证）
	r.GET("/captcha", handler.GetCaptcha)
	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.POST("/refresh-token", handler.RefreshToken)

	// 需要认证的接口（用户身份）
	auth := r.Group("")
	auth.Use(middleware.Auth(), middleware.UserOnly())
	{
		auth.POST("/logout", handler.Logout)
		auth.GET("/profile", handler.GetProfile)
		auth.POST("/update", handler.UpdateProfile)
		auth.POST("/password/change", handler.ChangePassword)
	}
}

// RegisterAdminRoutes 注册用户管理路由（管理员侧）
func RegisterAdminRoutes(r *gin.RouterGroup) {
	initHandler()

	// 管理员认证
	auth := r.Group("")
	auth.Use(middleware.Auth(), middleware.AdminTypeOnly())
	{
		auth.GET("/list", handler.AdminListUsers)
		auth.GET("/:id", handler.AdminGetUser)
		auth.POST("/update/:id", handler.AdminUpdateUser)
		auth.POST("/ban/:id", handler.AdminBanUser)
		auth.POST("/unban/:id", handler.AdminUnbanUser)
		auth.POST("/password/reset/:id", handler.AdminResetPassword)
	}
}
