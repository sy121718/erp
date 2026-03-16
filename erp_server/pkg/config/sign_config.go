package config

// SignConfig 签名配置
type SignConfig struct {
	Secret     string `mapstructure:"secret"`      // 签名密钥
	ExpireTime int    `mapstructure:"expire_time"` // 签名有效期(秒)
}

// GetSecret 获取签名密钥
func (c *SignConfig) GetSecret() string {
	if c.Secret == "" {
		return "erp-default-sign-secret-key-2024"
	}
	return c.Secret
}

// GetExpireTime 获取签名有效期(秒)
func (c *SignConfig) GetExpireTime() int {
	if c.ExpireTime <= 0 {
		return 300 // 默认5分钟
	}
	return c.ExpireTime
}

// DefaultSignConfig 默认签名配置
func DefaultSignConfig() *SignConfig {
	return &SignConfig{
		Secret:     "erp-default-sign-secret-key-2024",
		ExpireTime: 300,
	}
}