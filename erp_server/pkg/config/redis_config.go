package config

import "time"

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`      // Redis地址
	Port     int    `mapstructure:"port"`      // Redis端口
	Password string `mapstructure:"password"`  // Redis密码
	DB       int    `mapstructure:"db"`        // Redis数据库
	PoolSize int    `mapstructure:"pool_size"` // 连接池大小
}

// GetAddr 获取Redis地址
func (c *RedisConfig) GetAddr() string {
	if c.Host == "" {
		return "127.0.0.1:6379"
	}
	return c.Host
}

// GetPassword 获取密码
func (c *RedisConfig) GetPassword() string {
	return c.Password
}

// GetDB 获取数据库
func (c *RedisConfig) GetDB() int {
	return c.DB
}

// GetPoolSize 获取连接池大小
func (c *RedisConfig) GetPoolSize() int {
	if c.PoolSize <= 0 {
		return 10
	}
	return c.PoolSize
}

// GetRefreshTokenTTL 获取refresh_token过期时间（7天）
func (c *RedisConfig) GetRefreshTokenTTL() time.Duration {
	return 7 * 24 * time.Hour
}

// GetAccessTokenTTL 获取access_token过期时间（5分钟）
func (c *RedisConfig) GetAccessTokenTTL() time.Duration {
	return 5 * time.Minute
}

// DefaultRedisConfig 默认Redis配置
func DefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "",
		DB:       0,
		PoolSize: 10,
	}
}