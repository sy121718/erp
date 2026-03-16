package config

import "time"

// JWTConfig JWT配置
type JWTConfig struct {
	Secret           string `mapstructure:"secret"`            // 签名密钥
	AccessTokenTime  int    `mapstructure:"access_token_time"` // access_token过期时间(分钟)
	RefreshTokenTime int    `mapstructure:"refresh_token_time"` // refresh_token过期时间(天)
	Issuer           string `mapstructure:"issuer"`            // 签发者
}

// GetAccessTokenDuration 获取access_token过期时间
func (c *JWTConfig) GetAccessTokenDuration() time.Duration {
	if c.AccessTokenTime <= 0 {
		return 5 * time.Minute // 默认5分钟
	}
	return time.Duration(c.AccessTokenTime) * time.Minute
}

// GetRefreshTokenDuration 获取refresh_token过期时间
func (c *JWTConfig) GetRefreshTokenDuration() time.Duration {
	if c.RefreshTokenTime <= 0 {
		return 7 * 24 * time.Hour // 默认7天
	}
	return time.Duration(c.RefreshTokenTime) * 24 * time.Hour
}

// GetSecret 获取密钥，如果为空则使用默认值
func (c *JWTConfig) GetSecret() string {
	if c.Secret == "" {
		return "erp-default-jwt-secret-key-2024"
	}
	return c.Secret
}

// GetIssuer 获取签发者
func (c *JWTConfig) GetIssuer() string {
	if c.Issuer == "" {
		return "erp-server"
	}
	return c.Issuer
}

// DefaultJWTConfig 默认JWT配置
func DefaultJWTConfig() *JWTConfig {
	return &JWTConfig{
		Secret:           "erp-default-jwt-secret-key-2024",
		AccessTokenTime:  5,  // 5分钟
		RefreshTokenTime: 7,  // 7天
		Issuer:           "erp-server",
	}
}