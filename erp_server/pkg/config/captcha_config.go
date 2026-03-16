package config

import "time"

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	Length     int `mapstructure:"length"`      // 验证码长度
	ExpireTime int `mapstructure:"expire_time"` // 过期时间(秒)
	Width      int `mapstructure:"width"`       // 图片宽度
	Height     int `mapstructure:"height"`      // 图片高度
}

// GetExpireDuration 获取过期时间duration
func (c *CaptchaConfig) GetExpireDuration() time.Duration {
	if c.ExpireTime <= 0 {
		return 5 * time.Minute
	}
	return time.Duration(c.ExpireTime) * time.Second
}

// GetLength 获取验证码长度
func (c *CaptchaConfig) GetLength() int {
	if c.Length <= 0 {
		return 6
	}
	return c.Length
}

// GetWidth 获取图片宽度
func (c *CaptchaConfig) GetWidth() int {
	if c.Width <= 0 {
		return 120
	}
	return c.Width
}

// GetHeight 获取图片高度
func (c *CaptchaConfig) GetHeight() int {
	if c.Height <= 0 {
		return 40
	}
	return c.Height
}

// DefaultCaptchaConfig 默认验证码配置
func DefaultCaptchaConfig() *CaptchaConfig {
	return &CaptchaConfig{
		Length:     6,
		ExpireTime: 300,
		Width:      120,
		Height:     40,
	}
}