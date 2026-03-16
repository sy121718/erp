package middleware

import (
	"strconv"
	"strings"

	"erp-server/pkg/errors"
	"erp-server/pkg/jwt"
	"erp-server/pkg/log"
	"erp-server/pkg/sign"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 白名单路由（不需要认证的路径）
var whiteList = []string{
	"/health",
	"/api/ping",
	"/api/admin/login",
	"/api/admin/captcha",
	"/api/admin/refresh-token",
}

// Auth 认证中间件（JWT + 签名验证）
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 检查是否在白名单中
		if isWhiteListed(path) {
			c.Next()
			return
		}

		// 1. 验证签名（测试模式下跳过）
		if gin.Mode() != gin.TestMode {
			if !verifySign(c) {
				c.Error(errors.ErrSignatureInvalid)
				c.Abort()
				return
			}
		}

		// 2. 验证JWT Token
		token := getTokenFromHeader(c)
		if token == "" {
			c.Error(errors.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := jwt.Get().ParseToken(token)
		if err != nil {
			log.Warn("JWT解析失败", zap.Error(err))
			c.Error(errors.ErrTokenInvalid.WithError(err))
			c.Abort()
			return
		}

		// 3. 将用户信息存入上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("name", claims.Name)
		c.Set("is_admin", claims.IsAdmin)

		c.Next()
	}
}

// AdminOnly 仅管理员可访问
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			c.Error(errors.ErrNotAdmin)
			c.Abort()
			return
		}
		c.Next()
	}
}

// isWhiteListed 检查路径是否在白名单中
func isWhiteListed(path string) bool {
	for _, p := range whiteList {
		if path == p || strings.HasPrefix(path, p+"/") {
			return true
		}
	}
	return false
}

// getTokenFromHeader 从Header获取Token
func getTokenFromHeader(c *gin.Context) string {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		return ""
	}

	// Bearer token格式
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}

	return auth
}

// verifySign 验证签名
func verifySign(c *gin.Context) bool {
	// 从Header获取签名信息
	timestampStr := c.GetHeader("X-Timestamp")
	nonce := c.GetHeader("X-Nonce")
	signature := c.GetHeader("X-Sign")

	if timestampStr == "" || nonce == "" || signature == "" {
		return false
	}

	// 解析时间戳
	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return false
	}

	// 获取请求参数
	params := make(map[string]string)
	if c.Request.Method == "GET" {
		query := c.Request.URL.Query()
		for k, v := range query {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
	}

	// 验证签名
	return sign.Get().Verify(params, timestamp, nonce, signature)
}

// GetUserID 从上下文获取用户ID
func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(int64)
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}

// GetIsAdmin 从上下文获取是否管理员
func GetIsAdmin(c *gin.Context) bool {
	isAdmin, exists := c.Get("is_admin")
	if !exists {
		return false
	}
	return isAdmin.(bool)
}

// AddToWhiteList 添加路由到白名单
func AddToWhiteList(paths ...string) {
	whiteList = append(whiteList, paths...)
}