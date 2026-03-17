package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	"/api/admin/logout",
	"/api/user/login",
	"/api/user/register",
	"/api/user/refresh-token",
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
		userType := claims.UserType
		if userType == "" {
			userType = jwt.UserTypeAdmin
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("name", claims.Name)
		c.Set("is_admin", claims.IsAdmin)
		c.Set("user_type", userType)

		// 4. 检查token是否需要续期
		if jwt.Get().ShouldRefreshToken(claims) {
			newToken, _, err := jwt.Get().GenerateAccessToken(
				claims.UserID,
				claims.Username,
				claims.Name,
				claims.IsAdmin,
				userType,
			)
			if err == nil {
				// 在响应头中返回新的token
				c.Header("X-New-Access-Token", newToken)
			}
		}

		c.Next()
	}
}

// AdminOnly 仅超级管理员可访问
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

// UserOnly 仅普通用户可访问
func UserOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType := GetUserType(c)
		if userType != jwt.UserTypeUser {
			c.Error(errors.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}

// AdminTypeOnly 仅管理员类型可访问（包括超管和普通管理员）
func AdminTypeOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType := GetUserType(c)
		if userType != jwt.UserTypeAdmin {
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
		// GET 请求：从 query 参数获取
		query := c.Request.URL.Query()
		for k, v := range query {
			if len(v) > 0 {
				params[k] = v[0]
			}
		}
	} else if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
		// POST/PUT/DELETE 请求：从 body 读取参数
		// 注意：需要先读取 body，然后重新设置回去
		bodyBytes, err := c.GetRawData()
		if err == nil && len(bodyBytes) > 0 {
			// 将 body 重新设置回去，供后续处理使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// 解析 JSON body
			var bodyMap map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &bodyMap); err == nil {
				// 递归展平嵌套对象
				flattenMap(bodyMap, "", params)
			}
		}
	}

	// 验证签名
	return sign.Get().Verify(params, timestamp, nonce, signature)
}

// flattenMap 递归展平嵌套的 map
func flattenMap(obj map[string]interface{}, prefix string, result map[string]string) {
	for key, value := range obj {
		newKey := prefix
		if prefix != "" {
			newKey += "."
		}
		newKey += key

		// 跳过 nil 值
		if value == nil {
			continue
		}

		// 如果是嵌套对象（map）
		if nested, ok := value.(map[string]interface{}); ok {
			flattenMap(nested, newKey, result)
		} else {
			// 转换为字符串
			result[newKey] = fmt.Sprintf("%v", value)
		}
	}
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

// GetUserType 从上下文获取用户类型
func GetUserType(c *gin.Context) string {
	userType, exists := c.Get("user_type")
	if !exists {
		return ""
	}
	return userType.(string)
}

// AddToWhiteList 添加路由到白名单
func AddToWhiteList(paths ...string) {
	whiteList = append(whiteList, paths...)
}
