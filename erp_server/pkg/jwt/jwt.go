package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"erp-server/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired     = errors.New("token已过期")
	ErrTokenInvalid     = errors.New("token无效")
	ErrTokenMalformed   = errors.New("token格式错误")
	ErrTokenNotValidYet = errors.New("token尚未生效")
)

// UserType 常量
const (
	UserTypeAdmin = "admin" // 管理员
	UserTypeUser  = "user"  // 普通用户
)

// Claims 自定义Claims (access_token使用)
type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"is_admin"`
	UserType string `json:"user_type"` // "admin" 或 "user"
	TokenID  string `json:"jti"`
	jwt.RegisteredClaims
}

// RefreshClaims 刷新token的Claims (refresh_token使用)
type RefreshClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	UserType string `json:"user_type"`
	TokenID  string `json:"jti"`
	jwt.RegisteredClaims
}

// JWT JWT工具
type JWT struct {
	config *config.JWTConfig
}

var jwtInstance *JWT

// Init 初始化JWT
func Init(cfg *config.JWTConfig) {
	jwtInstance = &JWT{config: cfg}
}

// Get 获取JWT实例
func Get() *JWT {
	if jwtInstance == nil {
		jwtInstance = &JWT{
			config: config.DefaultJWTConfig(),
		}
	}
	return jwtInstance
}

// generateTokenID 生成唯一的Token ID
func generateTokenID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// GenerateAccessToken 生成access_token (短期，5分钟)
func (j *JWT) GenerateAccessToken(userID int64, username, name string, isAdmin bool, userType string) (tokenString, tokenID string, err error) {
	if userType == "" {
		userType = UserTypeAdmin
	}
	tokenID = generateTokenID()
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Name:     name,
		IsAdmin:  isAdmin,
		UserType: userType,
		TokenID:  tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(now.Add(j.config.GetAccessTokenDuration())),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.config.GetIssuer(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(j.config.GetSecret()))
	return
}

// GenerateRefreshToken 生成refresh_token (长期，7天)
func (j *JWT) GenerateRefreshToken(userID int64, username string, isAdmin bool, userType string) (tokenString, tokenID string, err error) {
	if userType == "" {
		userType = UserTypeAdmin
	}
	tokenID = generateTokenID()
	now := time.Now()
	claims := RefreshClaims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		UserType: userType,
		TokenID:  tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(now.Add(j.config.GetRefreshTokenDuration())),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.config.GetIssuer(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString([]byte(j.config.GetSecret()))
	return
}

// GenerateToken 生成token (兼容旧接口)
func (j *JWT) GenerateToken(userID int64, username, name string, isAdmin bool) (string, error) {
	tokenString, _, err := j.GenerateAccessToken(userID, username, name, isAdmin, UserTypeAdmin)
	return tokenString, err
}

// ParseToken 解析access_token
func (j *JWT) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.GetSecret()), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// ParseRefreshToken 解析refresh_token
func (j *JWT) ParseRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.config.GetSecret()), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		}
		if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// GetTokenIDFromAccessToken 从access_token获取token_id
func (j *JWT) GetTokenIDFromAccessToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.TokenID, nil
}

// GetTokenIDFromRefreshToken 从refresh_token获取token_id
func (j *JWT) GetTokenIDFromRefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseRefreshToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.TokenID, nil
}

// RefreshToken 使用refresh_token刷新，返回新的access_token和refresh_token
func (j *JWT) RefreshToken(refreshTokenString string) (accessToken, newRefreshToken string, err error) {
	claims, err := j.ParseRefreshToken(refreshTokenString)
	if err != nil {
		return "", "", err
	}

	userType := claims.UserType
	if userType == "" {
		userType = UserTypeAdmin
	}

	accessToken, _, err = j.GenerateAccessToken(claims.UserID, claims.Username, "", claims.IsAdmin, userType)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, _, err = j.GenerateRefreshToken(claims.UserID, claims.Username, claims.IsAdmin, userType)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

// GenerateRandomSecret 生成随机密钥
func GenerateRandomSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// ShouldRefreshToken 检查token是否需要续期
// 如果token剩余时间少于总有效期的30%，则应该续期
func (j *JWT) ShouldRefreshToken(claims *Claims) bool {
	if claims.ExpiresAt == nil {
		return false
	}

	// 计算剩余时间
	now := time.Now()
	expiresAt := claims.ExpiresAt.Time
	remaining := expiresAt.Sub(now)

	// 计算总有效期
	totalDuration := j.config.GetAccessTokenDuration()

	// 如果剩余时间小于总有效期的30%，则需要续期
	return remaining < time.Duration(float64(totalDuration.Seconds())*0.3)
}

// RefreshAccessToken 为现有的claims生成新的access_token
func (j *JWT) RefreshAccessToken(claims *Claims) (string, error) {
	tokenID := generateTokenID()
	now := time.Now()
	userType := claims.UserType
	if userType == "" {
		userType = UserTypeAdmin
	}
	newClaims := Claims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Name:     claims.Name,
		IsAdmin:  claims.IsAdmin,
		UserType: userType,
		TokenID:  tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID,
			ExpiresAt: jwt.NewNumericDate(now.Add(j.config.GetAccessTokenDuration())),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    j.config.GetIssuer(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return token.SignedString([]byte(j.config.GetSecret()))
}