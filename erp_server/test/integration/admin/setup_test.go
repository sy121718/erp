package admin_test

import (
	"time"

	"erp-server/internal/handler/router"
	"erp-server/pkg/config"
	"erp-server/pkg/database"
	"erp-server/pkg/jwt"
	"erp-server/pkg/log"
	"erp-server/pkg/redis"
	"erp-server/pkg/token"
	"erp-server/pkg/validator"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// 测试用的全局变量
var (
	testRouter   *gin.Engine
	adminToken   string
	refreshToken string
	testAdminID  int64
)

// 超级管理员配置
const (
	superAdminUsername = "sky"
	superAdminPassword = "123456"
	superAdminName     = "超级管理员"
)

// TestMain 测试入口
func TestMain(m *testing.M) {
	// 设置测试模式
	gin.SetMode(gin.TestMode)

	// 初始化测试环境
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	// 初始化配置
	if err := config.Init(env); err != nil {
		panic(err)
	}

	// 初始化日志
	if err := log.Init(&config.Get().Log); err != nil {
		panic(err)
	}

	// 初始化验证器
	validator.Init()

	// 初始化数据库
	if err := database.Init(&config.Get().Database); err != nil {
		panic(err)
	}

	// 初始化Redis（测试环境可能没有Redis，跳过）
	if err := redis.Init(&config.Get().Redis); err != nil {
		log.Warn("Redis初始化失败，部分功能可能不可用", zap.Error(err))
	}

	// 初始化JWT
	jwt.Init(&config.Get().JWT)

	// 初始化Token存储
	token.Init(&config.Get().Redis)

	// 初始化路由
	testRouter = router.SetupRouter()

	// 确保超级管理员存在且密码正确
	ensureSuperAdmin()

	// 运行测试
	code := m.Run()

	// 清理
	cleanupTestData()
	database.Close()
	redis.Close()
	log.Sync()

	os.Exit(code)
}

// ensureSuperAdmin 确保超级管理员存在且密码正确
func ensureSuperAdmin() {
	db := database.Get()

	// 检查超级管理员是否存在
	var count int64
	db.Raw("SELECT COUNT(*) FROM sys_admin WHERE username = ?", superAdminUsername).Scan(&count)

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(superAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	now := time.Now()

	if count == 0 {
		// 创建超级管理员
		result := db.Exec(`
			INSERT INTO sys_admin (username, password, name, status, is_admin, create_by, update_by, create_time, update_time)
			VALUES (?, ?, ?, 1, 1, 1, 1, ?, ?)
		`, superAdminUsername, string(hashedPassword), superAdminName, now, now)

		if result.Error != nil {
			panic(result.Error)
		}
	} else {
		// 更新密码和重置登录失败次数
		db.Exec(`
			UPDATE sys_admin 
			SET password = ?, login_failure_count = 0, locked_until_time = NULL, update_time = ?
			WHERE username = ?
		`, string(hashedPassword), now, superAdminUsername)
	}
}

// cleanupTestData 清理测试数据（不删除超级管理员 sky）
func cleanupTestData() {
	db := database.Get()
	// 只删除测试创建的用户，不删除超级管理员 sky
	if testAdminID > 0 {
		db.Exec("DELETE FROM sys_admin WHERE id = ? AND username != ?", testAdminID, superAdminUsername)
	}
	// 清理测试过程中创建的其他测试用户
	db.Exec("DELETE FROM sys_admin WHERE username LIKE 'test%' AND username != ?", superAdminUsername)
	db.Exec("DELETE FROM sys_admin WHERE username = 'temp_delete_test'")
	db.Exec("DELETE FROM sys_admin WHERE username LIKE 'special_%'")
	db.Exec("DELETE FROM sys_admin WHERE username LIKE 'xss_%'")
}