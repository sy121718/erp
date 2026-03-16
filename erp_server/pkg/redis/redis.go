package redis

import (
	"context"
	"fmt"

	"erp-server/pkg/config"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

// Init 初始化Redis连接
func Init(cfg *config.RedisConfig) error {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// 测试连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis连接失败: %w", err)
	}

	return nil
}

// Get 获取Redis客户端
func Get() *redis.Client {
	if client == nil {
		panic("Redis未初始化")
	}
	return client
}

// Close 关闭Redis连接
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}