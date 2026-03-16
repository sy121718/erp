package main

import (
	"fmt"
	"os"

	"erp-server/internal/handler/router"
	"erp-server/pkg/captcha"
	"erp-server/pkg/config"
	"erp-server/pkg/database"
	"erp-server/pkg/jwt"
	"erp-server/pkg/log"
	"erp-server/pkg/redis"
	"erp-server/pkg/sign"
	"erp-server/pkg/token"
	"erp-server/pkg/validator"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 获取运行环境，默认为 dev
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// 1. 初始化配置
	if err := config.Init(env); err != nil {
		fmt.Printf("初始化配置失败: %v\n", err)
		os.Exit(1)
	}
	cfg := config.Get()

	// 2. 初始化日志
	if err := log.Init(&cfg.Log); err != nil {
		fmt.Printf("初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	// 3. 初始化验证器
	validator.Init()

	// 4. 初始化数据库
	if err := database.Init(&cfg         .Database); err != nil {
		log.Fatal("初始化数据库失败", zap.Error(err))
	}
	defer database.Close()

	// 5. 初始化Redis
	if err := redis.Init(&cfg.Redis); err != nil {
		log.Fatal("初始化Redis失败", zap.Error(err))
	}
	defer redis.Close()

	// 6. 初始化JWT
	jwt.Init(&cfg.JWT)

	// 7. 初始化Token存储
	token.Init(&cfg.Redis)

	// 8. 初始化验证码
	captcha.Init(&captcha.Config{
		Length:     cfg.Captcha.GetLength(),
		ExpireTime: cfg.Captcha.GetExpireDuration(),
		Width:      cfg.Captcha.GetWidth(),
		Height:     cfg.Captcha.GetHeight(),
	})

	// 9. 初始化签名
	sign.Init(&cfg.Sign)

	// 10. 设置 Gin 运行模式
	gin.SetMode(cfg.Server.Mode)

	// 11. 初始化路由
	r := router.SetupRouter()

	// 12. 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info(fmt.Sprintf("服务启动，监听端口 %s", addr))
	if err := r.Run(addr); err != nil {
		log.Fatal("服务启动失败", zap.Error(err))
	}
}